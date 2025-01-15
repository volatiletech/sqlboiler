package main

import (
	"github.com/twitter-payments/sqlboiler/v4/drivers"
	"github.com/twitter-payments/sqlboiler/v4/drivers/sqlboiler-mssql/driver"
)

func main() {
	drivers.DriverMain(&driver.MSSQLDriver{})
}
