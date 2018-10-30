package v1alpha1

import (
	v1alpha1 "github.com/zenjoy/hcloud-floating-ip-operator/apis/hcloud/v1alpha1"
	scheme "github.com/zenjoy/hcloud-floating-ip-operator/client/k8s/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// FloatingIPPoolsGetter has a method to return a FloatingIPPoolInterface.
// A group's client should implement this interface.
type FloatingIPPoolsGetter interface {
	FloatingIPPools() FloatingIPPoolInterface
}

// FloatingIPPoolInterface has methods to work with FloatingIPPool resources.
type FloatingIPPoolInterface interface {
	Create(*v1alpha1.FloatingIPPool) (*v1alpha1.FloatingIPPool, error)
	Update(*v1alpha1.FloatingIPPool) (*v1alpha1.FloatingIPPool, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha1.FloatingIPPool, error)
	List(opts v1.ListOptions) (*v1alpha1.FloatingIPPoolList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.FloatingIPPool, err error)
	FloatingIPPoolExpansion
}

// floatingIPPools implements FloatingIPPoolInterface
type floatingIPPools struct {
	client rest.Interface
}

// newFloatingIPPools returns a FloatingIPPools
func newFloatingIPPools(c *HcloudV1alpha1Client) *floatingIPPools {
	return &floatingIPPools{
		client: c.RESTClient(),
	}
}

// Get takes name of the floatingIPPool, and returns the corresponding floatingIPPool object, and an error if there is any.
func (c *floatingIPPools) Get(name string, options v1.GetOptions) (result *v1alpha1.FloatingIPPool, err error) {
	result = &v1alpha1.FloatingIPPool{}
	err = c.client.Get().
		Resource("floatingippools").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of FloatingIPPools that match those selectors.
func (c *floatingIPPools) List(opts v1.ListOptions) (result *v1alpha1.FloatingIPPoolList, err error) {
	result = &v1alpha1.FloatingIPPoolList{}
	err = c.client.Get().
		Resource("floatingippools").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested floatingIPPools.
func (c *floatingIPPools) Watch(opts v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Resource("floatingippools").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a floatingIPPool and creates it.  Returns the server's representation of the floatingIPPool, and an error, if there is any.
func (c *floatingIPPools) Create(floatingIPPool *v1alpha1.FloatingIPPool) (result *v1alpha1.FloatingIPPool, err error) {
	result = &v1alpha1.FloatingIPPool{}
	err = c.client.Post().
		Resource("floatingippools").
		Body(floatingIPPool).
		Do().
		Into(result)
	return
}

// Update takes the representation of a floatingIPPool and updates it. Returns the server's representation of the floatingIPPool, and an error, if there is any.
func (c *floatingIPPools) Update(floatingIPPool *v1alpha1.FloatingIPPool) (result *v1alpha1.FloatingIPPool, err error) {
	result = &v1alpha1.FloatingIPPool{}
	err = c.client.Put().
		Resource("floatingippools").
		Name(floatingIPPool.Name).
		Body(floatingIPPool).
		Do().
		Into(result)
	return
}

// Delete takes name of the floatingIPPool and deletes it. Returns an error if one occurs.
func (c *floatingIPPools) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Resource("floatingippools").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *floatingIPPools) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	return c.client.Delete().
		Resource("floatingippools").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched floatingIPPool.
func (c *floatingIPPools) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.FloatingIPPool, err error) {
	result = &v1alpha1.FloatingIPPool{}
	err = c.client.Patch(pt).
		Resource("floatingippools").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
