package ratelimiter

import (
	"github.com/giantswarm/microerror"
)

var invalidConfigError = &microerror.Error{
	Kind: "invalidConfigError",
}

// IsInvalidConfig asserts invalidConfigError.
func IsInvalidConfig(err error) bool {
	return microerror.Cause(err) == invalidConfigError
}

var limitReachedError = &microerror.Error{
	Kind: "limitReachedError",
}

// IsLimitReached asserts limitReachedError.
func IsLimitReached(err error) bool {
	return microerror.Cause(err) == limitReachedError
}
