package main

import (
	"github.com/razor-1/sqlboiler/v4/drivers"
	"github.com/razor-1/sqlboiler/v4/drivers/sqlboiler-sqlite3/driver"
)

func main() {
	drivers.DriverMain(&driver.SQLiteDriver{})
}
