package v1

import "k8s.io/apimachinery/pkg/runtime"

// DeepCopyInto copies all properties of this object into another object of the
// same type that is provided as a pointer.
func (in *CatalogItem) DeepCopyInto(out *CatalogItem) {
	out.TypeMeta = in.TypeMeta
	out.ObjectMeta = in.ObjectMeta
	// out.Spec = CatalogItemSpec{
	// 	Replicas: in.Spec.Replicas,
	// }
}

// DeepCopyObject returns a generically typed copy of an object
func (in *CatalogItem) DeepCopyObject() runtime.Object {
	out := CatalogItem{}
	in.DeepCopyInto(&out)

	return &out
}

// DeepCopyObject returns a generically typed copy of an object
func (in *CatalogItemList) DeepCopyObject() runtime.Object {
	out := CatalogItemList{}
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta

	if in.Items != nil {
		out.Items = make([]CatalogItem, len(in.Items))
		for i := range in.Items {
			in.Items[i].DeepCopyInto(&out.Items[i])
		}
	}

	return &out
}
