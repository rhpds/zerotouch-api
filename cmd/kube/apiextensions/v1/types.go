package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// -----------------------------------------------------------------------------
// CatalogItem
// -----------------------------------------------------------------------------
type CatalogItemLifespan struct {
	Default         string `json:"default"`
	Maximum         string `json:"maximum"`
	RelativeMaximum string `json:"relativeMaximum"`
}

type CatalogItemSpec struct {
	Lifespan CatalogItemLifespan `json:"lifespan"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type CatalogItem struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec CatalogItemSpec `json:"spec"`
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
	StopTimeStamp  string `json:"stop_timestamp,omitempty"`
}

type ResourceClaimProvider struct {
	Name            string                       `json:"name"`
	ParameterValues ResourceClaimParameterValues `json:"parameterValues"`
}

type ResourceClaimProvisionData struct {
	GUID   string `json:"GUID"`
	LabURL string `json:"lab_ui_url"`
}

type ResourceClaimStatusSummary struct {
	ProvisionData  ResourceClaimProvisionData `json:"provision_data"`
	RuntimeDefault string                     `json:"runtime_default"`
	RuntimeMaximum string                     `json:"runtime_maximum"`
	State          string                     `json:"state"`
}

type ResourceClaimStatusLifespan struct {
	End   string `json:"end"`
	Start string `json:"start"`
}

type ResourceClaimSpecLifespan struct {
	End string `json:"end"`
}

// +k8s:deepcopy-gen=true
type ResourceClaimStatus struct {
	Summary  ResourceClaimStatusSummary  `json:"summary"`
	Lifespan ResourceClaimStatusLifespan `json:"lifespan"`
}

// +k8s:deepcopy-gen=true
type ResourceClaimSpec struct {
	Lifespan *ResourceClaimSpecLifespan `json:"lifespan,omitempty"`
	Provider ResourceClaimProvider      `json:"provider"`
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
