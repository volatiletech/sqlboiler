package main

import (
	"github.com/volatiletech/sqlboiler/v5/drivers"
	"github.com/volatiletech/sqlboiler/v5/drivers/sqlboiler-sqlite3/driver"
)

func main() {
	drivers.DriverMain(&driver.SQLiteDriver{})
}
