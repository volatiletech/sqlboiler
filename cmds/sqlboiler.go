package cmds

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"github.com/pobri19/sqlboiler/dbdrivers"
	"github.com/spf13/cobra"
)

const (
	templatesDirectory         = "/cmds/templates"
	templatesTestDirectory     = "/cmds/templates_test"
	templatesTestMainDirectory = "/cmds/templates_test/main_test"
)

// LoadTemplates loads all template folders into the cmdData object.
func initTemplates(cmdData *CmdData) error {
	var err error
	cmdData.Templates, err = loadTemplates(templatesDirectory)
	if err != nil {
		return err
	}

	cmdData.TestTemplates, err = loadTemplates(templatesTestDirectory)
	if err != nil {
		return err
	}

	filename := cmdData.DriverName + "_main.tpl"
	cmdData.TestMainTemplate, err = loadTemplate(templatesTestMainDirectory, filename)
	if err != nil {
		return err
	}

	return nil
}

// loadTemplates loads all of the template files in the specified directory.
func loadTemplates(dir string) ([]*template.Template, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	pattern := filepath.Join(wd, dir, "*.tpl")
	tpl, err := template.New("").Funcs(sqlBoilerTemplateFuncs).ParseGlob(pattern)

	if err != nil {
		return nil, err
	}

	templates := templater(tpl.Templates())
	sort.Sort(templates)

	return templates, err
}

// loadTemplate loads a single template file.
func loadTemplate(dir string, filename string) (*template.Template, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	pattern := filepath.Join(wd, dir, filename)
	tpl, err := template.New("").Funcs(sqlBoilerTemplateFuncs).ParseFiles(pattern)

	if err != nil {
		return nil, err
	}

	return tpl.Lookup(filename), err
}

// SQLBoilerPostRun cleans up the output file and db connection once all cmds are finished.
func (c *CmdData) SQLBoilerPostRun(cmd *cobra.Command, args []string) error {
	c.Interface.Close()
	return nil
}

// SQLBoilerPreRun performs the initialization tasks before the root command is run
func (c *CmdData) SQLBoilerPreRun(cmd *cobra.Command, args []string) error {
	// Initialize package name
	pkgName := cmd.PersistentFlags().Lookup("pkgname").Value.String()

	// Retrieve driver flag
	driverName := cmd.PersistentFlags().Lookup("driver").Value.String()
	if driverName == "" {
		return errors.New("Must supply a driver flag.")
	}

	tableName := cmd.PersistentFlags().Lookup("table").Value.String()

	outFolder := cmd.PersistentFlags().Lookup("folder").Value.String()
	if outFolder == "" {
		return fmt.Errorf("No output folder specified.")
	}

	return c.initCmdData(pkgName, driverName, tableName, outFolder)
}

// SQLBoilerRun is a proxy method for the run function
func (c *CmdData) SQLBoilerRun(cmd *cobra.Command, args []string) error {
	return c.run(true)
}

// run executes the sqlboiler templates and outputs them to files.
func (c *CmdData) run(includeTests bool) error {
	if includeTests {
		if err := generateTestMainOutput(c); err != nil {
			return fmt.Errorf("Unable to generate TestMain output: %s", err)
		}
	}

	for _, table := range c.Tables {
		data := &tplData{
			Table:   table,
			PkgName: c.PkgName,
		}

		// Generate the regular templates
		if err := generateOutput(c, data); err != nil {
			return fmt.Errorf("Unable to generate test output: %s", err)
		}

		// Generate the test templates
		if includeTests {
			if err := generateTestOutput(c, data); err != nil {
				return fmt.Errorf("Unable to generate output: %s", err)
			}
		}
	}

	return nil
}

func (c *CmdData) initCmdData(pkgName, driverName, tableName, outFolder string) error {
	c.OutFolder = outFolder
	c.PkgName = pkgName

	err := initInterface(driverName, c)
	if err != nil {
		return err
	}

	// Connect to the driver database
	if err = c.Interface.Open(); err != nil {
		return fmt.Errorf("Unable to connect to the database: %s", err)
	}

	err = initTables(tableName, c)
	if err != nil {
		return fmt.Errorf("Unable to initialize tables: %s", err)
	}

	err = initOutFolder(c)
	if err != nil {
		return fmt.Errorf("Unable to initialize the output folder: %s", err)
	}

	err = initTemplates(c)
	if err != nil {
		return fmt.Errorf("Unable to initialize templates: %s", err)
	}

	return nil
}

// initInterface attempts to set the cmdData Interface based off the passed in
// driver flag value. If an invalid flag string is provided an error is returned.
func initInterface(driverName string, cmdData *CmdData) error {
	// Create a driver based off driver flag
	switch driverName {
	case "postgres":
		cmdData.Interface = dbdrivers.NewPostgresDriver(
			cmdData.Config.Postgres.User,
			cmdData.Config.Postgres.Pass,
			cmdData.Config.Postgres.DBName,
			cmdData.Config.Postgres.Host,
			cmdData.Config.Postgres.Port,
		)
	}

	if cmdData.Interface == nil {
		return errors.New("An invalid driver name was provided")
	}

	cmdData.DriverName = driverName
	return nil
}

// initTables will create a string slice out of the passed in table flag value
// if one is provided. If no flag is provided, it will attempt to connect to the
// database to retrieve all "public" table names, and build a slice out of that result.
func initTables(tableName string, cmdData *CmdData) error {
	var tableNames []string

	if len(tableName) != 0 {
		tableNames = strings.Split(tableName, ",")
		for i, name := range tableNames {
			tableNames[i] = strings.TrimSpace(name)
		}
	}

	var err error
	cmdData.Tables, err = dbdrivers.Tables(cmdData.Interface, tableNames...)
	if err != nil {
		return fmt.Errorf("Unable to get all table names: %s", err)
	}

	if len(cmdData.Tables) == 0 {
		return errors.New("No tables found in database, migrate some tables first")
	}

	return nil
}

// initOutFolder creates the folder that will hold the generated output.
func initOutFolder(cmdData *CmdData) error {
	if err := os.MkdirAll(cmdData.OutFolder, os.ModePerm); err != nil {
		return fmt.Errorf("Unable to make output folder: %s", err)
	}

	return nil
}
