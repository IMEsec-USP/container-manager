package applications

import (
	"regexp"
)

type Application interface {
	ContainerAvaliator
	ImageReference() string
}

type ContainerAvaliator interface {
	Matches(containerName string) bool
}

type application struct {
	imageReference string
	regexes        []*regexp.Regexp
}

func (a *application) ImageReference() string {
	return a.imageReference
}

func (a *application) Matches(containerName string) bool {
	for _, regex := range a.regexes {
		if regex.MatchString(containerName) {
			return true
		}
	}
	return false
}
