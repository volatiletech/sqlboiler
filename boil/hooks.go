package boil

// HookPoint is the point in time at which we hook
type HookPoint int

// the hook point constants
const (
	HookAfterCreate HookPoint = iota + 1
	HookAfterUpdate
	HookAfterUpsert
	HookBeforeCreate
	HookBeforeUpdate
	HookBeforeUpsert
)
