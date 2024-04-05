package main

import (
	"github.com/IOTechSystems/sqlboiler/v4/drivers"
	"github.com/IOTechSystems/sqlboiler/v4/drivers/sqlboiler-mssql/driver"
)

func main() {
	drivers.DriverMain(&driver.MSSQLDriver{})
}
