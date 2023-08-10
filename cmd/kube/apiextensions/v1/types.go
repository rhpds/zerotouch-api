package v1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

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
