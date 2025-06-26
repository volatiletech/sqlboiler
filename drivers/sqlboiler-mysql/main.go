package main

import (
	"github.com/aarondl/sqlboiler/v4/drivers"
	"github.com/aarondl/sqlboiler/v4/drivers/sqlboiler-mysql/driver"
)

func main() {
	drivers.DriverMain(&driver.MySQLDriver{})
}
