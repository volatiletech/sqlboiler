{{- define "relationship_to_one_helper"}}
// {{.Function.Name}} pointed to by the foreign key.
func ({{.Function.Receiver}} *{{.LocalTable.NameGo}}) {{.Function.Name}}(selectCols ...string) (*{{.ForeignTable.NameGo}}, error) {
  return {{.Function.Receiver}}.{{.Function.Name}}X(boil.GetDB(), selectCols...)
}

// {{.Function.Name}}P pointed to by the foreign key. Panics on error.
func ({{.Function.Receiver}} *{{.LocalTable.NameGo}}) {{.Function.Name}}P(selectCols ...string) *{{.ForeignTable.NameGo}} {
  o, err := {{.Function.Receiver}}.{{.Function.Name}}X(boil.GetDB(), selectCols...)
  if err != nil {
    panic(boil.WrapErr(err))
  }

  return o
}

// {{.Function.Name}}XP pointed to by the foreign key with exeuctor. Panics on error.
func ({{.Function.Receiver}} *{{.LocalTable.NameGo}}) {{.Function.Name}}XP(exec boil.Executor, selectCols ...string) *{{.ForeignTable.NameGo}} {
  o, err := {{.Function.Receiver}}.{{.Function.Name}}X(exec, selectCols...)
  if err != nil {
    panic(boil.WrapErr(err))
  }

  return o
}

// {{.Function.Name}}X pointed to by the foreign key.
func ({{.Function.Receiver}} *{{.LocalTable.NameGo}}) {{.Function.Name}}X(exec boil.Executor, selectCols ...string) (*{{.ForeignTable.NameGo}}, error) {
  {{.Function.Varname}} := &{{.ForeignTable.NameGo}}{}

  selectColumns := `*`
  if len(selectCols) != 0 {
    selectColumns = fmt.Sprintf(`"%s"`, strings.Join(selectCols, `","`))
  }

  query := fmt.Sprintf(`select %s from {{.ForeignTable.Name}} where "{{.ForeignTable.ColumnName}}" = $1`, selectColumns)
  err := exec.QueryRow(query, {{.Function.Receiver}}.{{.LocalTable.ColumnNameGo}}).Scan(boil.GetStructPointers({{.Function.Varname}}, selectCols...)...)
  if err != nil {
    return nil, fmt.Errorf(`{{.Function.PackageName}}: unable to select from {{.ForeignTable.Name}}: %v`, err)
  }

  return {{.Function.Varname}}, nil
}

{{end -}}
{{- if .Table.IsJoinTable -}}
{{- else -}}
  {{- $dot := . -}}
  {{- range .Table.FKeys -}}
    {{- $rel := textsFromForeignKey $dot.PkgName $dot.Tables $dot.Table . -}}
{{- template "relationship_to_one_helper" $rel -}}
{{end -}}
{{- end -}}
