package nodepoolid

import (
	"context"
)

type key string

var nodePoolIDKey key = "nodePoolID"

func NewContext(ctx context.Context, v string) context.Context {
	if v == "" {
		// In case the node pool ID is empty we do not add it, but only return the
		// given context as it is. That way the existence check when reading the
		// context works as expected when no node pool ID or an empty node pool ID
		// was tried to be added.
		return ctx
	}

	return context.WithValue(ctx, nodePoolIDKey, v)
}

func FromContext(ctx context.Context) string {
	v, _ := ctx.Value(nodePoolIDKey).(string)
	return v
}
