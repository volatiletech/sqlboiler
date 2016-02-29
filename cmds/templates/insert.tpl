{{- $tableName := .TableName -}}
// {{makeGoColName $tableName}}Insert inserts a single record.
func {{makeGoColName $tableName}}Insert(db *sqlx.DB, o *{{makeGoColName $tableName}}) (int, error) {
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
