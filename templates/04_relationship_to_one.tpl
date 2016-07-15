{{- if .Table.IsJoinTable -}}
{{- else -}}
  {{- $dot := . -}}
  {{- range .Table.FKeys -}}
    {{- $rel := textsFromForeignKey $dot.Tables $dot.Table . -}}
// {{$rel.LocalTable.ColumnNameGo}} pointed to by the foreign key.
func ({{$rel.Function.Receiver}} *{{$rel.LocalTable.NameGo}}) {{$rel.LocalTable.ColumnNameGo}}(selectCols ...string) (*{{$rel.ForeignTable.NameGo}}, error) {
  return {{$rel.Function.Receiver}}.{{$rel.LocalTable.ColumnNameGo}}X(boil.GetDB(), selectCols...)
}

// {{$rel.LocalTable.ColumnNameGo}} pointed to by the foreign key.
func ({{$rel.Function.Receiver}} *{{$rel.LocalTable.NameGo}}) {{$rel.LocalTable.ColumnNameGo}}X(exec boil.Executor, selectCols ...string) (*{{$rel.ForeignTable.NameGo}}, error) {
  {{$rel.Function.Varname}} := &{{$rel.ForeignTable.NameGo}}{}

  selectColumns := `*`
  if len(selectCols) != 0 {
    selectColumns = fmt.Sprintf(`"%s"`, strings.Join(selectCols, `","`))
  }

  query := fmt.Sprintf(`select %s from {{.ForeignTable}} where "{{.ForeignColumn}}" = $1`, selectColumns)
  err := exec.QueryRow(query, {{$rel.Function.Receiver}}.{{titleCase .Column}}).Scan(boil.GetStructPointers({{$rel.Function.Varname}}, selectCols...)...)
  if err != nil {
    return nil, fmt.Errorf(`{{$dot.PkgName}}: unable to select from {{.ForeignTable}}: %v`, err)
  }

  return {{$rel.Function.Varname}}, nil
}

{{end -}}
{{- end -}}
