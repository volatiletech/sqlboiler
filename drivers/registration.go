package drivers

import "fmt"

// registeredDrivers are all the drivers which are currently registered
var registeredDrivers = map[string]Interface{}

// RegisterBinary is used to register drivers that are binaries.
// Panics if a driver with the same name has been previously loaded.
func RegisterBinary(name, path string) {
	register(name, binaryDriver(path))
}

// RegisterFromInit is typically called by a side-effect loaded driver
// during init time.
// Panics if a driver with the same name has been previously loaded.
func RegisterFromInit(name string, driver Interface) {
	register(name, driver)
}

// GetDriver retrieves the driver by name
func GetDriver(name string) Interface {
	if d, ok := registeredDrivers[name]; ok {
		return d
	}

	panic(fmt.Sprintf("drivers: sqlboiler driver %s has not been registered", name))
}

func register(name string, driver Interface) {
	if _, ok := registeredDrivers[name]; ok {
		panic(fmt.Sprintf("drivers: sqlboiler driver %s already loaded", name))
	}

	registeredDrivers[name] = driver
}
