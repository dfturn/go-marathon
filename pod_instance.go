package marathon

type PodInstanceState string

const (
	PodInstanceStatePending  PodInstanceState = "PENDING"
	PodInstanceStateStaging  PodInstanceState = "STAGING"
	PodInstanceStateStable   PodInstanceState = "STABLE"
	PodInstanceStateDegraded PodInstanceState = "DEGRADED"
	PodInstanceStateTerminal PodInstanceState = "TERMINAL"
)

type PodInstanceStatus struct {
	AgentHostname string              `json:"agentHostname,omitempty"`
	Conditions    []*StatusCondition  `json:"conditions,omitempty"`
	Containers    []*ContainerStatus  `json:"containers,omitempty"`
	ID            string              `json:"id,omitempty"`
	LastChanged   string              `json:"lastChanged,omitempty"`
	LastUpdated   string              `json:"lastUpdated,omitempty"`
	Message       string              `json:"message,omitempty"`
	Networks      []*PodNetworkStatus `json:"networks,omitempty"`
	Resources     *Resources          `json:"resources,omitempty"`
	SpecReference string              `json:"specReference,omitempty"`
	Status        PodInstanceState    `json:"status,omitempty"`
	StatusSince   string              `json:"statusSince,omitempty"`
}

type PodNetworkStatus struct {
	Addresses []string `json:"addresses,omitempty"`
	Name      string   `json:"name,omitempty"`
}

type StatusCondition struct {
	Name        string `json:"name,omitempty"`
	Value       string `json:"value,omitempty"`
	Reason      string `json:"reason,omitempty"`
	LastChanged string `json:"lastChanged,omitempty"`
	LastUpdated string `json:"lastUpdated,omitempty"`
}

type ContainerStatus struct {
	Conditions  []*StatusCondition         `json:"conditions,omitempty"`
	ContainerID string                     `json:"containerId,omitempty"`
	Endpoints   []*PodEndpoint             `json:"endpoints,omitempty"`
	LastChanged string                     `json:"lastChanged,omitempty"`
	LastUpdated string                     `json:"lastUpdated,omitempty"`
	Message     string                     `json:"message,omitempty"`
	Name        string                     `json:"name,omitempty"`
	Resources   *Resources                 `json:"resources,omitempty"`
	Status      string                     `json:"status,omitempty"`
	StatusSince string                     `json:"statusSince,omitempty"`
	Termination *ContainerTerminationState `json:"termination,omitempty"`
}

type ContainerTerminationState struct {
	ExitCode int    `json:"exitCode,omitempty"`
	Message  string `json:"message,omitempty"`
}
