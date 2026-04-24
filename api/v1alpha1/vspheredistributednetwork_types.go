// Copyright (c) 2020-2026 Broadcom. All Rights Reserved.
// Broadcom Confidential. The term "Broadcom" refers to Broadcom Inc.
// and/or its subsidiaries.

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type VSphereDistributedNetworkConditionType string

const (
	// VSphereDistributedNetworkPortGroupFailure is added when PortGroupID specified either doesn't exist, or
	// there was an error in communicating with vCenter Server.
	VSphereDistributedNetworkPortGroupFailure VSphereDistributedNetworkConditionType = "PortGroupFailure"
	// VSphereDistributedNetworkIPPoolInvalid is added when no valid IPPool references exists.
	VSphereDistributedNetworkIPPoolInvalid VSphereDistributedNetworkConditionType = "IPPoolInvalid"
	// VsphereDistributedNetworkIPPoolPressure condition status is set to True when IPPool is low on free IPs.
	VsphereDistributedNetworkIPPoolPressure VSphereDistributedNetworkConditionType = "IPPoolPressure"
)

type IPAssignmentModeType string

const (
	// IPAssignmentModeDHCP indicates IP address is assigned dynamically using DHCP.
	IPAssignmentModeDHCP IPAssignmentModeType = "dhcp"
	// IPAssignmentModeStaticPool indicates IP address is assigned from a static pool of IP addresses.
	IPAssignmentModeStaticPool IPAssignmentModeType = "staticpool"
	// IPAssignmentModeNone indicates that no IP assignment will be performed.
	// The operator will not assign an IP and no DHCP client will be configured.
	IPAssignmentModeNone IPAssignmentModeType = "none"
)

// VSphereDistributedNetworkCondition describes the state of a VSphereDistributedNetwork at a certain point.
type VSphereDistributedNetworkCondition struct {
	// type is the type of VSphereDistributedNetwork condition.
	// +required
	//nolint:kubeapilinter // Stable v1alpha1: KAL maxlength (wire shape unchanged).
	Type VSphereDistributedNetworkConditionType `json:"type"`
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
	// lastTransitionTime is the timestamp for when the VSphereDistributedNetwork object last transitioned from one status to another.
	// +optional
	//nolint:kubeapilinter // Stable v1alpha1: preserve API wire type and markers; new fields should comply with KAL.
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty" patchStrategy:"replace"`
}

// VSphereDistributedNetworkSpec defines the desired state of VSphereDistributedNetwork.
type VSphereDistributedNetworkSpec struct {
	// portGroupID is an existing vSphere Distributed PortGroup identifier.
	// +required
	//nolint:kubeapilinter // Stable v1alpha1: KAL maxlength (wire shape unchanged).
	PortGroupID string `json:"portGroupID"`

	// ipAssignmentMode to use for network interfaces. If unset, defaults to IPAssignmentModeStaticPool.
	// For IPAssignmentModeDHCP and IPAssignmentModeNone, the IPPools, Gateway and SubnetMask
	// fields should be empty/unset. When using IPAssignmentModeNone, no IP will be assigned
	// and no DHCP client will be configured.
	// +optional
	//nolint:kubeapilinter // Stable v1alpha1: preserve API wire type and markers; new fields should comply with KAL.
	IPAssignmentMode IPAssignmentModeType `json:"ipAssignmentMode,omitempty"`

	// ipPools references list of IPPool objects. This field should only be set when using
	// IPAssignmentModeStaticPool. For all other modes (IPAssignmentModeDHCP, IPAssignmentModeNone), this should be set
	// to an empty list.
	// +listType=atomic
	// +optional
	//nolint:kubeapilinter // Stable v1alpha1: preserve API wire type and markers; new fields should comply with KAL.
	IPPools []IPPoolReference `json:"ipPools"`

	// gateway setting to use for network interfaces. This field should only be set when using
	// IPAssignmentModeStaticPool. For all other modes (IPAssignmentModeDHCP, IPAssignmentModeNone), this should be set
	// to an empty string.
	// +optional
	//nolint:kubeapilinter // Stable v1alpha1: preserve API wire type and markers; new fields should comply with KAL.
	Gateway string `json:"gateway"`

	// subnetMask setting to use for network interfaces. This field should only be set when using
	// IPAssignmentModeStaticPool. For all other modes (IPAssignmentModeDHCP, IPAssignmentModeNone), this should be set
	// to an empty string.
	// +required
	//nolint:kubeapilinter // Stable v1alpha1: KAL maxlength (wire shape unchanged).
	SubnetMask string `json:"subnetMask"`
}

// VLANType represents the type of VLAN configuration
type VLANType string

const (
	// VLANTypeStandard represents a standard VLAN configuration with a single VLAN ID
	VLANTypeStandard VLANType = "standard"
	// VLANTypeTrunk represents a VLAN trunk configuration that allows multiple VLANs
	VLANTypeTrunk VLANType = "trunk"
	// VLANTypePrivate represents a private VLAN configuration
	VLANTypePrivate VLANType = "private"
)

// VLANTrunkRange represents a range of VLAN IDs for trunk configuration
type VLANTrunkRange struct {
	// start represents the beginning of the VLAN ID range (inclusive).
	// +required
	//nolint:kubeapilinter // Stable v1alpha1: KAL requiredfields (wire shape unchanged).
	Start int32 `json:"start"`

	// end represents the end of the VLAN ID range (inclusive).
	// +required
	//nolint:kubeapilinter // Stable v1alpha1: KAL requiredfields (wire shape unchanged).
	End int32 `json:"end"`
}

// VlanSpec represents the VLAN configuration.
type VlanSpec struct {
	// type indicates the type of VLAN configuration (standard, trunk, or private).
	// +required
	//nolint:kubeapilinter // Stable v1alpha1: KAL maxlength (wire shape unchanged).
	Type VLANType `json:"type"`

	// vlanID specifies the VLAN ID when Type is VLANTypeStandard.
	// This field is ignored for other VLAN types.
	// Possible values:
	// - A value of 0 indicates there is no VLAN configuration for the port.
	// - A value from 1 to 4094 specifies a VLAN ID for the port.
	// +optional
	VlanID *int32 `json:"vlanID,omitempty"`

	// trunkRange specifies the ranges of allowed VLANs when Type is VLANTypeTrunk.
	// This field is ignored for other VLAN types.
	// Each range's Start and End values must be between 0 and 4094 inclusive.
	// Overlapping ranges are allowed.
	// +listType=atomic
	// +optional
	//nolint:kubeapilinter // Stable v1alpha1: KAL maxlength (wire shape unchanged).
	TrunkRange []VLANTrunkRange `json:"trunkRange,omitempty"`

	// privateVlanID specifies the private VLAN ID when Type is VLANTypePrivate.
	// This field is ignored for other VLAN types.
	// +optional
	PrivateVlanID *int32 `json:"privateVlanID,omitempty"`
}

// VSphereDistributedPortConfig represents the port-level configuration for a vSphere Distributed Network
type VSphereDistributedPortConfig struct {
	// vlan represents the VLAN configuration for this port.
	// If unset, indicates that no VLAN configuration has been retrieved yet for this port.
	// +optional
	Vlan *VlanSpec `json:"vlan,omitempty"`
}

// VSphereDistributedNetworkStatus defines the observed state of VSphereDistributedNetwork.
type VSphereDistributedNetworkStatus struct {
	// conditions is an array of current observed vSphere Distributed network conditions.
	// +listType=atomic
	// +optional
	//nolint:kubeapilinter // Stable v1alpha1: KAL conditions (wire shape unchanged).
	Conditions []VSphereDistributedNetworkCondition `json:"conditions,omitempty"`

	// defaultPortConfig represents the default port-level configuration that applies to all ports
	// unless overridden at the individual port level.
	// +optional
	DefaultPortConfig *VSphereDistributedPortConfig `json:"defaultPortConfig,omitempty"`
}

// +genclient
// +genclient:nonNamespaced
// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Cluster

// VSphereDistributedNetwork represents schema for a network backed by a vSphere Distributed PortGroup on vSphere
// Distributed switch.
// +kubebuilder:subresource:status
type VSphereDistributedNetwork struct {
	metav1.TypeMeta `json:",inline"`
	// metadata is the standard object's metadata.
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// spec describes the desired vSphere distributed network configuration.
	// +optional
	//nolint:kubeapilinter // Stable v1alpha1: preserve API wire type and markers; new fields should comply with KAL.
	Spec VSphereDistributedNetworkSpec `json:"spec,omitempty"`
	// status reflects the observed state of the vSphere distributed network.
	// +optional
	//nolint:kubeapilinter // Stable v1alpha1: preserve API wire type and markers; new fields should comply with KAL.
	Status VSphereDistributedNetworkStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// VSphereDistributedNetworkList contains a list of VSphereDistributedNetwork
type VSphereDistributedNetworkList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	// +required
	Items []VSphereDistributedNetwork `json:"items"`
}

func init() {
	RegisterTypeWithScheme(&VSphereDistributedNetwork{}, &VSphereDistributedNetworkList{})
}
