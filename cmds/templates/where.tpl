{{- $tableNameSingular := titleCaseSingular .Table.Name -}}
{{- $dbName := singular .Table.Name -}}
{{- $tableNamePlural := titleCasePlural .Table.Name -}}
{{- $varNamePlural := camelCasePlural .Table.Name -}}
// {{$tableNamePlural}}Where retrieves all records with the specified column values.
func {{$tableNamePlural}}Where(db boil.DB, columns map[string]interface{}) ([]*{{$tableNameSingular}}, error) {
  var {{$varNamePlural}} []*{{$tableNameSingular}}
  query := fmt.Sprintf(`SELECT {{selectParamNames $dbName .Table.Columns}} FROM {{.Table.Name}} WHERE %s`, boil.Where(columns))
  err := db.Select(&{{$varNamePlural}}, query, boil.WhereParams(columns)...)

  if err != nil {
    return nil, fmt.Errorf("{{.PkgName}}: unable to select from {{.Table.Name}}: %s", err)
  }

  return {{$varNamePlural}}, nil
}
