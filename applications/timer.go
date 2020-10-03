package applications

type LHTimer struct{}

func (l LHTimer) ContainerName() string {
	return "docker-compose_timer_1"
}

func (l LHTimer) ImageReference() string {
	return "docker.io/lightninghacks/timer"
}
