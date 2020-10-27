package applications

import (
	"errors"
)

func ParseApplication(maybeAppID string) (Application, error) {
	if app, ok := applications[maybeAppID]; ok {
		return app, nil
	}

	return nil, errors.New("could not find an application with such ID")
}
