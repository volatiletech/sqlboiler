// +build linux,go1.8

package boilingcore

import (
	"plugin"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/vattle/sqlboiler/bdb"
	"github.com/vattle/sqlboiler/bdb/drivers"
)

type driverPlugin struct {
	*plugin.Plugin
}

func (d *driverPlugin) getDriver() (bdb.Interface, error) {
	sym, err := d.Lookup("InitDriver")
	if err != nil {
		return nil, errors.Wrap(err, "could not find symbol InitDriver in plugin")
	}

	initializer, ok := sym.(func() (bdb.Interface, error))
	if !ok {
		return nil, errors.New("symbol InitDriver is not `func() (bdb.Interface, error)`")
	}

	driver, err := initializer()
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize driver")
	}

	return driver, nil
}

func (d *driverPlugin) injectConfig() error {
	sym, err := d.Lookup("SetConfig")
	if err != nil {
		return nil
	}

	injector, ok := sym.(func(interface{}) error)
	if !ok {
		return errors.New("symbol SetConfig is not `func(interface{}) error")
	}

	injector(viper.GetViper())

	return nil
}

// initDriver attempts to set the state Interface based off the passed in
// driver flag value. If an invalid flag string is provided an error is returned.
func (s *State) initDriver(driverName string) error {
	// Create a driver based off driver flag
	switch driverName {
	case "postgres":
		s.Driver = drivers.NewPostgresDriver(
			s.Config.Postgres.User,
			s.Config.Postgres.Pass,
			s.Config.Postgres.DBName,
			s.Config.Postgres.Host,
			s.Config.Postgres.Port,
			s.Config.Postgres.SSLMode,
		)
	case "mysql":
		s.Driver = drivers.NewMySQLDriver(
			s.Config.MySQL.User,
			s.Config.MySQL.Pass,
			s.Config.MySQL.DBName,
			s.Config.MySQL.Host,
			s.Config.MySQL.Port,
			s.Config.MySQL.SSLMode,
		)
	case "mock":
		s.Driver = &drivers.MockDriver{}
	default:
		plg, err := plugin.Open(driverName)
		if err != nil {
			return errors.Wrap(err, "unable to open plugin")
		}

		driverPlg := &driverPlugin{Plugin: plg}
		if err := driverPlg.injectConfig(); err != nil {
			return err
		}

		driver, err := driverPlg.getDriver()
		if err != nil {
			return err
		}

		s.Driver = driver
	}

	if s.Driver == nil {
		return errors.New("An invalid driver name was provided")
	}

	s.Dialect.LQ = s.Driver.LeftQuote()
	s.Dialect.RQ = s.Driver.RightQuote()
	s.Dialect.IndexPlaceholders = s.Driver.IndexPlaceholders()

	return nil
}
