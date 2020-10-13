package middleware

import (
	"github.com/giantswarm/azure-imds-agent-app/flag/middleware/ratelimiter"
)

type Middleware struct {
	RateLimiter ratelimiter.RateLimiter
}
