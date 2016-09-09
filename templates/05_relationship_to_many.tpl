{{- /* Begin execution of template for many-to-one or many-to-many relationship helper */ -}}
{{- if .Table.IsJoinTable -}}
{{- else -}}
  {{- $dot := . -}}
  {{- $table := .Table -}}
  {{- range .Table.ToManyRelationships -}}
    {{- $varNameSingular := .ForeignTable | singular | camelCase -}}
    {{- if (and .ForeignColumnUnique (not .ToJoinTable)) -}}
      {{- /* Begin execution of template for many-to-one relationship. */ -}}
      {{- $txt := textsFromOneToOneRelationship $dot.PkgName $dot.Tables $table . -}}
      {{- template "relationship_to_one_helper" (preserveDot $dot $txt) -}}
    {{- else -}}
      {{- /* Begin execution of template for many-to-many relationship. */ -}}
      {{- $rel := textsFromRelationship $dot.Tables $table . -}}
// {{$rel.Function.Name}}G retrieves all the {{$rel.LocalTable.NameSingular}}'s {{$rel.ForeignTable.NameHumanReadable}}
{{- if not (eq $rel.Function.Name $rel.ForeignTable.NamePluralGo)}} via {{.ForeignColumn}} column{{- end}}.
func ({{$rel.Function.Receiver}} *{{$rel.LocalTable.NameGo}}) {{$rel.Function.Name}}G(mods ...qm.QueryMod) {{$varNameSingular}}Query {
  return {{$rel.Function.Receiver}}.{{$rel.Function.Name}}(boil.GetDB(), mods...)
}

// {{$rel.Function.Name}} retrieves all the {{$rel.LocalTable.NameSingular}}'s {{$rel.ForeignTable.NameHumanReadable}} with an executor
{{- if not (eq $rel.Function.Name $rel.ForeignTable.NamePluralGo)}} via {{.ForeignColumn}} column{{- end}}.
func ({{$rel.Function.Receiver}} *{{$rel.LocalTable.NameGo}}) {{$rel.Function.Name}}(exec boil.Executor, mods ...qm.QueryMod) {{$varNameSingular}}Query {
  queryMods := []qm.QueryMod{
    qm.Select(`"{{id 0}}".*`),
  }

  if len(mods) != 0 {
    queryMods = append(queryMods, mods...)
  }

    {{if .ToJoinTable -}}
  queryMods = append(queryMods,
    qm.InnerJoin(`{{schemaTable $dot.Dialect.LQ $dot.Dialect.RQ $dot.DriverName $dot.Schema .JoinTable}} as "{{id 1}}" on "{{id 0}}"."{{.ForeignColumn}}" = "{{id 1}}"."{{.JoinForeignColumn}}"`),
    qm.Where(`"{{id 1}}"."{{.JoinLocalColumn}}"=$1`, {{$rel.Function.Receiver}}.{{$rel.LocalTable.ColumnNameGo}}),
  )
    {{else -}}
  queryMods = append(queryMods,
    qm.Where(`"{{id 0}}"."{{.ForeignColumn}}"=$1`, {{$rel.Function.Receiver}}.{{$rel.LocalTable.ColumnNameGo}}),
  )
    {{end}}

  query := {{$rel.ForeignTable.NamePluralGo}}(exec, queryMods...)
  boil.SetFrom(query.Query, `{{schemaTable $dot.Dialect.LQ $dot.Dialect.RQ $dot.DriverName $dot.Schema .ForeignTable}} as "{{id 0}}"`)
  return query
}

{{end -}}{{- /* if unique foreign key */ -}}
{{- end -}}{{- /* range relationships */ -}}
{{- end -}}{{- /* if isJoinTable */ -}}
