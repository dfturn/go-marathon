/*
Copyright 2014 Rohith All rights reserved.

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

import (
	"encoding/json"
	"fmt"
)

// Pod is the definition for an pod in marathon
type Pod struct {
	Containers  []*PodContainer         `json:"containers,omitempty"`
	Environment *map[string]interface{} `json:"environment,omitempty"`
	ID          string                  `json:"id,omitempty"`
	Labels      *map[string]string      `json:"labels,omitempty"`
	Networks    []*PodNetwork           `json:"networks,omitempty"`
	Scaling     *PodScalingPolicy       `json:"scaling,omitempty"`
	Scheduling  *PodSchedulingPolicy    `json:"scheduling,omitempty"`
	User        string                  `json:"user,omitempty"`
	Version     string                  `json:"version,omitempty"`
	Volumes     []*PodVolume            `json:"volumes,omitempty"`

	Secrets *map[string]Secret `json:"secrets,omitempty"`
}

func (p *Pod) String() string {
	s, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return fmt.Sprintf(`{"error": "error decoding type into json: %s"}`, err)
	}

	return string(s)
}

type Secret struct {
	Source string `json:"source"`
}

type PodContainerImage struct {
	Kind      string `json:"kind"`
	ID        string `json:"id"`
	ForcePull bool   `json:"forcePull,omitempty"`
}

type PodResources struct {
	Cpus float64 `json:"cpus"`
	Mem  int     `json:"mem"`
	Disk int     `json:"disk"`
	Gpus int     `json:"gpus,omitempty"`
}

type PodHTTPEndpoint struct {
	Endpoint string `json:"endpoint,omitempty"`
	Path     string `json:"path,omitempty"`
	Scheme   string `json:"scheme,omitempty"`
}

type PodHealthCheck struct {
	HTTP                   *PodHTTPEndpoint `json:"http,omitempty"`
	GracePeriodSeconds     int              `json:"gracePeriodSeconds,omitempty"`
	IntervalSeconds        int              `json:"intervalSeconds,omitempty"`
	MaxConsecutiveFailures int              `json:"maxConsecutiveFailures,omitempty"`
	TimeoutSeconds         int              `json:"timeoutSeconds,omitempty"`
	DelaySeconds           int              `json:"delaySeconds,omitempty"`
}

type PodCommand struct {
	Shell string `json:"shell,omitempty"`
}
type PodExec struct {
	Command PodCommand `json:"command,omitempty"`
}

type PodEndpoints struct {
	Name          string             `json:"name"`
	ContainerPort int                `json:"containerPort"`
	HostPort      int                `json:"hostPort"`
	Protocol      []string           `json:"protocol"`
	Labels        *map[string]string `json:"labels,omitempty"`
}

type PodVolumeMounts struct {
	Name      string `json:"name"`
	MountPath string `json:"mountPath"`
}

type PodArtifact struct {
	Uri        string `json:"uri"`
	Extract    bool   `json:"extract,omitempty"`
	Executable bool   `json:"executable,omitempty"`
	Cache      bool   `json:"cache,omitempty"`
	DestPath   string `json:"destPath,omitempty"`
}

type PodContainer struct {
	Name         string                  `json:"name"`
	Exec         *PodExec                `json:"exec,omitempty"`
	Resources    *PodResources           `json:"resources"`
	Endpoints    []*PodEndpoints         `json:"endpoints,omitempty"`
	Image        *PodContainerImage      `json:"image"`
	Environment  *map[string]interface{} `json:"environment,omitempty"`
	User         string                  `json:"user,omitempty"`
	HealthCheck  *PodHealthCheck         `json:"healthCheck,omitempty"`
	VolumeMounts []*PodVolumeMounts      `json:"volumeMounts,omitempty"`
	Artifacts    []PodArtifact           `json:"artifacts,omitempty"`
	Labels       *map[string]string      `json:"labels,omitempty"`
	//Lifecycle interface{} `json:"lifecycle"`
}

type PodNetwork struct {
	Name   string             `json:"name"`
	Mode   string             `json:"mode"`
	Labels *map[string]string `json:"labels"`
}

type PodScalingPolicy struct {
	Kind         string `json:"kind"`
	Instances    int    `json:"instances"`
	MaxInstances int    `json:"maxInstances,omitempty"`
}

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

type PodVolume struct {
	Name string `json:"name"`
	Host string `json:"host"`
}

// CreatePod creates a new application in Marathon
// 		application:		the structure holding the application configuration
func (r *marathonClient) CreatePod(pod *Pod) (*Pod, error) {
	result := new(Pod)
	if err := r.apiPost(marathonAPIPods, &pod, result); err != nil {
		return nil, err
	}

	return result, nil
}

// DeleteApplication deletes an application from marathon
// 		name: 		the id used to identify the application
//		force:		used to force the delete operation in case of blocked deployment
func (r *marathonClient) DeletePod(name string) (*DeploymentID, error) {
	uri := buildPodURI(name)
	// step: check of the application already exists
	deployID := new(DeploymentID)
	if err := r.apiDelete(uri, nil, deployID); err != nil {
		return nil, nil // TODO: Get headers and get the deployment ID
	}

	return deployID, nil
}

func buildPodURI(path string) string {
	return fmt.Sprintf("%s/%s", marathonAPIPods, trimRootPath(path))
}
