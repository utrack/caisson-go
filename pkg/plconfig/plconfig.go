package plconfig

import (
	"runtime/debug"
	"strings"

	"github.com/utrack/caisson-go/errors"
	"github.com/utrack/envconfig"
)

type Config struct {
	ServiceName string
	Otel        TelemetryConfig
}

type TelemetryConfig struct {
	Enable            bool `required:"true"` // required so that the telemetry isn't accidentally off on prod (explicit v implicit)
	CollectorEndpoint string
	CollectorInsecure bool
}

func read() (*Config, error) {
	var c Config
	err := envconfig.ProcessWithOptions("", &c, envconfig.Options{SplitWords: true})
	if err != nil {
		return nil, errors.Wrap(err, "failed to read caisson config")
	}
	if c.ServiceName == "" {
		dInfo, ok := debug.ReadBuildInfo()
		if !ok {
			return nil, errors.Wrap(err, "failed to read Go build info and service_name is empty")
		}
		c.ServiceName = dInfo.Main.Path
	}

	c.ServiceName = strings.ReplaceAll(c.ServiceName, "/", "-")

	if c.Otel.Enable && c.Otel.CollectorEndpoint == "" {
		return nil, errors.Errorf("caisson/baseconfig: OTEL_COLLECTOR_ENDPOINT is required when OTEL_ENABLE is true")
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
