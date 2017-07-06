package queries

import (
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/lbryio/null.go"
	"github.com/lbryio/sqlboiler/boil"

	"github.com/go-errors/errors"
)

// joinKind is the type of join
type joinKind int

// Join type constants
const (
	JoinInner joinKind = iota
	JoinOuterLeft
	JoinOuterRight
	JoinNatural
)

// Query holds the state for the built up query
type Query struct {
	executor   boil.Executor
	dialect    *Dialect
	rawSQL     rawSQL
	load       []string
	delete     bool
	update     map[string]interface{}
	selectCols []string
	count      bool
	from       []string
	joins      []join
	where      []where
	in         []in
	groupBy    []string
	orderBy    []string
	having     []having
	limit      int
	offset     int
	forlock    string
}

// Dialect holds values that direct the query builder
// how to build compatible queries for each database.
// Each database driver needs to implement functions
// that provide these values.
type Dialect struct {
	// The left quote character for SQL identifiers
	LQ byte
	// The right quote character for SQL identifiers
	RQ byte
	// Bool flag indicating whether indexed
	// placeholders ($1) are used, or ? placeholders.
	IndexPlaceholders bool
	// Bool flag indicating whether "TOP" or "LIMIT" clause
	// must be used for rows limitation
	UseTopClause bool
}

type where struct {
	clause      string
	orSeparator bool
	args        []interface{}
}

type in struct {
	clause      string
	orSeparator bool
	args        []interface{}
}

type having struct {
	clause string
	args   []interface{}
}

type rawSQL struct {
	sql  string
	args []interface{}
}

type join struct {
	kind   joinKind
	clause string
	args   []interface{}
}

// Raw makes a raw query, usually for use with bind
func Raw(exec boil.Executor, query string, args ...interface{}) *Query {
	return &Query{
		executor: exec,
		rawSQL: rawSQL{
			sql:  query,
			args: args,
		},
	}
}

// RawG makes a raw query using the global boil.Executor, usually for use with bind
func RawG(query string, args ...interface{}) *Query {
	return Raw(boil.GetDB(), query, args...)
}

// Exec executes a query that does not need a row returned
func (q *Query) Exec() (sql.Result, error) {
	qs, args := buildQuery(q)
	if boil.DebugMode {
		qStr, err := interpolateParams(qs, args...)
		if err != nil {
			return nil, err
		}
		fmt.Fprintln(boil.DebugWriter, qStr)
	}
	return q.executor.Exec(qs, args...)
}

// QueryRow executes the query for the One finisher and returns a row
func (q *Query) QueryRow() *sql.Row {
	qs, args := buildQuery(q)
	if boil.DebugMode {
		qStr, err := interpolateParams(qs, args...)
		if err != nil {
			panic(err)
		}
		fmt.Fprintln(boil.DebugWriter, qStr)
	}
	return q.executor.QueryRow(qs, args...)
}

// Query executes the query for the All finisher and returns multiple rows
func (q *Query) Query() (*sql.Rows, error) {
	qs, args := buildQuery(q)
	if boil.DebugMode {
		qStr, err := interpolateParams(qs, args...)
		if err != nil {
			return nil, err
		}
		fmt.Fprintln(boil.DebugWriter, qStr)
	}
	return q.executor.Query(qs, args...)
}

// ExecP executes a query that does not need a row returned
// It will panic on error
func (q *Query) ExecP() sql.Result {
	res, err := q.Exec()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return res
}

// QueryP executes the query for the All finisher and returns multiple rows
// It will panic on error
func (q *Query) QueryP() *sql.Rows {
	rows, err := q.Query()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return rows
}

// SetExecutor on the query.
func SetExecutor(q *Query, exec boil.Executor) {
	q.executor = exec
}

// GetExecutor on the query.
func GetExecutor(q *Query) boil.Executor {
	return q.executor
}

// SetDialect on the query.
func SetDialect(q *Query, dialect *Dialect) {
	q.dialect = dialect
}

// SetSQL on the query.
func SetSQL(q *Query, sql string, args ...interface{}) {
	q.rawSQL = rawSQL{sql: sql, args: args}
}

// SetLoad on the query.
func SetLoad(q *Query, relationships ...string) {
	q.load = append([]string(nil), relationships...)
}

// AppendLoad on the query.
func AppendLoad(q *Query, relationships ...string) {
	q.load = append(q.load, relationships...)
}

// SetSelect on the query.
func SetSelect(q *Query, sel []string) {
	q.selectCols = sel
}

// SetCount on the query.
func SetCount(q *Query) {
	q.count = true
}

// SetDelete on the query.
func SetDelete(q *Query) {
	q.delete = true
}

// SetLimit on the query.
func SetLimit(q *Query, limit int) {
	q.limit = limit
}

// SetOffset on the query.
func SetOffset(q *Query, offset int) {
	q.offset = offset
}

// SetFor on the query.
func SetFor(q *Query, clause string) {
	q.forlock = clause
}

// SetUpdate on the query.
func SetUpdate(q *Query, cols map[string]interface{}) {
	q.update = cols
}

// AppendSelect on the query.
func AppendSelect(q *Query, columns ...string) {
	q.selectCols = append(q.selectCols, columns...)
}

// AppendFrom on the query.
func AppendFrom(q *Query, from ...string) {
	q.from = append(q.from, from...)
}

// SetFrom replaces the current from statements.
func SetFrom(q *Query, from ...string) {
	q.from = append([]string(nil), from...)
}

// AppendInnerJoin on the query.
func AppendInnerJoin(q *Query, clause string, args ...interface{}) {
	q.joins = append(q.joins, join{clause: clause, kind: JoinInner, args: args})
}

// AppendHaving on the query.
func AppendHaving(q *Query, clause string, args ...interface{}) {
	q.having = append(q.having, having{clause: clause, args: args})
}

// AppendWhere on the query.
func AppendWhere(q *Query, clause string, args ...interface{}) {
	q.where = append(q.where, where{clause: clause, args: args})
}

// AppendIn on the query.
func AppendIn(q *Query, clause string, args ...interface{}) {
	q.in = append(q.in, in{clause: clause, args: args})
}

// SetLastWhereAsOr sets the or separator for the tail "WHERE" in the slice
func SetLastWhereAsOr(q *Query) {
	if len(q.where) == 0 {
		return
	}

	q.where[len(q.where)-1].orSeparator = true
}

// SetLastInAsOr sets the or separator for the tail "IN" in the slice
func SetLastInAsOr(q *Query) {
	if len(q.in) == 0 {
		return
	}

	q.in[len(q.in)-1].orSeparator = true
}

// AppendGroupBy on the query.
func AppendGroupBy(q *Query, clause string) {
	q.groupBy = append(q.groupBy, clause)
}

// AppendOrderBy on the query.
func AppendOrderBy(q *Query, clause string) {
	q.orderBy = append(q.orderBy, clause)
}

// duplicated in boil_queries.tpl
func interpolateParams(query string, args ...interface{}) (string, error) {
	for i := 0; i < len(args); i++ {
		field := reflect.ValueOf(args[i])

		if value, ok := field.Interface().(time.Time); ok {
			query = strings.Replace(query, "?", `"`+value.Format("2006-01-02 15:04:05")+`"`, 1)
		} else if nullable, ok := field.Interface().(null.Nullable); ok {
			if nullable.IsNull() {
				query = strings.Replace(query, "?", "NULL", 1)
			} else {
				switch field.Type() {
				case reflect.TypeOf(null.Time{}):
					query = strings.Replace(query, "?", `"`+field.Interface().(null.Time).Time.Format("2006-01-02 15:04:05")+`"`, 1)
				case reflect.TypeOf(null.Int{}):
					query = strings.Replace(query, "?", strconv.FormatInt(int64(field.Interface().(null.Int).Int), 10), 1)
				case reflect.TypeOf(null.Int8{}):
					query = strings.Replace(query, "?", strconv.FormatInt(int64(field.Interface().(null.Int8).Int8), 10), 1)
				case reflect.TypeOf(null.Int16{}):
					query = strings.Replace(query, "?", strconv.FormatInt(int64(field.Interface().(null.Int16).Int16), 10), 1)
				case reflect.TypeOf(null.Int32{}):
					query = strings.Replace(query, "?", strconv.FormatInt(int64(field.Interface().(null.Int32).Int32), 10), 1)
				case reflect.TypeOf(null.Int64{}):
					query = strings.Replace(query, "?", strconv.FormatInt(field.Interface().(null.Int64).Int64, 10), 1)
				case reflect.TypeOf(null.Uint{}):
					query = strings.Replace(query, "?", strconv.FormatUint(uint64(field.Interface().(null.Uint).Uint), 10), 1)
				case reflect.TypeOf(null.Uint8{}):
					query = strings.Replace(query, "?", strconv.FormatUint(uint64(field.Interface().(null.Uint8).Uint8), 10), 1)
				case reflect.TypeOf(null.Uint16{}):
					query = strings.Replace(query, "?", strconv.FormatUint(uint64(field.Interface().(null.Uint16).Uint16), 10), 1)
				case reflect.TypeOf(null.Uint32{}):
					query = strings.Replace(query, "?", strconv.FormatUint(uint64(field.Interface().(null.Uint32).Uint32), 10), 1)
				case reflect.TypeOf(null.Uint64{}):
					query = strings.Replace(query, "?", strconv.FormatUint(field.Interface().(null.Uint64).Uint64, 10), 1)
				case reflect.TypeOf(null.String{}):
					query = strings.Replace(query, "?", `"`+field.Interface().(null.String).String+`"`, 1)
				case reflect.TypeOf(null.Bool{}):
					if field.Interface().(null.Bool).Bool {
						query = strings.Replace(query, "?", "1", 1)
					} else {
						query = strings.Replace(query, "?", "0", 1)
					}
				}
			}
		} else {
			switch field.Kind() {
			case reflect.Bool:
				boolString := "0"
				if field.Bool() {
					boolString = "1"
				}
				query = strings.Replace(query, "?", boolString, 1)
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				query = strings.Replace(query, "?", strconv.FormatInt(field.Int(), 10), 1)
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				query = strings.Replace(query, "?", strconv.FormatUint(field.Uint(), 10), 1)
			case reflect.String:
				query = strings.Replace(query, "?", `"`+field.String()+`"`, 1)
			default:
				return "", errors.New("Dont know how to interpolate type " + field.Type().String())
			}
		}
	}
	return query, nil
}
