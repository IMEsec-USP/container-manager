package applications

type LHTimer struct{}

func (l LHTimer) ContainerName() string {
	return "test_timer"
}

func (l LHTimer) ImageReference() string {
	return "docker.io/lightninghacks/timer"
}
