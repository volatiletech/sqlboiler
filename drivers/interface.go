// Package drivers talks to various database backends and retrieves table,
// column, type, and foreign key information
package drivers

import (
	"sort"

	"github.com/friendsofgo/errors"
	"github.com/volatiletech/sqlboiler/v4/importers"
	"github.com/volatiletech/strmangle"
)

// These constants are used in the config map passed into the driver
const (
	ConfigBlacklist = "blacklist"
	ConfigWhitelist = "whitelist"
	ConfigSchema    = "schema"

	ConfigUser    = "user"
	ConfigPass    = "pass"
	ConfigHost    = "host"
	ConfigPort    = "port"
	ConfigDBName  = "dbname"
	ConfigSSLMode = "sslmode"
)

// Interface abstracts either a side-effect imported driver or a binary
// that is called in order to produce the data required for generation.
type Interface interface {
	// Assemble the database information into a nice struct
	Assemble(config Config) (*DBInfo, error)
	// Templates to add/replace for generation
	Templates() (map[string]string, error)
	// Imports to merge for generation
	Imports() (importers.Collection, error)
}

// DBInfo is the database's table data and dialect.
type DBInfo struct {
	Schema  string  `json:"schema"`
	Tables  []Table `json:"tables"`
	Dialect Dialect `json:"dialect"`
}

// Dialect describes the databases requirements in terms of which features
// it speaks and what kind of quoting mechanisms it uses.
//
// WARNING: When updating this struct there is a copy of it inside
// the boil_queries template that is used for users to create queries
// without having to figure out what their dialect is.
type Dialect struct {
	LQ rune `json:"lq"`
	RQ rune `json:"rq"`

	UseIndexPlaceholders bool `json:"use_index_placeholders"`
	UseLastInsertID      bool `json:"use_last_insert_id"`
	UseSchema            bool `json:"use_schema"`
	UseDefaultKeyword    bool `json:"use_default_keyword"`

	// The following is mostly for T-SQL/MSSQL, what a show
	UseAutoColumns          bool `json:"use_auto_columns"`
	UseTopClause            bool `json:"use_top_clause"`
	UseOutputClause         bool `json:"use_output_clause"`
	UseCaseWhenExistsClause bool `json:"use_case_when_exists_clause"`
}

// Constructor breaks down the functionality required to implement a driver
// such that the drivers.Tables method can be used to reduce duplication in driver
// implementations.
type Constructor interface {
	TableNames(schema string, whitelist, blacklist []string) ([]string, error)
	Columns(schema, tableName string, whitelist, blacklist []string) ([]Column, error)
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

		if t.Columns, err = c.Columns(schema, name, whitelist, blacklist); err != nil {
			return nil, errors.Wrapf(err, "unable to fetch table column info (%s)", name)
		}

		for i, col := range t.Columns {
			t.Columns[i] = c.TranslateColumnType(col)
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
