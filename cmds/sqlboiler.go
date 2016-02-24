package cmds

import (
	"errors"
	"os"
	"strings"

	"github.com/pobri19/sqlboiler/dbdrivers"
	"github.com/spf13/cobra"
)

type CmdData struct {
	TablesInfo [][]dbdrivers.DBTable
	TableNames []string
	DBDriver   dbdrivers.DBDriver
	OutFile    *os.File
}

var cmdData *CmdData

func init() {
	SQLBoiler.PersistentFlags().StringP("driver", "d", "", "The name of the driver in your config.toml")
	SQLBoiler.PersistentFlags().StringP("table", "t", "", "A comma seperated list of table names")
	SQLBoiler.PersistentFlags().StringP("out", "o", "", "The name of the output file")
	SQLBoiler.PersistentPreRun = sqlBoilerPreRun
	SQLBoiler.PersistentPostRun = sqlBoilerPostRun
}

var SQLBoiler = &cobra.Command{
	Use:   "sqlboiler",
	Short: "SQL Boiler generates boilerplate structs and statements",
	Long: "SQL Boiler generates boilerplate structs and statements.\n" +
		`Complete documentation is available at http://github.com/pobri19/sqlboiler`,
}

func sqlBoilerPostRun(cmd *cobra.Command, args []string) {
	cmdData.OutFile.Close()
	cmdData.DBDriver.Close()
}

func sqlBoilerPreRun(cmd *cobra.Command, args []string) {
	var err error
	cmdData = &CmdData{}

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

	// Connect to the driver database
	if err = cmdData.DBDriver.Open(); err != nil {
		errorQuit(err)
	}

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
		cmdData.TableNames, err = cmdData.DBDriver.GetAllTableNames()
		if err != nil {
			errorQuit(err)
		}

		if len(cmdData.TableNames) == 0 {
			errorQuit(errors.New("No tables found in database, migrate some tables first"))
		}
	}

	// loop over table Names and build TablesInfo
	for i := 0; i < len(cmdData.TableNames); i++ {
		tInfo, err := cmdData.DBDriver.GetTableInfo(cmdData.TableNames[i])
		if err != nil {
			errorQuit(err)
		}

		cmdData.TablesInfo = append(cmdData.TablesInfo, tInfo)
	}

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
