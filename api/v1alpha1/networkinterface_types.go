// Copyright (c) 2020-2026 Broadcom. All Rights Reserved.
// Broadcom Confidential. The term "Broadcom" refers to Broadcom Inc.
// and/or its subsidiaries.

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// NetworkInterfaceFinalizer allows the Controller to clean up resources associated
	// with a NetworkInterface before removing it from the API Server.
	NetworkInterfaceFinalizer = "networkinterface.netoperator.vmware.com"

	// NetworkInterfaceClientManagedAnnotation annotations means the NetworkInterface is
	// client managed and the Controller will not reconcile it. The value does not need
	// to be truthy; the presence of the key is what disables reconciliation.
	NetworkInterfaceClientManagedAnnotation = "networkinterface.netoperator.vmware.com/client-managed"
)

// IPConfig represents an IP configuration.
type IPConfig struct {
	// ip setting.
	// +required
	//nolint:kubeapilinter // Stable v1alpha1: KAL maxlength (wire shape unchanged).
	IP string `json:"ip"`
	// ipFamily specifies the IP family (IPv4 vs IPv6) the IP belongs to.
	// +required
	//nolint:kubeapilinter // Stable v1alpha1: KAL requiredfields (wire shape unchanged).
	IPFamily corev1.IPFamily `json:"ipFamily"`
	// gateway setting.
	// +required
	//nolint:kubeapilinter // Stable v1alpha1: KAL maxlength (wire shape unchanged).
	Gateway string `json:"gateway"`
	// subnetMask setting.
	// +required
	//nolint:kubeapilinter // Stable v1alpha1: KAL maxlength (wire shape unchanged).
	SubnetMask string `json:"subnetMask"`
}

// NetworkInterfaceProviderReference contains info to locate a network interface provider object.
type NetworkInterfaceProviderReference struct {
	// apiGroup is the group for the resource being referenced.
	// +required
	//nolint:kubeapilinter // Stable v1alpha1: KAL maxlength (wire shape unchanged).
	APIGroup string `json:"apiGroup"`
	// kind is the type of resource being referenced
	// +required
	//nolint:kubeapilinter // Stable v1alpha1: KAL maxlength (wire shape unchanged).
	Kind string `json:"kind"`
	// name is the name of resource being referenced
	// +required
	//nolint:kubeapilinter // Stable v1alpha1: KAL maxlength (wire shape unchanged).
	Name string `json:"name"`
	// apiVersion is the API version of the referent.
	// +optional
	//nolint:kubeapilinter // Stable v1alpha1: preserve API wire type and markers; new fields should comply with KAL.
	APIVersion string `json:"apiVersion,omitempty"`
}

type NetworkInterfaceConditionType string

const (
	// NetworkInterfaceReady is added when all network settings have been updated and the network
	// interface is ready to be used.
	NetworkInterfaceReady NetworkInterfaceConditionType = "Ready"
	// NetworkInterfaceFailure is added when network provider plugin returns an error.
	NetworkInterfaceFailure NetworkInterfaceConditionType = "Failure"
)

type NetworkInterfaceConditionReason string

const (
	// NetworkInterfaceFailureReasonCannotAllocIP indicates NetworkInterface is in failed state because an
	// IPConfig cannot be allocated.
	NetworkInterfaceFailureReasonCannotAllocIP NetworkInterfaceConditionReason = "CannotAllocIP"
	// NetworkInterfaceFailureReasonCannotAllocPort indicates NetworkInterface is in failed state because
	// port cannot be allocated for network interface on the network.
	NetworkInterfaceFailureReasonCannotAllocPort NetworkInterfaceConditionReason = "CannotAllocPort"
	// NetworkInterfaceFailureReasonNetworkDeleted indicates NetworkInterface is in failed state because
	// the underlying Network resource has been deleted.
	NetworkInterfaceFailureReasonNetworkDeleted NetworkInterfaceConditionReason = "NetworkDeleted"
)

// NetworkInterfaceCondition describes the state of a NetworkInterface at a certain point.
type NetworkInterfaceCondition struct {
	// type is the type of network interface condition.
	// +required
	//nolint:kubeapilinter // Stable v1alpha1: KAL maxlength (wire shape unchanged).
	Type NetworkInterfaceConditionType `json:"type"`
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
	Reason NetworkInterfaceConditionReason `json:"reason,omitempty"`
	// message is a human-readable message indicating details about last transition.
	// +optional
	//nolint:kubeapilinter // Stable v1alpha1: preserve API wire type and markers; new fields should comply with KAL.
	Message string `json:"message,omitempty"`
}

// NetworkInterfaceStatus defines the observed state of NetworkInterface.
// Once NetworkInterfaceReady condition is True, it should contain configuration to use to place
// a VM/Pod/Container's nic on the specified network.
type NetworkInterfaceStatus struct {
	// conditions is an array of current observed network interface conditions.
	// +listType=atomic
	// +optional
	//nolint:kubeapilinter // Stable v1alpha1: KAL conditions (wire shape unchanged).
	Conditions []NetworkInterfaceCondition `json:"conditions,omitempty"`
	// ipConfigs is an array of IP configurations for the network interface.
	// +listType=atomic
	// +optional
	//nolint:kubeapilinter // Stable v1alpha1: KAL maxlength (wire shape unchanged).
	IPConfigs []IPConfig `json:"ipConfigs,omitempty"`
	// macAddress setting for the network interface.
	// +optional
	//nolint:kubeapilinter // Stable v1alpha1: preserve API wire type and markers; new fields should comply with KAL.
	MacAddress string `json:"macAddress,omitempty"`
	// externalID is a network provider specific identifier assigned to the network interface.
	// +optional
	//nolint:kubeapilinter // Stable v1alpha1: preserve API wire type and markers; new fields should comply with KAL.
	ExternalID string `json:"externalID,omitempty"`
	// networkID is an network provider specific identifier for the network backing the network
	// interface.
	// +optional
	//nolint:kubeapilinter // Stable v1alpha1: preserve API wire type and markers; new fields should comply with KAL.
	NetworkID string `json:"networkID,omitempty"`
	// portID is a network provider specific port identifier allocated for this network interface on
	// the backing network. It is only valid on requested node and is set only if port allocation
	// was requested.
	// +optional
	//nolint:kubeapilinter // Stable v1alpha1: preserve API wire type and markers; new fields should comply with KAL.
	PortID string `json:"portID,omitempty"`
	// connectionID is a network provider specific port connection identifier allocated for this
	// network interface on the backing network. It is only valid on requested node and is set
	// only if port allocation was requested.
	// +optional
	//nolint:kubeapilinter // Stable v1alpha1: preserve API wire type and markers; new fields should comply with KAL.
	ConnectionID string `json:"connectionID,omitempty"`
	// ipAssignmentMode indicates how IP addresses are assigned to this interface.
	// When unset:
	// - If IP is assigned, it is assumed to be NetworkInterfaceIPAssignmentModeStaticPool.
	// - If IP is unassigned, it is assumed to be NetworkInterfaceIPAssignmentModeDHCP.
	// When set to NetworkInterfaceIPAssignmentModeStaticPool, indicates IP is assigned from a static pool.
	// When set to NetworkInterfaceIPAssignmentModeDHCP, indicates IP should be obtained via DHCP.
	// When set to NetworkInterfaceIPAssignmentModeNone, indicates no IP assignment should be performed.
	// +optional
	//nolint:kubeapilinter // Stable v1alpha1: preserve API wire type and markers; new fields should comply with KAL.
	IPAssignmentMode NetworkInterfaceIPAssignmentMode `json:"ipAssignmentMode,omitempty"`
}

type NetworkInterfaceType string

const (
	// NetworkInterfaceTypeVMXNet3 is for a VMXNET3 device.
	NetworkInterfaceTypeVMXNet3 = NetworkInterfaceType("vmxnet3")
)

// NetworkInterfaceIPAssignmentMode defines how IP addresses are assigned to a network interface
type NetworkInterfaceIPAssignmentMode string

const (
	// NetworkInterfaceIPAssignmentModeStaticPool indicates IP address is assigned from a static pool.
	NetworkInterfaceIPAssignmentModeStaticPool NetworkInterfaceIPAssignmentMode = "staticpool"

	// NetworkInterfaceIPAssignmentModeDHCP indicates IP address should be obtained via DHCP.
	NetworkInterfaceIPAssignmentModeDHCP NetworkInterfaceIPAssignmentMode = "dhcp"

	// NetworkInterfaceIPAssignmentModeNone indicates no IP assignment should be performed.
	NetworkInterfaceIPAssignmentModeNone NetworkInterfaceIPAssignmentMode = "none"
)

// NetworkInterfacePortAllocation describes the settings for network interface port allocation request.
type NetworkInterfacePortAllocation struct {
	// nodeName is the node where port must be allocated for this network interface.
	// +required
	//nolint:kubeapilinter // Stable v1alpha1: KAL maxlength (wire shape unchanged).
	NodeName string `json:"nodeName"`
}

// NetworkInterfaceSpec defines the desired state of NetworkInterface.
type NetworkInterfaceSpec struct {
	// networkName refers to a NetworkObject in the same namespace.
	// +optional
	//nolint:kubeapilinter // Stable v1alpha1: preserve API wire type and markers; new fields should comply with KAL.
	NetworkName string `json:"networkName,omitempty"`
	// type is the type of NetworkInterface. Supported values are vmxnet3.
	// +optional
	//nolint:kubeapilinter // Stable v1alpha1: preserve API wire type and markers; new fields should comply with KAL.
	Type NetworkInterfaceType `json:"type,omitempty"`
	// providerRef is a reference to a provider specific network interface object
	// that specifies the network interface configuration.
	// If unset, default configuration is assumed.
	// +optional
	ProviderRef *NetworkInterfaceProviderReference `json:"providerRef,omitempty"`
	// portAllocation is a request to allocate a port for this network interface on the backing network.
	// This feature is currently supported only if backing network type is NetworkTypeVDS. In all other
	// cases this field is ignored. Typically this is done implicitly by vCenter Server at the time
	// of attaching a network interface to a network and should be left unset. This is used primarily when
	// attachment of network interface to the network is done without vCenter Server's knowledge.
	// +optional
	PortAllocation *NetworkInterfacePortAllocation `json:"portAllocation,omitempty"`
	// externalID describes a value that will be surfaced as status.externalID.
	// If this field is omitted, then it is up to the underlying network
	// provider to surface any information in status.externalID.
	// +optional
	//nolint:kubeapilinter // Stable v1alpha1: preserve API wire type and markers; new fields should comply with KAL.
	ExternalID string `json:"externalID,omitempty"`
}

// NetworkInterfaceReference is an object that points to a NetworkInterface.
type NetworkInterfaceReference struct {
	// kind is the type of resource being referenced.
	// +required
	//nolint:kubeapilinter // Stable v1alpha1: KAL maxlength (wire shape unchanged).
	Kind string `json:"kind"`
	// name is the name of resource being referenced.
	// +required
	//nolint:kubeapilinter // Stable v1alpha1: KAL maxlength (wire shape unchanged).
	Name string `json:"name"`
	// apiVersion is the API version of the referent.
	//
	// +optional
	//nolint:kubeapilinter // Stable v1alpha1: preserve API wire type and markers; new fields should comply with KAL.
	APIVersion string `json:"apiVersion,omitempty"`
}

// +genclient
// +kubebuilder:object:root=true

// NetworkInterface is the Schema for the networkinterfaces API.
// A NetworkInterface represents a user's request for network configuration to use to place a
// VM/Pod/Container's nic on a specified network.
// +kubebuilder:subresource:status
type NetworkInterface struct {
	metav1.TypeMeta `json:",inline"`
	// metadata is the standard object's metadata.
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// spec describes the desired network interface configuration.
	// +optional
	//nolint:kubeapilinter // Stable v1alpha1: preserve API wire type and markers; new fields should comply with KAL.
	Spec NetworkInterfaceSpec `json:"spec,omitempty"`
	// status reflects the observed state of the network interface.
	// +optional
	//nolint:kubeapilinter // Stable v1alpha1: preserve API wire type and markers; new fields should comply with KAL.
	Status NetworkInterfaceStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// NetworkInterfaceList contains a list of NetworkInterface
type NetworkInterfaceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	// +required
	Items []NetworkInterface `json:"items"`
}

func init() {
	RegisterTypeWithScheme(&NetworkInterface{}, &NetworkInterfaceList{})
}
