package endpoint

import (
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/spf13/viper"

	"github.com/giantswarm/azure-imds-agent-app/flag"
	"github.com/giantswarm/azure-imds-agent-app/server/endpoint/imds"
	"github.com/giantswarm/azure-imds-agent-app/server/middleware"
	"github.com/giantswarm/azure-imds-agent-app/service"
)

// Config represents the configuration used to create a endpoint.
type Config struct {
	// Dependencies.
	Flag       *flag.Flag
	Logger     micrologger.Logger
	Middleware *middleware.Middleware
	Service    *service.Service
	Viper      *viper.Viper
}

// DefaultConfig provides a default configuration to create a new endpoint by
// best effort.
func DefaultConfig() Config {
	return Config{
		// Dependencies.
		Flag:       nil,
		Logger:     nil,
		Middleware: nil,
		Service:    nil,
	}
}

// Endpoint is the endpoint collection.
type Endpoint struct {
	IMDS *imds.Endpoint
}

// New creates a new configured endpoint.
func New(config Config) (*Endpoint, error) {
	if config.Flag == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Flag must not be empty", config)
	}
	if config.Viper == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Viper must not be empty", config)
	}

	var err error

	var imdsEndpoint *imds.Endpoint
	{
		c := imds.Config{
			Logger:     config.Logger,
			Middleware: config.Middleware,
			Service:    config.Service,
		}
		imdsEndpoint, err = imds.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	endpoint := &Endpoint{
		IMDS: imdsEndpoint,
	}

	return endpoint, nil
}
