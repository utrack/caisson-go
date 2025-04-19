package caiapp

import "github.com/utrack/caisson-go/pkg/http/hhandler"


type Handlers struct {
	http hhandler.Server
}

func (c *Handlers) HTTP() hhandler.Configurer {
	return c.http
}