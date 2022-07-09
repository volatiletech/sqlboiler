package drivers

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

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

// RegisterBinaryFromCmdArg is used to register drivers from a command line argument
// The argument is either just the driver name or a path to a specific driver
// Panics if a driver with the same name has been previously loaded.
func RegisterBinaryFromCmdArg(arg string) (name, path string, err error) {
	path, err = getFullPath(arg)
	if err != nil {
		return name, path, err
	}

	name = getNameFromPath(path)

	RegisterBinary(name, path)

	return name, path, nil
}

// Get the full path to the driver binary from the given path
// the path can also be just the driver name e.g. "psql"
func getFullPath(path string) (string, error) {
	var err error

	if strings.ContainsRune(path, os.PathSeparator) {
		return path, nil
	}

	path, err = exec.LookPath("sqlboiler-" + path)
	if err != nil {
		return path, fmt.Errorf("could not find driver executable: %w", err)
	}

	path, err = filepath.Abs(path)
	if err != nil {
		return path, fmt.Errorf("could not find absolute path to driver: %w", err)
	}

	return path, nil
}

// Get the driver name from the path.
// strips the "sqlboiler-" prefix if it exists
// strips the ".exe" suffix if it exits
func getNameFromPath(name string) string {
	name = strings.Replace(filepath.Base(name), "sqlboiler-", "", 1)
	name = strings.Replace(name, ".exe", "", 1)

	return name
}
