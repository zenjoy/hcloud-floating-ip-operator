package operator

import (
	"context"
	"fmt"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"

	hcloudv1alpha1 "github.com/zenjoy/hcloud-floating-ip-operator/apis/hcloud/v1alpha1"
	"github.com/zenjoy/hcloud-floating-ip-operator/pkg/log"
	"github.com/zenjoy/hcloud-floating-ip-operator/pkg/service"
)

// Handler is the floating ip assignment handler that will handle the
// events received from kubernetes.
type handler struct {
	service service.Syncer
	logger  log.Logger
}

// newHandler returns a new handler.
func newHandler(k8sCli kubernetes.Interface, hcloudCli *hcloud.Client, logger log.Logger) *handler {
	return &handler{
		service: service.NewService(k8sCli, hcloudCli, logger),
		logger:  logger,
	}
}

// Add will ensure that the required assigner is running.
func (h *handler) Add(ctx context.Context, obj runtime.Object) error {
	fip, ok := obj.(*hcloudv1alpha1.FloatingIPPool)
	if !ok {
		return fmt.Errorf("%v is not a floating ip object", obj.GetObjectKind())
	}

	return h.service.EnsureFloatingIPPool(fip)
}

// Delete will ensure the reuited pod terminator is not running.
func (h *handler) Delete(ctx context.Context, name string) error {
	return h.service.DeleteFloatingIPPool(name)
}
