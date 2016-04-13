package boil

// HookPoint is the point in time at which we hook
type HookPoint int

const (
	HookAfterCreate HookPoint = iota + 1
	HookAfterUpdate
	HookBeforeCreate
	HookBeforeUpdate
)
