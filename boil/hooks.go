package boil

// HookPoint is the point in time at which we hook
type HookPoint int

// the hook point constants
const (
	HookAfterInsert HookPoint = iota + 1
	HookAfterSelect
	HookAfterUpdate
	HookAfterDelete
	HookAfterUpsert
	HookBeforeInsert
	HookBeforeSelect
	HookBeforeUpdate
	HookBeforeDelete
	HookBeforeUpsert
)
