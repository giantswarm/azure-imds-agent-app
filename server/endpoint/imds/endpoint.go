package imds

import (
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"

	"github.com/giantswarm/azure-imds-agent-app/server/endpoint/imds/instance"
	"github.com/giantswarm/azure-imds-agent-app/server/middleware"
	"github.com/giantswarm/azure-imds-agent-app/service"
)

// Config represents the configuration used to create a cluster endpoint.
type Config struct {
	Logger     micrologger.Logger
	Middleware *middleware.Middleware
	Service    *service.Service
}

// Endpoint contains a collection of endpoints
type Endpoint struct {
	Instance *instance.Endpoint
}

// New creates a new configured IMDS endpoint.
func New(config Config) (*Endpoint, error) {
	var err error

	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}
	if config.Middleware == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Middleware must not be empty", config)
	}
	if config.Service == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Service must not be empty", config)
	}

	var instanceEndpoint *instance.Endpoint
	{
		c := instance.Config{
			Logger:     config.Logger,
			Middleware: config.Middleware,
			Service:    config.Service,
		}
		instanceEndpoint, err = instance.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	endpoint := &Endpoint{
		Instance: instanceEndpoint,
	}

	return endpoint, nil
}
