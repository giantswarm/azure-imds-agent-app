// Package method stores and accesses the HTTP method in and from
// context.Context.
package method

import (
	"context"
)

// key is an unexported type for keys defined in this package. This prevents
// collisions with keys defined in other packages.
type key string

// methodKey is the key for HTTP method values in context.Context.
// Clients use method.NewContext and method.FromContext instead of using this
// key directly.
var methodKey key = "method"

// NewContext returns a new context.Context that carries value h.
func NewContext(ctx context.Context, h string) context.Context {
	if h == "" {
		// In case the method is empty we do not add it, but only return the given
		// context as it is. That way the existence check when reading the context
		// works as expected when no method or an empty method was tried to be
		// added.
		return ctx
	}

	return context.WithValue(ctx, methodKey, h)
}

// FromContext returns the HTTP method value stored in ctx, if
// any.
func FromContext(ctx context.Context) string {
	h, _ := ctx.Value(methodKey).(string)
	return h
}
