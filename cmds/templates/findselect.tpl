{{- $tableName := .Table -}}
// {{titleCase $tableName}}FindSelect retrieves the specified columns for a single record by ID.
// Pass in a pointer to an object with `db` tags that match the column names you wish to retrieve.
// For example: friendName string `db:"friend_name"`
func {{titleCase $tableName}}FindSelect(db boil.DB, id int, results interface{}) error {
  if id == 0 {
    return nil, errors.New("{{.PkgName}}: no id provided for {{$tableName}} select")
  }

  query := fmt.Sprintf(`SELECT %s FROM {{$tableName}} WHERE id=$1`, boil.SelectNames(results))
  err := db.Select(results, query, id)

  if err != nil {
    return fmt.Errorf("{{.PkgName}}: unable to select from {{$tableName}}: %s", err)
  }

  return nil
}
