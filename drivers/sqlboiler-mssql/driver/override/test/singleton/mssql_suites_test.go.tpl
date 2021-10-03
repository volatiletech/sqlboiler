func TestUpsert(t *testing.T) {
  {{- range $index, $table := .Tables}}
  {{- if $table.IsJoinTable -}}
  {{- else -}}
  {{- $alias := $.Aliases.Table $table.Name}}
  t.Run("{{$alias.UpPlural}}", test{{$alias.UpPlural}}Upsert)
  {{end -}}
  {{- end -}}
}
