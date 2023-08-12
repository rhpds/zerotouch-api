package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// -----------------------------------------------------------------------------
// CatalogItem
// -----------------------------------------------------------------------------

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type CatalogItem struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type CatalogItemList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []CatalogItem `json:"items"`
}

// -----------------------------------------------------------------------------
// ResourceClaim
// -----------------------------------------------------------------------------

type ResourceClaimParameterValues struct {
	Purpose        string `json:"purpose"`
	StartTimeStamp string `json:"start_timestamp"`
	EndTimeStamp   string `json:"end_timestamp"`
}

type ResourceClaimProvider struct {
	Name            string                       `json:"name"`
	ParameterValues ResourceClaimParameterValues `json:"parameterValues"`
}

type ResourceClaimLifespan struct {
	End string `json:"end"`
}

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

type ResourceClaimSpec struct {
	Lifespan ResourceClaimLifespan `json:"lifespan"`
	Provider ResourceClaimProvider `json:"provider"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ResourceClaim struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ResourceClaimSpec   `json:"spec"`
	Status ResourceClaimStatus `json:"status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ResourceClaimList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []ResourceClaim `json:"items"`
}
