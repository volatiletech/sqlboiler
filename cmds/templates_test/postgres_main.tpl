type PostgresCfg struct {
	User   string `toml:"user"`
	Pass   string `toml:"pass"`
	Host   string `toml:"host"`
	Port   int    `toml:"port"`
	DBName string `toml:"dbname"`
}

type Config struct {
	Postgres PostgresCfg `toml:"postgres"`
	TestPostgres *PostgresCfg `toml:"postgres_test"`
}

var cfg *Config

var dbConn *sql.DB

func DBConnect(user, pass, dbname, host string, port int) error {
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d",
		user, pass, dbname, host, port)

		var err error
		dbConn, err = sql.Open("postgres", connStr)
		if err != nil {
			return err
		}

		return nil
}

func LoadConfigFile(filename string) error {
	_, err := toml.DecodeFile(filename, &cfg)

	if os.IsNotExist(err) {
		return fmt.Errorf("Failed to find the toml configuration file %s: %s", filename, err)
	}

	if err != nil {
		return fmt.Errorf("Failed to decode toml configuration file:", err)
	}

	if cfg.TestPostgres != nil {
		if cfg.TestPostgres.User == "" || cfg.TestPostgres.Pass == "" ||
			cfg.TestPostgres.Host == "" || cfg.TestPostgres.Port == 0 ||
			cfg.TestPostgres.DBName == "" || cfg.Postgres.DBName == cfg.TestPostgres.DBName {
			cfg.TestPostgres = nil
		}
	}

  if cfg.TestPostgres == nil {
    return errors.New("Failed to load config.toml postgres_test config")
  }

	return nil
}

func TestMain(m *testing.M) {
  err := setup()
	if err != nil {
		os.Exit(-1)
	}
  code := m.Run()
  // shutdown
  os.Exit(code)
}

func setup() error {
  err := LoadConfigFile("../config.toml")
	if err != nil {
		return fmt.Errorf("Unable to load config file: %s", err)
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
	pwBytes := []byte(fmt.Sprintf("*:*:*:*:%s", cfg.Postgres.Pass))
	passFilePath := passDir + "/pwfile"

	err = ioutil.WriteFile(passFilePath, pwBytes, 0600)
	if err != nil {
		return fmt.Errorf("Unable to create pwfile in passDir: %s", err)
	}

	params := []string{
		fmt.Sprintf(`--host=%s`, cfg.Postgres.Host),
		fmt.Sprintf(`--port=%d`, cfg.Postgres.Port),
		fmt.Sprintf(`--username=%s`, cfg.Postgres.User),
		"--schema-only",
		cfg.Postgres.DBName,
	}

	errBuf := bytes.Buffer{}
	cmd := exec.Command("pg_dump", params...)
	cmd.Stderr = &errBuf
	cmd.Stdout = fhSchema
	cmd.Env = append(os.Environ(), fmt.Sprintf(`PGPASSFILE=%s`, passFilePath))

	if err := cmd.Run(); err != nil {
		fmt.Printf("pg_dump exec failed: %s\n\n%s\n", err, errBuf.String())
	}

	err = DBConnect(cfg.Postgres.User, cfg.Postgres.Pass, cfg.Postgres.DBName, cfg.Postgres.Host, cfg.Postgres.Port)
	_, err = dbConn.Exec(fmt.Sprintf(`CREATE DATABASE %s WITH ENCODING 'UTF8'`, cfg.TestPostgres.DBName))
	return nil
}
