package main

import (
	"github.com/twitter-payments/sqlboiler/v4/drivers"
	"github.com/twitter-payments/sqlboiler/v4/drivers/sqlboiler-psql/driver"
)

func main() {
	drivers.DriverMain(&driver.PostgresDriver{})
}
