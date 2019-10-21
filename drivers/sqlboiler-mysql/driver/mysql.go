package driver

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/friendsofgo/errors"
	"github.com/volatiletech/sqlboiler/drivers"
	"github.com/volatiletech/sqlboiler/importers"
)

func init() {
	drivers.RegisterFromInit("mysql", &MySQLDriver{})
}

//go:generate go-bindata -nometadata -pkg driver -prefix override override/...

// Assemble is more useful for calling into the library so you don't
// have to instantiate an empty type.
func Assemble(config drivers.Config) (dbinfo *drivers.DBInfo, err error) {
	driver := MySQLDriver{}
	return driver.Assemble(config)
}

// MySQLDriver holds the database connection string and a handle
// to the database connection.
type MySQLDriver struct {
	connStr string
	conn    *sql.DB

	tinyIntAsInt bool
}

// Templates that should be added/overridden
func (MySQLDriver) Templates() (map[string]string, error) {
	names := AssetNames()
	tpls := make(map[string]string)
	for _, n := range names {
		b, err := Asset(n)
		if err != nil {
			return nil, err
		}

		tpls[n] = base64.StdEncoding.EncodeToString(b)
	}

	return tpls, nil
}

// Assemble all the information we need to provide back to the driver
func (m *MySQLDriver) Assemble(config drivers.Config) (dbinfo *drivers.DBInfo, err error) {
	defer func() {
		if r := recover(); r != nil && err == nil {
			dbinfo = nil
			err = r.(error)
		}
	}()

	user := config.MustString(drivers.ConfigUser)
	pass, _ := config.String(drivers.ConfigPass)
	dbname := config.MustString(drivers.ConfigDBName)
	host := config.MustString(drivers.ConfigHost)
	port := config.DefaultInt(drivers.ConfigPort, 3306)
	sslmode := config.DefaultString(drivers.ConfigSSLMode, "true")

	schema := dbname
	whitelist, _ := config.StringSlice(drivers.ConfigWhitelist)
	blacklist, _ := config.StringSlice(drivers.ConfigBlacklist)

	tinyIntAsIntIntf, ok := config["tinyint_as_int"]
	if ok {
		if b, ok := tinyIntAsIntIntf.(bool); ok {
			m.tinyIntAsInt = b
		}
	}

	m.connStr = MySQLBuildQueryString(user, pass, dbname, host, port, sslmode)
	m.conn, err = sql.Open("mysql", m.connStr)
	if err != nil {
		return nil, errors.Wrap(err, "sqlboiler-mysql failed to connect to database")
	}

	defer func() {
		if e := m.conn.Close(); e != nil {
			dbinfo = nil
			err = e
		}
	}()

	dbinfo = &drivers.DBInfo{
		Dialect: drivers.Dialect{
			LQ: '`',
			RQ: '`',

			UseLastInsertID: true,
			UseSchema:       false,
		},
	}

	dbinfo.Tables, err = drivers.Tables(m, schema, whitelist, blacklist)
	if err != nil {
		return nil, err
	}

	return dbinfo, err
}

// MySQLBuildQueryString builds a query string for MySQL.
func MySQLBuildQueryString(user, pass, dbname, host string, port int, sslmode string) string {
	config := mysql.NewConfig()

	config.User = user
	if len(pass) != 0 {
		config.Passwd = pass
	}
	config.DBName = dbname
	config.Net = "tcp"
	config.Addr = host
	if port == 0 {
		port = 3306
	}
	config.Addr += ":" + strconv.Itoa(port)
	config.TLSConfig = sslmode

	// MySQL is a bad, and by default reads date/datetime into a []byte
	// instead of a time.Time. Tell it to stop being a bad.
	config.ParseTime = true

	return config.FormatDSN()
}

// TableNames connects to the mysql database and
// retrieves all table names from the information_schema where the
// table schema is public.
func (m *MySQLDriver) TableNames(schema string, whitelist, blacklist []string) ([]string, error) {
	var names []string

	query := fmt.Sprintf(`select table_name from information_schema.tables where table_schema = ? and table_type = 'BASE TABLE'`)
	args := []interface{}{schema}
	if len(whitelist) > 0 {
		tables := drivers.TablesFromList(whitelist)
		if len(tables) > 0 {
			query += fmt.Sprintf(" and table_name in (%s)", strings.Repeat(",?", len(tables))[1:])
			for _, w := range tables {
				args = append(args, w)
			}
		}
	} else if len(blacklist) > 0 {
		tables := drivers.TablesFromList(blacklist)
		if len(tables) > 0 {
			query += fmt.Sprintf(" and table_name not in (%s)", strings.Repeat(",?", len(tables))[1:])
			for _, b := range tables {
				args = append(args, b)
			}
		}
	}

	query += ` order by table_name;`

	rows, err := m.conn.Query(query, args...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		names = append(names, name)
	}

	return names, nil
}

// Columns takes a table name and attempts to retrieve the table information
// from the database information_schema.columns. It retrieves the column names
// and column types and returns those as a []Column after TranslateColumnType()
// converts the SQL types to Go types, for example: "varchar" to "string"
func (m *MySQLDriver) Columns(schema, tableName string, whitelist, blacklist []string) ([]drivers.Column, error) {
	var columns []drivers.Column
	args := []interface{}{tableName, schema}

	query := `
	select
	c.column_name,
	c.column_type,
	if(c.data_type = 'enum', c.column_type, c.data_type),
	if(extra = 'auto_increment','auto_increment', c.column_default),
	c.is_nullable = 'YES',
		exists (
			select c.column_name
			from information_schema.table_constraints tc
			inner join information_schema.key_column_usage kcu
				on tc.constraint_name = kcu.constraint_name and tc.table_name = kcu.table_name and tc.table_schema = kcu.table_schema
			where c.table_name = tc.table_name and c.column_name = kcu.column_name and c.table_schema = kcu.constraint_schema and 
				(tc.constraint_type = 'PRIMARY KEY' or tc.constraint_type = 'UNIQUE') and
				(select count(*) from information_schema.key_column_usage where table_schema = kcu.table_schema and constraint_schema = kcu.table_schema and table_name = tc.table_name and constraint_name = tc.constraint_name) = 1
		) as is_unique
	from information_schema.columns as c
	where table_name = ? and table_schema = ? and c.extra not like '%VIRTUAL%'`

	if len(whitelist) > 0 {
		cols := drivers.ColumnsFromList(whitelist, tableName)
		if len(cols) > 0 {
			query += fmt.Sprintf(" and c.column_name in (%s)", strings.Repeat(",?", len(cols))[1:])
			for _, w := range cols {
				args = append(args, w)
			}
		}
	} else if len(blacklist) > 0 {
		cols := drivers.ColumnsFromList(blacklist, tableName)
		if len(cols) > 0 {
			query += fmt.Sprintf(" and c.column_name not in (%s)", strings.Repeat(",?", len(cols))[1:])
			for _, w := range cols {
				args = append(args, w)
			}
		}
	}

	query += ` order by c.ordinal_position;`

	rows, err := m.conn.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var colName, colType, colFullType string
		var nullable, unique bool
		var defaultValue *string
		if err := rows.Scan(&colName, &colFullType, &colType, &defaultValue, &nullable, &unique); err != nil {
			return nil, errors.Wrapf(err, "unable to scan for table %s", tableName)
		}

		column := drivers.Column{
			Name:       colName,
			FullDBType: colFullType, // example: tinyint(1) instead of tinyint
			DBType:     colType,
			Nullable:   nullable,
			Unique:     unique,
		}

		if defaultValue != nil && *defaultValue != "NULL" {
			column.Default = *defaultValue
		}

		columns = append(columns, column)
	}

	return columns, nil
}

// PrimaryKeyInfo looks up the primary key for a table.
func (m *MySQLDriver) PrimaryKeyInfo(schema, tableName string) (*drivers.PrimaryKey, error) {
	pkey := &drivers.PrimaryKey{}
	var err error

	query := `
	select tc.constraint_name
	from information_schema.table_constraints as tc
	where tc.table_name = ? and tc.constraint_type = 'PRIMARY KEY' and tc.table_schema = ?;`

	row := m.conn.QueryRow(query, tableName, schema)
	if err = row.Scan(&pkey.Name); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	queryColumns := `
	select kcu.column_name
	from   information_schema.key_column_usage as kcu
	where  table_name = ? and constraint_name = ? and table_schema = ?;`

	var rows *sql.Rows
	if rows, err = m.conn.Query(queryColumns, tableName, pkey.Name, schema); err != nil {
		return nil, err
	}
	defer rows.Close()

	var columns []string
	for rows.Next() {
		var column string

		err = rows.Scan(&column)
		if err != nil {
			return nil, err
		}

		columns = append(columns, column)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	pkey.Columns = columns

	return pkey, nil
}

// ForeignKeyInfo retrieves the foreign keys for a given table name.
func (m *MySQLDriver) ForeignKeyInfo(schema, tableName string) ([]drivers.ForeignKey, error) {
	var fkeys []drivers.ForeignKey

	query := `
	select constraint_name, table_name, column_name, referenced_table_name, referenced_column_name
	from information_schema.key_column_usage
	where table_schema = ? and referenced_table_schema = ? and table_name = ?
	order by constraint_name, table_name, column_name, referenced_table_name, referenced_column_name
	`

	var rows *sql.Rows
	var err error
	if rows, err = m.conn.Query(query, schema, schema, tableName); err != nil {
		return nil, err
	}

	for rows.Next() {
		var fkey drivers.ForeignKey
		var sourceTable string

		fkey.Table = tableName
		err = rows.Scan(&fkey.Name, &sourceTable, &fkey.Column, &fkey.ForeignTable, &fkey.ForeignColumn)
		if err != nil {
			return nil, err
		}

		fkeys = append(fkeys, fkey)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return fkeys, nil
}

// TranslateColumnType converts mysql database types to Go types, for example
// "varchar" to "string" and "bigint" to "int64". It returns this parsed data
// as a Column object.
func (m *MySQLDriver) TranslateColumnType(c drivers.Column) drivers.Column {
	unsigned := strings.Contains(c.FullDBType, "unsigned")
	if c.Nullable {
		switch c.DBType {
		case "tinyint":
			// map tinyint(1) to bool if TinyintAsBool is true
			if !m.tinyIntAsInt && c.FullDBType == "tinyint(1)" {
				c.Type = "null.Bool"
			} else if unsigned {
				c.Type = "null.Uint8"
			} else {
				c.Type = "null.Int8"
			}
		case "smallint":
			if unsigned {
				c.Type = "null.Uint16"
			} else {
				c.Type = "null.Int16"
			}
		case "mediumint":
			if unsigned {
				c.Type = "null.Uint32"
			} else {
				c.Type = "null.Int32"
			}
		case "int", "integer":
			if unsigned {
				c.Type = "null.Uint"
			} else {
				c.Type = "null.Int"
			}
		case "bigint":
			if unsigned {
				c.Type = "null.Uint64"
			} else {
				c.Type = "null.Int64"
			}
		case "float":
			c.Type = "null.Float32"
		case "double", "double precision", "real":
			c.Type = "null.Float64"
		case "boolean", "bool":
			c.Type = "null.Bool"
		case "date", "datetime", "timestamp":
			c.Type = "null.Time"
		case "binary", "varbinary", "tinyblob", "blob", "mediumblob", "longblob":
			c.Type = "null.Bytes"
		case "numeric", "decimal", "dec", "fixed":
			c.Type = "types.NullDecimal"
		case "json":
			c.Type = "null.JSON"
		default:
			c.Type = "null.String"
		}
	} else {
		switch c.DBType {
		case "tinyint":
			// map tinyint(1) to bool if TinyintAsBool is true
			if !m.tinyIntAsInt && c.FullDBType == "tinyint(1)" {
				c.Type = "bool"
			} else if unsigned {
				c.Type = "uint8"
			} else {
				c.Type = "int8"
			}
		case "smallint":
			if unsigned {
				c.Type = "uint16"
			} else {
				c.Type = "int16"
			}
		case "mediumint":
			if unsigned {
				c.Type = "uint32"
			} else {
				c.Type = "int32"
			}
		case "int", "integer":
			if unsigned {
				c.Type = "uint"
			} else {
				c.Type = "int"
			}
		case "bigint":
			if unsigned {
				c.Type = "uint64"
			} else {
				c.Type = "int64"
			}
		case "float":
			c.Type = "float32"
		case "double", "double precision", "real":
			c.Type = "float64"
		case "boolean", "bool":
			c.Type = "bool"
		case "date", "datetime", "timestamp":
			c.Type = "time.Time"
		case "binary", "varbinary", "tinyblob", "blob", "mediumblob", "longblob":
			c.Type = "[]byte"
		case "numeric", "decimal", "dec", "fixed":
			c.Type = "types.Decimal"
		case "json":
			c.Type = "types.JSON"
		default:
			c.Type = "string"
		}
	}

	return c
}

// Imports returns important imports for the driver
func (MySQLDriver) Imports() (col importers.Collection, err error) {
	col.All = importers.Set{
		Standard: importers.List{
			`"strconv"`,
		},
	}

	col.Singleton = importers.Map{
		"mysql_upsert": {
			Standard: importers.List{
				`"fmt"`,
				`"strings"`,
			},
			ThirdParty: importers.List{
				`"github.com/volatiletech/sqlboiler/strmangle"`,
				`"github.com/volatiletech/sqlboiler/drivers"`,
			},
		},
	}

	col.TestSingleton = importers.Map{
		"mysql_suites_test": {
			Standard: importers.List{
				`"testing"`,
			},
		},
		"mysql_main_test": {
			Standard: importers.List{
				`"bytes"`,
				`"database/sql"`,
				`"fmt"`,
				`"io"`,
				`"io/ioutil"`,
				`"os"`,
				`"os/exec"`,
				`"regexp"`,
				`"strings"`,
			},
			ThirdParty: importers.List{
				`"github.com/kat-co/vala"`,
				`"github.com/friendsofgo/errors"`,
				`"github.com/spf13/viper"`,
				`"github.com/volatiletech/sqlboiler/drivers/sqlboiler-mysql/driver"`,
				`"github.com/volatiletech/sqlboiler/randomize"`,
				`_ "github.com/go-sql-driver/mysql"`,
			},
		},
	}

	col.BasedOnType = importers.Map{
		"null.Float32": {
			ThirdParty: importers.List{`"github.com/volatiletech/null"`},
		},
		"null.Float64": {
			ThirdParty: importers.List{`"github.com/volatiletech/null"`},
		},
		"null.Int": {
			ThirdParty: importers.List{`"github.com/volatiletech/null"`},
		},
		"null.Int8": {
			ThirdParty: importers.List{`"github.com/volatiletech/null"`},
		},
		"null.Int16": {
			ThirdParty: importers.List{`"github.com/volatiletech/null"`},
		},
		"null.Int32": {
			ThirdParty: importers.List{`"github.com/volatiletech/null"`},
		},
		"null.Int64": {
			ThirdParty: importers.List{`"github.com/volatiletech/null"`},
		},
		"null.Uint": {
			ThirdParty: importers.List{`"github.com/volatiletech/null"`},
		},
		"null.Uint8": {
			ThirdParty: importers.List{`"github.com/volatiletech/null"`},
		},
		"null.Uint16": {
			ThirdParty: importers.List{`"github.com/volatiletech/null"`},
		},
		"null.Uint32": {
			ThirdParty: importers.List{`"github.com/volatiletech/null"`},
		},
		"null.Uint64": {
			ThirdParty: importers.List{`"github.com/volatiletech/null"`},
		},
		"null.String": {
			ThirdParty: importers.List{`"github.com/volatiletech/null"`},
		},
		"null.Bool": {
			ThirdParty: importers.List{`"github.com/volatiletech/null"`},
		},
		"null.Time": {
			ThirdParty: importers.List{`"github.com/volatiletech/null"`},
		},
		"null.Bytes": {
			ThirdParty: importers.List{`"github.com/volatiletech/null"`},
		},
		"null.JSON": {
			ThirdParty: importers.List{`"github.com/volatiletech/null"`},
		},

		"time.Time": {
			Standard: importers.List{`"time"`},
		},
		"types.JSON": {
			ThirdParty: importers.List{`"github.com/volatiletech/sqlboiler/types"`},
		},
		"types.Decimal": {
			ThirdParty: importers.List{`"github.com/volatiletech/sqlboiler/types"`},
		},
		"types.NullDecimal": {
			ThirdParty: importers.List{`"github.com/volatiletech/sqlboiler/types"`},
		},
	}
	return col, err
}
