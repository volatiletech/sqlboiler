var rgxMSSQLkey = regexp.MustCompile(`(?m)^ALTER TABLE .*ADD\s+CONSTRAINT .* FOREIGN KEY.*?.*\n?REFERENCES.*`)

type mssqlTester struct {
	dbConn     *sql.DB
	dbName     string
	host       string
	user       string
	pass       string
	sslmode    string
	port       int
	testDBName string
	skipSQLCmd bool
}

func init() {
	dbMain = &mssqlTester{}
}

func (m *mssqlTester) setup() error {
	var err error

	viper.SetDefault("mssql.schema", "dbo")
	viper.SetDefault("mssql.sslmode", "true")
	viper.SetDefault("mssql.port", 1433)

	m.dbName = viper.GetString("mssql.dbname")
	m.host = viper.GetString("mssql.host")
	m.user = viper.GetString("mssql.user")
	m.pass = viper.GetString("mssql.pass")
	m.port = viper.GetInt("mssql.port")
	m.sslmode = viper.GetString("mssql.sslmode")
	m.testDBName = viper.GetString("mssql.testdbname")
	m.skipSQLCmd = viper.GetBool("mssql.skipsqlcmd")

	err = vala.BeginValidation().Validate(
		vala.StringNotEmpty(viper.GetString("mssql.user"), "mssql.user"),
		vala.StringNotEmpty(viper.GetString("mssql.host"), "mssql.host"),
		vala.Not(vala.Equals(viper.GetInt("mssql.port"), 0, "mssql.port")),
		vala.StringNotEmpty(viper.GetString("mssql.dbname"), "mssql.dbname"),
		vala.StringNotEmpty(viper.GetString("mssql.sslmode"), "mssql.sslmode"),
	).Check()

	if err != nil {
		return err
	}

	// Create a randomized db name.
	if len(m.testDBName) == 0 {
		m.testDBName = randomize.StableDBName(m.dbName)
	}

	if !m.skipSQLCmd {
		if err = m.dropTestDB(); err != nil {
			return err
		}
		if err = m.createTestDB(); err != nil {
			return err
		}

		createCmd := exec.Command("sqlcmd", "-S", m.host, "-U", m.user, "-P", m.pass, "-d", m.testDBName)

		f, err := os.Open("tables_schema.sql")
		if err != nil {
			return errors.Wrap(err, "failed to open tables_schema.sql file")
		}

		defer func() { _ = f.Close() }()

		stderr := &bytes.Buffer{}
		createCmd.Stdin = newFKeyDestroyer(rgxMSSQLkey, f)
		createCmd.Stderr = stderr

		if err = createCmd.Start(); err != nil {
			return errors.Wrap(err, "failed to start sqlcmd command")
		}

		if err = createCmd.Wait(); err != nil {
			fmt.Println(err)
			fmt.Println(stderr.String())
			return errors.Wrap(err, "failed to wait for sqlcmd command")
		}
	}

	return nil
}

func (m *mssqlTester) sslMode(mode string) string {
	switch mode {
	case "true":
		return "true"
	case "false":
		return "false"
	default:
		return "disable"
	}
}

func (m *mssqlTester) createTestDB() error {
	sql := fmt.Sprintf(`
	CREATE DATABASE %s;
	GO
	ALTER DATABASE %[1]s
	SET READ_COMMITTED_SNAPSHOT ON;
	GO`, m.testDBName)
	return m.runCmd(sql, "sqlcmd", "-S", m.host, "-U", m.user, "-P", m.pass)
}

func (m *mssqlTester) dropTestDB() error {
	// Since MS SQL 2016 it can be done with
	// DROP DATABASE [ IF EXISTS ] { database_name | database_snapshot_name } [ ,...n ] [;]
	sql := fmt.Sprintf(`
	IF EXISTS(SELECT name FROM sys.databases 
		WHERE name = '%s')
		DROP DATABASE %s
	GO`, m.testDBName, m.testDBName)
	return m.runCmd(sql, "sqlcmd", "-S", m.host, "-U", m.user, "-P", m.pass)
}

func (m *mssqlTester) teardown() error {
	if m.dbConn != nil {
		m.dbConn.Close()
	}

	if !m.skipSQLCmd {
		if err := m.dropTestDB(); err != nil {
			return err
		}
	}

	return nil
}

func (m *mssqlTester) runCmd(stdin, command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdin = strings.NewReader(stdin)

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("failed running:", command, args)
		fmt.Println(stdout.String())
		fmt.Println(stderr.String())
		return err
	}

	return nil
}

func (m *mssqlTester) conn() (*sql.DB, error) {
	if m.dbConn != nil {
		return m.dbConn, nil
	}

	var err error
	m.dbConn, err = sql.Open("mssql", driver.MSSQLBuildQueryString(m.user, m.pass, m.testDBName, m.host, m.port, m.sslmode))
	if err != nil {
		return nil, err
	}

	return m.dbConn, nil
}
