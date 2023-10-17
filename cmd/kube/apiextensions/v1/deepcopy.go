package v1

func (in *ResourceState) DeepCopyInto(out *ResourceState) {
	out.State = make(map[string]interface{})
	for key, value := range in.State {
		out.State[key] = value
	}
}
