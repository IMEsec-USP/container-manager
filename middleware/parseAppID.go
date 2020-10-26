package middleware

import (
	"net/http"

	"github.com/IMEsec-USP/container-manager/applications"
	"github.com/go-martini/martini"
	"github.com/rs/zerolog"
)

func ParseAppID(
	c martini.Context,
	params martini.Params,
	res http.ResponseWriter,
	logger zerolog.Logger,
) {
	appID := params["appID"]
	app, err := applications.ParseApplication(appID)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		_, werr := res.Write([]byte(err.Error()))
		if werr != nil {
			logger.Error().Err(werr).Msg("could not write to response stream")
		}
	}

	logger = logger.With().Str("appID", appID).Logger()
	c.Map(app).Map(logger)
}
