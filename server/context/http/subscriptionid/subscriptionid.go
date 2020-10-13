package subscriptionid

import (
	"context"
)

type key string

var subscriptionIDKey key = "subscriptionID"

func NewContext(ctx context.Context, v string) context.Context {
	if v == "" {
		// In case the subscription ID is empty we do not add it, but only
		// return the given context as it is. That way the existence check when
		// reading the context works as expected when no subscription ID or an
		// empty subscription ID was tried to be added.
		return ctx
	}

	return context.WithValue(ctx, subscriptionIDKey, v)
}

func FromContext(ctx context.Context) string {
	v, _ := ctx.Value(subscriptionIDKey).(string)
	return v
}
