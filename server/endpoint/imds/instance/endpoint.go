package instance

import (
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"

	"github.com/giantswarm/azure-imds-agent-app/server/endpoint/imds/instance/getter"
	"github.com/giantswarm/azure-imds-agent-app/server/middleware"
	"github.com/giantswarm/azure-imds-agent-app/service"
)

// Config represents the configuration used to create a cluster endpoint.
type Config struct {
	Logger     micrologger.Logger
	Middleware *middleware.Middleware
	Service    *service.Service
}

type Endpoint struct {
	Getter *getter.Endpoint
}

// New creates a new configured app config endpoint.
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

	var getterEndpoint *getter.Endpoint
	{
		c := getter.Config{
			Logger:          config.Logger,
			Middleware:      config.Middleware,
			InstanceService: config.Service.IMDSService.InstanceService,
		}

		getterEndpoint, err = getter.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	endpoint := &Endpoint{
		Getter: getterEndpoint,
	}

	return endpoint, nil
}
