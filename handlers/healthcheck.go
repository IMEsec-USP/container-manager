package handlers

import (
	"net/http"
)

func HealthCheck(res http.ResponseWriter) {
	res.WriteHeader(http.StatusOK)
}

func (h *HTTPHandler) RegisterHealthCheck() {
	h.Get("/_healthcheck", HealthCheck)
}
