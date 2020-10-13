package instanceid

import (
	"context"
)

type key string

var instanceIDKey key = "instanceID"

func NewContext(ctx context.Context, v string) context.Context {
	if v == "" {
		// In case the instance ID is empty we do not add it, but only return the
		// given context as it is. That way the existence check when reading the
		// context works as expected when no instance ID or an empty instance ID was
		// tried to be added.
		return ctx
	}

	return context.WithValue(ctx, instanceIDKey, v)
}

func FromContext(ctx context.Context) string {
	v, _ := ctx.Value(instanceIDKey).(string)
	return v
}
