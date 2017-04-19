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
}

type PodScalingPolicy struct {
	Kind         string `json:"kind"`
	Instances    int    `json:"instances"`
	MaxInstances int    `json:"maxInstances,omitempty"`
}

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

func (p *Pod) EmptyLabels() *Pod {
	p.Labels = make(map[string]string)
	return p
}

func (p *Pod) AddLabel(key, value string) *Pod {
	p.Labels[key] = value
	return p
}

func (p *Pod) SetLabels(labels map[string]string) *Pod {
	p.Labels = labels
	return p
}

func (p *Pod) EmptyEnvironment() *Pod {
	p.Environment = make(map[string]interface{})
	return p
}

func (p *Pod) AddEnvironment(name, value string) *Pod {
	p.Environment[name] = value
	return p
}

func (p *Pod) SetEnvironment(env map[string]string) *Pod {
	for k, v := range env {
		p.AddEnvironment(k, v)
	}
	return p
}

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

func (p *Pod) AddEnvironmentSecret(name, secretName string) *Pod {
	p.Environment[name] = EnvironmentSecret{
		Secret: secretName,
	}
	return p
}

func (p *Pod) AddContainer(container *PodContainer) *Pod {
	p.Containers = append(p.Containers, container)
	return p
}

func (p *Pod) EmptySecrets() *Pod {
	p.Secrets = make(map[string]SecretSource)
	return p
}

func (p *Pod) GetSecretSource(name string) (string, error) {
	if val, ok := p.Secrets[name]; ok {
		return val.Source, nil
	} else {
		return "", fmt.Errorf("secret does not exist")
	}
}

func (p *Pod) AddSecret(name, value string) *Pod {
	p.Secrets[name] = SecretSource{
		Source: value,
	}
	return p
}

func (p *Pod) AddVolume(vol *PodVolume) *Pod {
	p.Volumes = append(p.Volumes, vol)
	return p
}

func (p *Pod) AddNetwork(net *PodNetwork) *Pod {
	p.Networks = append(p.Networks, net)
	return p
}

func (p *Pod) Count(count int) *Pod {
	p.Scaling = &PodScalingPolicy{
		Kind:      "fixed",
		Instances: count,
	}
	return p
}

func (p *Pod) SetPodSchedulingPolicy(policy *PodSchedulingPolicy) *Pod {
	p.Scheduling = policy
	return p
}

func (p *Pod) SetExecutorResources(resources *ExecutorResources) *Pod {
	p.ExecutorResources = resources
	return p
}

// CreatePod creates a new pod in Marathon
// 		application:		the structure holding the pod configuration
func (r *marathonClient) CreatePod(pod *Pod) (*Pod, error) {
	result := new(Pod)
	if err := r.apiPost(marathonAPIPods, &pod, result); err != nil {
		return nil, err
	}

	return result, nil
}

// DeletePod deletes a pod from marathon
// 		name: 		the id used to identify the pod
func (r *marathonClient) DeletePod(name string) (*DeploymentID, error) {
	uri := buildPodURI(name)
	// step: check of the pod already exists
	deployID := new(DeploymentID)
	if err := r.apiDelete(uri, nil, deployID); err != nil {
		return nil, nil // TODO: Get headers and get the deployment ID
	}

	return deployID, nil
}

// CreatePod creates a new pod in Marathon
// 		application:		the structure holding the pod configuration
func (r *marathonClient) UpdatePod(pod *Pod) (*Pod, error) {
	uri := buildPodURI(pod.ID)
	result := new(Pod)
	// step: check of the pod already exists
	if err := r.apiPut(uri, nil, result); err != nil {
		return nil, nil // TODO: Get headers and get the deployment ID
	}

	return result, nil
}

func buildPodURI(path string) string {
	return fmt.Sprintf("%s/%s", marathonAPIPods, trimRootPath(path))
}

// TODO: Add other APIs
