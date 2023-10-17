package models

import (
	"encoding/json"

	v1 "github.com/rhpds/zerotouch-api/cmd/kube/apiextensions/v1"
)

// This implementation is not efficient and normally should be done
// with the reflection. But it's temprorary until Rates Database will be migrated
// and start to use the ResourceClaim.metadata.uuids
//
// After migration this function will not be needed
// TODO: delete this function after migration
func findUUIDs(rcs *v1.ResourceClaimStatus) []string {
	uuids := make([]string, 0)
	for _, value := range rcs.Resources {
		jsonString, err := json.Marshal(value)
		if err != nil {
			continue
		}

		state := struct {
			State struct {
				Spec struct {
					Vars struct {
						JobVars struct {
							Uuid string `json:"uuid"`
						} `json:"job_vars"`
					} `json:"vars"`
				} `json:"spec"`
			} `json:"state"`
		}{}

		err = json.Unmarshal(jsonString, &state)
		if err != nil {
			continue
		}

		uuids = append(uuids, state.State.Spec.Vars.JobVars.Uuid)
	}
	return uuids
}
