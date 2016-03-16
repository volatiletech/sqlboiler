{{- $tableName := .Table -}}
// {{titleCase $tableName}}SelectWhere retrieves all records with the specified column values.
func {{titleCase $tableName}}SelectWhere(db boil.DB, results interface{}, columns map[string]interface{}) error {
  query := fmt.Sprintf(`SELECT %s FROM {{$tableName}} WHERE %s`, boil.SelectNames(results), boil.Where(columns))
  err := db.Select(results, query, boil.WhereParams(columns)...)

  if err != nil {
    return fmt.Errorf("{{.PkgName}}: unable to select from {{$tableName}}: %s", err)
  }

  return nil
}
