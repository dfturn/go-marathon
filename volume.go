package marathon

type PodVolume struct {
	Name string `json:"name,omitempty"`
	Host string `json:"host,omitempty"`
}

type PodVolumeMount struct {
	Name      string `json:"name,omitempty"`
	MountPath string `json:"mountPath,omitempty"`
}

func NewPodVolume(name, path string) *PodVolume {
	return &PodVolume{
		Name: name,
		Host: path,
	}
}

func NewPodVolumeMount(name, mount string) *PodVolumeMount {
	return &PodVolumeMount{
		Name:      name,
		MountPath: mount,
	}
}
