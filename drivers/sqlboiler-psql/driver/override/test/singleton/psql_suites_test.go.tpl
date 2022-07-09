func TestUpsert(t *testing.T) {
  {{- range $index, $table := .Tables}}
  {{- if or $table.IsJoinTable $table.IsView -}}
  {{- else -}}
  {{- $alias := $.Aliases.Table $table.Name}}
  t.Run("{{$alias.UpPlural}}", test{{$alias.UpPlural}}Upsert)
  {{end -}}
  {{- end -}}
}
