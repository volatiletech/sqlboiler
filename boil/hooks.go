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

// QueryHookPoint is the point in time at which we hook query
type QueryHookPoint int

// the query hook point constants
const (
	InsertHook QueryHookPoint = iota + 1
	SelectHook
	UpdateHook
	DeleteHook
	UpsertHook
)
