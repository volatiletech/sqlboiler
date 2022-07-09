var ViewNames = struct {
	{{range $table := .Tables}}{{if $table.IsView -}}
	{{titleCase $table.Name}} string
	{{end}}{{end -}}
}{
	{{range $table := .Tables}}{{if $table.IsView -}}
	{{titleCase $table.Name}}: "{{$table.Name}}",
	{{end}}{{end -}}
}

