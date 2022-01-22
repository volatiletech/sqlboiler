var ViewNames = struct {
	{{range $view := .Views -}}
	{{titleCase $view.Name}} string
	{{end -}}
}{
	{{range $view := .Views -}}
	{{titleCase $view.Name}}: "{{$view.Name}}",
	{{end -}}
}

