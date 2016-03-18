{{- $tableNameSingular := titleCaseSingular .Table -}}
{{- $dbName := singular .Table -}}
{{- $tableNamePlural := titleCasePlural .Table -}}
{{- $varNamePlural := camelCasePlural .Table -}}
// {{$tableNamePlural}}Where retrieves all records with the specified column values.
func {{$tableNamePlural}}Where(db boil.DB, columns map[string]interface{}) ([]*{{$tableNameSingular}}, error) {
  var {{$varNamePlural}} []*{{$tableNameSingular}}
  query := fmt.Sprintf(`SELECT {{selectParamNames $dbName .Columns}} FROM {{.Table}} WHERE %s`, boil.Where(columns))
  err := db.Select(&{{$varNamePlural}}, query, boil.WhereParams(columns)...)

  if err != nil {
    return nil, fmt.Errorf("{{.PkgName}}: unable to select from {{.Table}}: %s", err)
  }

  return {{$varNamePlural}}, nil
}
