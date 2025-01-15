package main

import (
	"github.com/twitter-payments/sqlboiler/v4/drivers"
	"github.com/twitter-payments/sqlboiler/v4/drivers/sqlboiler-sqlite3/driver"
)

func main() {
	drivers.DriverMain(&driver.SQLiteDriver{})
}
