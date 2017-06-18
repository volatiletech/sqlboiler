// Package bdb supplies the sql(b)oiler (d)ata(b)ase abstractions.
package bdb

import (
	"sort"

	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/strmangle"
)

// Interface abstracts either a side-effect imported driver or a binary
// that is called in order to produce the data required for generation.
type Interface interface {
	Assemble(config map[string]interface{}) (*Assembly, error)
}

// Assembly is the database's properties and the table data from within.
type Assembly struct {
	Tables []Tables
	Props  Properties
}

// Properties describes the databases requirements in terms of which features
// it supports and what kind of quoting mechanisms it uses.
type Properties struct {
	LQ rune
	RQ rune

	UseLastInsertID      bool
	UseIndexPlaceholders bool
	UseTopClause         bool
}

// Constructor breaks down the functionality required to implement a driver
// such that the bdb.Tables method can be used to reduce duplication in driver
// implementations.
type Constructor interface {
	TableNames(schema string, whitelist, blacklist []string) ([]string, error)
	Columns(schema, tableName string) ([]Column, error)
	PrimaryKeyInfo(schema, tableName string) (*PrimaryKey, error)
	ForeignKeyInfo(schema, tableName string) ([]ForeignKey, error)

	// TranslateColumnType takes a Database column type and returns a go column type.
	TranslateColumnType(Column) Column
}

// Tables returns the metadata for all tables, minus the tables
// specified in the blacklist.
func Tables(c Constructor, schema string, whitelist, blacklist []string) ([]Table, error) {
	var err error

	names, err := c.TableNames(schema, whitelist, blacklist)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get table names")
	}

	sort.Strings(names)

	var tables []Table
	for _, name := range names {
		t := Table{
			Name: name,
		}

		if t.Columns, err = c.Columns(schema, name); err != nil {
			return nil, errors.Wrapf(err, "unable to fetch table column info (%s)", name)
		}

		for i, c := range t.Columns {
			t.Columns[i] = c.TranslateColumnType(c)
		}

		if t.PKey, err = c.PrimaryKeyInfo(schema, name); err != nil {
			return nil, errors.Wrapf(err, "unable to fetch table pkey info (%s)", name)
		}

		if t.FKeys, err = c.ForeignKeyInfo(schema, name); err != nil {
			return nil, errors.Wrapf(err, "unable to fetch table fkey info (%s)", name)
		}

		filterForeignKeys(&t, whitelist, blacklist)

		setIsJoinTable(&t)

		tables = append(tables, t)
	}

	// Relationships have a dependency on foreign key nullability.
	for i := range tables {
		tbl := &tables[i]
		setForeignKeyConstraints(tbl, tables)
	}
	for i := range tables {
		tbl := &tables[i]
		setRelationships(tbl, tables)
	}

	return tables, nil
}

// filterForeignKeys filter FK whose ForeignTable is not in whitelist or in blacklist
func filterForeignKeys(t *Table, whitelist, blacklist []string) {
	var fkeys []ForeignKey
	for _, fkey := range t.FKeys {
		if (len(whitelist) == 0 || strmangle.SetInclude(fkey.ForeignTable, whitelist)) &&
		(len(blacklist) == 0 || !strmangle.SetInclude(fkey.ForeignTable, blacklist)) {
			fkeys = append(fkeys, fkey)
		}
	}
	t.FKeys = fkeys
}

// setIsJoinTable if there are:
// A composite primary key involving two columns
// Both primary key columns are also foreign keys
func setIsJoinTable(t *Table) {
	if t.PKey == nil || len(t.PKey.Columns) != 2 || len(t.FKeys) < 2 || len(t.Columns) > 2 {
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

func setForeignKeyConstraints(t *Table, tables []Table) {
	for i, fkey := range t.FKeys {
		localColumn := t.GetColumn(fkey.Column)
		foreignTable := GetTable(tables, fkey.ForeignTable)
		foreignColumn := foreignTable.GetColumn(fkey.ForeignColumn)

		t.FKeys[i].Nullable = localColumn.Nullable
		t.FKeys[i].Unique = localColumn.Unique
		t.FKeys[i].ForeignColumnNullable = foreignColumn.Nullable
		t.FKeys[i].ForeignColumnUnique = foreignColumn.Unique
	}
}

func setRelationships(t *Table, tables []Table) {
	t.ToOneRelationships = toOneRelationships(*t, tables)
	t.ToManyRelationships = toManyRelationships(*t, tables)
}
