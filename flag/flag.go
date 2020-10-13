package flag

import (
	"github.com/giantswarm/microkit/flag"

	"github.com/giantswarm/azure-imds-agent-app/flag/middleware"
	"github.com/giantswarm/azure-imds-agent-app/flag/server"
)

type Flag struct {
	Middleware middleware.Middleware
	Server     server.Server
}

func New() *Flag {
	f := &Flag{}
	flag.Init(f)
	return f
}
