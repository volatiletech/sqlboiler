{{- $tableNamePlural := titleCasePlural .Table.Name -}}
// {{$tableNamePlural}}SelectWhere retrieves all records with the specified column values.
func {{$tableNamePlural}}SelectWhere(results interface{}, columns map[string]interface{}) error {
  query := fmt.Sprintf(`SELECT %s FROM {{.Table.Name}} WHERE %s`, boil.SelectNames(results), boil.WhereClause(columns))
  err := boil.GetDB().Select(results, query, boil.WhereParams(columns)...)

  if err != nil {
    return fmt.Errorf("{{.PkgName}}: unable to select from {{.Table.Name}}: %s", err)
  }

  return nil
}
