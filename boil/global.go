package boil

import (
	"time"
)

var (
	// currentDB is a global database handle for the package
	currentDB        Executor
	currentContextDB ContextExecutor
	// timestampLocation is the timezone used for the
	// automated setting of created_at/updated_at columns
	timestampLocation = time.UTC
)

// SetDB initializes the database handle for all template db interactions
func SetDB(db Executor) {
	currentDB = db
	if c, ok := currentDB.(ContextExecutor); ok {
		currentContextDB = c
	}
}

// GetDB retrieves the global state database handle
func GetDB() Executor {
	return currentDB
}

// GetContextDB retrieves the global state database handle as a context executor
func GetContextDB() ContextExecutor {
	return currentContextDB
}

// SetLocation sets the global timestamp Location.
// This is the timezone used by the generated package for the
// automated setting of created_at and updated_at columns.
// If the package was generated with the --no-auto-timestamps flag
// then this function has no effect.
func SetLocation(loc *time.Location) {
	timestampLocation = loc
}

// GetLocation retrieves the global timestamp Location.
// This is the timezone used by the generated package for the
// automated setting of created_at and updated_at columns
// if the package was not generated with the --no-auto-timestamps flag.
func GetLocation() *time.Location {
	return timestampLocation
}
