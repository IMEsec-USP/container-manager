package handlers

import (
	"github.com/go-martini/martini"
)

type HTTPHandler struct {
	*martini.ClassicMartini
}

func NewHTTPHandler() *HTTPHandler {
	return &HTTPHandler{
		martini.Classic(),
	}
}
