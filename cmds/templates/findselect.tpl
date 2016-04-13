{{- $tableNameSingular := titleCaseSingular .Table.Name -}}
// {{$tableNameSingular}}FindSelect retrieves the specified columns for a single record by ID.
// Pass in a pointer to an object with `db` tags that match the column names you wish to retrieve.
// For example: friendName string `db:"friend_name"`
func {{$tableNameSingular}}FindSelect(id int, results interface{}) error {
  if id == 0 {
    return errors.New("{{.PkgName}}: no id provided for {{.Table.Name}} select")
  }

  query := fmt.Sprintf(`SELECT %s FROM {{.Table.Name}} WHERE id=$1`, boil.SelectNames(results))
  err := boil.GetDB().Select(results, query, id)

  if err != nil {
    return fmt.Errorf("{{.PkgName}}: unable to select from {{.Table.Name}}: %s", err)
  }

  return nil
}
