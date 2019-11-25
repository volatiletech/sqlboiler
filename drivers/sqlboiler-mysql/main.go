package main

import (
	"github.com/razor-1/sqlboiler/v3/drivers"
	"github.com/razor-1/sqlboiler/v3/drivers/sqlboiler-mysql/driver"
)

func main() {
	drivers.DriverMain(&driver.MySQLDriver{})
}
