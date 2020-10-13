package middleware

import (
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/spf13/viper"

	"github.com/giantswarm/azure-imds-agent-app/flag"
	"github.com/giantswarm/azure-imds-agent-app/server/middleware/ratelimiter"
	"github.com/giantswarm/azure-imds-agent-app/service"
)

// Config represents the configuration used to create a middleware.
type Config struct {
	// Dependencies.
	Flag    *flag.Flag
	Logger  micrologger.Logger
	Service *service.Service
	Viper   *viper.Viper
}

type Middleware struct {
	RateLimiter *ratelimiter.Middleware
}

// New creates a new configured middleware.
func New(config Config) (*Middleware, error) {
	var err error

	var ratelimiterMiddleware *ratelimiter.Middleware
	{
		c := ratelimiter.Config{
			Logger: config.Logger,

			DefaultPeriod: config.Viper.GetDuration(config.Flag.Middleware.RateLimiter.DefaultPeriod),
			DefaultLimit:  config.Viper.GetInt64(config.Flag.Middleware.RateLimiter.DefaultLimit),
		}

		ratelimiterMiddleware, err = ratelimiter.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	middleware := &Middleware{
		RateLimiter: ratelimiterMiddleware,
	}

	return middleware, nil
}
