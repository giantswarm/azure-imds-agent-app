package vmssname

import (
	"context"
)

type key string

var vmssNameKey key = "vmssName"

func NewContext(ctx context.Context, v string) context.Context {
	if v == "" {
		// In case the VMSS name is empty we do not add it, but only return the
		// given context as it is. That way the existence check when reading the
		// context works as expected when no VMSS name or an empty VMSS name was
		// tried to be added.
		return ctx
	}

	return context.WithValue(ctx, vmssNameKey, v)
}

func FromContext(ctx context.Context) string {
	v, _ := ctx.Value(vmssNameKey).(string)
	return v
}
