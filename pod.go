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

import (
	"encoding/json"
	"fmt"
)

// Pod is the definition for an pod in marathon
type Pod struct {
	ID                string                  `json:"id,omitempty"`
	Labels            map[string]string       `json:"labels,omitempty"`
	Version           string                  `json:"version,omitempty"`
	User              string                  `json:"user,omitempty"`
	Environment       map[string]interface{}  `json:"environment,omitempty"`
	Containers        []*PodContainer         `json:"containers,omitempty"`
	Secrets           map[string]SecretSource `json:"secrets,omitempty"`
	Volumes           []*PodVolume            `json:"volumes,omitempty"`
	Networks          []*PodNetwork           `json:"networks,omitempty"`
	Scaling           *PodScalingPolicy       `json:"scaling,omitempty"`
	Scheduling        *PodSchedulingPolicy    `json:"scheduling,omitempty"`
	ExecutorResources *ExecutorResources      `json:"executorResources,omitempty"`

	marathonClient Marathon
}

// PodScalingPolicy is the scaling policy of the pod
type PodScalingPolicy struct {
	Kind         string `json:"kind"`
	Instances    int    `json:"instances"`
	MaxInstances int    `json:"maxInstances,omitempty"`
}

// NewPod create an empty pod
func NewPod() *Pod {
	return &Pod{
		Labels:      map[string]string{},
		Environment: map[string]interface{}{},
		Containers:  []*PodContainer{},
		Secrets:     map[string]SecretSource{},
		Volumes:     []*PodVolume{},
		Networks:    []*PodNetwork{},
	}
}

// String marshals the pod as an indented string
func (p *Pod) String() string {
	s, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return fmt.Sprintf(`{"error": "error decoding type into json: %s"}`, err)
	}

	return string(s)
}

// Name sets the name / ID of the pod i.e. the identifier for this pod
func (p *Pod) Name(id string) *Pod {
	p.ID = validateID(id)
	return p
}

// SetUser sets the user to run the pod as
func (p *Pod) SetUser(user string) *Pod {
	p.User = user
	return p
}

// EmptyLabels empties the labels in a pod
func (p *Pod) EmptyLabels() *Pod {
	p.Labels = make(map[string]string)
	return p
}

// AddLabel adds a label to a pod
func (p *Pod) AddLabel(key, value string) *Pod {
	p.Labels[key] = value
	return p
}

// SetLabels sets the labels for a pod
func (p *Pod) SetLabels(labels map[string]string) *Pod {
	p.Labels = labels
	return p
}

// EmptyEnvironment empties the environment variables for a pod
func (p *Pod) EmptyEnvironment() *Pod {
	p.Environment = make(map[string]interface{})
	return p
}

// AddEnvironment adds an environment variable to a pod
func (p *Pod) AddEnvironment(name, value string) *Pod {
	p.Environment[name] = value
	return p
}

// ExtendEnvironment extends the environment with the new environment variables
func (p *Pod) ExtendEnvironment(env map[string]string) *Pod {
	for k, v := range env {
		p.AddEnvironment(k, v)
	}
	return p
}

// GetEnvironmentVariable gets the string contained in an environment variable
func (p *Pod) GetEnvironmentVariable(name string) (string, error) {
	str := ""
	var err error

	if val, ok := p.Environment[name]; ok {
		switch val.(type) {
		case string:
			str = val.(string)
			err = nil
		case EnvironmentSecret:
			err = fmt.Errorf("environment variable refers to a secret")
		default:
			err = fmt.Errorf("environment variable refers to unknown type")
		}
	} else {
		err = fmt.Errorf("environment variable not found")
	}
	return str, err
}

// AddEnvironmentSecret adds a secret to a pod
func (p *Pod) AddEnvironmentSecret(name, secretName string) *Pod {
	p.Environment[name] = EnvironmentSecret{
		Secret: secretName,
	}
	return p
}

// AddContainer adds a container to a pod
func (p *Pod) AddContainer(container *PodContainer) *Pod {
	p.Containers = append(p.Containers, container)
	return p
}

// EmptySecrets empties the secret sources in a pod
func (p *Pod) EmptySecrets() *Pod {
	p.Secrets = make(map[string]SecretSource)
	return p
}

// GetSecretSource gets the source of the named secret
func (p *Pod) GetSecretSource(name string) (string, error) {
	if val, ok := p.Secrets[name]; ok {
		return val.Source, nil
	}
	return "", fmt.Errorf("secret does not exist")
}

// AddSecret adds a secret source to a pod
func (p *Pod) AddSecret(name, value string) *Pod {
	p.Secrets[name] = SecretSource{
		Source: value,
	}
	return p
}

// AddVolume adds a volume to a pod
func (p *Pod) AddVolume(vol *PodVolume) *Pod {
	p.Volumes = append(p.Volumes, vol)
	return p
}

// AddNetwork adds a PodNetwork to a pod
func (p *Pod) AddNetwork(net *PodNetwork) *Pod {
	p.Networks = append(p.Networks, net)
	return p
}

// Count sets the count of the pod
func (p *Pod) Count(count int) *Pod {
	p.Scaling = &PodScalingPolicy{
		Kind:      "fixed",
		Instances: count,
	}
	return p
}

// SetPodSchedulingPolicy sets the PodSchedulingPolicy of a pod
func (p *Pod) SetPodSchedulingPolicy(policy *PodSchedulingPolicy) *Pod {
	p.Scheduling = policy
	return p
}

// SetExecutorResources sets the resources for the pod executor
func (p *Pod) SetExecutorResources(resources *ExecutorResources) *Pod {
	p.ExecutorResources = resources
	return p
}

// SupportsPods determines if this version of marathon supports pods
// If HEAD returns 200 it does
func (r *marathonClient) SupportsPods() bool {
	if err := r.apiHead(marathonAPIPods, nil, nil); err != nil {
		return false
	}

	return true
}

// GetPod gets a pod from marathon
// 		name: 		the id used to identify the pod
func (r *marathonClient) GetPod(name string) (*Pod, error) {
	uri := buildPodURI(name)
	result := new(Pod)
	if err := r.apiGet(uri, nil, result); err != nil {
		return nil, err
	}
	result.marathonClient = r

	return result, nil
}

// GetAllPods gets all pods from marathon
func (r *marathonClient) GetAllPods() ([]*Pod, error) {
	var result []*Pod
	if err := r.apiGet(marathonAPIPods, nil, &result); err != nil {
		return nil, err
	}

	for _, p := range result {
		p.marathonClient = r
	}

	return result, nil
}

// CreatePod creates a new pod in Marathon
// 		pod:		the structure holding the pod configuration
func (r *marathonClient) CreatePod(pod *Pod) (*Pod, error) {
	result := new(Pod)
	if err := r.apiPost(marathonAPIPods, &pod, result); err != nil {
		return nil, err
	}

	return result, nil
}

// DeletePod deletes a pod from marathon
// 		name: 		the id used to identify the pod
// 		force: 		whether
func (r *marathonClient) DeletePod(name string, force bool) (*DeploymentID, error) {
	uri := fmt.Sprintf("%s?force=%v", buildPodURI(name), force)
	// step: check of the pod already exists
	deployID := new(DeploymentID)
	if err := r.apiDelete(uri, nil, deployID); err != nil {
		return nil, err
	}

	return deployID, nil
}

// CreatePod creates a new pod in Marathon
// 		pod:		the structure holding the pod configuration
func (r *marathonClient) UpdatePod(pod *Pod, force bool) (*Pod, error) {
	uri := fmt.Sprintf("%s?force=%v", buildPodURI(pod.ID), force)
	result := new(Pod)
	// step: check of the pod already exists
	if err := r.apiPut(uri, pod, result); err != nil {
		return nil, err
	}

	return result, nil
}

// GetVersions gets the versions of a pod
// 		name:		the id of the pod
func (r *marathonClient) GetVersions(name string) ([]string, error) {
	uri := buildPodVersionURI(name)
	var result []string
	if err := r.apiGet(uri, nil, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// GetPodByVersion gets a pod by a version
// 		name:		the id of the pod
// 		version:	the version of the pod
func (r *marathonClient) GetPodByVersion(name, version string) (*Pod, error) {
	uri := fmt.Sprintf("%s/%s", buildPodVersionURI(name), version)
	result := new(Pod)
	if err := r.apiGet(uri, nil, result); err != nil {
		return nil, err
	}

	return result, nil
}

func buildPodVersionURI(name string) string {
	return fmt.Sprintf("%s/%s::versions", marathonAPIPods, trimRootPath(name))
}

func buildPodURI(path string) string {
	return fmt.Sprintf("%s/%s", marathonAPIPods, trimRootPath(path))
}

func (p *Pod) Versions() (*ApplicationVersions, error) {
	versions, err := p.marathonClient.GetVersions(p.ID)
	if err != nil {
		return nil, err
	}

	appVersions := new(ApplicationVersions)
	appVersions.Versions = versions
	return appVersions, nil
}

func (p *Pod) Update(force bool) error {
	_, err := p.marathonClient.UpdatePod(p, force)

	return err
}
