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

// CmdData holds the table schema a slice of (column name, column type) slices.
// It also holds a slice of all of the table names sqlboiler is generating against,
// the database driver chosen by the driver flag at runtime, and a pointer to the
// output file, if one is specified with a flag.
type CmdData struct {
	TablesInfo [][]dbdrivers.DBTable
	TableNames []string
	DBDriver   dbdrivers.DBDriver
	OutFile    *os.File
}

// cmdData is used globally by all commands to access the table schema data,
// the database driver and the output file. cmdData is initialized by
// the root SQLBoiler cobra command at run time, before other commands execute.
var cmdData *CmdData

// templates holds a slice of pointers to all templates in the templates directory.
var templates []*template.Template

// init initializes the sqlboiler flags, such as driver, table, and output file.
// It also sets the global preRun hook and postRun hook. Every command will execute
// these hooks before and after running to initialize the shared state.
func init() {
	SQLBoiler.PersistentFlags().StringP("driver", "d", "", "The name of the driver in your config.toml (mandatory)")
	SQLBoiler.PersistentFlags().StringP("table", "t", "", "A comma seperated list of table names")
	SQLBoiler.PersistentFlags().StringP("out", "o", "", "The name of the output file")
	SQLBoiler.PersistentPreRun = sqlBoilerPreRun
	SQLBoiler.PersistentPostRun = sqlBoilerPostRun
}

// SQLBoiler is the root app command
var SQLBoiler = &cobra.Command{
	Use:   "sqlboiler",
	Short: "SQL Boiler generates boilerplate structs and statements",
	Long: "SQL Boiler generates boilerplate structs and statements.\n" +
		`Complete documentation is available at http://github.com/pobri19/sqlboiler`,
}

// sqlBoilerPostRun cleans up the output file and database connection once
// all commands are finished running.
func sqlBoilerPostRun(cmd *cobra.Command, args []string) {
	cmdData.OutFile.Close()
	cmdData.DBDriver.Close()
}

// sqlBoilerPreRun executes before all commands start running. Its job is to
// initialize the shared state object (cmdData). Initialization tasks include
// assigning the database driver based off the driver flag and opening the database connection,
// retrieving a list of all the tables in the database (if specific tables are not provided
// via a flag), building the table schema for use in the templates, and opening the output file
// handle if one is specified with a flag.
func sqlBoilerPreRun(cmd *cobra.Command, args []string) {
	var err error
	cmdData = &CmdData{}

	// Initialize the cmdData.DBDriver
	initDBDriver()

	// Connect to the driver database
	if err = cmdData.DBDriver.Open(); err != nil {
		errorQuit(err)
	}

	// Initialize the cmdData.TableNames
	initTableNames()

	// Initialize the cmdData.TablesInfo
	initTablesInfo()

	// Initialize the cmdData.OutFile
	initOutFile()

	// Initialize the templates
	templates, err = initTemplates()
	if err != nil {
		errorQuit(fmt.Errorf("Unable to initialize templates: %s", err))
	}
}

// initDBDriver attempts to set the cmdData DBDriver based off the passed in
// driver flag value. If an invalid flag string is provided the program will exit.
func initDBDriver() {
	// Retrieve driver flag
	driverName := SQLBoiler.PersistentFlags().Lookup("driver").Value.String()
	if driverName == "" {
		errorQuit(errors.New("Must supply a driver flag."))
	}

	// Create a driver based off driver flag
	switch driverName {
	case "postgres":
		cmdData.DBDriver = dbdrivers.NewPostgresDriver(
			cfg.Postgres.User,
			cfg.Postgres.Pass,
			cfg.Postgres.DBName,
			cfg.Postgres.Host,
			cfg.Postgres.Port,
		)
	}

	if cmdData.DBDriver == nil {
		errorQuit(errors.New("An invalid driver name was provided"))
	}
}

// initTableNames will create a string slice out of the passed in table flag value
// if one is provided. If no flag is provided, it will attempt to connect to the
// database to retrieve all "public" table names, and build a slice out of that result.
func initTableNames() {
	// Retrieve the list of tables
	tn := SQLBoiler.PersistentFlags().Lookup("table").Value.String()

	if len(tn) != 0 {
		cmdData.TableNames = strings.Split(tn, ",")
		for i, name := range cmdData.TableNames {
			cmdData.TableNames[i] = strings.TrimSpace(name)
		}
	}

	// If no table names are provided attempt to process all tables in database
	if len(cmdData.TableNames) == 0 {
		// get all table names
		var err error
		cmdData.TableNames, err = cmdData.DBDriver.GetAllTableNames()
		if err != nil {
			errorQuit(err)
		}

		if len(cmdData.TableNames) == 0 {
			errorQuit(errors.New("No tables found in database, migrate some tables first"))
		}
	}
}

// initTablesInfo builds a description of each table (column name, column type)
// and assigns it to cmdData.TablesInfo, the slice of dbdrivers.DBTable slices.
func initTablesInfo() {
	// loop over table Names and build TablesInfo
	for i := 0; i < len(cmdData.TableNames); i++ {
		tInfo, err := cmdData.DBDriver.GetTableInfo(cmdData.TableNames[i])
		if err != nil {
			errorQuit(err)
		}

		cmdData.TablesInfo = append(cmdData.TablesInfo, tInfo)
	}
}

// initOutFile opens a file handle to the file name specified by the out flag.
// If no file name is provided, out will remain nil and future output will be
// piped to Stdout instead of to a file.
func initOutFile() {
	// open the out file filehandle
	outf := SQLBoiler.PersistentFlags().Lookup("out").Value.String()
	if outf != "" {
		var err error
		cmdData.OutFile, err = os.Create(outf)
		if err != nil {
			errorQuit(err)
		}
	}
}

func initTemplates() ([]*template.Template, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	pattern := filepath.Join(wd, "templates", "*.tpl")

	tpl, err := template.New("").Funcs(template.FuncMap{
		"makeGoColName":          makeGoColName,
		"makeGoVarName":          makeGoVarName,
		"makeDBColName":          makeDBColName,
		"makeSelectParamNames":   makeSelectParamNames,
		"makeGoInsertParamNames": makeGoInsertParamNames,
		"makeGoInsertParamFlags": makeGoInsertParamFlags,
	}).ParseGlob(pattern)

	if err != nil {
		return nil, err
	}

	return tpl.Templates(), err
}
