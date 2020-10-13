package service

import (
	"net/http"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/spf13/viper"

	"github.com/giantswarm/azure-imds-agent-app/flag"
	"github.com/giantswarm/azure-imds-agent-app/service/imds"
)

type Config struct {
	Logger micrologger.Logger

	Description string
	Flag        *flag.Flag
	GitCommit   string
	Name        string
	Source      string
	Viper       *viper.Viper

	HttpClient *http.Client
}

type Service struct {
	Logger micrologger.Logger
	Flag   *flag.Flag
	Viper  *viper.Viper

	IMDSService *imds.Service
}

func New(config Config) (*Service, error) {
	if config.Flag == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Flag must not be empty", config)
	}
	if config.Viper == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Viper must not be empty", config)
	}
	if config.HttpClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.HttpClient must not be empty", config)
	}

	s, err := NewMiddleware(config.Logger, config.HttpClient, config.Description, config.GitCommit, config.Name, config.Source, config.Flag, config.Viper)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return s, nil
}

func NewMiddleware(logger micrologger.Logger, httpClient *http.Client, description, gitCommit, name, source string, f *flag.Flag, v *viper.Viper) (*Service, error) {
	var err error

	var imdsService *imds.Service
	{
		c := imds.Config{
			Logger:     logger,
			HttpClient: httpClient,
		}
		imdsService, err = imds.NewService(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	service := &Service{
		Logger:      logger,
		Flag:        f,
		Viper:       v,
		IMDSService: imdsService,
	}

	return service, nil
}
