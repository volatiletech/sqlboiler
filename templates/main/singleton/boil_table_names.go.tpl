var TableNames = struct {
	{{range $table := .Tables}}{{if not $table.IsView -}}
	{{$tblAlias := index $.Aliases.Tables $table.Name -}}
	{{$tblAlias.UpPlural}} string
	{{end}}{{end -}}
}{
	{{range $table := .Tables}}{{if not $table.IsView -}}
	{{$tblAlias := index $.Aliases.Tables $table.Name -}}
	{{$tblAlias.UpPlural}}: "{{$table.Name}}",
	{{end}}{{end -}}
}
