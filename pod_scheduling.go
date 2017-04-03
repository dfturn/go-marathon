package marathon

type PodBackoff struct {
	Backoff        int     `json:"backoff"`
	BackoffFactor  float64 `json:"backoffFactor"`
	MaxLaunchDelay int     `json:"maxLaunchDelay"`
}

type PodUpgrade struct {
	MinimumHealthCapacity int `json:"minimumHealthCapacity"`
	MaximumOverCapacity   int `json:"maximumOverCapacity"`
}

type PodPlacement struct {
	Constraints           *[][]string `json:"constraints"`
	AcceptedResourceRoles []string    `json:"acceptedResourceRoles"`
}

type PodSchedulingPolicy struct {
	Backoff   *PodBackoff   `json:"backoff,omitempty"`
	Upgrade   *PodUpgrade   `json:"upgrade,omitempty"`
	Placement *PodPlacement `json:"placement,omitempty"`
}

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

func NewPodSchedulingPolicy() *PodSchedulingPolicy {
	return &PodSchedulingPolicy{
		Placement: NewPodPlacement(),
	}
}
