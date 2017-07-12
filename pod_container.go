package marathon

type PodContainer struct {
	Name         string                 `json:"name,omitempty"`
	Exec         *PodExec               `json:"exec,omitempty"`
	Resources    *Resources             `json:"resources,omitempty"`
	Endpoints    []*PodEndpoint         `json:"endpoints,omitempty"`
	Image        *PodContainerImage     `json:"image,omitempty"`
	Environment  map[string]interface{} `json:"environment,omitempty"`
	User         string                 `json:"user,omitempty"`
	HealthCheck  *PodHealthCheck        `json:"healthCheck,omitempty"`
	VolumeMounts []*PodVolumeMount      `json:"volumeMounts,omitempty"`
	Artifacts    []*PodArtifact         `json:"artifacts,omitempty"`
	Labels       map[string]string      `json:"labels,omitempty"`
	Lifecycle    PodLifecycle           `json:"lifecycle,omitempty"`
}

type PodLifecycle struct {
	KillGracePeriodSeconds float64 `json:"killGracePeriodSeconds,omitempty"`
}

type PodCommand struct {
	Shell string `json:"shell,omitempty"`
}

type PodExec struct {
	Command PodCommand `json:"command,omitempty"`
}

type PodArtifact struct {
	Uri        string `json:"uri,omitempty"`
	Extract    bool   `json:"extract,omitempty"`
	Executable bool   `json:"executable,omitempty"`
	Cache      bool   `json:"cache,omitempty"`
	DestPath   string `json:"destPath,omitempty"`
}

func NewPodContainer() *PodContainer {
	return &PodContainer{
		Endpoints:    []*PodEndpoint{},
		Environment:  map[string]interface{}{},
		VolumeMounts: []*PodVolumeMount{},
		Artifacts:    []*PodArtifact{},
		Labels:       map[string]string{},
		Resources:    NewResources(),
	}
}

func (p *PodContainer) SetName(name string) *PodContainer {
	p.Name = name
	return p
}

func (p *PodContainer) SetCommand(name string) *PodContainer {
	p.Exec = &PodExec{
		Command: PodCommand{
			Shell: name,
		},
	}
	return p
}

func (r *PodContainer) CPUs(cpu float64) *PodContainer {
	r.Resources.Cpus = cpu
	return r
}

func (r *PodContainer) Memory(memory float64) *PodContainer {
	r.Resources.Mem = memory
	return r
}

func (r *PodContainer) Storage(disk float64) *PodContainer {
	r.Resources.Disk = disk
	return r
}

func (r *PodContainer) GPUs(gpu int32) *PodContainer {
	r.Resources.Gpus = gpu
	return r
}

func (r *PodContainer) AddEndpoint(endpoint *PodEndpoint) *PodContainer {
	r.Endpoints = append(r.Endpoints, endpoint)
	return r
}

func (r *PodContainer) SetImage(image *PodContainerImage) *PodContainer {
	r.Image = image
	return r
}

func (p *PodContainer) AddEnvironment(name, value string) *PodContainer {
	p.Environment[name] = value
	return p
}

func (p *PodContainer) SetEnvironment(env map[string]string) *PodContainer {
	for k, v := range env {
		p.AddEnvironment(k, v)
	}
	return p
}

func (p *PodContainer) AddEnvironmentSecret(name, secretName string) *PodContainer {
	p.Environment[name] = EnvironmentSecret{
		Secret: secretName,
	}
	return p
}

// SetUser sets the user to run the pod as
func (p *PodContainer) SetUser(user string) *PodContainer {
	p.User = user
	return p
}

func (r *PodContainer) SetHealthCheck(healthcheck *PodHealthCheck) *PodContainer {
	r.HealthCheck = healthcheck
	return r
}

func (r *PodContainer) AddVolumeMount(mount *PodVolumeMount) *PodContainer {
	r.VolumeMounts = append(r.VolumeMounts, mount)
	return r
}

func (r *PodContainer) AddArtifact(artifact *PodArtifact) *PodContainer {
	r.Artifacts = append(r.Artifacts, artifact)
	return r
}

func (p *PodContainer) AddLabel(key, value string) *PodContainer {
	p.Labels[key] = value
	return p
}

func (p *PodContainer) SetLifecycle(lifecycle PodLifecycle) *PodContainer {
	p.Lifecycle = lifecycle
	return p
}
