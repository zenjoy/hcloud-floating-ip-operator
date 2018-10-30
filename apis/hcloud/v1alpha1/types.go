package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type FloatingIPPool struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec FloatinIPPoolSpec `json:"spec"`
}

// FloatinIPPoolSpec defines a floating ip resource
type FloatinIPPoolSpec struct {
	// Floating IP from Hetzner that will be assigned to nodes matching the
	// nodeSelector
	Ips []string `json:"ips"`

	// Query to select a pool of nodes that
	NodeSelector map[string]string `json:"nodeSelector"`

	// Frequency for reconcilation loops
	IntervalSeconds Seconds `json:"intervalSeconds,omitempty"`
}

// Seconds is an duration in seconds
type Seconds int64

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type FloatingIPPoolList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []FloatingIPPool `json:"items"`
}
