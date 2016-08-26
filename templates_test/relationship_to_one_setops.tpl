{{- define "relationship_to_one_setops_test_helper" -}}
{{- $varNameSingular := .ForeignKey.ForeignTable | singular | camelCase -}}
func test{{.LocalTable.NameGo}}ToOneSetOp{{.ForeignTable.NameGo}}_{{.Function.Name}}(t *testing.T) {
}
{{- if .ForeignKey.Nullable}}

func test{{.LocalTable.NameGo}}ToOneRemoveOp{{.ForeignTable.NameGo}}_{{.Function.Name}}(t *testing.T) {
}
{{end -}}
{{- end -}}
{{- if .Table.IsJoinTable -}}
{{- else -}}
  {{- $dot := . -}}
  {{- range .Table.FKeys -}}
    {{- $rel := textsFromForeignKey $dot.PkgName $dot.Tables $dot.Table .}}

{{template "relationship_to_one_setops_test_helper" $rel -}}
{{- end -}}
{{- end -}}
