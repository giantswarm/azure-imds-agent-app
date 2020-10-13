package organizationid

import (
	"context"
)

type key string

var organizationIDKey key = "organizationID"

func NewContext(ctx context.Context, v string) context.Context {
	if v == "" {
		// In case the organization ID is empty we do not add it, but only return
		// the given context as it is. That way the existence check when reading the
		// context works as expected when no organization ID or an empty
		// organization ID was tried to be added.
		return ctx
	}

	return context.WithValue(ctx, organizationIDKey, v)
}

func FromContext(ctx context.Context) string {
	v, _ := ctx.Value(organizationIDKey).(string)
	return v
}
