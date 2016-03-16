{{- $tableName := .Table -}}
// {{titleCase $tableName}}Select retrieves the specified columns for all records.
// Pass in a pointer to an object with `db` tags that match the column names you wish to retrieve.
// For example: friendName string `db:"friend_name"`
func {{titleCase $tableName}}Select(db boil.DB, results interface{}) error {
  query := fmt.Sprintf(`SELECT %s FROM {{$tableName}}`, boil.SelectNames(results))
  err := db.Select(results, query)

  if err != nil {
    return fmt.Errorf("{{.PkgName}}: unable to select from {{$tableName}}: %s", err)
  }

  return nil
}
