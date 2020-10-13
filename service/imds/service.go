package imds

import (
	"net/http"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"

	"github.com/giantswarm/azure-imds-agent-app/service/imds/instance"
)

type Config struct {
	// Dependencies.
	Logger     micrologger.Logger
	HttpClient *http.Client
}

type Service struct {
	InstanceService *instance.Service
}

func NewService(config Config) (*Service, error) {
	var err error
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}
	if config.HttpClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.HttpClient must not be empty", config)
	}

	var instanceService *instance.Service
	{
		c := instance.Config{
			Logger:     config.Logger,
			HttpClient: config.HttpClient,
		}
		instanceService, err = instance.NewService(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	service := &Service{
		InstanceService: instanceService,
	}

	return service, nil
}
