// Copyright (c) 2020-2026 Broadcom. All Rights Reserved.
// Broadcom Confidential. The term "Broadcom" refers to Broadcom Inc.
// and/or its subsidiaries.

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// IPAddressAllocationFinalizer is a finalizer that allows the controller to perform cleanup
// of resources associated with an IPAddressAllocation before it is removed from the API Server.
const IPAddressAllocationFinalizer = "ipaddressallocation.netoperator.vmware.com"

// IPAddressAllocationConditionType is a string type for the condition types of an IPAddressAllocation.
type IPAddressAllocationConditionType string

const (
	// IPAddressAllocationReady indicates the IP has been successfully allocated.
	IPAddressAllocationReady IPAddressAllocationConditionType = "Ready"
	// IPAddressAllocationFail indicates an error was encountered during allocation.
	IPAddressAllocationFail IPAddressAllocationConditionType = "Failure"
)

// IPAddressAllocationConditionReason describes the reason for the last transition of a condition.
type IPAddressAllocationConditionReason string

const (
	// IPAddressAllocationConditionInvalidRequestedIP is used when the IPAddressAllocation fails due to an invalid RequestedIP.
	IPAddressAllocationConditionInvalidRequestedIP IPAddressAllocationConditionReason = "InvalidRequestedIP"
	// IPAddressAllocationConditionFailureReasonCannotAllocIP is used when the IPAddressAllocation fails because an IP cannot be allocated.
	IPAddressAllocationConditionFailureReasonCannotAllocIP IPAddressAllocationConditionReason = "CannotAllocIP"
	// IPAddressAllocationConditionFailureReasonIPPoolRefRetrievalFailed is used when retrieval of the IPPoolRef has failed.
	IPAddressAllocationConditionFailureReasonIPPoolRefRetrievalFailed IPAddressAllocationConditionReason = "IPPoolRefRetrievalFailed"
)

// IPAddressAllocationCondition describes the state of an IPAddressAllocation at a specific point in time.
type IPAddressAllocationCondition struct {
	// type is the type of the condition.
	// +required
	//nolint:kubeapilinter // Stable v1alpha1: KAL maxlength (wire shape unchanged).
	Type IPAddressAllocationConditionType `json:"type"`
	// status reflects whether the condition is True, False, or Unknown.
	// +required
	//nolint:kubeapilinter // Stable v1alpha1: KAL requiredfields (wire shape unchanged).
	Status corev1.ConditionStatus `json:"status"`
	// lastTransitionTime is the timestamp of the last change to the condition's status.
	// +optional
	//nolint:kubeapilinter // Stable v1alpha1: preserve API wire type and markers; new fields should comply with KAL.
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty"`
	// reason provides a machine-readable explanation for the last status transition.
	// +optional
	//nolint:kubeapilinter // Stable v1alpha1: preserve API wire type and markers; new fields should comply with KAL.
	Reason IPAddressAllocationConditionReason `json:"reason,omitempty"`
	// message provides a human-readable explanation for the last status transition.
	// +optional
	//nolint:kubeapilinter // Stable v1alpha1: preserve API wire type and markers; new fields should comply with KAL.
	Message string `json:"message,omitempty"`
}

// IPAddressAllocationSpec defines the desired state of an IPAddressAllocation, including the pool reference and an optional requested IP.
type IPAddressAllocationSpec struct {
	// poolRef is the reference to the network's IP pool within the namespace.
	// It currently only supports reference to a Network.
	// +required
	//nolint:kubeapilinter // Stable v1alpha1: KAL requiredfields (wire shape unchanged).
	PoolRef corev1.TypedLocalObjectReference `json:"poolRef"`
	// requestedIP is an optional field for a user to specify a particular IP they want to request.
	// If omitted, the system will allocate a single IP address.
	// +optional
	//nolint:kubeapilinter // Stable v1alpha1: preserve API wire type and markers; new fields should comply with KAL.
	RequestedIP string `json:"requestedIP,omitempty"`
}

// IPAddressAllocationStatus contains the current status of an IPAddressAllocation, including the allocated IP address and conditions.
type IPAddressAllocationStatus struct {
	// ipaddress is the actually allocated IP address.
	// +optional
	//nolint:kubeapilinter // Stable v1alpha1: preserve API wire type and markers; new fields should comply with KAL.
	IPAddress string `json:"ipaddress,omitempty"`
	// conditions provide detailed information about the status of the allocation.
	// +listType=atomic
	// +optional
	//nolint:kubeapilinter // Stable v1alpha1: KAL conditions (wire shape unchanged).
	Conditions []IPAddressAllocationCondition `json:"conditions,omitempty"`
}

// +genclient
// +kubebuilder:object:root=true

// IPAddressAllocation represents a request for IP address allocation, including the desired state and current status.
// +kubebuilder:subresource:status
type IPAddressAllocation struct {
	metav1.TypeMeta `json:",inline"`
	// metadata is the standard object's metadata.
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`
	// spec describes the desired IP address allocation.
	// +optional
	//nolint:kubeapilinter // Stable v1alpha1: preserve API wire type and markers; new fields should comply with KAL.
	Spec IPAddressAllocationSpec `json:"spec,omitempty"`
	// status reflects the observed state of the IP address allocation.
	// +optional
	//nolint:kubeapilinter // Stable v1alpha1: preserve API wire type and markers; new fields should comply with KAL.
	Status IPAddressAllocationStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// IPAddressAllocationList is a list of IPAddressAllocation objects.
type IPAddressAllocationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	// +required
	Items []IPAddressAllocation `json:"items"`
}

// init function registers the IPAddressAllocation type with the scheme.
func init() {
	RegisterTypeWithScheme(&IPAddressAllocation{}, &IPAddressAllocationList{})
}
