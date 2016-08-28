package boil

// HookPoint is the point in time at which we hook
type HookPoint int

// the hook point constants
const (
	HookBeforeInsert HookPoint = iota + 1
	HookBeforeUpdate
	HookBeforeDelete
	HookBeforeUpsert
	HookAfterInsert
	HookAfterSelect
	HookAfterUpdate
	HookAfterDelete
	HookAfterUpsert
)
