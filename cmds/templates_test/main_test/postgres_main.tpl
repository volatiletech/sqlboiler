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

var cfg *Config
var testCfg *Config

var dbConn *sql.DB

func TestMain(m *testing.M) {
	rand.Seed(time.Now().UnixNano())

  err := setup()
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

  code := m.Run()

	err = teardown()
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

  os.Exit(code)
}

// teardown switches its connection to the template1 database temporarily
// so that it can drop the test database and the test user.
func teardown() error {
	err := dbConn.Close()
	if err != nil {
		return err
	}

	dbConn, err = DBConnect(cfg.Postgres.User, cfg.Postgres.Pass, "template1", cfg.Postgres.Host, cfg.Postgres.Port)
	if err != nil {
		return err
	}

	_, err = dbConn.Exec(fmt.Sprintf(`DROP DATABASE %s;`, testCfg.Postgres.DBName))
	if err != nil {
		return err
	}

	_, err = dbConn.Exec(fmt.Sprintf(`DROP USER %s;`, testCfg.Postgres.User))
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

func randSeq(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyz")

  randStr := make([]rune, n)
  for i := range randStr {
    randStr[i] = letters[rand.Intn(len(letters))]
  }
  return string(randStr)
}

func LoadConfigFile(filename string) error {
	_, err := toml.DecodeFile(filename, &cfg)

	if os.IsNotExist(err) {
		return fmt.Errorf("Failed to find the toml configuration file %s: %s", filename, err)
	}

	if err != nil {
		return fmt.Errorf("Failed to decode toml configuration file: %s", err)
	}

	return nil
}

// setup dumps the database schema and imports it into a temporary randomly
// generated test database so that tests can be run against it using the
// generated sqlboiler ORM package.
func setup() error {
	// Load the config file in the parent directory.
  err := LoadConfigFile("../config.toml")
	if err != nil {
		return fmt.Errorf("Unable to load config file: %s", err)
	}

	// Create a randomized test configuration object.
	testCfg = &Config{}
	testCfg.Postgres.Host = cfg.Postgres.Host
	testCfg.Postgres.Port = cfg.Postgres.Port
	testCfg.Postgres.User = randSeq(20)
	testCfg.Postgres.Pass = randSeq(20)
	testCfg.Postgres.DBName = cfg.Postgres.DBName + "_" + randSeq(10)

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
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.DBName,
		cfg.Postgres.User,
		cfg.Postgres.Pass,
	))

	passFilePath := passDir + "/pwfile"

	err = ioutil.WriteFile(passFilePath, pwBytes, 0600)
	if err != nil {
		return fmt.Errorf("Unable to create pwfile in passDir: %s", err)
	}

	// The params for the pg_dump command to dump the database schema
	params := []string{
		fmt.Sprintf(`--host=%s`, cfg.Postgres.Host),
		fmt.Sprintf(`--port=%d`, cfg.Postgres.Port),
		fmt.Sprintf(`--username=%s`, cfg.Postgres.User),
		"--schema-only",
		cfg.Postgres.DBName,
	}

	// Dump the database schema into the sqlboilerschema tmp file
	errBuf := bytes.Buffer{}
	cmd := exec.Command("pg_dump", params...)
	cmd.Stderr = &errBuf
	cmd.Stdout = fhSchema
	cmd.Env = append(os.Environ(), fmt.Sprintf(`PGPASSFILE=%s`, passFilePath))

	if err := cmd.Run(); err != nil {
		fmt.Printf("pg_dump exec failed: %s\n\n%s\n", err, errBuf.String())
	}

	dbConn, err = DBConnect(cfg.Postgres.User, cfg.Postgres.Pass, cfg.Postgres.DBName, cfg.Postgres.Host, cfg.Postgres.Port)
	if err != nil {
		return err
	}

	// Create the randomly generated database test user
	if err = createTestUser(dbConn); err != nil {
		return err
	}

	// Create the randomly generated database
	_, err = dbConn.Exec(fmt.Sprintf(`CREATE DATABASE %s WITH ENCODING 'UTF8'`, testCfg.Postgres.DBName))
	if err != nil {
		return err
	}

	// Assign the randomly generated db test user to the generated test db
	_, err = dbConn.Exec(fmt.Sprintf(`ALTER DATABASE %s OWNER TO %s;`, testCfg.Postgres.DBName, testCfg.Postgres.User))
	if err != nil {
		return err
	}

	// Close the old connection so we can reconnect with the restricted access generated user
	if err = dbConn.Close(); err != nil {
		return err
	}

	// Connect to the generated test db with the restricted privilege generated user
	dbConn, err = DBConnect(testCfg.Postgres.User, testCfg.Postgres.Pass, testCfg.Postgres.DBName, testCfg.Postgres.Host, testCfg.Postgres.Port)
	if err != nil {
		return err
	}

	// Write the generated user password to a tmp file for pg_dump
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

// createTestUser creates a temporary database user with restricted privileges
func createTestUser(db *sql.DB) error {
	now := time.Now().Add(time.Hour * 24 * 2)
	valid := now.Format("2006-1-2")

	query := fmt.Sprintf(`CREATE USER %s WITH PASSWORD '%s' VALID UNTIL '%s';`,
		testCfg.Postgres.User,
		testCfg.Postgres.Pass,
		valid,
	)

	_, err := dbConn.Exec(query)
	return err
}
