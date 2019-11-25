package main

import (
	"github.com/razor-1/sqlboiler/drivers"
	"github.com/razor-1/sqlboiler/drivers/sqlboiler-mysql/driver"
)

func main() {
	drivers.DriverMain(&driver.MySQLDriver{})
}
