package applications

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

var applications = make(map[string]Application)

func BuildApplicationsFromConfig(logger zerolog.Logger) error {
	configApps := viper.Sub("applications")
	if configApps == nil {
		return errors.New("could not find configuration for apps")
	}

	// FIXME AllSettings() reads the whole configuration and puts it into a map,
	// instead of reading just the keys from one level deep. We don't need it to read
	// all the configuration, just to enumerate all the keys.
	// (AllKeys() doesn't work. It goes all layers deep as well.)
	for appName := range configApps.AllSettings() {
		logger.Info().Str("appName", appName).Msg("loading config from app")
		appConfig := configApps.Sub(appName)
		if appConfig == nil {
			return fmt.Errorf("something went terribly wrong. appConfig is nil")
		}
		app := &application{}

		// TODO we could use reflection over here with struct tags.
		imageReference := appConfig.GetString("image-reference")
		if imageReference == "" {
			return fmt.Errorf("app with name %s does not have key image-reference", appName)
		}
		app.imageReference = imageReference

		regexStrings := appConfig.GetStringSlice("matches-regexes")
		if len(regexStrings) == 0 {
			return fmt.Errorf("app with name %s does not define matches-regexes as a list of strings", appName)
		}
		app.regexes = make([]*regexp.Regexp, 0, len(regexStrings))
		for _, regexStr := range regexStrings {
			regexStr = "/" + regexStr
			regex, err := regexp.Compile(regexStr)
			if err != nil {
				return fmt.Errorf("error while processing regex from app %s: %w", regexStr, err)
			}
			logger.Info().Str("regexStr", regexStr).Msg("loaded regex")
			app.regexes = append(app.regexes, regex)
		}

		applications[appName] = app
	}

	return nil
}
