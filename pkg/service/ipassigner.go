package service

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"reflect"
	"sync"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"

	"github.com/hetznercloud/hcloud-go/hcloud"

	hcloudv1alpha1 "github.com/zenjoy/hcloud-floating-ip-operator/apis/hcloud/v1alpha1"
	"github.com/zenjoy/hcloud-floating-ip-operator/pkg/log"
)

const (
	// MinimalIntervalSeconds is the lower bound for refreshing floating ips
	MinimalIntervalSeconds hcloudv1alpha1.Seconds = 5
)

// TimeWrapper is a wrapper around time so it can be mocked
type TimeWrapper interface {
	// After is the same as time.After
	After(d time.Duration) <-chan time.Time
	// Now is the same as Now.
	Now() time.Time
}

type timeStd struct{}

func (t *timeStd) After(d time.Duration) <-chan time.Time { return time.After(d) }
func (t *timeStd) Now() time.Time                         { return time.Now() }

// IPAssigner will verify ip assignment at regular intervals.
type IPAssigner struct {
	fip       *hcloudv1alpha1.FloatingIPPool
	k8sCli    kubernetes.Interface
	hcloudCli *hcloud.Client
	logger    log.Logger
	time      TimeWrapper

	running bool
	mutex   sync.Mutex
	stopC   chan struct{}
}

// NewIPAssigner returns a new ip assigner.
func NewIPAssigner(fip *hcloudv1alpha1.FloatingIPPool, k8sCli kubernetes.Interface, hcloudCli *hcloud.Client, logger log.Logger) *IPAssigner {
	return &IPAssigner{
		fip:       fip,
		k8sCli:    k8sCli,
		hcloudCli: hcloudCli,
		logger:    logger,
		time:      &timeStd{},
	}
}

// NewCustomIPAssigner is a constructor that lets you customize everything on the object construction.
func NewCustomIPAssigner(fip *hcloudv1alpha1.FloatingIPPool, k8sCli kubernetes.Interface, hcloudCli *hcloud.Client, time TimeWrapper, logger log.Logger) *IPAssigner {
	return &IPAssigner{
		fip:       fip,
		k8sCli:    k8sCli,
		hcloudCli: hcloudCli,
		logger:    logger,
		time:      time,
	}
}

// SameSpec checks if the ip assigner has the same spec.
func (p *IPAssigner) SameSpec(fip *hcloudv1alpha1.FloatingIPPool) bool {
	return reflect.DeepEqual(p.fip.Spec, fip.Spec)
}

// Start will run the ip assigner at regular intervals.
func (p *IPAssigner) Start() error {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if p.running {
		return fmt.Errorf("already running")
	}

	p.stopC = make(chan struct{})
	p.running = true

	go func() {
		p.logger.Infof("started %s ip assigner", p.fip.Name)
		if err := p.run(); err != nil {
			p.logger.Errorf("error executing ip assigner: %s", err)
		}
	}()

	return nil
}

// Stop stops the ip assigner.
func (p *IPAssigner) Stop() error {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if p.running {
		close(p.stopC)
		p.logger.Infof("stopped %s ip assigner", p.fip.Name)
	}

	p.running = false
	return nil
}

// run will run the loop that will kill eventually the required pods.
func (p *IPAssigner) run() error {
	for {
		select {
		case <-p.time.After(time.Duration(max(p.fip.Spec.IntervalSeconds, MinimalIntervalSeconds)) * time.Second):
			if err := p.assign(); err != nil {
				p.logger.Errorf("error assigning ip: %s", err)
			}
		case <-p.stopC:
			return nil
		}
	}
}

// asign will verify current assignment of the floating ip and change
// assignment to a node matching the nodeSelector in case the floating
// ip is currently not correctly assigned
func (p *IPAssigner) assign() error {
	// Get all probable targets.
	nodes, err := p.getProbableNodes()
	if err != nil {
		return err
	}

	total := len(nodes.Items)
	if total == 0 {
		p.logger.Errorf("0 nodes probable targets")
		return fmt.Errorf("%s ip assigner: 0 nodes probable targets", p.fip.Name)
	}

	// Get nodes in random order.
	targets, targetNames := p.getRandomNodes(nodes)

	// Get available Hetzner IPs
	hetznerIps, err := p.findHCloudFloatingIps()
	if err != nil {
		return err
	}

	// Get all Hetzner servers #TODO: support pagination or filter to get only nodes in probableNodes?
	hetznerServers, err := p.findServers()
	if err != nil {
		return err
	}
	var hetznerServersByID = make(map[int]*hcloud.Server)
	var hetznerServersByName = make(map[string]*hcloud.Server)

	for i := range hetznerServers {
		server := hetznerServers[i]
		hetznerServersByID[server.ID] = server
		hetznerServersByName[server.Name] = server
	}

	var hetznerIpsToAssign = make([]*hcloud.FloatingIP, 0)
	var assignedServers = make([]string, 0)
	var assignments = map[string]([]*hcloud.FloatingIP){}

	// Check which ips needs to be assigned
	for i := range hetznerIps {
		if hetznerIps[i].Server == nil {
			p.logger.Infof("%s ip %s is not assigned to any node", p.fip.Name, hetznerIps[i].IP.String())
			hetznerIpsToAssign = append(hetznerIpsToAssign, hetznerIps[i])
		} else {
			var found = false
			var serverName = "---"
			var server = hetznerServersByID[hetznerIps[i].Server.ID]

			if server != nil {
				serverName = server.Name
			}

			for j := range targets {
				if targets[j].Name == serverName {
					found = true
				}
			}
			if found {
				assignedServers = append(assignedServers, serverName)
				assignments[serverName] = append(assignments[serverName], hetznerIps[i])
			} else {
				p.logger.Infof("%s ip %s is assigned to unknown node", p.fip.Name, hetznerIps[i].IP.String())
				hetznerIpsToAssign = append(hetznerIpsToAssign, hetznerIps[i])
			}
		}
	}

	if (len(hetznerIps)-len(hetznerIpsToAssign)) >= len(targets) &&
		len(assignments) < len(targets) &&
		len(assignments) < (len(hetznerIps)-len(hetznerIpsToAssign)) {
		i := 0
		j := 0
		nbIpsToReassign := ((len(hetznerIps) - len(hetznerIpsToAssign)) - len(assignments))

		p.logger.Infof("%s ips are not equally spread over possible targets", p.fip.Name)

		for i < nbIpsToReassign && j < 100 { // prevent endless loop
			for k, v := range assignments {
				if len(v) > 1 {
					ip, v := v[0], v[1:]
					assignments[k] = v
					hetznerIpsToAssign = append(hetznerIpsToAssign, ip)
					p.logger.Infof("%s ip %s will be reassigned", p.fip.Name, ip.IP.String())
					i++
				}
				if i >= nbIpsToReassign {
					break
				}
			}
			j++
		}
	}

	unassignedServers := p.difference(targetNames, assignedServers)

	for i := range hetznerIpsToAssign {
		var nodeName string

		// first use any servers with no ip assigned yet
		if len(unassignedServers) > 0 {
			nodeName, unassignedServers = unassignedServers[0], unassignedServers[1:]
		} else {
			nodeName = targetNames[i%len(targetNames)]
		}

		server := hetznerServersByName[nodeName]

		_, _, err = p.hcloudCli.FloatingIP.Assign(context.TODO(), hetznerIpsToAssign[i], server)
		if err != nil {
			return err
		}

		p.logger.Infof("%s ip %s assigned to node %s", p.fip.Name, hetznerIpsToAssign[i].IP.String(), nodeName)
	}

	return nil
}

// Gets all the pods filtered that can be a target of termination.
func (p *IPAssigner) getProbableNodes() (*corev1.NodeList, error) {
	set := labels.Set(p.fip.Spec.NodeSelector)
	slc := set.AsSelector()
	opts := metav1.ListOptions{
		LabelSelector: slc.String(),
	}
	return p.k8sCli.CoreV1().Nodes().List(opts)
}

// getRandomNodes will shuffle the list of available nodes.
func (p *IPAssigner) getRandomNodes(nodes *corev1.NodeList) ([]corev1.Node, []string) {
	items := nodes.Items
	nodeNames := make([]string, len(items))

	r := rand.New(rand.NewSource(time.Now().Unix()))
	// We start at the end of the items, inserting our random
	// values one at a time.
	for n := len(items); n > 0; n-- {
		randIndex := r.Intn(n)
		// We swap the value at index n-1 and the random index
		// to move our randomly chosen value to the end of the
		// items, and to move the value that was at n-1 into our
		// unshuffled portion of the items.
		items[n-1], items[randIndex] = items[randIndex], items[n-1]
	}

	for i := range items {
		nodeNames[i] = items[i].Name
	}

	return items, nodeNames
}

// findHCloudFloatingIP will return a hcloud FloatingIP resource that matches
// the ip specified in the FloatingIP CRD resource
func (p *IPAssigner) findHCloudFloatingIps() ([]*hcloud.FloatingIP, error) {
	ips := make([]net.IP, len(p.fip.Spec.Ips))

	for i := range p.fip.Spec.Ips {
		ips[i] = net.ParseIP(p.fip.Spec.Ips[i])
		if ips[i] == nil {
			return nil, fmt.Errorf("error parsing ip from spec: %s", p.fip.Spec.Ips[i])
		}
	}

	fips, err := p.hcloudCli.FloatingIP.All(context.TODO())
	if err != nil {
		return nil, err
	}

	hetznerIps := make([]*hcloud.FloatingIP, len(ips))

	for i := range fips {
		for j := range ips {
			if fips[i].IP.Equal(ips[j]) {
				hetznerIps[j] = fips[i]
			}
		}
	}

	for i := range hetznerIps {
		if hetznerIps[i] == nil {
			return nil, fmt.Errorf("ip %s does not match any floating ip resource", ips[i].String())
		}
	}

	return hetznerIps, nil
}

func (p *IPAssigner) findServer(nodeName string) (*hcloud.Server, error) {
	server, _, err := p.hcloudCli.Server.GetByName(context.TODO(), nodeName)

	return server, err
}

func (p *IPAssigner) findServers() ([]*hcloud.Server, error) {
	servers, err := p.hcloudCli.Server.All(context.TODO())

	return servers, err
}

func (p *IPAssigner) difference(slice1 []string, slice2 []string) []string {
	var diff []string

	// Loop two times, first to find slice1 strings not in slice2,
	// second loop to find slice2 strings not in slice1
	for i := 0; i < 2; i++ {
		for _, s1 := range slice1 {
			found := false
			for _, s2 := range slice2 {
				if s1 == s2 {
					found = true
					break
				}
			}
			// String not found. We add it to return slice
			if !found {
				diff = append(diff, s1)
			}
		}
		// Swap the slices, only if it was the first loop
		if i == 0 {
			slice1, slice2 = slice2, slice1
		}
	}

	return diff
}
