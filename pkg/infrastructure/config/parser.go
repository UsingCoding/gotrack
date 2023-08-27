package config

import (
	"encoding/json"
	"os"

	"github.com/pkg/errors"

	appconfig "gotrack/pkg/app/config"
)

type Parser struct{}

func (p Parser) Parse(path string) (appconfig.Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return appconfig.Config{}, errors.Wrap(err, "failed to read config data")
	}

	var c config
	err = json.Unmarshal(data, &c)
	if err != nil {
		return appconfig.Config{}, errors.Wrap(err, "failed to unmarshal config data")
	}

	return appconfig.Config{
		YouTrackHost: c.YouTrackHost,
		Token:        c.Token,
	}, err
}

type config struct {
	YouTrackHost string `json:"youTrackHost"`
	Token        string `json:"token"`
}
