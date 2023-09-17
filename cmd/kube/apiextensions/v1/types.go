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
	StopTimeStamp  string `json:"stop_timestamp"`
}

type ResourceClaimProvider struct {
	Name            string                       `json:"name"`
	ParameterValues ResourceClaimParameterValues `json:"parameterValues"`
}

type ResourceClaimProvisionData struct {
	GUID   string `json:"guid"`
	LabURL string `json:"lab_ui_url"`
}

type ResourceClaimStatusSummary struct {
	ProvisionData  ResourceClaimProvisionData `json:"provision_data"`
	RuntimeDefault string                     `json:"runtime_default"`
	RuntimeMaximum string                     `json:"runtime_maximum"`
	State          string                     `json:"state"`
}

type ResourceClaimStatusLifespan struct {
	End string `json:"end"`
}

type ResourceClaimStatus struct {
	Summary  ResourceClaimStatusSummary  `json:"summary"`
	Lifespan ResourceClaimStatusLifespan `json:"lifespan"`
}

type ResourceClaimSpec struct {
	Provider ResourceClaimProvider `json:"provider"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ResourceClaim struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ResourceClaimSpec   `json:"spec"`
	Status ResourceClaimStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ResourceClaimList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []ResourceClaim `json:"items"`
}
