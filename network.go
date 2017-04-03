package marathon

type PodNetwork struct {
	Name   string            `json:"name,omitempty"`
	Mode   string            `json:"mode,omitempty"`
	Labels map[string]string `json:"labels,omitempty"`
}

type PodEndpoint struct {
	Name          string            `json:"name,omitempty"`
	ContainerPort int               `json:"containerPort,omitempty"`
	HostPort      int               `json:"hostPort,omitempty"`
	Protocol      []string          `json:"protocol,omitempty"`
	Labels        map[string]string `json:"labels,omitempty"`
}

func NewPodNetwork(name string) *PodNetwork {
	return &PodNetwork{
		Name:   name,
		Labels: map[string]string{},
	}
}

func NewPodEndpoint() *PodEndpoint {
	return &PodEndpoint{
		Protocol: []string{},
		Labels:   map[string]string{},
	}
}

func NewContainerPodNetwork(name string) *PodNetwork {
	pn := NewPodNetwork(name)
	return pn.SetMode("container")
}

func (n *PodNetwork) SetName(name string) *PodNetwork {
	n.Name = name
	return n
}

func (n *PodNetwork) SetMode(mode string) *PodNetwork {
	n.Mode = mode
	return n
}

func (n *PodNetwork) Label(key, value string) *PodNetwork {
	n.Labels[key] = value
	return n
}

func (e *PodEndpoint) SetName(name string) *PodEndpoint {
	e.Name = name
	return e
}

func (e *PodEndpoint) SetContainerPort(port int) *PodEndpoint {
	e.ContainerPort = port
	return e
}

func (e *PodEndpoint) SetHostPort(port int) *PodEndpoint {
	e.HostPort = port
	return e
}

func (e *PodEndpoint) AddProtocol(protocol string) *PodEndpoint {
	e.Protocol = append(e.Protocol, protocol)
	return e
}

func (e *PodEndpoint) Label(key, value string) *PodEndpoint {
	e.Labels[key] = value
	return e
}
