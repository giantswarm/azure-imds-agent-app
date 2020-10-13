package instance

import (
	"context"
	"net/http"

	"github.com/giantswarm/azure-imds-client/pkg/imds/instance"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
)

type Config struct {
	// Dependencies.
	Logger     micrologger.Logger
	HttpClient *http.Client
}

type Service struct {
	// Dependencies.
	logger     micrologger.Logger
	httpClient *http.Client

	instanceClient *instance.Client
}

func NewService(config Config) (*Service, error) {
	var err error

	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}
	if config.HttpClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.HttpClient must not be empty", config)
	}

	var instanceClient *instance.Client
	{
		c := instance.ClientConfig{
			Logger:     config.Logger,
			HttpClient: config.HttpClient,
		}
		instanceClient, err = instance.NewClient(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	service := &Service{
		logger:         config.Logger,
		httpClient:     config.HttpClient,
		instanceClient: instanceClient,
	}

	return service, nil
}

func (s *Service) GetMetadata(ctx context.Context) (metadata *instance.Metadata, err error) {
	instanceMetadata, err := s.instanceClient.GetMetadata(ctx)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return instanceMetadata, nil
}
