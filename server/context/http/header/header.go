// Package header stores and accesses the HTTP Authorization header in and from
// context.Context.
package header

import (
	"context"
)

// key is an unexported type for keys defined in this package. This prevents
// collisions with keys defined in other packages.
type key string

// headerKey is the key for HTTP Authorization header values in context.Context.
// Clients use header.NewContext and header.FromContext instead of using this
// key directly.
var headerKey key = "header"

// NewContext returns a new context.Context that carries value h.
func NewContext(ctx context.Context, h string) context.Context {
	if h == "" {
		// In case the header is empty we do not add it, but only return the given
		// context as it is. That way the existence check when reading the context
		// works as expected when no header or an empty header was tried to be
		// added.
		return ctx
	}

	return context.WithValue(ctx, headerKey, h)
}

// FromContext returns the HTTP Authorization header value stored in ctx, if
// any.
func FromContext(ctx context.Context) (string, bool) {
	h, ok := ctx.Value(headerKey).(string)
	return h, ok
}
