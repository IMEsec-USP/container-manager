package applications

type Application interface {
	ImageReference() string
	ContainerName() string
}

type ContainerAvaliator interface {
	Matches(containerName string) bool
}

type application struct {
	imageReference string
	containerName  string
	regexes        []string
}

func (a *application) ImageReference() string {
	return a.imageReference
}

func (a *application) ContainerName() string {
	return a.containerName
}
