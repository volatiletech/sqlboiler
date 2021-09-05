func TestUpsert(t *testing.T) {
  {{- range $index, $table := .Tables}}
  {{- if $table.IsJoinTable -}}
  {{- else -}}
  {{- $alias := $.Aliases.Table $table.Name}}
  {{if $.AddStrictUpsert -}}
  t.Run("{{$alias.UpPlural}}", test{{$alias.UpPlural}}UpsertBy{{$table.PKey.TitleCase}})
  {{range $ukey := $table.UKeys -}}
  t.Run("{{$alias.UpPlural}}", test{{$alias.UpPlural}}UpsertBy{{$ukey.TitleCase}})
  {{end -}}
  t.Run("{{$alias.UpPlural}}", test{{$alias.UpPlural}}UpsertDoNothing)
  {{- end -}}
  {{end -}}
  {{- end -}}
}
