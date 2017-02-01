// +build plugin

package main

import (
	"github.com/spf13/viper"
	"github.com/vattle/sqlboiler/bdb"
	"github.com/vattle/sqlboiler/bdb/drivers"
)

// InitDriver enables the default mysql driver to be used as a plugin. You could
// also implement your own custom driver if you wanted.
func InitDriver() (bdb.Interface, error) {
	return drivers.NewMySQLDriver(
		viper.GetString("mysql.user"),
		viper.GetString("mysql.pass"),
		viper.GetString("mysql.dbname"),
		viper.GetString("mysql.host"),
		viper.GetInt("mysql.port"),
		viper.GetString("mysql.sslmode"),
	), nil
}
