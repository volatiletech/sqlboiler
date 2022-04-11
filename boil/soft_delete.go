package boil

import "context"

// HardDelete modifies a context to ignore soft deletion
func HardDelete(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxHardDelete, true)
}

// DoHardDelete returns true if the context hard deletes
func DoHardDelete(ctx context.Context) bool {
	hardDelete := ctx.Value(ctxHardDelete)
	return hardDelete != nil && hardDelete.(bool)
}
