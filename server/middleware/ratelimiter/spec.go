package ratelimiter

import (
	"context"
)

// LookupFunc is the middleware configuration required to lookup the rate limit
// key of the user making the current request. This has to be configured
// independently since the rate limit key is provided differently with each
// endpoint.
type LookupFunc func(ctx context.Context, request interface{}) (string, error)
