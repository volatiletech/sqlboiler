package bdb

import (
	"fmt"
	"regexp"
	"strings"
)

var rgxAutoIncColumn = regexp.MustCompile(`^nextval\(.*\)`)

// PrimaryKey represents a primary key constraint in a database
type PrimaryKey struct {
	Name    string
	Columns []string
}

// ForeignKey represents a foreign key constraint in a database
type ForeignKey struct {
	Name   string
	Column string

	ForeignTable  string
	ForeignColumn string
}

// SQLColumnDef formats a column name and type like an SQL column definition.
type SQLColumnDef struct {
	Name string
	Type string
}

func (s SQLColumnDef) String() string {
	return fmt.Sprintf("%s %s", s.Name, s.Type)
}

// SQLColDefinitions creates a definition in sql format for a column
// example: id int64, thingName string
func SQLColDefinitions(cols []Column, names []string) []SQLColumnDef {
	ret := make([]SQLColumnDef, len(names))

	for i, n := range names {
		for _, c := range cols {
			if n != c.Name {
				continue
			}

			ret[i] = SQLColumnDef{Name: n, Type: c.Type}
		}
	}

	return ret
}

// AutoIncPrimaryKey returns the auto-increment primary key column name or an
// empty string.
func AutoIncPrimaryKey(cols []Column, pkey *PrimaryKey) (col Column, ok bool) {
	if pkey == nil {
		return col, false
	}

	for _, pkeyColumn := range pkey.Columns {
		for _, c := range cols {
			if c.Name != pkeyColumn {
				continue
			}

			if !rgxAutoIncColumn.MatchString(c.Default) || c.IsNullable ||
				!(strings.HasPrefix(c.Type, "int") || strings.HasPrefix(c.Type, "uint")) {
				continue
			}

			return c, true
		}
	}

	return col, false
}
