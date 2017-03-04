// +build plugin

package main

import (
	"github.com/vattle/sqlboiler/bdb"
	"github.com/vattle/sqlboiler/bdb/drivers"
)

type configGetter interface {
	GetString(string) string
	GetInt(string) int
}

var config configGetter

func SetConfig(c interface{}) error {
	config = c.(configGetter)
	return nil
}

// InitDriver enables the default mysql driver to be used as a plugin. You could
// also implement your own custom driver if you wanted.
func InitDriver() (bdb.Interface, error) {
	return drivers.NewMySQLDriver(
		config.GetString("mysql.user"),
		config.GetString("mysql.pass"),
		config.GetString("mysql.dbname"),
		config.GetString("mysql.host"),
		config.GetInt("mysql.port"),
		config.GetString("mysql.sslmode"),
	), nil
}
