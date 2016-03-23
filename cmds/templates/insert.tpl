{{- $tableNameSingular := titleCaseSingular .Table.Name -}}
// {{$tableNameSingular}}Insert inserts a single record.
func {{$tableNameSingular}}Insert(db boil.DB, o *{{$tableNameSingular}}) (int, error) {
  if o == nil {
    return 0, errors.New("{{.PkgName}}: no {{.Table.Name}} provided for insertion")
  }

  var rowID int
  err := db.QueryRow(`INSERT INTO {{.Table.Name}} ({{insertParamNames .Table.Columns}}) VALUES({{insertParamFlags .Table.Columns}}) RETURNING id`, {{insertParamVariables "o." .Table.Columns}}).Scan(&rowID)

  if err != nil {
    return 0, fmt.Errorf("{{.PkgName}}: unable to insert {{.Table.Name}}: %s", err)
  }

  return rowID, nil
}
