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

  p.dbName = viper.GetString("postgres.dbname")
  p.host = viper.GetString("postgres.host")
  p.user = viper.GetString("postgres.user")
  p.pass = viper.GetString("postgres.pass")
  p.port = viper.GetInt("postgres.port")
  p.sslmode = viper.GetString("postgres.sslmode")
  // Create a randomized db name.
  p.testDBName = randomize.StableDBName(p.dbName)

  if err = p.makePGPassFile(); err != nil {
    return err
  }

  if err = p.dropTestDB(); err != nil {
    return err
  }
  if err = p.createTestDB(); err != nil {
    return err
  }

  dumpCmd := exec.Command("pg_dump", "--schema-only", p.dbName)
  dumpCmd.Env = append(os.Environ(), p.pgEnv()...)
  createCmd := exec.Command("psql", p.testDBName)
  createCmd.Env = append(os.Environ(), p.pgEnv()...)

  r, w := io.Pipe()
  dumpCmd.Stdout = w
  createCmd.Stdin = newFKeyDestroyer(rgxPGFkey, r)

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
  cmd.Env = append(os.Environ(), p.pgEnv()...)

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

func (p *pgTester) pgEnv() []string {
  return []string{
    fmt.Sprintf("PGHOST=%s", p.host),
    fmt.Sprintf("PGPORT=%d", p.port),
    fmt.Sprintf("PGUSER=%s", p.user),
    fmt.Sprintf("PGPASSFILE=%s", p.pgPassFile),
  }
}

func (p *pgTester) makePGPassFile() error {
  tmp, err := ioutil.TempFile("", "pgpass")
  if err != nil {
    return errors.Wrap(err, "failed to create option file")
  }

  fmt.Fprintf(tmp, "%s:%d:postgres:%s", p.host, p.port, p.user)
  if len(p.pass) != 0 {
    fmt.Fprintf(tmp, ":%s", p.pass)
  }
  fmt.Fprintln(tmp)

  fmt.Fprintf(tmp, "%s:%d:%s:%s", p.host, p.port, p.dbName, p.user)
  if len(p.pass) != 0 {
    fmt.Fprintf(tmp, ":%s", p.pass)
  }
  fmt.Fprintln(tmp)

  fmt.Fprintf(tmp, "%s:%d:%s:%s", p.host, p.port, p.testDBName, p.user)
  if len(p.pass) != 0 {
    fmt.Fprintf(tmp, ":%s", p.pass)
  }
  fmt.Fprintln(tmp)

  p.pgPassFile = tmp.Name()
  return tmp.Close()
}

func (p *pgTester) createTestDB() error {
  return p.runCmd("", "createdb", p.testDBName)
}

func (p *pgTester) dropTestDB() error {
  return p.runCmd("", "dropdb", "--if-exists", p.testDBName)
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

  return os.Remove(p.pgPassFile)
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

