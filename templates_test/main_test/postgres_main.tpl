type PostgresCfg struct {
	User   string `toml:"user"`
	Pass   string `toml:"pass"`
	Host   string `toml:"host"`
	Port   int    `toml:"port"`
	DBName string `toml:"dbname"`
}

type Config struct {
	Postgres PostgresCfg `toml:"postgres"`
}

func TestMain(m *testing.M) {
	rand.Seed(time.Now().UnixNano())

	// Set DebugMode so we can see generated sql statements
	boil.DebugMode = false

	var err error
	if err = setup(); err != nil {
		fmt.Println("Unable to execute setup:", err)
		os.Exit(-2)
	}

	var code int
	if err = disableTriggers(); err != nil {
		fmt.Println("Unable to disable triggers:", err)
	} else {
		boil.SetDB(dbConn)
	  code = m.Run()
	}

	if err = teardown(); err != nil {
		fmt.Println("Unable to execute teardown:", err)
		os.Exit(-3)
	}

  os.Exit(code)
}

// disableTriggers is used to disable foreign key constraints for every table.
// If this is not used we cannot test inserts due to foreign key constraint errors.
func disableTriggers() error {
	var stmts []string

	{{range .Tables}}
	stmts = append(stmts, `ALTER TABLE {{.Name}} DISABLE TRIGGER ALL;`)
	{{- end}}

	if len(stmts) == 0 {
		return nil
	}

	var err error
	for _, s := range stmts {
		_, err = dbConn.Exec(s)
		if err != nil {
			return err
		}
	}

	return nil
}

// teardown executes cleanup tasks when the tests finish running
func teardown() error {
	err := dropTestDB()
	return err
}

// dropTestDB switches its connection to the template1 database temporarily
// so that it can drop the test database without causing "in use" conflicts.
// The template1 database should be present on all default postgres installations.
func dropTestDB() error {
	var err error
	if dbConn != nil {
		if err = dbConn.Close(); err != nil {
			return err
		}
	}

	dbConn, err = DBConnect(testCfg.Postgres.User, testCfg.Postgres.Pass, "template1", testCfg.Postgres.Host, testCfg.Postgres.Port)
	if err != nil {
		return err
	}

	_, err = dbConn.Exec(fmt.Sprintf(`DROP DATABASE IF EXISTS %s;`, testCfg.Postgres.DBName))
	if err != nil {
		return err
	}

	return dbConn.Close()
}

// DBConnect connects to a database and returns the handle.
func DBConnect(user, pass, dbname, host string, port int) (*sql.DB, error) {
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d",
		user, pass, dbname, host, port)

		return sql.Open("postgres", connStr)
}

// setup dumps the database schema and imports it into a temporary randomly
// generated test database so that tests can be run against it using the
// generated sqlboiler ORM package.
func setup() error {
	var err error

	// Initialize Viper and load the config file
	err = InitViper()
	if err != nil {
		return fmt.Errorf("Unable to load config file: %s", err)
	}

	// Create a randomized test configuration object.
	testCfg.Postgres.Host = viper.GetString("postgres.host")
	testCfg.Postgres.Port = viper.GetInt("postgres.port")
	testCfg.Postgres.User = viper.GetString("postgres.user")
	testCfg.Postgres.Pass = viper.GetString("postgres.pass")
	testCfg.Postgres.DBName = getDBNameHash(viper.GetString("postgres.dbname"))

	err = vala.BeginValidation().Validate(
		vala.StringNotEmpty(testCfg.Postgres.User, "postgres.user"),
		vala.StringNotEmpty(testCfg.Postgres.Pass, "postgres.pass"),
		vala.StringNotEmpty(testCfg.Postgres.Host, "postgres.host"),
		vala.Not(vala.Equals(testCfg.Postgres.Port, 0, "postgres.port")),
		vala.StringNotEmpty(testCfg.Postgres.DBName, "postgres.dbname"),
	).Check()

	if err != nil {
		return fmt.Errorf("Unable to load testCfg: %s", err.Error())
	}

	err = dropTestDB()
	if err != nil {
		fmt.Printf("%#v\n", err)
		return err
	}

	fhSchema, err := ioutil.TempFile(os.TempDir(), "sqlboilerschema")
	if err != nil {
		return fmt.Errorf("Unable to create sqlboiler schema tmp file: %s", err)
	}
	defer os.Remove(fhSchema.Name())

	passDir, err := ioutil.TempDir(os.TempDir(), "sqlboiler")
	if err != nil {
		return fmt.Errorf("Unable to create sqlboiler tmp dir for postgres pw file: %s", err)
	}
	defer os.RemoveAll(passDir)

	// Write the postgres user password to a tmp file for pg_dump
	pwBytes := []byte(fmt.Sprintf("%s:%d:%s:%s:%s",
		viper.GetString("postgres.host"),
		viper.GetInt("postgres.port"),
		viper.GetString("postgres.dbname"),
		viper.GetString("postgres.user"),
		viper.GetString("postgres.pass"),
	))

	passFilePath := passDir + "/pwfile"

	err = ioutil.WriteFile(passFilePath, pwBytes, 0600)
	if err != nil {
		return fmt.Errorf("Unable to create pwfile in passDir: %s", err)
	}

	// The params for the pg_dump command to dump the database schema
	params := []string{
		fmt.Sprintf(`--host=%s`, viper.GetString("postgres.host")),
		fmt.Sprintf(`--port=%d`, viper.GetInt("postgres.port")),
		fmt.Sprintf(`--username=%s`, viper.GetString("postgres.user")),
		"--schema-only",
		viper.GetString("postgres.dbname"),
	}

	// Dump the database schema into the sqlboilerschema tmp file
	errBuf := bytes.Buffer{}
	cmd := exec.Command("pg_dump", params...)
	cmd.Stderr = &errBuf
	cmd.Stdout = fhSchema
	cmd.Env = append(os.Environ(), fmt.Sprintf(`PGPASSFILE=%s`, passFilePath))

	if err := cmd.Run(); err != nil {
		fmt.Printf("pg_dump exec failed: %s\n\n%s\n", err, errBuf.String())
		return err
	}

	dbConn, err = DBConnect(
		viper.GetString("postgres.user"),
		viper.GetString("postgres.pass"),
		viper.GetString("postgres.dbname"),
		viper.GetString("postgres.host"),
		viper.GetInt("postgres.port"),
	)
	if err != nil {
		return err
	}

	// Create the randomly generated database
	_, err = dbConn.Exec(fmt.Sprintf(`CREATE DATABASE %s WITH ENCODING 'UTF8'`, testCfg.Postgres.DBName))
	if err != nil {
		return err
	}

	// Close the old connection so we can reconnect to the test database
	if err = dbConn.Close(); err != nil {
		return err
	}

	// Connect to the generated test db
	dbConn, err = DBConnect(testCfg.Postgres.User, testCfg.Postgres.Pass, testCfg.Postgres.DBName, testCfg.Postgres.Host, testCfg.Postgres.Port)
	if err != nil {
		return err
	}

	// Write the test config credentials to a tmp file for pg_dump
	testPwBytes := []byte(fmt.Sprintf("%s:%d:%s:%s:%s",
		testCfg.Postgres.Host,
		testCfg.Postgres.Port,
		testCfg.Postgres.DBName,
		testCfg.Postgres.User,
		testCfg.Postgres.Pass,
	))

	testPassFilePath := passDir + "/testpwfile"

	err = ioutil.WriteFile(testPassFilePath, testPwBytes, 0600)
	if err != nil {
		return fmt.Errorf("Unable to create testpwfile in passDir: %s", err)
	}

	// The params for the psql schema import command
	params = []string{
		fmt.Sprintf(`--dbname=%s`, testCfg.Postgres.DBName),
		fmt.Sprintf(`--host=%s`, testCfg.Postgres.Host),
		fmt.Sprintf(`--port=%d`, testCfg.Postgres.Port),
		fmt.Sprintf(`--username=%s`, testCfg.Postgres.User),
		fmt.Sprintf(`--file=%s`, fhSchema.Name()),
	}

	// Import the database schema into the generated database.
	// It is now ready to be used by the generated ORM package for testing.
	outBuf := bytes.Buffer{}
	cmd = exec.Command("psql", params...)
	cmd.Stderr = &errBuf
	cmd.Stdout = &outBuf
	cmd.Env = append(os.Environ(), fmt.Sprintf(`PGPASSFILE=%s`, testPassFilePath))

	if err = cmd.Run(); err != nil {
		fmt.Printf("psql schema import exec failed: %s\n\n%s\n", err, errBuf.String())
	}

	return nil
}
