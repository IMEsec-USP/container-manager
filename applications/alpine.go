package applications

type Alpine struct{}

func (a Alpine) ImageReference() string {
	return "docker.io/library/alpine"
}

func (a Alpine) ContainerName() string {
	return "alpine_tosco"
}
