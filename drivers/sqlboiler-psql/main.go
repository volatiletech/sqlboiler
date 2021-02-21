package main

import (
	"github.com/razor-1/sqlboiler/v4/drivers"
	"github.com/razor-1/sqlboiler/v4/drivers/sqlboiler-psql/driver"
)

func main() {
	drivers.DriverMain(&driver.PostgresDriver{})
}
