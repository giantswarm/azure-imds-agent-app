package clusterid

import (
	"context"
)

type key string

var clusterIDKey key = "clusterID"

func NewContext(ctx context.Context, v string) context.Context {
	if v == "" {
		// In case the cluster ID is empty we do not add it, but only return the
		// given context as it is. That way the existence check when reading the
		// context works as expected when no cluster ID or an empty cluster ID was
		// tried to be added.
		return ctx
	}

	return context.WithValue(ctx, clusterIDKey, v)
}

func FromContext(ctx context.Context) string {
	v, _ := ctx.Value(clusterIDKey).(string)
	return v
}
