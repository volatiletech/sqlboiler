{{- if .Table.IsJoinTable -}}
{{- else -}}
  {{- $dot := . -}}
  {{- $table := .Table -}}
  {{- range .Table.ToManyRelationships -}}
    {{- $varNameSingular := .ForeignTable | singular | camelCase -}}
    {{- if (and .ForeignColumnUnique (not .ToJoinTable)) -}}
{{- template "relationship_to_one_setops_test_helper" (textsFromOneToOneRelationship $dot.PkgName $dot.Tables $table .) -}}
    {{- else -}}
    {{- $rel := textsFromRelationship $dot.Tables $table .}}

func test{{$rel.LocalTable.NameGo}}ToManyAddOp{{$rel.Function.Name}}(t *testing.T) {
}
{{if .ForeignColumnNullable}}

func test{{$rel.LocalTable.NameGo}}ToManySetOp{{$rel.Function.Name}}(t *testing.T) {
}

func test{{$rel.LocalTable.NameGo}}ToManyRemoveOp{{$rel.Function.Name}}(t *testing.T) {
}
{{end -}}
{{- end -}}{{- /* if unique foreign key */ -}}
{{- end -}}{{- /* range relationships */ -}}
{{- end -}}{{- /* outer if join table */ -}}
