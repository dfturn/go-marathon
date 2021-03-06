/*
Copyright 2017 Devin All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package marathon

// PodBackoff describes the backoff for re-run attempts of a pod
type PodBackoff struct {
	Backoff        int     `json:"backoff"`
	BackoffFactor  float64 `json:"backoffFactor"`
	MaxLaunchDelay int     `json:"maxLaunchDelay"`
}

// PodUpgrade describes the policy for upgrading a pod in-place
type PodUpgrade struct {
	MinimumHealthCapacity int `json:"minimumHealthCapacity"`
	MaximumOverCapacity   int `json:"maximumOverCapacity"`
}

// PodPlacement supports constraining which hosts a pod is placed on
type PodPlacement struct {
	Constraints           *[][]string `json:"constraints"`
	AcceptedResourceRoles []string    `json:"acceptedResourceRoles,omitempty"`
}

// PodSchedulingPolicy is the overarching pod scheduling policy
type PodSchedulingPolicy struct {
	Backoff   *PodBackoff   `json:"backoff,omitempty"`
	Upgrade   *PodUpgrade   `json:"upgrade,omitempty"`
	Placement *PodPlacement `json:"placement,omitempty"`
}

// NewPodPlacement creates an empty PodPlacement
func NewPodPlacement() *PodPlacement {
	return &PodPlacement{
		Constraints:           &[][]string{},
		AcceptedResourceRoles: []string{},
	}
}

// AddConstraint adds a new constraint
//		constraints:	the constraint definition, one constraint per array element
func (r *PodPlacement) AddConstraint(constraints ...string) *PodPlacement {
	c := *r.Constraints
	c = append(c, constraints)
	r.Constraints = &c

	return r
}

// NewPodSchedulingPolicy creates an empty PodSchedulingPolicy
func NewPodSchedulingPolicy() *PodSchedulingPolicy {
	return &PodSchedulingPolicy{
		Placement: NewPodPlacement(),
	}
}
