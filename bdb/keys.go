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
	Table    string
	Name     string
	Column   string
	Nullable bool
	Unique   bool

	ForeignTable          string
	ForeignColumn         string
	ForeignColumnNullable bool
	ForeignColumnUnique   bool
}

// SQLColumnDef formats a column name and type like an SQL column definition.
type SQLColumnDef struct {
	Name string
	Type string
}

// String for fmt.Stringer
func (s SQLColumnDef) String() string {
	return fmt.Sprintf("%s %s", s.Name, s.Type)
}

// SQLColumnDefs has small helper functions
type SQLColumnDefs []SQLColumnDef

// Names returns all the names
func (s SQLColumnDefs) Names() []string {
	names := make([]string, len(s))

	for i, sqlDef := range s {
		names[i] = sqlDef.Name
	}

	return names
}

// Types returns all the types
func (s SQLColumnDefs) Types() []string {
	types := make([]string, len(s))

	for i, sqlDef := range s {
		types[i] = sqlDef.Type
	}

	return types
}

// SQLColDefinitions creates a definition in sql format for a column
func SQLColDefinitions(cols []Column, names []string) SQLColumnDefs {
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
// empty string. Primary key columns with default values are presumed
// to be auto-increment, because pkeys need to be unique and a static
// default value would cause collisions.
func AutoIncPrimaryKey(cols []Column, pkey *PrimaryKey) *Column {
	if pkey == nil {
		return nil
	}

	for _, pkeyColumn := range pkey.Columns {
		for _, c := range cols {
			if c.Name != pkeyColumn {
				continue
			}

			if c.Default != "auto_increment" || c.Nullable ||
				!(strings.HasPrefix(c.Type, "int") || strings.HasPrefix(c.Type, "uint")) {
				continue
			}

			return &c
		}
	}

	return nil
}
