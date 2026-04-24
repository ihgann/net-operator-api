// Copyright (c) 2020-2026 Broadcom. All Rights Reserved.
// Broadcom Confidential. The term "Broadcom" refers to Broadcom Inc.
// and/or its subsidiaries.

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// IPAMDisabledAnnotationKeyName is the name of the annotation added to
// GatewayClass resources that do not participate in net-operator's IPAM.
// The value does not need to be truthy; the presence of the key is what
// disables net-operator's IPAM for that GatewayClass.
const IPAMDisabledAnnotationKeyName = "netoperator.vmware.com/ipam-disabled"

type IPPoolUsageLabelValue string

const (
	// IPPoolUsageLabelKeyName is the name of a label used to indicate how IP pools
	// should be used. To create an affinity, you must create a NetworkInterface with a
	// label matching the intended use. For example, if you create a NetworkInterface with
	// a label matching netoperator.vmware.com/ipam-usage=vip, then net operator
	// will only provision from IPPools matching that label falling back to the general
	// pool if needed unless IPPoolUsageAnnotationStrictKeyName is set.
	IPPoolUsageLabelKeyName = "netoperator.vmware.com/ipam-usage"

	// IPPoolUsageAnnotationStrictKeyName indicates that an interface should not attempt
	// to retrieve IPPools meant for general purpose consumption. For example, if "vip" is set,
	// only IPPools matching the "vip" label will be used and "general" will not be used as a pool.
	IPPoolUsageAnnotationStrictKeyName = "netoperator.vmware.com/ipam-strict-usage"

	// IPPoolUsageLabelGeneralValue indicates an IP pool can be used for any purpose.
	// If a usage label is omitted from an IPPool, this value is implied.
	IPPoolUsageLabelGeneralValue IPPoolUsageLabelValue = "general"

	// IPPoolUsageLabelVIPValue indicates an IP pool is reserved for a NetworkInterface
	// which provisions virtual IP addresses.
	IPPoolUsageLabelVIPValue IPPoolUsageLabelValue = "vip"
)

type IPPoolConditionType string

const (
	// IPPoolFull condition is added when no more IPs are free in the pool.
	IPPoolFull IPPoolConditionType = "full"
	// IPPoolReady condition is added when IPPool has been realized.
	IPPoolReady IPPoolConditionType = "ready"
	// IPPoolFail condition is added when an error was encountered in realizing.
	IPPoolFail IPPoolConditionType = "failure"
)

// IPPoolCondition describes the state of a IPPool at a certain point.
type IPPoolCondition struct {
	// type is the type of IPPool condition.
	// +required
	//nolint:kubeapilinter // Stable v1alpha1: KAL maxlength (wire shape unchanged).
	Type IPPoolConditionType `json:"type"`
	// status is the status of the condition.
	// Can be True, False, Unknown.
	// +required
	//nolint:kubeapilinter // Stable v1alpha1: KAL requiredfields (wire shape unchanged).
	Status corev1.ConditionStatus `json:"status"`
	// reason is a machine understandable string that gives the reason for condition's last transition.
	// +optional
	//nolint:kubeapilinter // Stable v1alpha1: preserve API wire type and markers; new fields should comply with KAL.
	Reason string `json:"reason,omitempty"`
	// message is a human-readable message indicating details about last transition.
	// +optional
	//nolint:kubeapilinter // Stable v1alpha1: preserve API wire type and markers; new fields should comply with KAL.
	Message string `json:"message,omitempty"`
}

// IPPoolSpec defines the desired state of IPPool.
type IPPoolSpec struct {
	// startingAddress represents the starting IP address of the pool.
	// +required
	//nolint:kubeapilinter // Stable v1alpha1: KAL maxlength (wire shape unchanged).
	StartingAddress string `json:"startingAddress"`
	// addressCount represents the number of IP addresses in the pool.
	// +required
	//nolint:kubeapilinter // Stable v1alpha1: KAL requiredfields (wire shape unchanged).
	AddressCount int64 `json:"addressCount"`
}

// IPPoolStatus defines the current state of IPPool.
type IPPoolStatus struct {
	// allocated represents the number of IP addresses currently allocated to services.
	// +optional
	//nolint:kubeapilinter // Stable v1alpha1: preserve API wire type and markers; new fields should comply with KAL.
	Allocated int64 `json:"allocated,omitempty"`
	// conditions is an array of current observed IPPool conditions.
	// +listType=atomic
	// +optional
	//nolint:kubeapilinter // Stable v1alpha1: KAL conditions (wire shape unchanged).
	Conditions []IPPoolCondition `json:"conditions,omitempty"`
}

// +genclient
// +genclient:nonNamespaced
// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Cluster

// IPPool is the Schema for the ippools API.
// It represents a pool of IP addresses that are owned and managed by the IPPool controller.
// Provider specific networks can associate themselves with IPPool objects to use
// network operator's IPAM implementation.
// +kubebuilder:subresource:status
type IPPool struct {
	metav1.TypeMeta `json:",inline"`
	// metadata is the standard object's metadata.
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// spec describes the desired IP pool configuration.
	// +optional
	//nolint:kubeapilinter // Stable v1alpha1: preserve API wire type and markers; new fields should comply with KAL.
	Spec IPPoolSpec `json:"spec,omitempty"`
	// status reflects the observed state of the IP pool.
	// +optional
	//nolint:kubeapilinter // Stable v1alpha1: preserve API wire type and markers; new fields should comply with KAL.
	Status IPPoolStatus `json:"status,omitempty"`
}

type IPPoolReference struct {
	// name of the IPPool resource being referenced.
	// +required
	//nolint:kubeapilinter // Stable v1alpha1: KAL maxlength (wire shape unchanged).
	Name string `json:"name"`
	// apiVersion is the API version of the referent.
	// +optional
	//nolint:kubeapilinter // Stable v1alpha1: preserve API wire type and markers; new fields should comply with KAL.
	APIVersion string `json:"apiVersion,omitempty"`
}

// +kubebuilder:object:root=true

// IPPoolList contains a list of IPPool
type IPPoolList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	// +required
	Items []IPPool `json:"items"`
}

func init() {
	RegisterTypeWithScheme(&IPPool{}, &IPPoolList{})
}
