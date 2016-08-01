{{- if .Table.IsJoinTable -}}
{{- else -}}
  {{- $dot := . }}
  {{- $table := .Table }}
  {{- range .Table.ToManyRelationships -}}
    {{- if .ForeignColumnUnique -}}
{{- template "relationship_to_one_helper" (textsFromOneToOneRelationship $dot.PkgName $dot.Tables $table .) -}}
    {{- else -}}
    {{- $rel := textsFromRelationship $dot.Tables $table . -}}
// {{$rel.Function.Name}}G retrieves all the {{$rel.LocalTable.NameSingular}}'s {{$rel.ForeignTable.NameHumanReadable}}
{{- if not (eq $rel.Function.Name $rel.ForeignTable.NamePluralGo)}} via {{.ForeignColumn}} column{{- end}}.
func ({{$rel.Function.Receiver}} *{{$rel.LocalTable.NameGo}}) {{$rel.Function.Name}}G(selectCols ...string) ({{$rel.ForeignTable.Slice}}, error) {
  return {{$rel.Function.Receiver}}.{{$rel.Function.Name}}(boil.GetDB(), selectCols...)
}

// {{$rel.Function.Name}}GP panics on error. Retrieves all the {{$rel.LocalTable.NameSingular}}'s {{$rel.ForeignTable.NameHumanReadable}}
{{- if not (eq $rel.Function.Name $rel.ForeignTable.NamePluralGo)}} via {{.ForeignColumn}} column{{- end}}.
func ({{$rel.Function.Receiver}} *{{$rel.LocalTable.NameGo}}) {{$rel.Function.Name}}GP(selectCols ...string) {{$rel.ForeignTable.Slice}} {
  o, err := {{$rel.Function.Receiver}}.{{$rel.Function.Name}}(boil.GetDB(), selectCols...)
  if err != nil {
    panic(boil.WrapErr(err))
  }

  return o
}

// {{$rel.Function.Name}}P panics on error. Retrieves all the {{$rel.LocalTable.NameSingular}}'s {{$rel.ForeignTable.NameHumanReadable}} with an executor
{{- if not (eq $rel.Function.Name $rel.ForeignTable.NamePluralGo)}} via {{.ForeignColumn}} column{{- end}}.
func ({{$rel.Function.Receiver}} *{{$rel.LocalTable.NameGo}}) {{$rel.Function.Name}}P(exec boil.Executor, selectCols ...string) {{$rel.ForeignTable.Slice}} {
  o, err := {{$rel.Function.Receiver}}.{{$rel.Function.Name}}(exec, selectCols...)
  if err != nil {
    panic(boil.WrapErr(err))
  }

  return o
}

// {{$rel.Function.Name}} retrieves all the {{$rel.LocalTable.NameSingular}}'s {{$rel.ForeignTable.NameHumanReadable}} with an executor
{{- if not (eq $rel.Function.Name $rel.ForeignTable.NamePluralGo)}} via {{.ForeignColumn}} column{{- end}}.
func ({{$rel.Function.Receiver}} *{{$rel.LocalTable.NameGo}}) {{$rel.Function.Name}}(exec boil.Executor, selectCols ...string) ({{$rel.ForeignTable.Slice}}, error) {
  var ret {{$rel.ForeignTable.Slice}}

  selectColumns := `"{{id 0}}".*`
  if len(selectCols) != 0 {
    selectColumns = `"{{id 0}}".` + strings.Join(selectCols, `","{{id 0}}"."`)
  }
    {{if .ToJoinTable -}}
  query := fmt.Sprintf(`select %s from {{.ForeignTable}} "{{id 0}}" inner join {{.JoinTable}} as "{{id 1}}" on "{{id 1}}"."{{.JoinForeignColumn}}" = "{{id 0}}"."{{.ForeignColumn}}" where "{{id 1}}"."{{.JoinLocalColumn}}"=$1`, selectColumns)
    {{else -}}
  query := fmt.Sprintf(`select %s from {{.ForeignTable}} "{{id 0}}" where "{{id 0}}"."{{.ForeignColumn}}"=$1`, selectColumns)
    {{end}}
  rows, err := exec.Query(query, {{.Column | titleCase | printf "%s.%s" $rel.Function.Receiver }})
  if err != nil {
    return nil, fmt.Errorf(`{{$dot.PkgName}}: unable to select from {{.ForeignTable}}: %v`, err)
  }
  defer rows.Close()

  for rows.Next() {
    next := new({{$rel.ForeignTable.NameGo}})

    err = rows.Scan(boil.GetStructPointers(next, selectCols...)...)
    if err != nil {
      return nil, fmt.Errorf(`{{$dot.PkgName}}: unable to scan into {{$rel.ForeignTable.NameGo}}: %v`, err)
    }

    ret = append(ret, next)
  }

  if err = rows.Err(); err != nil {
    return nil, fmt.Errorf(`{{$dot.PkgName}}: unable to select from {{$rel.ForeignTable.NameGo}}: %v`, err)
  }

  return ret, nil
}

{{end -}}{{- /* if unique foreign key */ -}}
{{- end -}}{{- /* range relationships */ -}}
{{- end -}}{{- /* outer if join table */ -}}
