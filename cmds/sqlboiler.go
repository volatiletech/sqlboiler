package cmds

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/pobri19/sqlboiler/dbdrivers"
	"github.com/spf13/cobra"
)

const (
	templatesDirectory     = "/cmds/templates"
	templatesTestDirectory = "/cmds/templates_test"
)

// LoadTemplates loads all template folders into the cmdData object.
func (cmdData *CmdData) LoadTemplates() error {
	var err error
	cmdData.Templates, err = loadTemplates(templatesDirectory)
	if err != nil {
		return err
	}

	cmdData.TestTemplates, err = loadTemplates(templatesTestDirectory)
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

	return tpl.Templates(), err
}

// SQLBoilerPostRun cleans up the output file and db connection once all cmds are finished.
func (cmdData *CmdData) SQLBoilerPostRun(cmd *cobra.Command, args []string) error {
	cmdData.Interface.Close()
	return nil
}

// SQLBoilerPreRun performs the initialization tasks before the root command is run
func (cmdData *CmdData) SQLBoilerPreRun(cmd *cobra.Command, args []string) error {
	var err error

	// Initialize package name
	cmdData.PkgName = cmd.PersistentFlags().Lookup("pkgname").Value.String()

	err = initInterface(cmd, cmdData.Config, cmdData)
	if err != nil {
		return err
	}

	// Connect to the driver database
	if err = cmdData.Interface.Open(); err != nil {
		return fmt.Errorf("Unable to connect to the database: %s", err)
	}

	err = initTables(cmd, cmdData)
	if err != nil {
		return fmt.Errorf("Unable to initialize tables: %s", err)
	}

	err = initOutFolder(cmd, cmdData)
	if err != nil {
		return fmt.Errorf("Unable to initialize the output folder: %s", err)
	}

	return nil
}

// SQLBoilerRun executes every sqlboiler template and outputs them to files.
func (cmdData *CmdData) SQLBoilerRun(cmd *cobra.Command, args []string) error {
	for _, table := range cmdData.Tables {
		data := &tplData{
			Table:   table,
			PkgName: cmdData.PkgName,
		}

		var out [][]byte
		var imps imports

		imps.standard = sqlBoilerImports.standard
		imps.thirdparty = sqlBoilerImports.thirdparty

		// Loop through and generate every command template (excluding skipTemplates)
		for _, template := range cmdData.Templates {
			imps = combineTypeImports(imps, sqlBoilerTypeImports, data.Table.Columns)
			resp, err := generateTemplate(template, data)
			if err != nil {
				return err
			}
			out = append(out, resp)
		}

		err := outHandler(cmdData, out, data, imps, false)
		if err != nil {
			return err
		}

		// Generate the test templates for all commands
		if len(cmdData.TestTemplates) != 0 {
			var testOut [][]byte
			var testImps imports

			testImps.standard = sqlBoilerTestImports.standard
			testImps.thirdparty = sqlBoilerTestImports.thirdparty

			testImps = combineImports(testImps, sqlBoilerDriverTestImports[cmdData.DriverName])

			// Loop through and generate every command test template (excluding skipTemplates)
			for _, template := range cmdData.TestTemplates {
				resp, err := generateTemplate(template, data)
				if err != nil {
					return err
				}
				testOut = append(testOut, resp)
			}

			err = outHandler(cmdData, testOut, data, testImps, true)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// initInterface attempts to set the cmdData Interface based off the passed in
// driver flag value. If an invalid flag string is provided an error is returned.
func initInterface(cmd *cobra.Command, cfg *Config, cmdData *CmdData) error {
	// Retrieve driver flag
	driverName := cmd.PersistentFlags().Lookup("driver").Value.String()
	if driverName == "" {
		return errors.New("Must supply a driver flag.")
	}

	// Create a driver based off driver flag
	switch driverName {
	case "postgres":
		cmdData.Interface = dbdrivers.NewPostgresDriver(
			cfg.Postgres.User,
			cfg.Postgres.Pass,
			cfg.Postgres.DBName,
			cfg.Postgres.Host,
			cfg.Postgres.Port,
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
func initTables(cmd *cobra.Command, cmdData *CmdData) error {
	var tableNames []string
	tn := cmd.PersistentFlags().Lookup("table").Value.String()

	if len(tn) != 0 {
		tableNames = strings.Split(tn, ",")
		for i, name := range tableNames {
			tableNames[i] = strings.TrimSpace(name)
		}
	}

	var err error
	cmdData.Tables, err = cmdData.Interface.Tables(tableNames...)
	if err != nil {
		return fmt.Errorf("Unable to get all table names: %s", err)
	}

	if len(cmdData.Tables) == 0 {
		return errors.New("No tables found in database, migrate some tables first")
	}

	return nil
}

// initOutFolder creates the folder that will hold the generated output.
func initOutFolder(cmd *cobra.Command, cmdData *CmdData) error {
	cmdData.OutFolder = cmd.PersistentFlags().Lookup("folder").Value.String()
	if cmdData.OutFolder == "" {
		return fmt.Errorf("No output folder specified.")
	}

	if err := os.MkdirAll(cmdData.OutFolder, os.ModePerm); err != nil {
		return fmt.Errorf("Unable to make output folder: %s", err)
	}

	return nil
}
