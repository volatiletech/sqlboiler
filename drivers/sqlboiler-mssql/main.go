package main

import (
	"github.com/razor-1/sqlboiler/v4/drivers"
	"github.com/razor-1/sqlboiler/v4/drivers/sqlboiler-mssql/driver"
)

func main() {
	drivers.DriverMain(&driver.MSSQLDriver{})
}
