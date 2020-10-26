package applications

import (
	"errors"
)

var applications = map[string]Application{
	"ada": &application{
		imageReference: "docker.io/imesec/ada",
		containerName:  "services_ada_1",
	},
	"lh-timer": &application{
		imageReference: "docker.io/lightninghacks/timer",
		containerName:  "test-timer",
	},
}

func ParseApplication(maybeAppID string) (Application, error) {
	if app, ok := applications[maybeAppID]; ok {
		return app, nil
	}

	return nil, errors.New("could not find an application with such ID")
}
