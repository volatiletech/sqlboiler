type mysqlTester struct {
	dbConn *sql.DB

	dbName	string
	host	string
	user	string
	pass	string
	sslmode	string
	port	int

	optionFile string

	testDBName string
}

func init() {
	dbMain = &mysqlTester{}
}

func (m *mysqlTester) setup() error {
	var err error

	m.dbName = viper.GetString("mysql.dbname")
	m.host = viper.GetString("mysql.host")
	m.user = viper.GetString("mysql.user")
	m.pass = viper.GetString("mysql.pass")
	m.port = viper.GetInt("mysql.port")
	m.sslmode = viper.GetString("mysql.sslmode")
	// Create a randomized db name.
	m.testDBName = randomize.StableDBName(m.dbName)

	if err = m.makeOptionFile(); err != nil {
		return errors.Prefix("couldn't make option file", err)
	}

	if err = m.dropTestDB(); err != nil {
		return errors.Err(err)
	}
	if err = m.createTestDB(); err != nil {
		return errors.Err(err)
	}

	dumpCmd := exec.Command("mysqldump", m.defaultsFile(), "--no-data", m.dbName)
	createCmd := exec.Command("mysql", m.defaultsFile(), "--database", m.testDBName)

	r, w := io.Pipe()
	dumpCmd.Stdout = w
	createCmd.Stdin = newFKeyDestroyer(rgxMySQLkey, r)

	if err = dumpCmd.Start(); err != nil {
		return errors.Prefix("failed to start mysqldump command", err)
	}
	if err = createCmd.Start(); err != nil {
		return errors.Prefix("failed to start mysql command", err)
	}

	if err = dumpCmd.Wait(); err != nil {
		fmt.Println(err)
		return errors.Prefix("failed to wait for mysqldump command", err)
	}

	w.Close() // After dumpCmd is done, close the write end of the pipe

	if err = createCmd.Wait(); err != nil {
		fmt.Println(err)
		return errors.Prefix("failed to wait for mysql command", err)
	}

	return nil
}

func (m *mysqlTester) sslMode(mode string) string {
	switch mode {
	case "true":
		return "REQUIRED"
	case "false":
		return "DISABLED"
	default:
		return "PREFERRED"
	}
}

func (m *mysqlTester) defaultsFile() string {
	return fmt.Sprintf("--defaults-file=%s", m.optionFile)
}

func (m *mysqlTester) makeOptionFile() error {
	tmp, err := ioutil.TempFile("", "optionfile")
	if err != nil {
		return errors.Prefix("failed to create option file", err)
	}

	isTCP := false
	_, err = os.Stat(m.host)
	if os.IsNotExist(err) {
		isTCP = true
	} else if err != nil {
		return errors.Prefix("could not stat m.host", err)
	}

	fmt.Fprintln(tmp, "[client]")
	fmt.Fprintf(tmp, "host=%s\n", m.host)
	fmt.Fprintf(tmp, "port=%d\n", m.port)
	fmt.Fprintf(tmp, "user=%s\n", m.user)
	fmt.Fprintf(tmp, "password=%s\n", m.pass)
	fmt.Fprintf(tmp, "ssl-mode=%s\n", m.sslMode(m.sslmode))
	if isTCP {
		fmt.Fprintln(tmp, "protocol=tcp")
	}

	fmt.Fprintln(tmp, "[mysqldump]")
	fmt.Fprintf(tmp, "host=%s\n", m.host)
	fmt.Fprintf(tmp, "port=%d\n", m.port)
	fmt.Fprintf(tmp, "user=%s\n", m.user)
	fmt.Fprintf(tmp, "password=%s\n", m.pass)
	fmt.Fprintf(tmp, "ssl-mode=%s\n", m.sslMode(m.sslmode))
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
		m.dbConn.Close()
	}

	if err := m.dropTestDB(); err != nil {
		return errors.Err(err)
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
	return errors.Err(err)
	}

	return nil
}

func (m *mysqlTester) conn() (*sql.DB, error) {
	if m.dbConn != nil {
	return m.dbConn, nil
	}

	var err error
	m.dbConn, err = sql.Open("mysql", drivers.MySQLBuildQueryString(m.user, m.pass, m.testDBName, m.host, m.port, m.sslmode))
	if err != nil {
	return nil, err
	}

	return m.dbConn, nil
}
