package main

import (
	"github.com/IOTechSystems/sqlboiler/v4/drivers"
	"github.com/IOTechSystems/sqlboiler/v4/drivers/sqlboiler-mysql/driver"
)

func main() {
	drivers.DriverMain(&driver.MySQLDriver{})
}
