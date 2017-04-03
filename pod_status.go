package marathon

import (
	"fmt"
	"time"
)

type PodState string

const (
	PodStateDegraded PodState = "DEGRADED"
	PodStateStable   PodState = "STABLE"
	PodStateTerminal PodState = "TERMINAL"
)

type PodStatus struct {
	ID                 string                   `json:"id,omitempty"`
	Spec               *Pod                     `json:"spec,omitempty"`
	Status             PodState                 `json:"status,omitempty"`
	StatusSince        string                   `json:"statusSince,omitempty"`
	Message            string                   `json:"message,omitempty"`
	Instances          []*PodInstanceStatus     `json:"instances,omitempty"`
	TerminationHistory []*PodTerminationHistory `json:"terminationHistory,omitempty"`
	LastUpdated        string                   `json:"lastUpdated,omitempty"`
	LastChanged        string                   `json:"lastChanged,omitempty"`
}

type PodTerminationHistory struct {
	InstanceId   string                         `json:"instanceId,omitempty"`
	StartedAt    string                         `json:"startedAt,omitempty"`
	TerminatedAt string                         `json:"terminatedAt,omitempty"`
	Message      string                         `json:"message,omitempty"`
	Containers   []*ContainerTerminationHistory `json:"containers,omitempty"`
}

type ContainerTerminationHistory struct {
	ContainerId    string                     `json:"containerId,omitempty"`
	LastKnownState string                     `json:"lastKnownState,omitempty"`
	Termination    *ContainerTerminationState `json:"termination,omitempty"`
}

// TODO: Add helper functions for anything that makes sense?

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
	ticker := time.NewTicker(r.config.PollingWaitTime)
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
