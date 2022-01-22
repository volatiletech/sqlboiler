var TableNames = struct {
	{{range $table := .Tables -}}
	{{titleCase $table.Name}} string
	{{end -}}
}{
	{{range $table := .Tables -}}
	{{titleCase $table.Name}}: "{{$table.Name}}",
	{{end -}}
}
