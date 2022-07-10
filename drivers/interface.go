// Package drivers talks to various database backends and retrieves table,
// column, type, and foreign key information
package drivers

import (
	"sort"
	"sync"

	"github.com/friendsofgo/errors"
	"github.com/volatiletech/sqlboiler/v4/importers"
	"github.com/volatiletech/strmangle"
)

// These constants are used in the config map passed into the driver
const (
	ConfigBlacklist      = "blacklist"
	ConfigWhitelist      = "whitelist"
	ConfigSchema         = "schema"
	ConfigAddEnumTypes   = "add-enum-types"
	ConfigEnumNullPrefix = "enum-null-prefix"
	ConfigConcurrency    = "concurrency"

	ConfigUser    = "user"
	ConfigPass    = "pass"
	ConfigHost    = "host"
	ConfigPort    = "port"
	ConfigDBName  = "dbname"
	ConfigSSLMode = "sslmode"

	// DefaultConcurrency defines the default amount of threads to use when loading tables info
	DefaultConcurrency = 10
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
	UseTopClause            bool `json:"use_top_clause"`
	UseOutputClause         bool `json:"use_output_clause"`
	UseCaseWhenExistsClause bool `json:"use_case_when_exists_clause"`

	// No longer used, left for backwards compatibility
	// should be removed in v5
	UseAutoColumns bool `json:"use_auto_columns"`
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

// Constructor breaks down the functionality required to implement a driver
// such that the drivers.Views method can be used to reduce duplication in driver
// implementations.
type ViewConstructor interface {
	ViewNames(schema string, whitelist, blacklist []string) ([]string, error)
	ViewCapabilities(schema, viewName string) (ViewCapabilities, error)
	ViewColumns(schema, tableName string, whitelist, blacklist []string) ([]Column, error)

	// TranslateColumnType takes a Database column type and returns a go column type.
	TranslateColumnType(Column) Column
}

type TableColumnTypeTranslator interface {
	// TranslateTableColumnType takes a Database column type and table name and returns a go column type.
	TranslateTableColumnType(c Column, tableName string) Column
}

// Tables returns the metadata for all tables, minus the tables
// specified in the blacklist.
func Tables(c Constructor, schema string, whitelist, blacklist []string) ([]Table, error) {
	return TablesConcurrently(c, schema, whitelist, blacklist, 1)
}

// TablesConcurrently is a concurrent version of Tables. It returns the
// metadata for all tables, minus the tables specified in the blacklist.
func TablesConcurrently(c Constructor, schema string, whitelist, blacklist []string, concurrency int) ([]Table, error) {
	var err error
	var ret []Table

	ret, err = tables(c, schema, whitelist, blacklist, concurrency)
	if err != nil {
		return nil, errors.Wrap(err, "unable to load tables")
	}

	if vc, ok := c.(ViewConstructor); ok {
		v, err := views(vc, schema, whitelist, blacklist, concurrency)
		if err != nil {
			return nil, errors.Wrap(err, "unable to load views")
		}
		ret = append(ret, v...)
	}

	return ret, nil
}

func tables(c Constructor, schema string, whitelist, blacklist []string, concurrency int) ([]Table, error) {
	var err error

	names, err := c.TableNames(schema, whitelist, blacklist)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get table names")
	}

	sort.Strings(names)

	ret := make([]Table, len(names))

	limiter := newConcurrencyLimiter(concurrency)
	wg := sync.WaitGroup{}
	errs := make(chan error, len(names))
	for i, name := range names {
		wg.Add(1)
		limiter.get()
		go func(i int, name string) {
			defer wg.Done()
			defer limiter.put()
			t, err := table(c, schema, name, whitelist, blacklist)
			if err != nil {
				errs <- err
				return
			}
			ret[i] = t
		}(i, name)
	}

	wg.Wait()

	// return first error occurred if any
	if len(errs) > 0 {
		return nil, <-errs
	}

	// Relationships have a dependency on foreign key nullability.
	for i := range ret {
		tbl := &ret[i]
		setForeignKeyConstraints(tbl, ret)
	}
	for i := range ret {
		tbl := &ret[i]
		setRelationships(tbl, ret)
	}

	return ret, nil
}

// table returns columns info for a given table
func table(c Constructor, schema string, name string, whitelist, blacklist []string) (Table, error) {
	var err error
	t := &Table{
		Name: name,
	}

	if t.Columns, err = c.Columns(schema, name, whitelist, blacklist); err != nil {
		return Table{}, errors.Wrapf(err, "unable to fetch table column info (%s)", name)
	}

	tr, ok := c.(TableColumnTypeTranslator)
	if ok {
		for i, col := range t.Columns {
			t.Columns[i] = tr.TranslateTableColumnType(col, name)
		}
	} else {
		for i, col := range t.Columns {
			t.Columns[i] = c.TranslateColumnType(col)
		}
	}

	if t.PKey, err = c.PrimaryKeyInfo(schema, name); err != nil {
		return Table{}, errors.Wrapf(err, "unable to fetch table pkey info (%s)", name)
	}

	if t.FKeys, err = c.ForeignKeyInfo(schema, name); err != nil {
		return Table{}, errors.Wrapf(err, "unable to fetch table fkey info (%s)", name)
	}

	filterPrimaryKey(t, whitelist, blacklist)
	filterForeignKeys(t, whitelist, blacklist)

	setIsJoinTable(t)

	return *t, nil
}

// views returns the metadata for all views, minus the views
// specified in the blacklist.
func views(c ViewConstructor, schema string, whitelist, blacklist []string, concurrency int) ([]Table, error) {
	var err error

	names, err := c.ViewNames(schema, whitelist, blacklist)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get view names")
	}

	sort.Strings(names)

	ret := make([]Table, len(names))

	limiter := newConcurrencyLimiter(concurrency)
	wg := sync.WaitGroup{}
	errs := make(chan error, len(names))
	for i, name := range names {
		wg.Add(1)
		limiter.get()
		go func(i int, name string) {
			defer wg.Done()
			defer limiter.put()
			t, err := view(c, schema, name, whitelist, blacklist)
			if err != nil {
				errs <- err
				return
			}
			ret[i] = t
		}(i, name)
	}

	wg.Wait()

	// return first error occurred if any
	if len(errs) > 0 {
		return nil, <-errs
	}

	return ret, nil
}

// view returns columns info for a given view
func view(c ViewConstructor, schema string, name string, whitelist, blacklist []string) (Table, error) {
	var err error
	t := Table{
		IsView: true,
		Name:   name,
	}

	if t.ViewCapabilities, err = c.ViewCapabilities(schema, name); err != nil {
		return Table{}, errors.Wrapf(err, "unable to fetch view capabilities info (%s)", name)
	}

	if t.Columns, err = c.ViewColumns(schema, name, whitelist, blacklist); err != nil {
		return Table{}, errors.Wrapf(err, "unable to fetch view column info (%s)", name)
	}

	tr, ok := c.(TableColumnTypeTranslator)
	if ok {
		for i, col := range t.Columns {
			t.Columns[i] = tr.TranslateTableColumnType(col, name)
		}
	} else {
		for i, col := range t.Columns {
			t.Columns[i] = c.TranslateColumnType(col)
		}
	}

	return t, nil
}

func knownColumn(table string, column string, whitelist, blacklist []string) bool {
	return (len(whitelist) == 0 ||
		strmangle.SetInclude(table, whitelist) ||
		strmangle.SetInclude(table+"."+column, whitelist) ||
		strmangle.SetInclude("*."+column, whitelist)) &&
		(len(blacklist) == 0 || (!strmangle.SetInclude(table, blacklist) &&
			!strmangle.SetInclude(table+"."+column, blacklist) &&
			!strmangle.SetInclude("*."+column, blacklist)))
}

// filterPrimaryKey filter columns from the primary key that are not in whitelist or in blacklist
func filterPrimaryKey(t *Table, whitelist, blacklist []string) {
	if t.PKey == nil {
		return
	}

	pkeyColumns := make([]string, 0, len(t.PKey.Columns))
	for _, c := range t.PKey.Columns {
		if knownColumn(t.Name, c, whitelist, blacklist) {
			pkeyColumns = append(pkeyColumns, c)
		}
	}
	t.PKey.Columns = pkeyColumns
}

// filterForeignKeys filter FK whose ForeignTable is not in whitelist or in blacklist
func filterForeignKeys(t *Table, whitelist, blacklist []string) {
	var fkeys []ForeignKey

	for _, fkey := range t.FKeys {
		if knownColumn(fkey.ForeignTable, fkey.ForeignColumn, whitelist, blacklist) &&
			knownColumn(fkey.Table, fkey.Column, whitelist, blacklist) {
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

// concurrencyCounter is a helper structure that can limit amount of concurrently processed requests
type concurrencyLimiter chan struct{}

func newConcurrencyLimiter(capacity int) concurrencyLimiter {
	ret := make(concurrencyLimiter, capacity)
	for i := 0; i < capacity; i++ {
		ret <- struct{}{}
	}

	return ret
}

func (c concurrencyLimiter) get() {
	<-c
}

func (c concurrencyLimiter) put() {
	c <- struct{}{}
}
