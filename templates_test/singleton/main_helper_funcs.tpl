var testCfg *Config

var dbConn *sql.DB

func InitViper() error {
  var err error
	testCfg = &Config{}

	viper.SetConfigName("sqlboiler")

	configHome := os.Getenv("XDG_CONFIG_HOME")
	homePath := os.Getenv("HOME")
	wd, err := os.Getwd()
	if err != nil {
		wd = "../"
	} else {
    wd = wd + "/.."
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

	err = viper.ReadInConfig()
  if err != nil {
    return err
  }

  return nil
}
