func TestUpsert(t *testing.T) {
  {{- range $index, $table := .Tables}}
  {{- if $table.IsJoinTable -}}
  {{- else -}}
  {{- $tableName := $table.Name | plural | titleCase -}}
  t.Run("{{$tableName}}", test{{$tableName}}Upsert)
  {{end -}}
  {{- end -}}
}
