{{- $tableNamePlural := titleCasePlural .Table.Name -}}
// {{$tableNamePlural}}Select retrieves the specified columns for all records.
// Pass in a pointer to an object with `db` tags that match the column names you wish to retrieve.
// For example: friendName string `db:"friend_name"`
func {{$tableNamePlural}}Select(db boil.DB, results interface{}) error {
  query := fmt.Sprintf(`SELECT %s FROM {{.Table.Name}}`, boil.SelectNames(results))
  err := db.Select(results, query)

  if err != nil {
    return fmt.Errorf("{{.PkgName}}: unable to select from {{.Table.Name}}: %s", err)
  }

  return nil
}
