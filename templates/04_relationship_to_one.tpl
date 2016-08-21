{{- define "relationship_to_one_helper" -}}
{{- $varNameSingular := .ForeignKey.ForeignTable | singular | camelCase -}}
// {{.Function.Name}}G pointed to by the foreign key.
func ({{.Function.Receiver}} *{{.LocalTable.NameGo}}) {{.Function.Name}}G(mods ...qm.QueryMod) {{$varNameSingular}}Query {
  return {{.Function.Receiver}}.{{.Function.Name}}(boil.GetDB(), mods...)
}

// {{.Function.Name}} pointed to by the foreign key.
func ({{.Function.Receiver}} *{{.LocalTable.NameGo}}) {{.Function.Name}}(exec boil.Executor, mods ...qm.QueryMod) ({{$varNameSingular}}Query) {
  queryMods := []qm.QueryMod{
    qm.Where("{{.ForeignTable.ColumnName}}=$1", {{.Function.Receiver}}.{{.LocalTable.ColumnNameGo}}),
  }

  queryMods = append(queryMods, mods...)

  query := {{.ForeignTable.NamePluralGo}}(exec, queryMods...)
  boil.SetFrom(query.Query, "{{.ForeignTable.Name}}")

  return query
}

{{end -}}
{{- if .Table.IsJoinTable -}}
{{- else -}}
  {{- $dot := . -}}
  {{- range .Table.FKeys -}}
    {{- $rel := textsFromForeignKey $dot.PkgName $dot.Tables $dot.Table . -}}
{{- template "relationship_to_one_helper" $rel -}}
{{- end -}}
{{- end -}}
