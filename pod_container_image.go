package marathon

type ImageType string

const (
	ImageTypeDocker ImageType = "DOCKER"
	ImageTypeAppC   ImageType = "APPC"
)

type PodContainerImage struct {
	Kind      ImageType `json:"kind,omitempty"`
	ID        string    `json:"id,omitempty"`
	ForcePull bool      `json:"forcePull,omitempty"`
}

func NewPodContainerImage() *PodContainerImage {
	return &PodContainerImage{}
}

func (i *PodContainerImage) SetKind(typ ImageType) *PodContainerImage {
	i.Kind = typ
	return i
}

func (i *PodContainerImage) SetID(id string) *PodContainerImage {
	i.ID = id
	return i
}

func NewDockerPodContainerImage() *PodContainerImage {
	return NewPodContainerImage().SetKind(ImageTypeDocker)
}
