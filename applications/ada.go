package applications

type Ada struct{}

func (a Ada) ImageReference() string {
	return "docker.io/imesec/ada"
}

func (a Ada) ContainerName() string {
	return "services_ada_1"
}
