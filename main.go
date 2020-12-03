package main

import (
	"os"

	"github.com/IMEsec-USP/container-manager/adapters"
	"github.com/IMEsec-USP/container-manager/applications"
	"github.com/IMEsec-USP/container-manager/handlers"
	"github.com/go-martini/martini"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

func main() {
	logger := zerolog.New(os.Stdout).Output(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Logger()
	logger.Info().Msg("started logger")

	readConfigurations(logger)
	err := applications.BuildApplicationsFromConfig(logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("cannot continue without correct application configuration")
	}

	logger.Info().Msg("loaded configs")

	dockerAdapter, err := adapters.NewDockerAdapter()
	if err != nil {
		logger.Fatal().Err(err).Msg("cannot continue without a docker client")
	}

	martini.Env = martini.Prod
	h := handlers.NewHTTPHandler()
	{
		h.Map(dockerAdapter)
		h.Map(logger)
		h.RegisterHealthCheck()
		h.RegisterRestart()
		h.RegisterPull()
	}
	h.RunOnAddr(viper.GetString("host"))
}

func readConfigurations(logger zerolog.Logger) {
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("/configs")
	{
		viper.SetDefault("host", "0.0.0.0")
		viper.SetDefault("port", "3000")
	}
	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			logger.Warn().Err(err).Msg("could not find any config files. Continuing with defaults")
		} else {
			logger.Fatal().Err(err).Msg("could not read from config files")
		}
	}
}
