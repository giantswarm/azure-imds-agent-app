// Package scheme stores and accesses the HTTP Authorization scheme in and from
// context.Context.
package scheme

import (
	"context"
)

// key is an unexported type for keys defined in this package. This prevents
// collisions with keys defined in other packages.
type key string

// schemeKey is the key for HTTP Authorization scheme values in context.Context.
// Clients use scheme.NewContext and scheme.FromContext instead of using this
// key directly.
var schemeKey key = "scheme"

// NewContext returns a new context.Context that carries value s.
func NewContext(ctx context.Context, s string) context.Context {
	if s == "" {
		// In case the scheme is empty we do not add it, but only return the given
		// context as it is. That way the existence check when reading the context
		// works as expected when no scheme or an empty scheme was tried to be
		// added.
		return ctx
	}

	return context.WithValue(ctx, schemeKey, s)
}

// FromContext returns the HTTP Authorization scheme value stored in ctx, if
// any.
func FromContext(ctx context.Context) (string, bool) {
	s, ok := ctx.Value(schemeKey).(string)
	return s, ok
}
