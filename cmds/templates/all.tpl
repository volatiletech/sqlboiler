{{- $tableNameSingular := titleCaseSingular .Table.Name -}}
{{- $dbName := singular .Table.Name -}}
{{- $tableNamePlural := titleCasePlural .Table.Name -}}
{{- $varNameSingular := camelCaseSingular .Table.Name -}}
{{- $varNamePlural := camelCasePlural .Table.Name -}}
type {{$varNameSingular}}Query struct {
  *boil.Query
}

// {{$tableNamePlural}}All retrieves all records.
func {{$tableNamePlural}}(mods ...QueryMod) {{$varNameSingular}}Query {
  var {{$varNamePlural}} []*{{$tableNameSingular}}

  rows, err := boil.GetDB().Query(`SELECT {{selectParamNames $dbName .Table.Columns}} FROM {{.Table.Name}}`)
  if err != nil {
    return nil, fmt.Errorf("{{.PkgName}}: failed to query: %v", err)
  }

  for rows.Next() {
    {{- $tmpVarName := (print $varNamePlural "Tmp") -}}
    {{$varNamePlural}}Tmp := {{$tableNameSingular}}{}

    if err := rows.Scan({{scanParamNames $tmpVarName .Table.Columns}}); err != nil {
      return nil, fmt.Errorf("{{.PkgName}}: failed to scan row: %v", err)
    }

    {{$varNamePlural}} = append({{$varNamePlural}}, &{{$varNamePlural}}Tmp)
  }

  if err := rows.Err(); err != nil {
    return nil, fmt.Errorf("{{.PkgName}}: failed to read rows: %v", err)
  }

  return {{$varNamePlural}}, nil
}

func {{$tableNamePlural}}X(exec boil.Executor, mods ...QueryMod) {{$tableNameSingular}}Query {


}
