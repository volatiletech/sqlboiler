var TableNames = struct {
	{{range $table := .Tables -}}
	{{titleCase $table.Name}} string
	{{end -}}
}{
	{{range $table := .Tables -}}
	{{titleCase $table.Name}}: "{{$table.Name}}",
	{{end -}}
}

var ViewNames = struct {
	{{range $view := .Views -}}
	{{titleCase $view.Name}} string
	{{end -}}
}{
	{{range $view := .Views -}}
	{{titleCase $view.Name}}: "{{$view.Name}}",
	{{end -}}
}
