type mysqlTester struct {
  dbConn *sql.DB
}

dbMain = mysqlTester{}

func (m mysqlTester) setup() error {
  return nil
}

func (m mysqlTester) teardown() error {
  return nil
}

func (m mysqlTester) conn() *sql.DB {
  return m.dbConn
}
