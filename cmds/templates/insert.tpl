{{- $tableName := .Table -}}
// {{titleCase $tableName}}Insert inserts a single record.
func {{titleCase $tableName}}Insert(db boil.DB, o *{{titleCase $tableName}}) (int, error) {
  if o == nil {
    return 0, errors.New("{{.PkgName}}: no {{$tableName}} provided for insertion")
  }

  var rowID int
  err := db.QueryRow(`INSERT INTO {{$tableName}} ({{insertParamNames .Columns}}) VALUES({{insertParamFlags .Columns}}) RETURNING id`)

  if err != nil {
    return 0, fmt.Errorf("{{.PkgName}}: unable to insert {{$tableName}}: %s", err)
  }

  return rowID, nil
}
