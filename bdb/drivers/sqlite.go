package drivers

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/vattle/sqlboiler/bdb"
)

// SQLiteDriver holds the database connection string and a handle
// to the database connection.
type SQLiteDriver struct {
	connStr string
	dbConn  *sql.DB
}

// NewSQLiteDriver takes the database connection details as parameters and
// returns a pointer to a SQLiteDriver object. Note that it is required to
// call SQLiteDriver.Open() and SQLiteDriver.Close() to open and close
// the database connection once an object has been obtained.
func NewSQLiteDriver(file string) *SQLiteDriver {
	driver := SQLiteDriver{
		connStr: SQLiteBuildQueryString(file),
	}

	return &driver
}

// SQLiteBuildQueryString builds a query string for SQLite.
func SQLiteBuildQueryString(file string) string {
	if !SQLITE_ENABLED {
		panic("SQLite is not enabled")
	}
	return "file:" + file + "?_loc=UTC"
}

// Open opens the database connection using the connection string
func (m *SQLiteDriver) Open() error {
	var err error

	m.dbConn, err = sql.Open("sqlite3", m.connStr)
	if err != nil {
		return err
	}

	return nil
}

// Close closes the database connection
func (m *SQLiteDriver) Close() {
	m.dbConn.Close()
}

// UseLastInsertID returns false for sqlite
func (m *SQLiteDriver) UseLastInsertID() bool {
	return true
}

// UseTopClause returns false to indicate SQLite doesnt support SQL TOP clause
func (m *SQLiteDriver) UseTopClause() bool {
	return false
}

// TableNames connects to the sqlite database and
// retrieves all table names from sqlite_master
func (m *SQLiteDriver) TableNames(schema string, whitelist, blacklist []string) ([]string, error) {
	var args []interface{}
	var names []string

	query := `SELECT name FROM sqlite_master WHERE type='table';`

	if len(whitelist) > 0 {
		query += fmt.Sprintf(" and name in (%s);", strings.Repeat(",?", len(whitelist))[1:])
		for _, w := range whitelist {
			args = append(args, w)
		}
	} else if len(blacklist) > 0 {
		query += fmt.Sprintf(" and name not in (%s);", strings.Repeat(",?", len(blacklist))[1:])
		for _, b := range blacklist {
			args = append(args, b)
		}
	}

	rows, err := m.dbConn.Query(query, args...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		if name != "sqlite_sequence" {
			names = append(names, name)
		}
	}

	return names, nil
}

type sqliteIndex struct {
	SeqNum  int
	Unique  int
	Partial int
	Name    string
	Origin  string
	Columns []string
}

type sqliteTableInfo struct {
	Cid          string
	Name         string
	Type         string
	NotNull      bool
	DefaultValue *string
	Pk           int
}

func (m *SQLiteDriver) tableInfo(tableName string) ([]*sqliteTableInfo, error) {
	var ret []*sqliteTableInfo
	rows, err := m.dbConn.Query(fmt.Sprintf("PRAGMA table_info('%s')", tableName))

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		tinfo := &sqliteTableInfo{}
		if err := rows.Scan(&tinfo.Cid, &tinfo.Name, &tinfo.Type, &tinfo.NotNull, &tinfo.DefaultValue, &tinfo.Pk); err != nil {
			return nil, errors.Wrapf(err, "unable to scan for table %s", tableName)
		}
		ret = append(ret, tinfo)
	}
	return ret, nil
}

func (m *SQLiteDriver) indexes(tableName string) ([]*sqliteIndex, error) {
	var ret []*sqliteIndex
	rows, err := m.dbConn.Query(fmt.Sprintf("PRAGMA index_list('%s')", tableName))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var idx = &sqliteIndex{}
		var columns []string
		if err := rows.Scan(&idx.SeqNum, &idx.Name, &idx.Unique, &idx.Origin, &idx.Partial); err != nil {
			return nil, err
		}
		// get all columns stored within the index
		rowsColumns, err := m.dbConn.Query(fmt.Sprintf("PRAGMA index_info('%s')", idx.Name))
		if err != nil {
			return nil, err
		}
		for rowsColumns.Next() {
			var rankIndex, rankTable int
			var colName string
			if err := rowsColumns.Scan(&rankIndex, &rankTable, &colName); err != nil {
				return nil, errors.Wrapf(err, "unable to scan for index %s", idx.Name)
			}
			columns = append(columns, colName)
		}
		rowsColumns.Close()
		idx.Columns = columns
		ret = append(ret, idx)
	}
	return ret, nil
}

// Columns takes a table name and attempts to retrieve the table information
// from the database. It retrieves the column names
// and column types and returns those as a []Column after TranslateColumnType()
// converts the SQL types to Go types, for example: "varchar" to "string"
func (m *SQLiteDriver) Columns(schema, tableName string) ([]bdb.Column, error) {
	var columns []bdb.Column

	// get all indexes
	idxs, err := m.indexes(tableName)
	if err != nil {
		return nil, err
	}

	// finally get the remaining information about the columns
	tinfo, err := m.tableInfo(tableName)
	if err != nil {
		return nil, err
	}

	query := "SELECT 1 FROM sqlite_master WHERE type = 'table' AND name = ? AND sql LIKE '%AUTOINCREMENT%'"
	result, err := m.dbConn.Query(query, tableName)
	if err != nil {
		return nil, err
	}
	autoIncr := result.Next()
	if err := result.Close(); err != nil {
		return nil, err
	}

	for _, column := range tinfo {
		bColumn := bdb.Column{
			Name:       column.Name,
			FullDBType: strings.ToUpper(column.Type),
			DBType:     strings.ToUpper(column.Type),
			Nullable:   !column.NotNull,
		}

		// also get a correct information for Unique
		for _, idx := range idxs {
			for _, name := range idx.Columns {
				if name == column.Name {
					bColumn.Unique = idx.Unique > 0
				}
			}
		}

		if column.DefaultValue != nil && *column.DefaultValue != "NULL" {
			bColumn.Default = *column.DefaultValue
		} else if autoIncr {
			bColumn.Default = "auto_increment"
		}

		columns = append(columns, bColumn)
	}

	return columns, nil
}

// PrimaryKeyInfo looks up the primary key for a table.
func (m *SQLiteDriver) PrimaryKeyInfo(schema, tableName string) (*bdb.PrimaryKey, error) {
	// lookup the columns affected by the PK
	tinfo, err := m.tableInfo(tableName)
	if err != nil {
		return nil, err
	}

	var columns []string
	for _, column := range tinfo {
		if column.Pk > 0 {
			columns = append(columns, column.Name)
		}
	}

	var pk *bdb.PrimaryKey
	if len(columns) > 0 {
		pk = &bdb.PrimaryKey{Columns: columns}
	}
	return pk, nil
}

// ForeignKeyInfo retrieves the foreign keys for a given table name.
func (m *SQLiteDriver) ForeignKeyInfo(schema, tableName string) ([]bdb.ForeignKey, error) {
	var fkeys []bdb.ForeignKey

	query := fmt.Sprintf("PRAGMA foreign_key_list('%s')", tableName)

	var rows *sql.Rows
	var err error
	if rows, err = m.dbConn.Query(query, tableName); err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var fkey bdb.ForeignKey
		var onu, ond, match string
		var id, seq int

		fkey.Table = tableName
		err = rows.Scan(&id, &seq, &fkey.ForeignTable, &fkey.Column, &fkey.ForeignColumn, &onu, &ond, &match)
		if err != nil {
			return nil, err
		}
		fkey.Name = fmt.Sprintf("FK_%d", id)

		fkeys = append(fkeys, fkey)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return fkeys, nil
}

// TranslateColumnType converts sqlite database types to Go types, for example
// "varchar" to "string" and "bigint" to "int64". It returns this parsed data
// as a Column object.
// https://sqlite.org/datatype3.html
func (m *SQLiteDriver) TranslateColumnType(c bdb.Column) bdb.Column {
	if c.Nullable {
		switch strings.Split(c.DBType, "(")[0] {
		case "INT", "INTEGER", "BIGINT":
			c.Type = "null.Int64"
		case "TINYINT", "INT8":
			c.Type = "null.Int8"
		case "SMALLINT", "INT2":
			c.Type = "null.Int16"
		case "MEDIUMINT":
			c.Type = "null.Int32"
		case "UNSIGNED BIG INT":
			c.Type = "null.Uint64"
		case "CHARACTER", "VARCHAR", "VARYING CHARACTER", "NCHAR",
			"NATIVE CHARACTER", "NVARCHAR", "TEXT", "CLOB":
			c.Type = "null.String"
		case "BLOB":
			c.Type = "null.Bytes"
		case "FLOAT":
			c.Type = "null.Float32"
		case "REAL", "DOUBLE", "DOUBLE PRECISION", "NUMERIC", "DECIMAL":
			c.Type = "null.Float64"
		case "BOOLEAN":
			c.Type = "null.Bool"
		case "DATE", "DATETIME":
			c.Type = "null.Time"

		default:
			c.Type = "null.String"
		}
	} else {
		switch c.DBType {
		case "INT", "INTEGER", "BIGINT":
			c.Type = "int64"
		case "TINYINT", "INT8":
			c.Type = "int8"
		case "SMALLINT", "INT2":
			c.Type = "int16"
		case "MEDIUMINT":
			c.Type = "int32"
		case "UNSIGNED BIG INT":
			c.Type = "uint64"
		case "CHARACTER", "VARCHAR", "VARYING CHARACTER", "NCHAR",
			"NATIVE CHARACTER", "NVARCHAR", "TEXT", "CLOB":
			c.Type = "string"
		case "BLOB":
			c.Type = "[]byte"
		case "FLOAT":
			c.Type = "float32"
		case "REAL", "DOUBLE", "DOUBLE PRECISION", "NUMERIC", "DECIMAL":
			c.Type = "float64"
		case "BOOLEAN":
			c.Type = "bool"
		case "DATE", "DATETIME":
			c.Type = "time.Time"

		default:
			c.Type = "string"
		}
	}

	return c
}

// RightQuote is the quoting character for the right side of the identifier
func (m *SQLiteDriver) RightQuote() byte {
	return '"'
}

// LeftQuote is the quoting character for the left side of the identifier
func (m *SQLiteDriver) LeftQuote() byte {
	return '"'
}

// IndexPlaceholders returns false to indicate SQLite doesnt support indexed placeholders
func (m *SQLiteDriver) IndexPlaceholders() bool {
	return false
}
