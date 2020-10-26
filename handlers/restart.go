package handlers

import (
	"net/http"
	"time"

	"github.com/IMEsec-USP/container-manager/adapters"
	"github.com/IMEsec-USP/container-manager/applications"
	"github.com/IMEsec-USP/container-manager/middleware"
	"github.com/rs/zerolog"
)

func RestartContainer(
	tctx middleware.TimeoutContext,
	adapter *adapters.DockerAdapter,
	logger zerolog.Logger,
	app applications.Application,
) (int, string) {
	logger.Info().Msg("getting app config")

	startClock := time.Now()

	containerConfig, err := adapter.GetContainerConfig(tctx, app)
	if err != nil {
		// TODO the container may already be dead. In that case, we could have
		// the initial settings of the containers so that we can restart them somehow.
		logger.Error().TimeDiff("duration", time.Now(), startClock).Err(err).Msg("could not get container config")
		return http.StatusInternalServerError, err.Error()
	}

	err = adapter.RemoveContainer(tctx, containerConfig)
	if err != nil {
		logger.Error().TimeDiff("duration", time.Now(), startClock).Err(err).Msg("could not kill container")
		return http.StatusInternalServerError, err.Error()
	}

	err = adapter.RunImage(tctx, app, containerConfig)
	if err != nil {
		logger.Error().TimeDiff("duration", time.Now(), startClock).Err(err).Msg("Could not run image with desired configurations")
		return http.StatusInternalServerError, err.Error()
	}

	logger.Info().TimeDiff("duration", time.Now(), startClock).Msg("restarted")

	return http.StatusOK, "restarted container"
}

// func logOutput(logger zerolog.Logger, output io.Reader) {
// 	_, err := io.Copy(logger, output)
// 	if err != nil {
// 		logger.Err(err)
// 	}
// }

func (h *HTTPHandler) RegisterRestart() {
	h.Post("/restart/:appID", middleware.WithFiveMinuteTimeout, middleware.ParseAppID, RestartContainer)
}
