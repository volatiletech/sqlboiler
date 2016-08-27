{{- define "relationship_to_one_setops_helper" -}}
{{- $varNameSingular := .ForeignKey.ForeignTable | singular | camelCase}}

// Set{{.Function.Name}} of the {{.ForeignKey.Table | singular}} to the related item.
// Sets R.{{.Function.Name}} to related.
func ({{.Function.Receiver}} *{{.LocalTable.NameGo}}) Set{{.Function.Name}}(exec boil.Executor, insert bool, related *{{.ForeignTable.NameGo}}) error {
  //{{.Function.Receiver}}.R.{{.Function.Name}} = related
  //{{.Function.Receiver}}.R.{{.Function.Name}}.{{.Function.ForeignAssignment}} = {{.Function.Receiver}}.{{.Function.LocalAssignment}}
  //if insert {
//    return related.Insert()
//  }

//  return exec.Exec(`update "{{.ForeignKey.Table}}" set "{{.ForeignKey.Column}}" = $1`, 5)
return nil
}
{{- if .ForeignKey.Nullable}}

// Remove{{.Function.Name}} relationship.
// Sets R.{{.Function.Name}} to nil.
func ({{.Function.Receiver}} *{{.LocalTable.NameGo}}) Remove{{.Function.Name}}(exec boil.Executor) error {
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
