package main

import (
	"github.com/volatiletech/sqlboiler/v5/drivers"
	"github.com/volatiletech/sqlboiler/v5/drivers/sqlboiler-mssql/driver"
)

func main() {
	drivers.DriverMain(&driver.MSSQLDriver{})
}
