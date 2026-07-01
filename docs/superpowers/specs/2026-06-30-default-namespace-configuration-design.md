# DefaultNamespaceConfiguration on WorkloadNetworkConfiguration

## Problem

`WorkloadNetworkConfigurationSpec.Providers[].SystemConfiguration` (a
`*NamespaceNetworkConfig`) describes the provider-specific template for the
currently-active system NNC. For the VPC provider, `systemConfiguration.vpcConfig`
uses `AutoCreateVPCConfig`, where `nsxProject` and `vpcConnectivityProfile` are
immutable once set and `privateCIDRs` is append-only (existing entries cannot be
removed) — see the CEL rules on `AutoCreateVPCConfig` in
`namespacenetworkconfiguration_types.go`.

This matches the "active NNC" semantics: the system configuration tracks what is
already provisioned and can only grow.

However, per the source-of-truth VMODL
(`bora/vpx/wcp/wcpsvc/vmodl/namespace_management/Clusters.vmodl`,
`VpcClusterNetworkSetSpec` / `VpcClusterNetworkUpdateSpec` at
[Clusters.vmodl#L1265-L1287](https://github-vcf.devops.broadcom.net/vcf/tera/blob/5a3943fb4487d4116551a59c02856e3382954728/bora/vpx/wcp/wcpsvc/vmodl/namespace_management/Clusters.vmodl#L1265-L1287)),
only `defaultPrivateCidrs` is independently mutable after initial cluster setup:

- `nsxProject` / `vpcConnectivityProfile` (`VpcConfig`/`VpcNetworkSummary`) are
  set once at cluster-enable time and never exposed on the update spec.
- `defaultPrivateCidrs` is exposed on both the set and update specs as
  `Optional<List<Ipv4Cidr>>` with full-replace semantics: "if unset, the current
  value will be retained." There is no append-only constraint. It establishes
  the default for *new* namespaces going forward; it does not affect namespaces
  that already exist.

There is currently no place in `WorkloadNetworkConfiguration` to express this
second, independently-mutable "default for new namespaces" value. This design
adds one, scoped to what VPC actually supports today.

## Shape

A new optional field on `NetworkProviderEntry`, structured as a union parallel
to `NamespaceNetworkConfig` so it can grow a sibling provider config later
(e.g. NSX Tier1) without restructuring:

```go
type NetworkProviderEntry struct {
	// type identifies the network provider for this entry.
	//
	// +required
	Type NetworkProvider `json:"type,omitempty"`

	// systemConfiguration holds the provider-specific NNC template for this provider.
	//
	// +required
	SystemConfiguration *NamespaceNetworkConfig `json:"systemConfiguration,omitempty"`

	// defaultNamespaceConfiguration optionally overrides configuration values
	// applied when provisioning new namespaces under this provider, going
	// forward. When unset, or when a given sub-field is unset, the equivalent
	// value from systemConfiguration is used instead.
	//
	// Unlike systemConfiguration, values here are not tied to the currently
	// active system NNC: they may be replaced freely and only take effect for
	// namespaces onboarded after the change. Namespaces already associated
	// under this provider are unaffected.
	//
	// +optional
	DefaultNamespaceConfiguration *NetworkProviderDefaultConfig `json:"defaultNamespaceConfiguration,omitempty"`
}

// NetworkProviderDefaultConfig holds the subset of provider configuration that
// may be independently overridden for new namespaces, distinct from and
// unconstrained by the append-only/immutable rules that apply to
// systemConfiguration. Currently only the vpc provider has such a divergence.
//
// +kubebuilder:validation:MinProperties=1
type NetworkProviderDefaultConfig struct {
	// vpcConfig overrides default values used when auto-creating VPCs for new
	// namespaces under the vpc provider.
	//
	// +optional
	VPCConfig *DefaultVPCConfig `json:"vpcConfig,omitempty"`
}

// DefaultVPCConfig holds VPC default values that may be changed independently
// of the vpc provider's systemConfiguration.
type DefaultVPCConfig struct {
	// privateCIDRs replaces the CIDR blocks used as private-subnet/private-pod-IP
	// defaults for VPCs auto-created for namespaces onboarding after this is
	// set. Already-created namespace VPCs are unaffected.
	//
	// This value fully replaces (not appends to) the previous default and is
	// not subject to the append-only constraint that applies to
	// systemConfiguration.vpcConfig.autoCreateConfig.privateCIDRs. When unset,
	// new namespaces use systemConfiguration's privateCIDRs instead.
	//
	// +optional
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=16
	// +kubebuilder:validation:items:MaxLength=64
	// +kubebuilder:validation:items:Pattern=`^([0-9]{1,3}\.){3}[0-9]{1,3}/[0-9]{1,2}$`
	// +listType=atomic
	PrivateCIDRs []string `json:"privateCIDRs,omitempty"`
}
```

`MinItems=1` on `privateCIDRs` only applies when the field is present (OpenAPI
validation skips absent optional fields), so an explicit empty list is
rejected — to inherit the system value, omit the field entirely rather than
setting `[]`.

## Validation (CEL)

One new rule on `NetworkProviderEntry`, mirroring the existing gating pattern
used for `systemConfiguration.vpcConfig`:

```
+kubebuilder:validation:XValidation:rule="self.type == 'vpc' || !has(self.defaultNamespaceConfiguration.vpcConfig)",message="defaultNamespaceConfiguration.vpcConfig may only be set when type is vpc"
```

Unlike `systemConfiguration.vpcConfig` there is no reverse "must be set when
type is vpc" rule — `defaultNamespaceConfiguration` is optional in its
entirety; a vpc provider entry may validly have no override at all, meaning
new namespaces continue to use `systemConfiguration`'s CIDRs.

`MinProperties=1` on `NetworkProviderDefaultConfig` prevents a pointless empty
`{}` override object.

## Extensibility (not implemented in this branch)

`NetworkProviderDefaultConfig` is deliberately shaped like `NamespaceNetworkConfig`
so that if NSX Tier1 gains an analogous "independently mutable default" concept,
it is added as one more optional sub-field (e.g. `NSXTier1Config`) plus one more
type-gating CEL rule, with no restructuring of this type.

## Out of scope

- No controller/reconciler changes — this repo (`net-operator-api`) contains
  only the API types, generated deepcopy/client, and CRD manifests. There is no
  reconciler here to consume this field.
- No NSX Tier1 fields in this branch.
- No change to `AutoCreateVPCConfig`'s existing append-only/immutable CEL rules.

## Files touched

- `api/v1alpha1/workloadnetworkconfiguration_types.go` — new field on
  `NetworkProviderEntry`, two new types, one new CEL rule.
- `api/v1alpha1/zz_generated.deepcopy.go` — regenerated via `make generate`.
- `config/crd/bases/supervisor.netoperator.vmware.com_workloadnetworkconfigurations.yaml` —
  regenerated via `make generate`.
- `pkg/client/...` — regenerated via `make generate` (generate-client target).

## Testing

The repo has no `*_test.go` files today. Verification is:

1. `make generate` produces a clean diff (deepcopy, CRD YAML, client all
   regenerate without manual edits needed).
2. Manual review of the generated CEL/schema in the CRD YAML for correctness.
3. Optionally, manual `kubectl apply --dry-run=server` of hand-built valid and
   invalid sample manifests against a cluster with the CRD installed, to
   exercise the new CEL rule.
