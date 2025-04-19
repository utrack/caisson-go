package cappconfig

import (
	"time"

	"github.com/utrack/caisson-go/errors"
	"github.com/utrack/envconfig"
)

type Config struct {
	Server           Server
	GracefulShutdown Grace
}

type Server struct {
	AddrHTTP string `default:"0.0.0.0"`
	PortHTTP int    `default:"8080"`

	AddrDebug string `default:"0.0.0.0"`
	PortDebug int    `default:"8082"`
}

type Grace struct {
	// Delay is the time between SIGTERM and the graceful shutdown commencement.
	// Used to give the ingress time to stop routing traffic to the server.
	Delay time.Duration `default:"5s"`
	// Timeout is the time between the graceful shutdown commencement and the server being stopped.
	Timeout time.Duration `default:"30s"`
}

func read() (*Config, error) {
	var c Config
	err := envconfig.ProcessWithOptions("", &c, envconfig.Options{SplitWords: true})
	if err != nil {
		return nil, errors.Wrap(err, "failed to read envvars")
	}
	return &c, nil
}

var c *Config

func Get() (*Config, error) {
	if c == nil {
		var err error
		c, err = read()
		if err != nil {
			return nil, errors.Wrap(err, "when reading caisson caiapp config")
		}
	}
	return c, nil
}
