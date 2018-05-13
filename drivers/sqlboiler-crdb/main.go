package main

import (
	"github.com/volatiletech/sqlboiler/drivers"
	"github.com/volatiletech/sqlboiler/drivers/sqlboiler-crdb/driver"
)

func main() {
	drivers.DriverMain(&driver.CockroachDBDriver{})
}
