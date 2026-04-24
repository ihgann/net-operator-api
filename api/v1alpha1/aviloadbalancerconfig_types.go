// Copyright (c) 2020-2026 Broadcom. All Rights Reserved.
// Broadcom Confidential. The term "Broadcom" refers to Broadcom Inc.
// and/or its subsidiaries.

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AviLoadBalancerLogLevel is a valid log level for the Avi Kubernetes Operator.
type AviLoadBalancerLogLevel string

const (
	// AviLoadBalancerLogLevelInfo is the INFO log level for AKO.
	AviLoadBalancerLogLevelInfo AviLoadBalancerLogLevel = "INFO"
	// AviLoadBalancerLogLevelDebug is the DEBUG log level for AKO.
	AviLoadBalancerLogLevelDebug AviLoadBalancerLogLevel = "DEBUG"
	// AviLoadBalancerLogLevelWarn is the WARN log level for AKO.
	AviLoadBalancerLogLevelWarn AviLoadBalancerLogLevel = "WARN"
	// AviLoadBalancerLogLevelError is the ERROR log level for AKO.
	AviLoadBalancerLogLevelError AviLoadBalancerLogLevel = "ERROR"
)

// AviLoadBalancerIPAMType is the type of IPAM used by Avi.
type AviLoadBalancerIPAMType string

const (
	// AviLoadBalancerSupervisorIPAM indicates that IPAM is provided by the
	// Supervisor cluster.
	AviLoadBalancerSupervisorIPAM AviLoadBalancerIPAMType = "supervisor"
	// AviLoadBalancerControllerIPAM indicates that IPAM is provided by the Avi
	// Controller.
	AviLoadBalancerControllerIPAM AviLoadBalancerIPAMType = "controller"
)

// AviLoadBalancerConfigSpec defines the configuration for an Avi load balancer.
// This specification is used to configure the resources the Avi Kubernetes
// Operator (AKO) requires in order to connect to the Avi load balancer.
type AviLoadBalancerConfigSpec struct {
	// server is the endpoint at which the Avi Controller REST API is available.
	// The format is [SCHEME://]ADDRESS[:PORT], ex. https://10.10.10.10
	//   * SCHEME may be http or https and defaults to https if the SCHEME is
	//     omitted
	//   * ADDRESS is the Avi Controller IP address or the Avi Cluster IP when
	//     two or more Avi Controllers are deployed in cluster mode.
	//   * PORT defaults to 80 when SCHEME is http and 443 when SCHEME is https.
	// +required
	//nolint:kubeapilinter // Stable v1alpha1: KAL maxlength (wire shape unchanged).
	Server string `json:"server"`

	// cloudName is used by the Avi Kubernetes Operator (AKO) when querying
	// properties via the Avi REST API, ex. /api/cloud/?name=CLOUD_NAME.
	// Defaults to Default-Cloud.
	// +kubebuilder:default:=Default-Cloud
	// +optional
	//nolint:kubeapilinter // Stable v1alpha1: keep optional string + kubebuilder:default (do not rewrite to *string / drop default).
	CloudName string `json:"cloudName,omitempty"`

	// advancedL4 is a flag that enables support for WCP in AKO.
	// Defaults to true.
	// +kubebuilder:default:=true
	// +optional
	//nolint:kubeapilinter // Stable v1alpha1: preserve API wire type and markers; new fields should comply with KAL.
	AdvancedL4 *bool `json:"advancedL4,omitempty"`

	// logLevel specifies the log level used by AKO.
	// +kubebuilder:default:=WARN
	// +kubebuilder:validation:Enum=INFO;DEBUG;WARN;ERROR
	// +optional
	//nolint:kubeapilinter // Stable v1alpha1: preserve API wire type and markers; new fields should comply with KAL.
	LogLevel AviLoadBalancerLogLevel `json:"logLevel,omitempty"`

	// ipamType is the type of IPAM used by the Avi Software Load Balancer.
	// +kubebuilder:default:=controller
	// +kubebuilder:validation:Enum=controller;supervisor
	// +optional
	//nolint:kubeapilinter // Stable v1alpha1: preserve API wire type and markers; new fields should comply with KAL.
	IPAMType AviLoadBalancerIPAMType `json:"ipamType,omitempty"`

	// credentialSecretRef points to a Secret resource used to access and
	// configure the Avi Controller.
	//
	// * certificateAuthorityData   PEM-encoded certificate authority
	//                              certificates
	// * username                   Username used with basic authentication for
	//                              the Avi REST API
	// * password                   Password used with basic authentication for
	//                              the Avi REST API
	//
	// The following YAML is an example secret:
	//
	// apiVersion: v1
	// kind: Secret
	// metadata:
	//   name: avi-lb-config
	//   namespace: vmware-system-netop
	// data:
	//   certificateAuthorityData: []byte
	//   username: []byte
	//   password: []byte
	// +required
	//nolint:kubeapilinter // Stable v1alpha1: KAL requiredfields (wire shape unchanged).
	CredentialSecretRef ClientSecretReference `json:"credentialSecretRef"`
}

// AviLoadBalancerConfigStatus is unused because AviLoadBalancerConfigSpec is
// purely a configuration resource.
type AviLoadBalancerConfigStatus struct {
}

// +genclient
// +genclient:nonNamespaced
// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Cluster

// AviLoadBalancerConfig is the Schema for the AviLoadBalancerConfigs API
// +kubebuilder:subresource:status
type AviLoadBalancerConfig struct {
	metav1.TypeMeta `json:",inline"`
	// metadata is the standard object's metadata.
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// spec describes the desired Avi load balancer configuration.
	// +optional
	//nolint:kubeapilinter // Stable v1alpha1: preserve API wire type and markers; new fields should comply with KAL.
	Spec AviLoadBalancerConfigSpec `json:"spec,omitempty"`
	// status reflects the observed state of the Avi load balancer configuration.
	// +optional
	//nolint:kubeapilinter // Stable v1alpha1: preserve API wire type and markers; new fields should comply with KAL.
	Status AviLoadBalancerConfigStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// AviLoadBalancerConfigList contains a list of AviLoadBalancerConfig
type AviLoadBalancerConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	// +kubebuilder:validation:Required
	Items []AviLoadBalancerConfig `json:"items"`
}

func init() {
	RegisterTypeWithScheme(&AviLoadBalancerConfig{}, &AviLoadBalancerConfigList{})
}
