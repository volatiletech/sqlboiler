// The reason that so many of these functions seem like duplication is simply for the better error
// messages so each driver doesn't need to deal with validation on these simple field types

package drivers

import (
	"os"
	"strconv"
	"strings"

	"github.com/friendsofgo/errors"
)

// Config is a map with helper functions
type Config map[string]interface{}

// MustString retrieves a string that must exist and must be a string, it must also not be the empty string
func (c Config) MustString(key string) string {
	s, ok := c[key]
	if !ok {
		panic(errors.Errorf("failed to find key %s in config", key))
	}

	str, ok := s.(string)
	if !ok {
		panic(errors.Errorf("found key %s in config, but it was not a string (%T)", key, s))
	}

	if len(str) == 0 {
		panic(errors.Errorf("found key %s in config, but it was an empty string", key))
	}

	return str
}

// MustInt retrieves a string that must exist and must be an int, and it must not be 0
func (c Config) MustInt(key string) int {
	i, ok := c[key]
	if !ok {
		panic(errors.Errorf("failed to find key %s in config", key))
	}

	var integer int
	switch t := i.(type) {
	case int:
		integer = t
	case float64:
		integer = int(t)
	case string:
		var err error
		integer, err = strconv.Atoi(t)
		if err != nil {
			panic(errors.Errorf("failed to parse key %s (%s) to int: %v", key, t, err))
		}
	default:
		panic(errors.Errorf("found key %s in config, but it was not an int (%T)", key, i))
	}

	if integer == 0 {
		panic(errors.Errorf("found key %s in config, but its value was 0", key))
	}

	return integer
}

// String retrieves a non-empty string, the bool says if it exists, is of appropriate type,
// and has a non-zero length or not.
func (c Config) String(key string) (string, bool) {
	s, ok := c[key]
	if !ok {
		return "", false
	}

	str, ok := s.(string)
	if !ok {
		return "", false
	}

	if len(str) == 0 {
		return "", false
	}

	return str, true
}

// DefaultString retrieves a non-empty string or the default value provided.
func (c Config) DefaultString(key, def string) string {
	str, ok := c.String(key)
	if !ok {
		return def
	}

	return str
}

// Int retrieves an int, the bool says if it exists, is of the appropriate type,
// and is non-zero. Coerces float64 to int because JSON and Javascript kinda suck.
func (c Config) Int(key string) (int, bool) {
	i, ok := c[key]
	if !ok {
		return 0, false
	}

	var integer int
	switch t := i.(type) {
	case int:
		integer = t
	case float64:
		integer = int(t)
	case string:
		var err error
		integer, err = strconv.Atoi(t)
		if err != nil {
			return 0, false
		}
	default:
		return 0, false
	}

	if integer == 0 {
		return 0, false
	}

	return integer, true
}

// DefaultInt retrieves a non-zero int or the default value provided.
func (c Config) DefaultInt(key string, def int) int {
	i, ok := c.Int(key)
	if !ok {
		return def
	}

	return i
}

// StringSlice retrieves an string slice, the bool says if it exists, is of the appropriate type,
// is non-nil and non-zero length
func (c Config) StringSlice(key string) ([]string, bool) {
	ss, ok := c[key]
	if !ok {
		return nil, false
	}

	var slice []string
	if intfSlice, ok := ss.([]interface{}); ok {
		for _, i := range intfSlice {
			slice = append(slice, i.(string))
		}
	} else if stringSlice, ok := ss.([]string); ok {
		slice = stringSlice
	} else {
		return nil, false
	}

	// Also detects nil
	if len(slice) == 0 {
		return nil, false
	}
	return slice, true
}

func (c Config) MustForeignKeys(key string) []ForeignKey {
	rawValue, ok := c[key]
	if !ok {
		return nil
	}

	switch v := rawValue.(type) {
	case nil:
		return nil
	case []ForeignKey:
		return v
	case []interface{}: // in case binary, config is pass to driver in json format, so this key will be []interface{}
		fks := make([]ForeignKey, 0, len(v))
		for _, item := range v {
			fk, ok := item.(map[string]interface{})
			if !ok {
				panic(errors.Errorf("found item of foreign keys, but it was not a map[string]interface{} (%T)", v))
			}

			configFK := Config(fk)

			fks = append(fks, ForeignKey{
				Name:          configFK.MustString("name"),
				Table:         configFK.MustString("table"),
				Column:        configFK.MustString("column"),
				ForeignTable:  configFK.MustString("foreign_table"),
				ForeignColumn: configFK.MustString("foreign_column"),
			})
		}
		return fks
	default:
		panic(errors.Errorf("found key %s in config, but it was invalid (%T)", key, v))
	}
}

// DefaultEnv grabs a value from the environment or a default.
// This is shared by drivers to get config for testing.
func DefaultEnv(key, def string) string {
	val := os.Getenv(key)
	if len(val) == 0 {
		val = def
	}
	return val
}

// TablesFromList takes a whitelist or blacklist and returns
// the table names.
func TablesFromList(list []string) []string {
	if len(list) == 0 {
		return nil
	}

	var tables []string
	for _, i := range list {
		splits := strings.Split(i, ".")

		if len(splits) == 1 {
			tables = append(tables, splits[0])
		}
	}

	return tables
}

// ColumnsFromList takes a whitelist or blacklist and returns
// the columns for a given table.
func ColumnsFromList(list []string, tablename string) []string {
	if len(list) == 0 {
		return nil
	}

	var columns []string
	for _, i := range list {
		splits := strings.Split(i, ".")

		if len(splits) != 2 {
			continue
		}

		if splits[0] == tablename || splits[0] == "*" {
			columns = append(columns, splits[1])
		}
	}

	return columns
}

// CombineConfigAndDBForeignKeys takes foreign keys from both config and db, filter by tableName and
// deduplicate by column name. If a foreign key is found in both config and db, the one in config will be used.
func CombineConfigAndDBForeignKeys(configForeignKeys []ForeignKey, tableName string, dbForeignKeys []ForeignKey) []ForeignKey {
	combinedForeignKeys := make([]ForeignKey, 0, len(configForeignKeys)+len(dbForeignKeys))
	appearedColumns := make(map[string]bool)
	fkNameCount := make(map[string]int)

	// Detect Composite Foreign Keys in the database by counting how many times they appear.
	// boiler doesn't support Composite FKs and should ignore those
	for _, fk := range dbForeignKeys {
		fkNameCount[fk.Name] += 1
	}

	for _, fk := range configForeignKeys {
		// need check table name here cause configForeignKeys contains all foreign keys of all tables
		if fk.Table != tableName {
			continue
		}
		if appearedColumns[fk.Column] {
			continue
		}

		combinedForeignKeys = append(combinedForeignKeys, fk)
		appearedColumns[fk.Column] = true
	}

	for _, fk := range dbForeignKeys {
		// no need check table here, because dbForeignKeys are already filtered by table name
		if appearedColumns[fk.Column] {
			continue
		}
		if fkNameCount[fk.Name] != 1 {
			continue
		}
		combinedForeignKeys = append(combinedForeignKeys, fk)
		appearedColumns[fk.Column] = true
	}

	return combinedForeignKeys
}
