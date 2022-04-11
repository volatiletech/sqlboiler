package helpers

import (
	"reflect"
	"time"

	"github.com/volatiletech/sqlboiler/v5/drivers"
)

type Table[T any] interface {
	New() T
	TableInfo() TableInfo

	// To set the deleted_at timestamp for a model with soft deletes
	// does nothing if not supported
	SetAsSoftDeleted(T, time.Time)
}

type TableInfo struct {
	Name              string
	Dialect           drivers.Dialect
	Type              reflect.Type
	Mapping           map[string]uint64
	PrimaryKeyMapping []uint64

	AllColumns            []string
	ColumnsWithDefault    []string
	ColumnsWithoutDefault []string
	PrimaryKeyColumns     []string
	GeneratedColumns      []string

	// For soft deletes
	// will be an empty string if disables or not supported
	DeletionColumnName string
}

type queryArgs interface {
	Values() []any
}
