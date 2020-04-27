package main

import (
	"github.com/volatiletech/sqlboiler/v4/drivers"
	"github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-mssql/driver"
)

func main() {
	drivers.DriverMain(&driver.MSSQLDriver{})
}
