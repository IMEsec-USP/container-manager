package handlers

import (
	"io"
	"net/http"
	"time"

	"github.com/IMEsec-USP/container-manager/adapters"
	"github.com/IMEsec-USP/container-manager/applications"
	"github.com/IMEsec-USP/container-manager/middleware"
	"github.com/rs/zerolog"
)

func PullImage(
	tctx middleware.TimeoutContext,
	adapter *adapters.DockerAdapter,
	logger zerolog.Logger,
	app applications.Application,
) (int, string) {
	startClock := time.Now()

	reader, err := adapter.PullImage(tctx, app)
	if err != nil {
		logger.Error().TimeDiff("duration", time.Now(), startClock).Err(err)
		return http.StatusInternalServerError, err.Error()
	}

	go writeToLogs(logger, reader)

	return http.StatusOK, "pulled image"
}

func writeToLogs(logger zerolog.Logger, reader io.Reader) {
	_, err := io.Copy(logger.Level(zerolog.InfoLevel), reader)
	if err != nil {
		logger.Error().Err(err).Msg("failed to write logs from pulling image")
	}
}

func (h *HTTPHandler) RegisterPull() {
	h.Post("/pull/:appID", middleware.WithFiveMinuteTimeout, middleware.ParseAppID, PullImage)
}
