package boil

type contextType int

const (
	ctxSkipHooks contextType = iota
	ctxSkipTimestamps
	ctxDebug
	ctxDebugWriter
)
