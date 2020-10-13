package resourcegroupname

import (
	"context"
)

type key string

var resourceGroupNameKey key = "resourceGroupName"

func NewContext(ctx context.Context, v string) context.Context {
	if v == "" {
		// In case the resource group name is empty we do not add it, but only
		// return the given context as it is. That way the existence check when
		// reading the context works as expected when no resource group name or
		// an empty resource group name was tried to be added.
		return ctx
	}

	return context.WithValue(ctx, resourceGroupNameKey, v)
}

func FromContext(ctx context.Context) string {
	v, _ := ctx.Value(resourceGroupNameKey).(string)
	return v
}
