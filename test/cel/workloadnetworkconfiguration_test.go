// Copyright (c) 2026 Broadcom. All Rights Reserved.
// Broadcom Confidential. The term "Broadcom" refers to Broadcom Inc.
// and/or its subsidiaries.

package cel_test

import (
	"testing"

	netv1alpha1 "github.com/vmware-tanzu/net-operator-api/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// validWNC returns a minimal valid WorkloadNetworkConfiguration: a single vpc
// provider entry with a valid autoCreateConfig, active on that provider.
func validWNC() *netv1alpha1.WorkloadNetworkConfiguration {
	return &netv1alpha1.WorkloadNetworkConfiguration{
		ObjectMeta: metav1.ObjectMeta{
			Name: netv1alpha1.WorkloadNetworkConfigurationName,
		},
		Spec: netv1alpha1.WorkloadNetworkConfigurationSpec{
			ActiveSystemProvider: netv1alpha1.NetworkProviderVPC,
			Providers: []netv1alpha1.NetworkProviderEntry{
				{
					Type: netv1alpha1.NetworkProviderVPC,
					SystemConfiguration: &netv1alpha1.NamespaceNetworkConfig{
						VPCConfig: netv1alpha1.VPCConfig{
							AutoCreateConfig: netv1alpha1.AutoCreateVPCConfig{
								NSXProject:             "/orgs/default/projects/default",
								VPCConnectivityProfile: "/orgs/default/projects/default/vpc-connectivity-profiles/default",
								PrivateCIDRs:           []string{"10.0.0.0/24"},
							},
						},
					},
				},
			},
		},
	}
}

func TestWorkloadNetworkConfiguration_ValidDefaultVPCConfig_Admitted(t *testing.T) {
	obj := validWNC()
	obj.Spec.Providers[0].DefaultNamespaceConfiguration = netv1alpha1.NetworkProviderDefaultConfig{
		VPCConfig: &netv1alpha1.DefaultVPCConfig{
			PrivateCIDRs: []string{"10.1.0.0/24"},
		},
	}
	if err := k8sClient.Create(testCtx, obj); err != nil {
		t.Fatalf("expected admission, got: %v", err)
	}
	defer func() { _ = k8sClient.Delete(testCtx, obj) }()
}

func TestWorkloadNetworkConfiguration_DefaultVPCConfigOnNonVPCType_Rejected(t *testing.T) {
	obj := &netv1alpha1.WorkloadNetworkConfiguration{
		ObjectMeta: metav1.ObjectMeta{
			Name: netv1alpha1.WorkloadNetworkConfigurationName,
		},
		Spec: netv1alpha1.WorkloadNetworkConfigurationSpec{
			ActiveSystemProvider: netv1alpha1.NetworkProviderVSphereDistributed,
			Providers: []netv1alpha1.NetworkProviderEntry{
				{
					Type: netv1alpha1.NetworkProviderVSphereDistributed,
					SystemConfiguration: &netv1alpha1.NamespaceNetworkConfig{
						VSphereDistributedConfig: netv1alpha1.VSphereDistributedConfig{
							Networks: []netv1alpha1.VSphereDistributedNetworkRef{
								{Name: "net-1"},
							},
							DefaultNetwork: "net-1",
						},
					},
					DefaultNamespaceConfiguration: netv1alpha1.NetworkProviderDefaultConfig{
						VPCConfig: &netv1alpha1.DefaultVPCConfig{
							PrivateCIDRs: []string{"10.1.0.0/24"},
						},
					},
				},
			},
		},
	}
	if err := k8sClient.Create(testCtx, obj); !isRejected(err) {
		t.Fatalf("expected rejection for defaultNamespaceConfiguration.vpcConfig on non-vpc type, got: %v", err)
	}
}

// validWNCUnstructured returns the same object as validWNC(), as Unstructured
// with the given defaultNamespaceConfiguration value spliced in on the wire.
// Use this instead of the typed client when the value under test would be
// dropped by omitempty/omitzero on marshal (e.g. an explicit empty map or
// empty list), since the typed client can't put those bytes on the wire.
func validWNCUnstructured(defaultNamespaceConfiguration map[string]interface{}) *unstructured.Unstructured {
	return &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "netoperator.vmware.com/v1alpha1",
			"kind":       "WorkloadNetworkConfiguration",
			"metadata": map[string]interface{}{
				"name": netv1alpha1.WorkloadNetworkConfigurationName,
			},
			"spec": map[string]interface{}{
				"activeSystemProvider": string(netv1alpha1.NetworkProviderVPC),
				"providers": []interface{}{
					map[string]interface{}{
						"type": string(netv1alpha1.NetworkProviderVPC),
						"systemConfiguration": map[string]interface{}{
							"vpcConfig": map[string]interface{}{
								"autoCreateConfig": map[string]interface{}{
									"nsxProject":             "/orgs/default/projects/default",
									"vpcConnectivityProfile": "/orgs/default/projects/default/vpc-connectivity-profiles/default",
									"privateCIDRs":           []interface{}{"10.0.0.0/24"},
								},
							},
						},
						"defaultNamespaceConfiguration": defaultNamespaceConfiguration,
					},
				},
			},
		},
	}
}

// TestWorkloadNetworkConfiguration_EmptyDefaultNamespaceConfiguration_Rejected sends the
// request as Unstructured because DefaultNamespaceConfiguration carries `omitzero`, which
// drops a Go zero-value struct entirely on marshal — the typed client would silently turn
// this into "field absent" instead of testing the MinProperties=1 rule.
func TestWorkloadNetworkConfiguration_EmptyDefaultNamespaceConfiguration_Rejected(t *testing.T) {
	obj := validWNCUnstructured(map[string]interface{}{})
	if err := k8sClient.Create(testCtx, obj); !isRejected(err) {
		t.Fatalf("expected rejection for empty defaultNamespaceConfiguration (MinProperties=1), got: %v", err)
	}
}

// TestWorkloadNetworkConfiguration_EmptyPrivateCIDRsList_Rejected sends the
// request as Unstructured because the typed client's `omitempty` drops an
// empty (but non-nil) []string entirely on marshal, which would silently
// turn this into "field absent" instead of testing the MinItems=1 rule.
func TestWorkloadNetworkConfiguration_EmptyPrivateCIDRsList_Rejected(t *testing.T) {
	obj := validWNCUnstructured(map[string]interface{}{
		"vpcConfig": map[string]interface{}{
			"privateCIDRs": []interface{}{},
		},
	})
	if err := k8sClient.Create(testCtx, obj); !isRejected(err) {
		t.Fatalf("expected rejection for explicit empty privateCIDRs list (MinItems=1), got: %v", err)
	}
}

func TestWorkloadNetworkConfiguration_DefaultPrivateCIDRsFullReplace_Admitted(t *testing.T) {
	obj := validWNC()
	obj.Spec.Providers[0].DefaultNamespaceConfiguration = netv1alpha1.NetworkProviderDefaultConfig{
		VPCConfig: &netv1alpha1.DefaultVPCConfig{
			PrivateCIDRs: []string{"10.1.0.0/24"},
		},
	}
	if err := k8sClient.Create(testCtx, obj); err != nil {
		t.Fatalf("create: %v", err)
	}
	defer func() { _ = k8sClient.Delete(testCtx, obj) }()

	latest := &netv1alpha1.WorkloadNetworkConfiguration{}
	if err := k8sClient.Get(testCtx, client.ObjectKeyFromObject(obj), latest); err != nil {
		t.Fatalf("get: %v", err)
	}
	// Full replace with a CIDR that was never in the original list. Unlike
	// systemConfiguration.vpcConfig.autoCreateConfig.privateCIDRs (append-only),
	// defaultNamespaceConfiguration.vpcConfig.privateCIDRs has no append-only
	// constraint, so dropping 10.1.0.0/24 entirely must be admitted.
	latest.Spec.Providers[0].DefaultNamespaceConfiguration.VPCConfig.PrivateCIDRs = []string{"10.2.0.0/24"}
	if err := k8sClient.Update(testCtx, latest); err != nil {
		t.Fatalf("expected admission for full-replace of defaultNamespaceConfiguration privateCIDRs, got: %v", err)
	}
}
