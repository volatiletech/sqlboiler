{{- $tableName := .TableName -}}
func Delete{{makeGoColName $tableName}}(id int, db *sqlx.DB) error {
  if id == nil {
    return nil, errors.New("No ID provided for {{makeGoColName $tableName}} delete")
  }

  err := db.Exec("DELETE FROM {{$tableName}} WHERE id=$1", id)
  if err != nil {
    return errors.New("Unable to delete from {{$tableName}}: %s", err)
  }

  return nil
}
