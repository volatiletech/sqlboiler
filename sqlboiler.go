// Package sqlboiler has types and methods useful for generating code that
// acts as a fully dynamic ORM might.
package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/nullbio/sqlboiler/dbdrivers"
)

const (
	templatesDirectory          = "cmds/templates"
	templatesSingletonDirectory = "cmds/templates/singleton"

	templatesTestDirectory          = "cmds/templates_test"
	templatesSingletonTestDirectory = "cmds/templates_test/singleton"
)

// State holds the global data needed by most pieces to run
type State struct {
	Config *Config

	Driver dbdrivers.Interface
	Tables []dbdrivers.Table

	Templates              templateList
	TestTemplates          templateList
	SingletonTemplates     templateList
	SingletonTestTemplates templateList

	TestMainTemplate *template.Template
}

// New creates a new state based off of the config
func New(config *Config) (*State, error) {
	s := &State{}

	err := s.initDriver(config.DriverName)
	if err != nil {
		return nil, err
	}

	// Connect to the driver database
	if err = s.Driver.Open(); err != nil {
		return nil, fmt.Errorf("Unable to connect to the database: %s", err)
	}

	err = s.initTables(config.TableName)
	if err != nil {
		return nil, fmt.Errorf("Unable to initialize tables: %s", err)
	}

	err = s.initOutFolder()
	if err != nil {
		return nil, fmt.Errorf("Unable to initialize the output folder: %s", err)
	}

	err = s.initTemplates()
	if err != nil {
		return nil, fmt.Errorf("Unable to initialize templates: %s", err)
	}

	return s, nil
}

// Run executes the sqlboiler templates and outputs them to files based on the
// state given.
func (s *State) Run(includeTests bool) error {
	singletonData := &templateData{
		Tables:     s.Tables,
		DriverName: s.Config.DriverName,
		PkgName:    s.Config.PkgName,
	}

	if err := generateSingletonOutput(s, singletonData); err != nil {
		return fmt.Errorf("Unable to generate singleton template output: %s", err)
	}

	if includeTests {
		if err := generateTestMainOutput(s, singletonData); err != nil {
			return fmt.Errorf("Unable to generate TestMain output: %s", err)
		}

		if err := generateSingletonTestOutput(s, singletonData); err != nil {
			return fmt.Errorf("Unable to generate singleton test template output: %s", err)
		}
	}

	for _, table := range s.Tables {
		if table.IsJoinTable {
			continue
		}

		data := &templateData{
			Table:      table,
			DriverName: s.Config.DriverName,
			PkgName:    s.Config.PkgName,
		}

		// Generate the regular templates
		if err := generateOutput(s, data); err != nil {
			return fmt.Errorf("Unable to generate output: %s", err)
		}

		// Generate the test templates
		if includeTests {
			if err := generateTestOutput(s, data); err != nil {
				return fmt.Errorf("Unable to generate test output: %s", err)
			}
		}
	}

	return nil
}

// Cleanup closes any resources that must be closed
func (s *State) Cleanup() error {
	s.Driver.Close()
	return nil
}

// initTemplates loads all template folders into the state object.
func (s *State) initTemplates() error {
	var err error

	s.Templates, err = loadTemplates(templatesDirectory)
	if err != nil {
		return err
	}

	s.SingletonTemplates, err = loadTemplates(templatesSingletonDirectory)
	if err != nil {
		return err
	}

	s.TestTemplates, err = loadTemplates(templatesTestDirectory)
	if err != nil {
		return err
	}

	s.SingletonTestTemplates, err = loadTemplates(templatesSingletonTestDirectory)
	if err != nil {
		return err
	}

	return nil
}

// initDriver attempts to set the state Interface based off the passed in
// driver flag value. If an invalid flag string is provided an error is returned.
func (s *State) initDriver(driverName string) error {
	// Create a driver based off driver flag
	switch driverName {
	case "postgres":
		s.Driver = dbdrivers.NewPostgresDriver(
			s.Config.Postgres.User,
			s.Config.Postgres.Pass,
			s.Config.Postgres.DBName,
			s.Config.Postgres.Host,
			s.Config.Postgres.Port,
		)
	}

	if s.Driver == nil {
		return errors.New("An invalid driver name was provided")
	}

	return nil
}

// initTables will create a string slice out of the passed in table flag value
// if one is provided. If no flag is provided, it will attempt to connect to the
// database to retrieve all "public" table names, and build a slice out of that
// result.
func (s *State) initTables(tableName string) error {
	var tableNames []string

	if len(tableName) != 0 {
		tableNames = strings.Split(tableName, ",")
		for i, name := range tableNames {
			tableNames[i] = strings.TrimSpace(name)
		}
	}

	var err error
	s.Tables, err = dbdrivers.Tables(s.Driver, tableNames...)
	if err != nil {
		return fmt.Errorf("Unable to get all table names: %s", err)
	}

	if len(s.Tables) == 0 {
		return errors.New("No tables found in database, migrate some tables first")
	}

	if err := checkPKeys(s.Tables); err != nil {
		return err
	}

	return nil
}

// initOutFolder creates the folder that will hold the generated output.
func (s *State) initOutFolder() error {
	if err := os.MkdirAll(s.Config.OutFolder, os.ModePerm); err != nil {
		return fmt.Errorf("Unable to make output folder: %s", err)
	}

	return nil
}

// checkPKeys ensures every table has a primary key column
func checkPKeys(tables []dbdrivers.Table) error {
	var missingPkey []string
	for _, t := range tables {
		if t.PKey == nil {
			missingPkey = append(missingPkey, t.Name)
		}
	}

	if len(missingPkey) != 0 {
		return fmt.Errorf("Cannot continue until the follow tables have PRIMARY KEY columns: %s", strings.Join(missingPkey, ", "))
	}

	return nil
}
