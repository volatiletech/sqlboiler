{{- if .Table.IsJoinTable -}}
{{- else -}}
  {{- $dot := . -}}
  {{- $table := .Table -}}
  {{- range .Table.ToManyRelationships -}}
    {{- if (and .ForeignColumnUnique (not .ToJoinTable)) -}}
{{- template "relationship_to_one_helper" (textsFromOneToOneRelationship $dot.PkgName $dot.Tables $table .) -}}
    {{- else -}}
    {{- $rel := textsFromRelationship $dot.Tables $table . -}}
// {{$rel.Function.Name}}G retrieves all the {{$rel.LocalTable.NameSingular}}'s {{$rel.ForeignTable.NameHumanReadable}}
{{- if not (eq $rel.Function.Name $rel.ForeignTable.NamePluralGo)}} via {{.ForeignColumn}} column{{- end}}.
func ({{$rel.Function.Receiver}} *{{$rel.LocalTable.NameGo}}) {{$rel.Function.Name}}G(mods ...qm.QueryMod) ({{$rel.ForeignTable.Slice}}, error) {
  return {{$rel.Function.Receiver}}.{{$rel.Function.Name}}(boil.GetDB(), mods...)
}

// {{$rel.Function.Name}}GP panics on error. Retrieves all the {{$rel.LocalTable.NameSingular}}'s {{$rel.ForeignTable.NameHumanReadable}}
{{- if not (eq $rel.Function.Name $rel.ForeignTable.NamePluralGo)}} via {{.ForeignColumn}} column{{- end}}.
func ({{$rel.Function.Receiver}} *{{$rel.LocalTable.NameGo}}) {{$rel.Function.Name}}GP(mods ...qm.QueryMod) {{$rel.ForeignTable.Slice}} {
  slice, err := {{$rel.Function.Receiver}}.{{$rel.Function.Name}}(boil.GetDB(), mods...)
  if err != nil {
    panic(boil.WrapErr(err))
  }

  return slice
}

// {{$rel.Function.Name}}P panics on error. Retrieves all the {{$rel.LocalTable.NameSingular}}'s {{$rel.ForeignTable.NameHumanReadable}} with an executor
{{- if not (eq $rel.Function.Name $rel.ForeignTable.NamePluralGo)}} via {{.ForeignColumn}} column{{- end}}.
func ({{$rel.Function.Receiver}} *{{$rel.LocalTable.NameGo}}) {{$rel.Function.Name}}P(exec boil.Executor, mods ...qm.QueryMod) {{$rel.ForeignTable.Slice}} {
  slice, err := {{$rel.Function.Receiver}}.{{$rel.Function.Name}}(exec, mods...)
  if err != nil {
    panic(boil.WrapErr(err))
  }

  return slice
}

// {{$rel.Function.Name}} retrieves all the {{$rel.LocalTable.NameSingular}}'s {{$rel.ForeignTable.NameHumanReadable}} with an executor
{{- if not (eq $rel.Function.Name $rel.ForeignTable.NamePluralGo)}} via {{.ForeignColumn}} column{{- end}}.
func ({{$rel.Function.Receiver}} *{{$rel.LocalTable.NameGo}}) {{$rel.Function.Name}}(exec boil.Executor, mods ...qm.QueryMod) ({{$rel.ForeignTable.Slice}}, error) {
  queryMods := []qm.QueryMod{
    qm.Select(`"{{id 0}}".*`),
  }

  if len(mods) != 0 {
    queryMods = append(queryMods, mods...)
  }

    {{if .ToJoinTable -}}
  queryMods = append(queryMods,
    qm.InnerJoin(`"{{.JoinTable}}" as "{{id 1}}" on "{{id 1}}"."{{.JoinForeignColumn}}" = "{{id 0}}"."{{.ForeignColumn}}"`),
    qm.Where(`"{{id 1}}"."{{.JoinLocalColumn}}"=$1`, {{.Column | titleCase | printf "%s.%s" $rel.Function.Receiver }}),
  )
    {{else -}}
  queryMods = append(queryMods,
    qm.Where(`"{{id 0}}"."{{.ForeignColumn}}"=$1`, {{.Column | titleCase | printf "%s.%s" $rel.Function.Receiver }}),
  )
    {{end}}

  query := {{$rel.ForeignTable.NamePluralGo}}(exec, queryMods...)
  boil.SetFrom(query.Query, `"{{.ForeignTable}}" as "{{id 0}}"`)
  return query.All()
}

{{end -}}{{- /* if unique foreign key */ -}}
{{- end -}}{{- /* range relationships */ -}}
{{- end -}}{{- /* outer if join table */ -}}
