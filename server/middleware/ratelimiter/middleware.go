package ratelimiter

import (
	"context"
	"time"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	kitendpoint "github.com/go-kit/kit/endpoint"
	"github.com/ulule/limiter"
	"github.com/ulule/limiter/drivers/store/memory"
)

type Config struct {
	Logger micrologger.Logger

	// DefaultLimit is the default rate limit for all rate limiters. It is
	// currently not possible to specify a rate limit per endpoint.
	DefaultLimit int64
	// DefaultPeriod is the default period for all rate limiters. It is currently
	// not possible to specify a period per endpoint.
	DefaultPeriod time.Duration
}

type Middleware struct {
	logger micrologger.Logger

	ratelimiter *limiter.Limiter
}

func New(config Config) (*Middleware, error) {
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}

	if config.DefaultLimit == 0 {
		return nil, microerror.Maskf(invalidConfigError, "%T.DefaultLimit must not be empty", config)
	}
	if config.DefaultPeriod == 0 {
		return nil, microerror.Maskf(invalidConfigError, "%T.DefaultPeriod must not be empty", config)
	}

	var rateLimiter *limiter.Limiter
	{
		r := limiter.Rate{
			Period: config.DefaultPeriod,
			Limit:  config.DefaultLimit,
		}

		rateLimiter = limiter.New(memory.NewStore(), r)
	}

	m := &Middleware{
		logger:      config.Logger,
		ratelimiter: rateLimiter,
	}

	return m, nil
}

// New takes a function and returns a middleware, unlike other middlewares, it
// needs to know how to extract some information from the request which might be
// in different places depending on the endpoint.
func (m *Middleware) New(lookupFunc LookupFunc) kitendpoint.Middleware {
	return func(next kitendpoint.Endpoint) kitendpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			// Get the rateKey using the function that was passed in when we created
			// this middleware in the endpoint.
			rateKey, err := lookupFunc(ctx, request)
			if err != nil {
				return nil, microerror.Mask(err)
			}

			// Get the rate limit for the given rateKey. This also 'increments' the
			// request count at the same time.
			context, err := m.ratelimiter.Get(ctx, rateKey)
			if err != nil {
				return nil, microerror.Mask(err)
			}

			// Check if the limit has been reached.
			if context.Reached {
				return nil, microerror.Mask(limitReachedError)
			}

			// Limit hasn't been reached, so proceed to the next middleware in the
			// chain.
			return next(ctx, request)
		}
	}
}
