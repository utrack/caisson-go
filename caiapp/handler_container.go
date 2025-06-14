package caiapp

import (
	"github.com/utrack/caisson-go/caiapp/internal/hchi"
	"github.com/utrack/caisson-go/pkg/http/hhandler"
)

type Handlers struct {
	http *hchi.ChiHandler
}

func (c *Handlers) HTTP() hhandler.Configurer {
	return c.http
}
