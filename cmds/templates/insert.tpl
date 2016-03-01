{{- $tableName := .TableName -}}
// {{makeGoName $tableName}}Insert inserts a single record.
func {{makeGoName $tableName}}Insert(db boil.DB, o *{{makeGoName $tableName}}) (int, error) {
  if o == nil {
    return 0, errors.New("model: no {{$tableName}} provided for insertion")
  }

  var rowID int
  err := db.QueryRow(`INSERT INTO {{$tableName}} ({{makeGoInsertParamNames .TableData}}) VALUES({{makeGoInsertParamFlags .TableData}}) RETURNING id`)

  if err != nil {
    return 0, fmt.Errorf("model: unable to insert {{$tableName}}: %s", err)
  }

  return rowID, nil
}
