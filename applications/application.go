package applications

type Application interface {
	ImageReference() string
	ContainerName() string
}

type ContainerAvaliator interface {
	Matches(containerName string) bool
}
