// Copyright (c) 2020-2026 Broadcom. All Rights Reserved.
// Broadcom Confidential. The term "Broadcom" refers to Broadcom Inc.
// and/or its subsidiaries.

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ClientSecretReference contains info to locate an object of Kind Secret
// which contains credential specifications for a load balancer.
type ClientSecretReference struct {
	// name is the name of resource being referenced.
	// +required
	//nolint:kubeapilinter // Stable v1alpha1: KAL maxlength (wire shape unchanged).
	Name string `json:"name"`
	// namespace of the resource being referenced. If empty, cluster scoped resource is assumed.
	// +kubebuilder:default:=default
	// +optional
	//nolint:kubeapilinter // Stable v1alpha1: preserve API wire type and markers; new fields should comply with KAL.
	Namespace string `json:"namespace,omitempty"`
}

// LoadBalancerConfigConditionType is used as a typed string for representing
// LoadBalancerConfig.Status.Conditions.
type LoadBalancerConfigConditionType string

const (
	// LoadBalancerConfigReady is added when the LoadBalancerConfig object has been successfully realized
	LoadBalancerConfigReady LoadBalancerConfigConditionType = "Ready"
	// LoadBalancerConfigFailure is added if any failure is encountered while realizing LoadBalancerConfig object
	LoadBalancerConfigFailure LoadBalancerConfigConditionType = "Failure"
	// LoadBalancerConfigIPPoolPressure condition status is set to True when IPPool is low on free IPs.
	LoadBalancerConfigIPPoolPressure LoadBalancerConfigConditionType = "IPPoolPressure"
)

// LoadBalancerConfigCondition describes the state of a LoadBalancerConfig at a certain point
type LoadBalancerConfigCondition struct {
	// type is the type of load balancer condition
	// Can be Ready or Failure
	// +required
	//nolint:kubeapilinter // Stable v1alpha1: KAL maxlength (wire shape unchanged).
	Type LoadBalancerConfigConditionType `json:"type"`
	// status is the status of the condition
	// Can be True, False, Unknown
	// +required
	//nolint:kubeapilinter // Stable v1alpha1: KAL requiredfields (wire shape unchanged).
	Status corev1.ConditionStatus `json:"status"`
	// reason is a machine understandable string that gives the reason for the condition's last transition.
	// +optional
	//nolint:kubeapilinter // Stable v1alpha1: preserve API wire type and markers; new fields should comply with KAL.
	Reason string `json:"reason,omitempty"`
	// message is a human-readable message indicating details about last transition.
	// +optional
	//nolint:kubeapilinter // Stable v1alpha1: preserve API wire type and markers; new fields should comply with KAL.
	Message string `json:"message,omitempty"`
	// lastTransitionTime is the timestamp for when the LoadBalancerConfig object last transitioned from one status to another.
	// +optional
	//nolint:kubeapilinter // Stable v1alpha1: preserve API wire type and markers; new fields should comply with KAL.
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty" patchStrategy:"replace"`
}

// LoadBalancerConfigProviderReference represents the specific load balancer instance that needs to be configured
type LoadBalancerConfigProviderReference struct {
	// apiGroup is the group for the resource being referenced
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

type LoadBalancerConfigType string

const (
	// LoadBalancerConfigTypeHAProxy is the LoadBalancerConfigType for HAProxy.
	LoadBalancerConfigTypeHAProxy LoadBalancerConfigType = "haproxy"

	// LoadBalancerConfigTypeAvi is the LoadBalancerConfigType for Avi.
	LoadBalancerConfigTypeAvi LoadBalancerConfigType = "avi"

	// LoadBalancerConfigTypeFoundation is the FoundationLoadBalancerConfigType for VCF Foundation Load Balancer.
	LoadBalancerConfigTypeFoundation LoadBalancerConfigType = "foundation"
)

// LoadBalancerConfigSpec defines the desired state of LoadBalancerConfig
type LoadBalancerConfigSpec struct {
	// type describes type of load balancer.
	// +kubebuilder:validation:Enum=haproxy;avi;foundation
	// +required
	//nolint:kubeapilinter // Stable v1alpha1: KAL requiredfields (wire shape unchanged).
	Type LoadBalancerConfigType `json:"type"`
	// providerRef is reference to a load balancer provider object that provides the details for this type of load balancer
	// +required
	//nolint:kubeapilinter // Stable v1alpha1: KAL requiredfields (wire shape unchanged).
	ProviderRef LoadBalancerConfigProviderReference `json:"providerRef"`
}

// LoadBalancerConfigStatus defines the observed state of LoadBalancerConfig
type LoadBalancerConfigStatus struct {
	// conditions is an array of current observed load balancer conditions
	// +listType=atomic
	// +optional
	//nolint:kubeapilinter // Stable v1alpha1: KAL conditions (wire shape unchanged).
	Conditions []LoadBalancerConfigCondition `json:"conditions,omitempty"`
}

// +genclient
// +genclient:nonNamespaced
// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Cluster

// LoadBalancerConfig is the Schema for the LoadBalancerConfigs API
// +kubebuilder:subresource:status
type LoadBalancerConfig struct {
	metav1.TypeMeta `json:",inline"`
	// metadata is the standard object's metadata.
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// spec describes the desired load balancer configuration.
	// +optional
	//nolint:kubeapilinter // Stable v1alpha1: preserve API wire type and markers; new fields should comply with KAL.
	Spec LoadBalancerConfigSpec `json:"spec,omitempty"`
	// status reflects the observed state of the load balancer configuration.
	// +optional
	//nolint:kubeapilinter // Stable v1alpha1: preserve API wire type and markers; new fields should comply with KAL.
	Status LoadBalancerConfigStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// LoadBalancerConfigList contains a list of LoadBalancerConfig
type LoadBalancerConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	// +required
	Items []LoadBalancerConfig `json:"items"`
}

func init() {
	RegisterTypeWithScheme(&LoadBalancerConfig{}, &LoadBalancerConfigList{})
}
