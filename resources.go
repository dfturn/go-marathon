package marathon

type ExecutorResources struct {
	Cpus float64 `json:"cpus,omitempty"`
	Mem  float64 `json:"mem,omitempty"`
	Disk float64 `json:"disk,omitempty"`
}

type Resources struct {
	Cpus float64 `json:"cpus"`
	Mem  float64 `json:"mem"`
	Disk float64 `json:"disk,omitempty"`
	Gpus int32   `json:"gpus,omitempty"`
}

func NewResources() *Resources {
	return &Resources{}
}
