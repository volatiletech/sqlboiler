{{- $tableNameSingular := titleCaseSingular .Table -}}
{{- $dbName := singular .Table -}}
{{- $tableNamePlural := titleCasePlural .Table -}}
{{- $varNamePlural := camelCasePlural .Table -}}
// {{$tableNamePlural}}All retrieves all records.
func {{$tableNamePlural}}All(db boil.DB) ([]*{{$tableNameSingular}}, error) {
  var {{$varNamePlural}} []*{{$tableNameSingular}}

  rows, err := db.Query(`SELECT {{selectParamNames $dbName .Columns}} FROM {{.Table}}`)
  if err != nil {
    return nil, fmt.Errorf("{{.PkgName}}: failed to query: %v", err)
  }

  for rows.Next() {
    {{- $tmpVarName := (print $varNamePlural "Tmp") -}}
    {{$varNamePlural}}Tmp := {{$tableNameSingular}}{}

    if err := rows.Scan({{scanParamNames $tmpVarName .Columns}}); err != nil {
      return nil, fmt.Errorf("{{.PkgName}}: failed to scan row: %v", err)
    }

    {{$varNamePlural}} = append({{$varNamePlural}}, &{{$varNamePlural}}Tmp)
  }

  if err := rows.Err(); err != nil {
    return nil, fmt.Errorf("{{.PkgName}}: failed to read rows: %v", err)
  }

  return {{$varNamePlural}}, nil
}
