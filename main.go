package main

import (
	"log"
	"os"

	"github.com/IMEsec-USP/container-manager/adapters"
	"github.com/IMEsec-USP/container-manager/handlers"
	"github.com/rs/zerolog"
)

func main() {
	dockerAdapter, err := adapters.NewDockerAdapter()
	if err != nil {
		log.Fatal(err)
	}

	logger := zerolog.New(os.Stdout).Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()
	logger.Info().Msg("started logger")

	h := handlers.NewHTTPHandler()
	{
		h.Map(dockerAdapter)
		h.Map(logger)
		h.RegisterHealthCheck()
		h.RegisterRestart()
		h.RegisterPull()
	}
	h.Run()
}
