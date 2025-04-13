package caisconfig

import (
	"runtime/debug"

	"github.com/utrack/caisson-go/errors"
	"github.com/utrack/envconfig"
)

type Config struct {
	ServiceName string
	Otel        TelemetryConfig
}

type TelemetryConfig struct {
	Enable       bool `required:"true"` // required so that the telemetry isn't accidentally off on prod (explicit v implicit)
	CollectorURL string
	Insecure     bool
}

func read() (*Config, error) {
	var c Config
	err := envconfig.Process("", &c)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read caisson config")
	}
	if c.ServiceName == "" {
		dInfo, ok := debug.ReadBuildInfo()
		if !ok {
			return nil, errors.Wrap(err, "failed to read Go build info and Caisson service_name is empty")
		}
		c.ServiceName = dInfo.Main.Path
	}
	return &c, nil
}

var c *Config

func Get() *Config {
	if c == nil {
		var err error
		c, err = read()
		if err != nil {
			panic(err)
		}
	}
	return c
}
