package qmhelper

import (
	"fmt"
	"reflect"

	"github.com/volatiletech/sqlboiler/v4/queries"
)

// Nullable object
type Nullable interface {
	IsZero() bool
}

// WhereQueryMod allows construction of where clauses
type WhereQueryMod struct {
	Clause string
	Args   []interface{}
}

// Apply implements QueryMod.Apply.
func (qm WhereQueryMod) Apply(q *queries.Query) {
	queries.AppendWhere(q, qm.Clause, qm.Args...)
}

// WhereNullEQ is a helper for doing equality with null types
func WhereNullEQ(name string, negated bool, value interface{}) WhereQueryMod {
	isNull := false
	if nullable, ok := value.(Nullable); ok {
		isNull = nullable.IsZero()
	} else {
		isNull = reflect.ValueOf(value).IsNil()
	}

	if isNull {
		var not string
		if negated {
			not = "not "
		}
		return WhereQueryMod{
			Clause: fmt.Sprintf("%s is %snull", name, not),
		}
	}

	op := "="
	if negated {
		op = "!="
	}

	return WhereQueryMod{
		Clause: fmt.Sprintf("%s %s ?", name, op),
		Args:   []interface{}{value},
	}
}

// WhereIsNull is a helper that just returns "name is null"
func WhereIsNull(name string) WhereQueryMod {
	return WhereQueryMod{Clause: fmt.Sprintf("%s is null", name)}
}

// WhereIsNotNull is a helper that just returns "name is not null"
func WhereIsNotNull(name string) WhereQueryMod {
	return WhereQueryMod{Clause: fmt.Sprintf("%s is not null", name)}
}

type operator string

// Supported operations
const (
	EQ  operator = "="
	NEQ operator = "!="
	LT  operator = "<"
	LTE operator = "<="
	GT  operator = ">"
	GTE operator = ">="
)

// Where is a helper for doing operations on primitive types
func Where(name string, operator operator, value interface{}) WhereQueryMod {
	return WhereQueryMod{
		Clause: fmt.Sprintf("%s %s ?", name, string(operator)),
		Args:   []interface{}{value},
	}
}
