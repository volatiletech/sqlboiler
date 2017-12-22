type pgTester struct {
  dbConn *sql.DB

  dbName  string
  host    string
  user    string
  pass    string
  sslmode string
  port    int

  pgPassFile string

  testDBName string
}

func init() {
  dbMain = &pgTester{}
}

// setup dumps the database schema and imports it into a temporary randomly
// generated test database so that tests can be run against it using the
// generated sqlboiler ORM package.
func (p *pgTester) setup() error {
  var err error

  p.dbName = viper.GetString("cockroach.dbname")
  p.host = viper.GetString("cockroach.host")
  p.user = viper.GetString("cockroach.user")
  p.pass = viper.GetString("cockroach.pass")
  p.port = viper.GetInt("cockroach.port")
  p.sslmode = viper.GetString("cockroach.sslmode")
  // Create a randomized db name.
  p.testDBName = randomize.StableDBName(p.dbName)

  if err = p.dropTestDB(); err != nil {
    return err
  }
  if err = p.createTestDB(); err != nil {
    return err
  }

  dumpCmd := exec.Command("cockroach", "dump", p.dbName, "--insecure", "--dump-mode=schema")
  createCmd := exec.Command("cockroach", "sql", "--database", p.testDBName, "--insecure")

  r, w := io.Pipe()
  dumpCmd.Stdout = w
  createCmd.Stdin = newFKeyDestroyer(rgxCDBFkey, r)

  if err = dumpCmd.Start(); err != nil {
      return errors.Wrap(err, "failed to start pg_dump command")
  }
  if err = createCmd.Start(); err != nil {
      return errors.Wrap(err, "failed to start psql command")
  }

  if err = dumpCmd.Wait(); err != nil {
      fmt.Println(err)
      return errors.Wrap(err, "failed to wait for pg_dump command")
  }

  w.Close() // After dumpCmd is done, close the write end of the pipe

  if err = createCmd.Wait(); err != nil {
      fmt.Println(err)
      return errors.Wrap(err, "failed to wait for psql command")
  }

  return nil
}

func (p *pgTester) runCmd(stdin, command string, args ...string) error {
  cmd := exec.Command(command, args...)
  cmd.Env = os.Environ()

  if len(stdin) != 0 {
    cmd.Stdin = strings.NewReader(stdin)
  }

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

func (p *pgTester) createTestDB() error {
  stmt := fmt.Sprintf("CREATE DATABASE %s", p.testDBName)
  return p.runCmd("", "cockroach", "sql", "--insecure", "--execute", stmt)
}

func (p *pgTester) dropTestDB() error {
  stmt := fmt.Sprintf("DROP DATABASE IF EXISTS %s CASCADE", p.testDBName)
  return p.runCmd("", "cockroach", "sql", "--insecure", "--execute", stmt)
}

// teardown executes cleanup tasks when the tests finish running
func (p *pgTester) teardown() error {
  var err error
  if err = p.dbConn.Close(); err != nil {
    return err
  }
  p.dbConn = nil

  if err = p.dropTestDB(); err != nil {
    return err
  }

  return nil
}

func (p *pgTester) conn() (*sql.DB, error) {
  if p.dbConn != nil {
    return p.dbConn, nil
  }

  var err error
  p.dbConn, err = sql.Open("postgres", drivers.PostgresBuildQueryString(p.user, p.pass, p.testDBName, p.host, p.port, p.sslmode))
  if err != nil {
    return nil, err
  }

  return p.dbConn, nil
}

