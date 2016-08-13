// M type is for providing columns and column values to UpdateAll.
type M map[string]interface{}

type upsertData struct {
  conflict  []string
  update    []string
  whitelist []string
  returning []string
}

// ErrSyncFail occurs during insert when the record could not be retrieved in
// order to populate default value information. This usually happens when LastInsertId
// fails or there was a primary key configuration that was not resolvable.
var ErrSyncFail = errors.New("{{.PkgName}}: failed to synchronize data after insert")
