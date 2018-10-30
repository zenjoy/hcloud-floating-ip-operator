package v1alpha1

import (
	v1alpha1 "github.com/zenjoy/hcloud-floating-ip-operator/apis/hcloud/v1alpha1"
	"github.com/zenjoy/hcloud-floating-ip-operator/client/k8s/clientset/versioned/scheme"
	serializer "k8s.io/apimachinery/pkg/runtime/serializer"
	rest "k8s.io/client-go/rest"
)

type HcloudV1alpha1Interface interface {
	RESTClient() rest.Interface
	FloatingIPPoolsGetter
}

// HcloudV1alpha1Client is used to interact with features provided by the hcloud.zenjoy.be group.
type HcloudV1alpha1Client struct {
	restClient rest.Interface
}

func (c *HcloudV1alpha1Client) FloatingIPPools() FloatingIPPoolInterface {
	return newFloatingIPPools(c)
}

// NewForConfig creates a new HcloudV1alpha1Client for the given config.
func NewForConfig(c *rest.Config) (*HcloudV1alpha1Client, error) {
	config := *c
	if err := setConfigDefaults(&config); err != nil {
		return nil, err
	}
	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}
	return &HcloudV1alpha1Client{client}, nil
}

// NewForConfigOrDie creates a new HcloudV1alpha1Client for the given config and
// panics if there is an error in the config.
func NewForConfigOrDie(c *rest.Config) *HcloudV1alpha1Client {
	client, err := NewForConfig(c)
	if err != nil {
		panic(err)
	}
	return client
}

// New creates a new HcloudV1alpha1Client for the given RESTClient.
func New(c rest.Interface) *HcloudV1alpha1Client {
	return &HcloudV1alpha1Client{c}
}

func setConfigDefaults(config *rest.Config) error {
	gv := v1alpha1.SchemeGroupVersion
	config.GroupVersion = &gv
	config.APIPath = "/apis"
	config.NegotiatedSerializer = serializer.DirectCodecFactory{CodecFactory: scheme.Codecs}

	if config.UserAgent == "" {
		config.UserAgent = rest.DefaultKubernetesUserAgent()
	}

	return nil
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *HcloudV1alpha1Client) RESTClient() rest.Interface {
	if c == nil {
		return nil
	}
	return c.restClient
}
