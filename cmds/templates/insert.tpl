{{- $tableNameSingular := titleCaseSingular .Table -}}
// {{$tableNameSingular}}Insert inserts a single record.
func {{$tableNameSingular}}Insert(db boil.DB, o *{{$tableNameSingular}}) (int, error) {
  if o == nil {
    return 0, errors.New("{{.PkgName}}: no {{.Table}} provided for insertion")
  }

  var rowID int
  err := db.QueryRow(`INSERT INTO {{.Table}} ({{insertParamNames .Columns}}) VALUES({{insertParamFlags .Columns}}) RETURNING id`, {{insertParamVariables "o." .Columns}}).Scan(&rowID)

  if err != nil {
    return 0, fmt.Errorf("{{.PkgName}}: unable to insert {{.Table}}: %s", err)
  }

  return rowID, nil
}
