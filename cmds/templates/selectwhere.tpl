{{- $tableNamePlural := titleCasePlural .Table -}}
// {{$tableNamePlural}}SelectWhere retrieves all records with the specified column values.
func {{$tableNamePlural}}SelectWhere(db boil.DB, results interface{}, columns map[string]interface{}) error {
  query := fmt.Sprintf(`SELECT %s FROM {{.Table}} WHERE %s`, boil.SelectNames(results), boil.Where(columns))
  err := db.Select(results, query, boil.WhereParams(columns)...)

  if err != nil {
    return fmt.Errorf("{{.PkgName}}: unable to select from {{.Table}}: %s", err)
  }

  return nil
}
