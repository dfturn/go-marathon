package marathon

import (
	"fmt"
	"time"
)

type PodState string

const (
	PodStateDegraded PodState = "degraded"
	PodStateStable   PodState = "stable"
	PodStateTerminal PodState = "terminal"
)

type PodStatus struct {
	ID                 string                   `json:"id"`
	Instances          []*PodInstanceStatus     `json:"instances"`
	LastChanged        string                   `json:"lastChanged"`
	LastUpdated        string                   `json:"lastUpdated"`
	Message            string                   `json:"message"`
	Spec               *Pod                     `json:"spec"`
	Status             PodState                 `json:"status"`
	StatusSince        string                   `json:"statusSince"`
	TerminationHistory []*PodTerminationHistory `json:"terminationHistory"`
}

type PodTerminationHistory struct {
	InstanceId   string                         `json:"instanceId"`
	StartedAt    string                         `json:"startedAt"` // TODO time
	TerminatedAt string                         `json:"terminatedAt"`
	Message      string                         `json:"message"`
	Containers   []*ContainerTerminationHistory `json:"containers"`
}

type ContainerTerminationHistory struct {
	ContainerId    string                     `json:"containerId"`
	LastKnownState string                     `json:"lastKnownState"`
	Termination    *ContainerTerminationState `json:"termination"`
}

type StatusCondition struct {
	Name        string `json:"name"`
	Value       string `json:"value"`
	Reason      string `json:"reason,omitempty"`
	LastChanged string `json:"lastChanged"` // TODO: datetime
	LastUpdated string `json:"lastUpdated"` // TODO: datetime
}

type ContainerTerminationState struct {
	ExitCode int    `json:"exitCode,omitempty"`
	Message  string `json:"message,omitempty"`
}

type ContainerStatus struct {
	Conditions  []*StatusCondition         `json:"conditions"`
	ContainerID string                     `json:"containerId"`
	Endpoints   []interface{}              `json:"endpoints"` // TODO
	LastChanged string                     `json:"lastChanged"`
	LastUpdated string                     `json:"lastUpdated"`
	Message     string                     `json:"message"`
	Name        string                     `json:"name"`
	Resources   *PodResources              `json:"resources"`
	Status      string                     `json:"status"`
	StatusSince string                     `json:"statusSince"`
	Termination *ContainerTerminationState `json:"termination,omitempty"`
}

type PodNetworkStatus struct {
	Addresses []string `json:"addresses"`
	Name      string   `json:"name"`
}

type PodInstanceStatus struct {
	AgentHostname string              `json:"agentHostname"`
	Conditions    []*StatusCondition  `json:"conditions"`
	Containers    []*ContainerStatus  `json:"containers"`
	ID            string              `json:"id"`
	LastChanged   string              `json:"lastChanged"`
	LastUpdated   string              `json:"lastUpdated"`
	Message       string              `json:"message"`
	Networks      []*PodNetworkStatus `json:"networks"`
	Resources     *PodResources       `json:"resources"`
	SpecReference string              `json:"specReference"`
	Status        string              `json:"status"` // TODO: This should be an enum (TASK_RUNNING etc)
	StatusSince   string              `json:"statusSince"`
}

// Application retrieves the application configuration from marathon
// 		name: 		the id used to identify the application
func (r *marathonClient) GetPodStatus(name string) (*PodStatus, error) {
	var podStatus PodStatus

	if err := r.apiGet(buildPodStatusURI(name), nil, &podStatus); err != nil {
		return nil, err
	}

	return &podStatus, nil
}

func buildPodStatusURI(path string) string {
	return fmt.Sprintf("%s/%s::status", marathonAPIPods, trimRootPath(path))
}

// WaitOnPod waits for a pod to be deployed
//		name:		the id of the pod
//		timeout:	a duration of time to wait for an pod to deploy
func (r *marathonClient) WaitOnPod(name string, timeout time.Duration) error {
	if r.podExistAndRunning(name) {
		return nil
	}

	timeoutTimer := time.After(timeout)
	ticker := time.NewTicker(r.pollingWaitTime)
	defer ticker.Stop()

	for {
		select {
		case <-timeoutTimer:
			return ErrTimeoutError
		case <-ticker.C:
			if r.podExistAndRunning(name) {
				return nil
			}
		}
	}
}

func (r *marathonClient) podExistAndRunning(name string) bool {
	podStatus, err := r.GetPodStatus(name)
	if apiErr, ok := err.(*APIError); ok && apiErr.ErrCode == ErrCodeNotFound {
		return false
	}
	if err == nil && podStatus.Status == PodStateStable {
		return true
	}
	return false
}
