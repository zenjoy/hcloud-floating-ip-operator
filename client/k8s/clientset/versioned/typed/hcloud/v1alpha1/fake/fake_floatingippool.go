package fake

import (
	v1alpha1 "github.com/zenjoy/hcloud-floating-ip-operator/apis/hcloud/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeFloatingIPPools implements FloatingIPPoolInterface
type FakeFloatingIPPools struct {
	Fake *FakeHcloudV1alpha1
}

var floatingippoolsResource = schema.GroupVersionResource{Group: "hcloud.zenjoy.be", Version: "v1alpha1", Resource: "floatingippools"}

var floatingippoolsKind = schema.GroupVersionKind{Group: "hcloud.zenjoy.be", Version: "v1alpha1", Kind: "FloatingIPPool"}

// Get takes name of the floatingIPPool, and returns the corresponding floatingIPPool object, and an error if there is any.
func (c *FakeFloatingIPPools) Get(name string, options v1.GetOptions) (result *v1alpha1.FloatingIPPool, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(floatingippoolsResource, name), &v1alpha1.FloatingIPPool{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.FloatingIPPool), err
}

// List takes label and field selectors, and returns the list of FloatingIPPools that match those selectors.
func (c *FakeFloatingIPPools) List(opts v1.ListOptions) (result *v1alpha1.FloatingIPPoolList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(floatingippoolsResource, floatingippoolsKind, opts), &v1alpha1.FloatingIPPoolList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.FloatingIPPoolList{}
	for _, item := range obj.(*v1alpha1.FloatingIPPoolList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested floatingIPPools.
func (c *FakeFloatingIPPools) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(floatingippoolsResource, opts))
}

// Create takes the representation of a floatingIPPool and creates it.  Returns the server's representation of the floatingIPPool, and an error, if there is any.
func (c *FakeFloatingIPPools) Create(floatingIPPool *v1alpha1.FloatingIPPool) (result *v1alpha1.FloatingIPPool, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(floatingippoolsResource, floatingIPPool), &v1alpha1.FloatingIPPool{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.FloatingIPPool), err
}

// Update takes the representation of a floatingIPPool and updates it. Returns the server's representation of the floatingIPPool, and an error, if there is any.
func (c *FakeFloatingIPPools) Update(floatingIPPool *v1alpha1.FloatingIPPool) (result *v1alpha1.FloatingIPPool, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(floatingippoolsResource, floatingIPPool), &v1alpha1.FloatingIPPool{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.FloatingIPPool), err
}

// Delete takes name of the floatingIPPool and deletes it. Returns an error if one occurs.
func (c *FakeFloatingIPPools) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteAction(floatingippoolsResource, name), &v1alpha1.FloatingIPPool{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeFloatingIPPools) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(floatingippoolsResource, listOptions)

	_, err := c.Fake.Invokes(action, &v1alpha1.FloatingIPPoolList{})
	return err
}

// Patch applies the patch and returns the patched floatingIPPool.
func (c *FakeFloatingIPPools) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.FloatingIPPool, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(floatingippoolsResource, name, data, subresources...), &v1alpha1.FloatingIPPool{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.FloatingIPPool), err
}
