// Copyright (c) 2020-2026 Broadcom. All Rights Reserved.
// Broadcom Confidential. The term "Broadcom" refers to Broadcom Inc.
// and/or its subsidiaries.

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// NetworkProtectionFinalizer allows the Controller to clean up resources and ensure
	// that no NetworkInterfaces are actively using this Network before deletion.
	NetworkProtectionFinalizer = "network.netoperator.vmware.com/network-protection"
)

type NetworkConditionType string

const (
	// NetworkDeletionBlocked indicates that the Network cannot be deleted, because
	// there may be some consumers (NetworkInterface) still actively using it.
	NetworkDeletionBlocked NetworkConditionType = "DeletionBlocked"
)

type NetworkConditionReason string

const (
	// NetworkDeletionBlockedReasonInUse indicates that the Network deletion is blocked
	// because there are NetworkInterfaces still actively using this Network.
	NetworkDeletionBlockedReasonInUse NetworkConditionReason = "NetworkInUse"
)

// NetworkProviderReference contains info to locate a network provider object.
type NetworkProviderReference struct {
	// apiGroup is the group for the resource being referenced.
	// +required
	//nolint:kubeapilinter // Stable v1alpha1: KAL maxlength (wire shape unchanged).
	APIGroup string `json:"apiGroup"`
	// kind is the type of resource being referenced.
	// +required
	//nolint:kubeapilinter // Stable v1alpha1: KAL maxlength (wire shape unchanged).
	Kind string `json:"kind"`
	// name is the name of resource being referenced.
	// +required
	//nolint:kubeapilinter // Stable v1alpha1: KAL maxlength (wire shape unchanged).
	Name string `json:"name"`
	// namespace of the resource being referenced. If empty, cluster scoped resource is assumed.
	// +optional
	//nolint:kubeapilinter // Stable v1alpha1: preserve API wire type and markers; new fields should comply with KAL.
	Namespace string `json:"namespace,omitempty"`
	// apiVersion is the API version of the referent.
	// +optional
	//nolint:kubeapilinter // Stable v1alpha1: preserve API wire type and markers; new fields should comply with KAL.
	APIVersion string `json:"apiVersion,omitempty"`
}

// NetworkType is used to type the constants describing possible network types.
type NetworkType string

const (
	// NetworkTypeNSXT is the network type describing NSX-T.
	NetworkTypeNSXT = NetworkType("nsx-t")

	// NetworkTypeVDS is the network type describing VSphere Distributed Switch.
	NetworkTypeVDS = NetworkType("vsphere-distributed")

	// NetworkTypeNSXTVPC is the network type describing NSX-T VPC.
	NetworkTypeNSXTVPC = NetworkType("nsx-t_vpc")
)

// NetworkSpec defines the state of Network.
type NetworkSpec struct {
	// type describes type of Network. Supported values are nsx-t, vsphere-distributed.
	// +required
	//nolint:kubeapilinter // Stable v1alpha1: KAL maxlength (wire shape unchanged).
	Type NetworkType `json:"type"`
	// providerRef is reference to a network provider object that provides this type of network.
	// +required
	//nolint:kubeapilinter // Stable v1alpha1: KAL requiredfields (wire shape unchanged).
	ProviderRef NetworkProviderReference `json:"providerRef"`
	// dns is a list of DNS server IPs to associate with network interfaces on this network.
	// +listType=atomic
	// +optional
	//nolint:kubeapilinter // Stable v1alpha1: KAL maxlength (wire shape unchanged).
	DNS []string `json:"dns,omitempty"`
	// dnsSearchDomains is a list of DNS search domains to associate with network interfaces on this network.
	// +listType=atomic
	// +optional
	//nolint:kubeapilinter // Stable v1alpha1: KAL maxlength (wire shape unchanged).
	DNSSearchDomains []string `json:"dnsSearchDomains,omitempty"`
	// ntp is a list of NTP server DNS names or IP addresses to use on this network.
	// +listType=atomic
	// +optional
	//nolint:kubeapilinter // Stable v1alpha1: KAL maxlength (wire shape unchanged).
	NTP []string `json:"ntp,omitempty"`
}

// NetworkCondition describes the state of a Network at a certain point.
type NetworkCondition struct {
	// type is the type of network condition.
	// +required
	//nolint:kubeapilinter // Stable v1alpha1: KAL maxlength (wire shape unchanged).
	Type NetworkConditionType `json:"type"`
	// status is the status of the condition.
	// Can be True, False, Unknown.
	// +required
	//nolint:kubeapilinter // Stable v1alpha1: KAL requiredfields (wire shape unchanged).
	Status corev1.ConditionStatus `json:"status"`
	// lastTransitionTime is the timestamp corresponding to the last status
	// change of this condition.
	// +optional
	//nolint:kubeapilinter // Stable v1alpha1: preserve API wire type and markers; new fields should comply with KAL.
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty"`
	// reason is a machine understandable string that gives the reason for condition's last transition.
	// +optional
	//nolint:kubeapilinter // Stable v1alpha1: preserve API wire type and markers; new fields should comply with KAL.
	Reason NetworkConditionReason `json:"reason,omitempty"`
	// message is a human-readable message indicating details about last transition.
	// +optional
	//nolint:kubeapilinter // Stable v1alpha1: preserve API wire type and markers; new fields should comply with KAL.
	Message string `json:"message,omitempty"`
}

// NetworkStatus defines the observed state of Network.
type NetworkStatus struct {
	// conditions is an array of current observed network conditions.
	// +listType=atomic
	// +optional
	//nolint:kubeapilinter // Stable v1alpha1: KAL conditions (wire shape unchanged).
	Conditions []NetworkCondition `json:"conditions,omitempty"`
}

// NetworkReference is an object that points to a Network.
type NetworkReference struct {
	// kind is the type of resource being referenced.
	// +required
	//nolint:kubeapilinter // Stable v1alpha1: KAL maxlength (wire shape unchanged).
	Kind string `json:"kind"`
	// name is the name of resource being referenced.
	// +required
	//nolint:kubeapilinter // Stable v1alpha1: KAL maxlength (wire shape unchanged).
	Name string `json:"name"`
	// apiVersion of the referent.
	//
	// +optional
	//nolint:kubeapilinter // Stable v1alpha1: preserve API wire type and markers; new fields should comply with KAL.
	APIVersion string `json:"apiVersion,omitempty"`
}

// +genclient
// +kubebuilder:object:root=true

// Network is the Schema for the networks API.
// A Network describes type, class and common attributes of a network available
// in a namespace. A NetworkInterface resource references a Network.
// +kubebuilder:subresource:status
type Network struct {
	metav1.TypeMeta `json:",inline"`
	// metadata is the standard object's metadata.
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// spec describes the desired network configuration.
	// +optional
	//nolint:kubeapilinter // Stable v1alpha1: preserve API wire type and markers; new fields should comply with KAL.
	Spec NetworkSpec `json:"spec,omitempty"`
	// status reflects the observed state of the network.
	// +optional
	//nolint:kubeapilinter // Stable v1alpha1: preserve API wire type and markers; new fields should comply with KAL.
	Status NetworkStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// NetworkList contains a list of Network
type NetworkList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	// +required
	Items []Network `json:"items"`
}

func init() {
	RegisterTypeWithScheme(&Network{}, &NetworkList{})
}
