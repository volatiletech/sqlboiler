{{- $tableName := titleCase .Table -}}
{{- $varName := camelCase .Table -}}
// {{$tableName}}All retrieves all records.
func {{$tableName}}All(db boil.DB) ([]*{{$tableName}}, error) {
  var {{$varName}} []*{{$tableName}}

  rows, err := db.Query(`SELECT {{selectParamNames .Table .Columns}} FROM {{.Table}}`)
  if err != nil {
    return nil, fmt.Errorf("{{.PkgName}}: failed to query: %v", err)
  }

  for rows.Next() {
    {{- $tmpVarName := (print $varName "Tmp") -}}
    {{$varName}}Tmp := {{$tableName}}{}

    if err := rows.Scan({{scanParamNames $tmpVarName .Columns}}); err != nil {
      return nil, fmt.Errorf("{{.PkgName}}: failed to scan row: %v", err)
    }

    {{$varName}} = append({{$varName}}, {{$varName}}Tmp)
  }

  if err := rows.Err(); err != nil {
    return nil, fmt.Errorf("{{.PkgName}}: failed to read rows: %v", err)
  }

  return {{$varName}}, nil
}
