// Copyright (c) 2026 Broadcom. All Rights Reserved.
// Broadcom Confidential. The term "Broadcom" refers to Broadcom Inc.
// and/or its subsidiaries.

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NetworkSettingsProvider is the type of network provider for a namespace.
type NetworkSettingsProvider string

const (
	// NetworkSettingsProviderVSphereDistributed indicates vSphere Distributed networking.
	NetworkSettingsProviderVSphereDistributed NetworkSettingsProvider = "vsphere-distributed"
	// NetworkSettingsProviderNSXTier1 indicates NSX-T Tier-1 networking.
	NetworkSettingsProviderNSXTier1 NetworkSettingsProvider = "nsx-tier1"
	// NetworkSettingsProviderVPC indicates VPC networking.
	NetworkSettingsProviderVPC NetworkSettingsProvider = "vpc"
)

// +genclient
// +kubebuilder:object:root=true
// NetworkSettings exposes read-only settings that present relevant network information for a namespace.
type NetworkSettings struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Provider is the network provider type backing the network configuration for this namespace.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Enum=vsphere-distributed;nsx-tier1;vpc
	Provider NetworkSettingsProvider `json:"provider"`
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
