package boil

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
