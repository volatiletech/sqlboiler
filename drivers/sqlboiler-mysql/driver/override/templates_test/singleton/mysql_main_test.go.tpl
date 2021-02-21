var rgxMySQLkey = regexp.MustCompile(`(?m)((,\n)?\s+CONSTRAINT.*?FOREIGN KEY.*?\n)+`)

type mysqlTester struct {
	dbConn *sql.DB

	dbName  string
	host    string
	user    string
	pass    string
	sslmode string
	port    int

	optionFile string

	testDBName string
	skipSQLCmd bool
}

func init() {
	dbMain = &mysqlTester{}
}

func (m *mysqlTester) setup() error {
	var err error

	viper.SetDefault("mysql.sslmode", "true")
	viper.SetDefault("mysql.port", 3306)

	m.dbName = viper.GetString("mysql.dbname")
	m.host = viper.GetString("mysql.host")
	m.user = viper.GetString("mysql.user")
	m.pass = viper.GetString("mysql.pass")
	m.port = viper.GetInt("mysql.port")
	m.sslmode = viper.GetString("mysql.sslmode")
	m.testDBName = viper.GetString("mysql.testdbname")
	m.skipSQLCmd = viper.GetBool("mysql.skipsqlcmd")

	err = vala.BeginValidation().Validate(
		vala.StringNotEmpty(m.user, "mysql.user"),
		vala.StringNotEmpty(m.host, "mysql.host"),
		vala.Not(vala.Equals(m.port, 0, "mysql.port")),
		vala.StringNotEmpty(m.dbName, "mysql.dbname"),
		vala.StringNotEmpty(m.sslmode, "mysql.sslmode"),
	).Check()

	if err != nil {
		return err
	}

	// Create a randomized db name.
	if len(m.testDBName) == 0 {
		m.testDBName = randomize.StableDBName(m.dbName)
	}

	if err = m.makeOptionFile(); err != nil {
		return errors.Wrap(err, "couldn't make option file")
	}

	if !m.skipSQLCmd {
		if err = m.dropTestDB(); err != nil {
			return err
		}
		if err = m.createTestDB(); err != nil {
			return err
		}

		dumpCmd := exec.Command("mysqldump", m.defaultsFile(), "--no-data", m.dbName)
		createCmd := exec.Command("mysql", m.defaultsFile(), "--database", m.testDBName)

		r, w := io.Pipe()
		dumpCmdStderr := &bytes.Buffer{}
		createCmdStderr := &bytes.Buffer{}

		dumpCmd.Stdout = w
		dumpCmd.Stderr = dumpCmdStderr

		createCmd.Stdin = newFKeyDestroyer(rgxMySQLkey, r)
		createCmd.Stderr = createCmdStderr

		if err = dumpCmd.Start(); err != nil {
			return errors.Wrap(err, "failed to start mysqldump command")
		}
		if err = createCmd.Start(); err != nil {
			return errors.Wrap(err, "failed to start mysql command")
		}

		if err = dumpCmd.Wait(); err != nil {
			fmt.Println(err)
			fmt.Println(dumpCmdStderr.String())
			return errors.Wrap(err, "failed to wait for mysqldump command")
		}

		_ = w.Close() // After dumpCmd is done, close the write end of the pipe

		if err = createCmd.Wait(); err != nil {
			fmt.Println(err)
			fmt.Println(createCmdStderr.String())
			return errors.Wrap(err, "failed to wait for mysql command")
		}
	}

	return nil
}

func (m *mysqlTester) sslMode(mode string) string {
	switch mode {
	case "true":
		return "REQUIRED"
	case "false":
		return "DISABLED"
	case "PREFERRED":
		return ""
	default:
		return mode
	}
}

func (m *mysqlTester) defaultsFile() string {
	return fmt.Sprintf("--defaults-file=%s", m.optionFile)
}

func (m *mysqlTester) makeOptionFile() error {
	tmp, err := ioutil.TempFile("", "optionfile")
	if err != nil {
		return errors.Wrap(err, "failed to create option file")
	}

	isTCP := false
	_, err = os.Stat(m.host)
	if os.IsNotExist(err) {
		isTCP = true
	} else if err != nil {
		return errors.Wrap(err, "could not stat m.host")
	}

	fmt.Fprintln(tmp, "[client]")
	fmt.Fprintf(tmp, "host=%s\n", m.host)
	fmt.Fprintf(tmp, "port=%d\n", m.port)
	fmt.Fprintf(tmp, "user=%s\n", m.user)
	if len(m.pass) != 0 {
		fmt.Fprintf(tmp, "password=%s\n", m.pass)
	}
	if mode := m.sslMode(m.sslmode); mode != "" {
		fmt.Fprintf(tmp, "ssl-mode=%s\n", mode)
	}
	if isTCP {
		fmt.Fprintln(tmp, "protocol=tcp")
	}

	fmt.Fprintln(tmp, "[mysqldump]")
	fmt.Fprintf(tmp, "host=%s\n", m.host)
	fmt.Fprintf(tmp, "port=%d\n", m.port)
	fmt.Fprintf(tmp, "user=%s\n", m.user)
	if len(m.pass) != 0 {
		fmt.Fprintf(tmp, "password=%s\n", m.pass)
	}
	if mode := m.sslMode(m.sslmode); mode != "" {
		fmt.Fprintf(tmp, "ssl-mode=%s\n", mode)
	}
	if isTCP {
		fmt.Fprintln(tmp, "protocol=tcp")
	}

	m.optionFile = tmp.Name()

	return tmp.Close()
}

func (m *mysqlTester) createTestDB() error {
	sql := fmt.Sprintf("create database %s;", m.testDBName)
	return m.runCmd(sql, "mysql")
}

func (m *mysqlTester) dropTestDB() error {
	sql := fmt.Sprintf("drop database if exists %s;", m.testDBName)
	return m.runCmd(sql, "mysql")
}

func (m *mysqlTester) teardown() error {
	if m.dbConn != nil {
		return m.dbConn.Close()
	}

	if !m.skipSQLCmd {
		if err := m.dropTestDB(); err != nil {
			return err
		}
	}

	return os.Remove(m.optionFile)
}

func (m *mysqlTester) runCmd(stdin, command string, args ...string) error {
	args = append([]string{m.defaultsFile()}, args...)

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

func (m *mysqlTester) conn() (*sql.DB, error) {
	if m.dbConn != nil {
	return m.dbConn, nil
	}

	var err error
	m.dbConn, err = sql.Open("mysql", driver.MySQLBuildQueryString(m.user, m.pass, m.testDBName, m.host, m.port, m.sslmode))
	if err != nil {
	return nil, err
	}

	return m.dbConn, nil
}
