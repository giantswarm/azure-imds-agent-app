// Package url stores and accesses the HTTP URL in and from
// context.Context.
package url

import (
	"context"
)

// key is an unexported type for keys defined in this package. This prevents
// collisions with keys defined in other packages.
type key string

// urlKey is the key for HTTP url values in context.Context.
// Clients use url.NewContext and url.FromContext instead of using this
// key directly.
var urlKey key = "url"

// NewContext returns a new context.Context that carries value h.
func NewContext(ctx context.Context, h string) context.Context {
	if h == "" {
		// In case the url is empty we do not add it, but only return the given
		// context as it is. That way the existence check when reading the context
		// works as expected when no url or an empty url was tried to be
		// added.
		return ctx
	}

	return context.WithValue(ctx, urlKey, h)
}

// FromContext returns the HTTP url value stored in ctx, if
// any.
func FromContext(ctx context.Context) string {
	h, _ := ctx.Value(urlKey).(string)
	return h
}
