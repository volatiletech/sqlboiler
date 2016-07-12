// Package bdb supplies the sql(b)oiler (d)ata(b)ase abstractions.
package bdb

import "github.com/pkg/errors"

// Interface for a database driver. Functionality required to support a specific
// database type (eg, MySQL, Postgres etc.)
type Interface interface {
	TableNames() ([]string, error)
	Columns(tableName string) ([]Column, error)
	PrimaryKeyInfo(tableName string) (*PrimaryKey, error)
	ForeignKeyInfo(tableName string) ([]ForeignKey, error)

	// TranslateColumnType takes a Database column type and returns a go column type.
	TranslateColumnType(Column) Column

	// Open the database connection
	Open() error

	// Close the database connection
	Close()
}

// Tables returns the table metadata for the given tables, or all tables if
// no tables are provided.
func Tables(db Interface, names ...string) ([]Table, error) {
	var err error
	if len(names) == 0 {
		if names, err = db.TableNames(); err != nil {
			return nil, errors.Wrap(err, "unable to get table names")
		}
	}

	var tables []Table
	for _, name := range names {
		t := Table{Name: name}

		if t.Columns, err = db.Columns(name); err != nil {
			return nil, errors.Wrapf(err, "unable to fetch table column info (%s)", name)
		}

		for i, c := range t.Columns {
			t.Columns[i] = db.TranslateColumnType(c)
		}

		if t.PKey, err = db.PrimaryKeyInfo(name); err != nil {
			return nil, errors.Wrapf(err, "unable to fetch table pkey info (%s)", name)
		}

		if t.FKeys, err = db.ForeignKeyInfo(name); err != nil {
			return nil, errors.Wrapf(err, "unable to fetch table fkey info (%s)", name)
		}

		setIsJoinTable(&t)
		setForeignKeyNullability(&t)

		tables = append(tables, t)
	}

	for i := range tables {
		tbl := &tables[i]
		setRelationships(tbl, tables)
	}

	return tables, nil
}

// setIsJoinTable iff there are:
// A composite primary key involving two columns
// Both primary key columns are also foreign keys
func setIsJoinTable(t *Table) {
	if t.PKey == nil || len(t.PKey.Columns) != 2 || len(t.FKeys) < 2 {
		return
	}

	for _, c := range t.PKey.Columns {
		found := false
		for _, f := range t.FKeys {
			if c == f.Column {
				found = true
				break
			}
		}
		if !found {
			return
		}
	}

	t.IsJoinTable = true
}

func setForeignKeyNullability(t *Table) {
	for i, fkey := range t.FKeys {

		found := -1
		for j, col := range t.Columns {
			if col.Name == fkey.Column {
				found = j
				break
			}
		}

		if found < 0 {
			panic("could not find foreign key column in table")
		}

		t.FKeys[i].Nullable = t.Columns[found].Nullable
	}
}

func setRelationships(t *Table, tables []Table) {
	t.ToManyRelationships = toManyRelationships(*t, tables)
}
