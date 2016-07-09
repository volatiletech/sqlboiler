{{- if .Table.IsJoinTable -}}
{{- else -}}
  {{- $pkg := .PkgName -}}
  {{- $localTable := .Table.Name | singular | titleCase -}}
  {{- range .Table.FKeys -}}
    {{- $localColumn := .Column | remove "_id" | singular | titleCase -}}
    {{- $foreignColumn := .ForeignColumn | remove "_id" | singular | titleCase -}}
    {{- $foreignTable := .ForeignTable | singular | titleCase -}}
    {{- $varname := .ForeignTable | singular | camelCase -}}
    {{- $receiver := $localTable | toLower | substring 0 1 -}}
// {{$localColumn}} fetches the {{$foreignTable}} pointed to by the foreign key.
func ({{$receiver}} *{{$localTable}}) {{$localColumn}}(selectCols ...string) (*{{$foreignTable}}, error) {
  return {{$receiver}}.{{$localColumn}}X(boil.GetDB(), selectCols...)
}

// {{$localColumn}} fetches the {{$foreignTable}} pointed to by the foreign key.
func ({{$receiver}} *{{$localTable}}) {{$localColumn}}X(exec boil.Executor, selectCols ...string) (*{{$foreignTable}}, error) {
  {{$varname}} := &{{$foreignTable}}{}

  query := fmt.Sprintf(`select "%s" from {{.ForeignTable}} where {{.Column}} = $1`, strings.Join(selectCols, `","`))
  err := exec.QueryRow(query, {{$receiver}}.{{titleCase .Column}}).Scan(boil.GetStructPointers({{$varname}}, selectCols...)...)
  if err != nil {
    return nil, fmt.Errorf(`{{$pkg}}: unable to select from {{.ForeignTable}}: %v`, err)
  }

  return {{$varname}}, nil
}

{{end -}}
{{- end -}}
