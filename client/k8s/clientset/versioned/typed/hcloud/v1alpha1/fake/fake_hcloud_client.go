package fake

import (
	v1alpha1 "github.com/zenjoy/hcloud-floating-ip-operator/client/k8s/clientset/versioned/typed/hcloud/v1alpha1"
	rest "k8s.io/client-go/rest"
	testing "k8s.io/client-go/testing"
)

type FakeHcloudV1alpha1 struct {
	*testing.Fake
}

func (c *FakeHcloudV1alpha1) FloatingIPPools() v1alpha1.FloatingIPPoolInterface {
	return &FakeFloatingIPPools{c}
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *FakeHcloudV1alpha1) RESTClient() rest.Interface {
	var ret *rest.RESTClient
	return ret
}
