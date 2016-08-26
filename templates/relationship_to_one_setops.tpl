{{- define "relationship_to_one_setops_helper" -}}
{{- $varNameSingular := .ForeignKey.ForeignTable | singular | camelCase}}

// Set{{.Function.Name}} of the {{.ForeignKey.Table | singular}} to the related item.
// Sets R.{{.Function.Name}} to related.
func (r *{{.LocalTable.NameGo}}Loaded) Set{{.Function.Name}}(exec boil.Executor, insert bool, related *{{.ForeignTable.NameGo}}) error {
  r.{{.Function.Name}} = related
  r.{{.Function.Name}}.{{.Function.ForeignAssignment}} = 5
  if insert {
    return related.Insert()
  }

  return exec.Exec(`update "{{.ForeignKey.Table}}" set "{{.ForeignKey.Column}}" = $1`, 5)
}
{{- if .ForeignKey.Nullable}}

// Remove{{.Function.Name}} relationship.
// Sets R.{{.Function.Name}} to nil.
func (r *{{.LocalTable.NameGo}}Loaded) Remove{{.Function.Name}}(exec boil.Executor) error {
  return nil
}
{{end -}}
{{- end -}}
{{- if .Table.IsJoinTable -}}
{{- else -}}
  {{- $dot := . -}}
  {{- range .Table.FKeys -}}
    {{- $rel := textsFromForeignKey $dot.PkgName $dot.Tables $dot.Table . -}}
{{- template "relationship_to_one_setops_helper" $rel -}}
{{- end -}}
{{- end -}}
