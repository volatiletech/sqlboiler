package main

import (
	"github.com/volatiletech/sqlboiler/drivers"
	"github.com/volatiletech/sqlboiler/drivers/sqlboiler-mysql/driver"
)

func main() {
	drivers.DriverMain(&driver.MySQLDriver{})
}
