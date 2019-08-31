package boil

import "context"

// SkipHooks modifies a context to prevent hooks from running for any query
// it encounters.
func SkipHooks(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxSkipHooks, true)
}

// HooksAreSkipped returns true if the context skips hooks
func HooksAreSkipped(ctx context.Context) bool {
	skip := ctx.Value(ctxSkipHooks)
	return skip != nil && skip.(bool)
}

// SkipTimestamps modifies a context to prevent hooks from running for any query
// it encounters.
func SkipTimestamps(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxSkipTimestamps, true)
}

// TimestampsAreSkipped returns true if the context skips hooks
func TimestampsAreSkipped(ctx context.Context) bool {
	skip := ctx.Value(ctxSkipTimestamps)
	return skip != nil && skip.(bool)
}

// HookPoint is the point in time at which we hook
type HookPoint int

// the hook point constants
const (
	BeforeInsertHook HookPoint = iota + 1
	BeforeUpdateHook
	BeforeDeleteHook
	BeforeUpsertHook
	AfterInsertHook
	AfterSelectHook
	AfterUpdateHook
	AfterDeleteHook
	AfterUpsertHook
)
