package v1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// -----------------------------------------------------------------------------
// CatalogItem
// -----------------------------------------------------------------------------

// type CatalogItemSpec struct {
// 	Replicas int `json:"replicas"`
// }

type CatalogItem struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec CatalogItemSpec `json:"spec"`
}

type CatalogItemList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []CatalogItem `json:"items"`
}

// -----------------------------------------------------------------------------
// ResourceClaim
// -----------------------------------------------------------------------------

type ResourceClaimProvisionData struct {
	RandomString string `json:"random_string"`
	GUID         string `json:"GUID"`
}

type ResourceClaimStatusSummary struct {
	ProvisionData  ResourceClaimProvisionData `json:"provision_data"`
	RuntimeDefault string                     `json:"runtime_default"`
	RuntimeMaximum string                     `json:"runtime_maximum"`
	State          string                     `json:"state"`
}

type ResourceClaimStatus struct {
	Summary ResourceClaimStatusSummary `json:"summary"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ResourceClaim struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec ResourceClaimStatus `json:"status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ResourceClaimList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []ResourceClaim `json:"items"`
}
