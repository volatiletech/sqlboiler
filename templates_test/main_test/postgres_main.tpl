type pgTester struct {
	dbConn *sql.DB

	dbName  string
	host    string
	user    string
	pass    string
	sslmode string
	port   int

	testDBName string
}

func init() {
	dbMain = &pgTester{}
}

// disableTriggers is used to disable foreign key constraints for every table.
// If this is not used we cannot test inserts due to foreign key constraint errors.
func (p *pgTester) disableTriggers() error {
	var stmts []string

	{{range .Tables -}}
	stmts = append(stmts, `ALTER TABLE {{.Name}} DISABLE TRIGGER ALL;`)
	{{end -}}

	if len(stmts) == 0 {
		return nil
	}

	var err error
	for _, s := range stmts {
		_, err = p.dbConn.Exec(s)
		if err != nil {
			return err
		}
	}

	return nil
}

// teardown executes cleanup tasks when the tests finish running
func (p *pgTester) teardown() error {
	return p.dropTestDB()
}

func (p *pgTester) conn() (*sql.DB, error) {
	return p.dbConn, nil
}

// dropTestDB switches its connection to the template1 database temporarily
// so that it can drop the test database without causing "in use" conflicts.
// The template1 database should be present on all default postgres installations.
func (p *pgTester) dropTestDB() error {
	var err error
	if p.dbConn != nil {
		if err = p.dbConn.Close(); err != nil {
			return err
		}
	}

	p.dbConn, err = DBConnect(p.user, p.pass, "template1", p.host, p.port, p.sslmode)
	if err != nil {
		return err
	}

	_, err = p.dbConn.Exec(fmt.Sprintf(`DROP DATABASE IF EXISTS %s;`, p.testDBName))
	if err != nil {
		return err
	}

	return p.dbConn.Close()
}

// DBConnect connects to a database and returns the handle.
func DBConnect(user, pass, dbname, host string, port int, sslmode string) (*sql.DB, error) {
	connStr := drivers.PostgresBuildQueryString(user, pass, dbname, host, port, sslmode)

	return sql.Open("postgres", connStr)
}

// setup dumps the database schema and imports it into a temporary randomly
// generated test database so that tests can be run against it using the
// generated sqlboiler ORM package.
func (p *pgTester) setup() error {
	var err error

	p.dbName = viper.GetString("postgres.dbname")
	p.host = viper.GetString("postgres.host")
	p.user = viper.GetString("postgres.user")
	p.pass = viper.GetString("postgres.pass")
	p.port = viper.GetInt("postgres.port")
	p.sslmode = viper.GetString("postgres.sslmode")
	// Create a randomized db name.
	p.testDBName = randomize.StableDBName(p.dbName)

	err = p.dropTestDB()
	if err != nil {
		fmt.Printf("%#v\n", err)
		return err
	}

	fhSchema, err := ioutil.TempFile(os.TempDir(), "sqlboilerschema")
	if err != nil {
		return errors.Wrap(err, "Unable to create sqlboiler schema tmp file")
	}
	defer os.Remove(fhSchema.Name())

	passDir, err := ioutil.TempDir(os.TempDir(), "sqlboiler")
	if err != nil {
		return errors.Wrap(err, "Unable to create sqlboiler tmp dir for postgres pw file")
	}
	defer os.RemoveAll(passDir)

	// Write the postgres user password to a tmp file for pg_dump
	pwBytes := []byte(fmt.Sprintf("%s:%d:%s:%s", p.host, p.port, p.dbName, p.user))

	if len(p.pass) > 0 {
		pwBytes = []byte(fmt.Sprintf("%s:%s", pwBytes, p.pass))
	}

	passFilePath := filepath.Join(passDir, "pwfile")

	err = ioutil.WriteFile(passFilePath, pwBytes, 0600)
	if err != nil {
		return errors.Wrap(err, "Unable to create pwfile in passDir")
	}

	// The params for the pg_dump command to dump the database schema
	params := []string{
		fmt.Sprintf(`--host=%s`, p.host),
		fmt.Sprintf(`--port=%d`, p.port),
		fmt.Sprintf(`--username=%s`, p.user),
		"--schema-only",
		p.dbName,
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

	p.dbConn, err = DBConnect(p.user, p.pass, p.dbName, p.host, p.port, p.sslmode)
	if err != nil {
		return err
	}

	// Create the randomly generated database
	_, err = p.dbConn.Exec(fmt.Sprintf(`CREATE DATABASE %s WITH ENCODING 'UTF8'`, p.testDBName))
	if err != nil {
		return err
	}

	// Close the old connection so we can reconnect to the test database
	if err = p.dbConn.Close(); err != nil {
		return err
	}

	// Connect to the generated test db
	p.dbConn, err = DBConnect(p.user, p.pass, p.testDBName, p.host, p.port, p.sslmode)
	if err != nil {
		return err
	}

	// Write the test config credentials to a tmp file for pg_dump
	testPwBytes := []byte(fmt.Sprintf("%s:%d:%s:%s", p.host, p.port, p.testDBName, p.user))

	if len(p.pass) > 0 {
		testPwBytes = []byte(fmt.Sprintf("%s:%s", testPwBytes, p.pass))
	}

	testPassFilePath := passDir + "/testpwfile"

	err = ioutil.WriteFile(testPassFilePath, testPwBytes, 0600)
	if err != nil {
		return errors.Wrapf(err, "Unable to create testpwfile in passDir")
	}

	// The params for the psql schema import command
	params = []string{
		fmt.Sprintf(`--dbname=%s`, p.testDBName),
		fmt.Sprintf(`--host=%s`, p.host),
		fmt.Sprintf(`--port=%d`, p.port),
		fmt.Sprintf(`--username=%s`, p.user),
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

	return p.disableTriggers()
}
