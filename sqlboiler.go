// Package sqlboiler has types and methods useful for generating code that
// acts as a fully dynamic ORM might.
package main

import (
	"encoding/json"
	"fmt"
	"go/build"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/pkg/errors"
	"github.com/vattle/sqlboiler/bdb"
	"github.com/vattle/sqlboiler/bdb/drivers"
	"github.com/vattle/sqlboiler/boil"
)

const (
	templatesDirectory          = "templates"
	templatesSingletonDirectory = "templates/singleton"

	templatesTestDirectory          = "templates_test"
	templatesSingletonTestDirectory = "templates_test/singleton"

	templatesTestMainDirectory = "templates_test/main_test"
)

// State holds the global data needed by most pieces to run
type State struct {
	Config *Config

	Driver bdb.Interface
	Tables []bdb.Table

	Templates              *templateList
	TestTemplates          *templateList
	SingletonTemplates     *templateList
	SingletonTestTemplates *templateList

	TestMainTemplate *template.Template
}

// New creates a new state based off of the config
func New(config *Config) (*State, error) {
	s := &State{
		Config: config,
	}

	err := s.initDriver(config.DriverName)
	if err != nil {
		return nil, err
	}

	// Connect to the driver database
	if err = s.Driver.Open(); err != nil {
		return nil, errors.Wrap(err, "unable to connect to the database")
	}

	err = s.initTables(config.Schema, config.WhitelistTables, config.BlacklistTables)
	if err != nil {
		return nil, errors.Wrap(err, "unable to initialize tables")
	}

	if s.Config.Debug {
		b, err := json.Marshal(s.Tables)
		if err != nil {
			return nil, errors.Wrap(err, "unable to json marshal tables")
		}
		boil.DebugWriter.Write(b)
		fmt.Fprintln(boil.DebugWriter)
	}

	err = s.initOutFolder()
	if err != nil {
		return nil, errors.Wrap(err, "unable to initialize the output folder")
	}

	err = s.initTemplates()
	if err != nil {
		return nil, errors.Wrap(err, "unable to initialize templates")
	}

	err = s.initTags(config.Tags)
	if err != nil {
		return nil, errors.Wrap(err, "unable to initialize struct tags")
	}

	return s, nil
}

// Run executes the sqlboiler templates and outputs them to files based on the
// state given.
func (s *State) Run(includeTests bool) error {
	singletonData := &templateData{
		Tables:           s.Tables,
		Schema:           s.Config.Schema,
		DriverName:       s.Config.DriverName,
		UseLastInsertID:  s.Driver.UseLastInsertID(),
		PkgName:          s.Config.PkgName,
		NoHooks:          s.Config.NoHooks,
		NoAutoTimestamps: s.Config.NoAutoTimestamps,

		StringFuncs: templateStringMappers,
	}

	if err := generateSingletonOutput(s, singletonData); err != nil {
		return errors.Wrap(err, "singleton template output")
	}

	if !s.Config.NoTests && includeTests {
		if err := generateTestMainOutput(s, singletonData); err != nil {
			return errors.Wrap(err, "unable to generate TestMain output")
		}

		if err := generateSingletonTestOutput(s, singletonData); err != nil {
			return errors.Wrap(err, "unable to generate singleton test template output")
		}
	}

	for _, table := range s.Tables {
		if table.IsJoinTable {
			continue
		}

		data := &templateData{
			Tables:           s.Tables,
			Table:            table,
			Schema:           s.Config.Schema,
			DriverName:       s.Config.DriverName,
			UseLastInsertID:  s.Driver.UseLastInsertID(),
			PkgName:          s.Config.PkgName,
			NoHooks:          s.Config.NoHooks,
			NoAutoTimestamps: s.Config.NoAutoTimestamps,
			Tags:             s.Config.Tags,

			StringFuncs: templateStringMappers,
		}

		// Generate the regular templates
		if err := generateOutput(s, data); err != nil {
			return errors.Wrap(err, "unable to generate output")
		}

		// Generate the test templates
		if !s.Config.NoTests && includeTests {
			if err := generateTestOutput(s, data); err != nil {
				return errors.Wrap(err, "unable to generate test output")
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

	basePath, err := getBasePath(s.Config.BaseDir)
	if err != nil {
		return err
	}

	s.Templates, err = loadTemplates(filepath.Join(basePath, templatesDirectory))
	if err != nil {
		return err
	}

	s.SingletonTemplates, err = loadTemplates(filepath.Join(basePath, templatesSingletonDirectory))
	if err != nil {
		return err
	}

	if !s.Config.NoTests {
		s.TestTemplates, err = loadTemplates(filepath.Join(basePath, templatesTestDirectory))
		if err != nil {
			return err
		}

		s.SingletonTestTemplates, err = loadTemplates(filepath.Join(basePath, templatesSingletonTestDirectory))
		if err != nil {
			return err
		}

		s.TestMainTemplate, err = loadTemplate(filepath.Join(basePath, templatesTestMainDirectory), s.Config.DriverName+"_main.tpl")
		if err != nil {
			return err
		}
	}

	return nil
}

var basePackage = "github.com/vattle/sqlboiler"

func getBasePath(baseDirConfig string) (string, error) {
	if len(baseDirConfig) > 0 {
		return baseDirConfig, nil
	}

	p, _ := build.Default.Import(basePackage, "", build.FindOnly)
	if p != nil && len(p.Dir) > 0 {
		return p.Dir, nil
	}

	return os.Getwd()
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
	}

	if s.Driver == nil {
		return errors.New("An invalid driver name was provided")
	}

	return nil
}

// initTables retrieves all "public" schema table names from the database.
func (s *State) initTables(schema string, whitelist, blacklist []string) error {
	var err error
	s.Tables, err = bdb.Tables(s.Driver, schema, whitelist, blacklist)
	if err != nil {
		return errors.Wrap(err, "unable to fetch table data")
	}

	if len(s.Tables) == 0 {
		return errors.New("no tables found in database")
	}

	if err := checkPKeys(s.Tables); err != nil {
		return err
	}

	return nil
}

// Tags must be in a format like: json, xml, etc.
var rgxValidTag = regexp.MustCompile(`[a-zA-Z_\.]+`)

// initTags removes duplicate tags and validates the format
// of all user tags are simple strings without quotes: [a-zA-Z_\.]+
func (s *State) initTags(tags []string) error {
	s.Config.Tags = removeDuplicates(s.Config.Tags)
	for _, v := range s.Config.Tags {
		if !rgxValidTag.MatchString(v) {
			return errors.New("Invalid tag format %q supplied, only specify name, eg: xml")
		}
	}

	return nil
}

// initOutFolder creates the folder that will hold the generated output.
func (s *State) initOutFolder() error {
	return os.MkdirAll(s.Config.OutFolder, os.ModePerm)
}

// checkPKeys ensures every table has a primary key column
func checkPKeys(tables []bdb.Table) error {
	var missingPkey []string
	for _, t := range tables {
		if t.PKey == nil {
			missingPkey = append(missingPkey, t.Name)
		}
	}

	if len(missingPkey) != 0 {
		return errors.Errorf("primary key missing in tables (%s)", strings.Join(missingPkey, ", "))
	}

	return nil
}
