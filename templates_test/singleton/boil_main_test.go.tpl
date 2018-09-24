var flagDebugMode = flag.Bool("test.sqldebug", false, "Turns on debug mode for SQL statements")
var flagConfigFile = flag.String("test.config", "", "Overrides the default config")

const outputDirDepth = {{.OutputDirDepth}}

var (
	dbMain tester
)

type tester interface {
	setup() error
	conn() (*sql.DB, error)
	teardown() error
}

func TestMain(m *testing.M) {
	if dbMain == nil {
		fmt.Println("no dbMain tester interface was ready")
		os.Exit(-1)
	}

	rand.Seed(time.Now().UnixNano())

	flag.Parse()

	var err error

	// Load configuration
	err = initViper()
	if err != nil {
		fmt.Println("unable to load config file")
		os.Exit(-2)
	}

	// Set DebugMode so we can see generated sql statements
	boil.DebugMode = *flagDebugMode

	if err = dbMain.setup(); err != nil {
		fmt.Println("Unable to execute setup:", err)
		os.Exit(-4)
	}

  conn, err := dbMain.conn()
  if err != nil {
    fmt.Println("failed to get connection:", err)
  }

	var code int
	boil.SetDB(conn)
	code = m.Run()

	if err = dbMain.teardown(); err != nil {
		fmt.Println("Unable to execute teardown:", err)
		os.Exit(-5)
	}

	os.Exit(code)
}

func initViper() error {
 	if flagConfigFile != nil && *flagConfigFile != "" {
		viper.SetConfigFile(*flagConfigFile)
		if err := viper.ReadInConfig(); err != nil {
			return err
		}
		return nil
	}

 	var err error

	viper.SetConfigName("sqlboiler")

	configHome := os.Getenv("XDG_CONFIG_HOME")
	homePath := os.Getenv("HOME")
	wd, err := os.Getwd()
	if err != nil {
		wd = strings.Repeat("../", outputDirDepth)
	} else {
		wd = wd + strings.Repeat("/..", outputDirDepth)
	}

	configPaths := []string{wd}
	if len(configHome) > 0 {
		configPaths = append(configPaths, filepath.Join(configHome, "sqlboiler"))
	} else {
		configPaths = append(configPaths, filepath.Join(homePath, ".config/sqlboiler"))
	}

	for _, p := range configPaths {
		viper.AddConfigPath(p)
	}

	// Ignore errors here, fall back to defaults and validation to provide errs
	_ = viper.ReadInConfig()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	return nil
}
