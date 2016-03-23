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

func LoadConfigFile(filename string) {
	_, err := toml.DecodeFile(filename, &cfg)

	if os.IsNotExist(err) {
		fmt.Fatalf("Failed to find the toml configuration file %s: %s", filename, err)
	}

	if err != nil {
		fmt.Fatalf("Failed to decode toml configuration file:", err)
	}

	if cfg.TestPostgres != nil {
		if cfg.TestPostgres.User == "" || cfg.TestPostgres.Pass == "" ||
			cfg.TestPostgres.Host == "" || cfg.TestPostgres.Port == 0 ||
			cfg.TestPostgres.DBName == "" || cfg.Postgres.DBName == cfg.TestPostgres.DBName {
			cfg.TestPostgres = nil
		}
	}

  if cfg.TestPostgres == nil {
    fmt.Fatalf("Failed to load config.toml postgres_test config")
  }
}

func TestMain(m *testing.M) {
  setup()
  code := m.Run()
  // shutdown
  os.Exit(code)
}

func setup() {
  LoadConfigFile("../config.toml")
}
