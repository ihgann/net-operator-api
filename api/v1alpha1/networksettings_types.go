// Copyright (c) 2026 Broadcom. All Rights Reserved.
// Broadcom Confidential. The term "Broadcom" refers to Broadcom Inc.
// and/or its subsidiaries.

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NetworkSettingsTopology is the active network topology for a namespace.
type NetworkSettingsTopology string

const (
	// NetworkSettingsTopologyVSphereDistributed is vSphere Distributed (VDS) network backing.
	NetworkSettingsTopologyVSphereDistributed NetworkSettingsTopology = "vsphere-distributed"
	// NetworkSettingsTopologyNSXTier1 is NSX-T Tier-1 (non-VPC) network backing.
	NetworkSettingsTopologyNSXTier1 NetworkSettingsTopology = "nsx-tier1"
	// NetworkSettingsTopologyVPC is VPC (NSX) network backing; see
	// crd.nsx.vmware.com/v1alpha1.NetworkInfo for more detail when this applies.
	NetworkSettingsTopologyVPC NetworkSettingsTopology = "vpc"
)

// +genclient
// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Namespaced
//
// NetworkSettings exposes read-only, operator-relevant information about the effective network
// configuration for a namespace.
//
// Consumers should treat it as observed, realized state, and expect it to track the network topology
// backing the namespace.
type NetworkSettings struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Topology is the active network topology in this namespace. Workloads and network-aware
	// components should use this to determine the network backing that is in effect, including
	// when choosing defaulting behavior or which provider-specific APIs to use when not specified
	// elsewhere. This value is mutable; it may change when the configuration changes.
	//
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Enum=vsphere-distributed;nsx-tier1;vpc
	Topology NetworkSettingsTopology `json:"topology"`
}

// +kubebuilder:object:root=true

// NetworkSettingsList is a list of NetworkSettings.
type NetworkSettingsList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NetworkSettings `json:"items"`
}

func init() {
	RegisterTypeWithScheme(&NetworkSettings{}, &NetworkSettingsList{})
}
