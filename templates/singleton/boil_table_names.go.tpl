var TableNames = struct {
	{{range $table := .Tables -}}
	{{titleCase $table.Name}} string
	{{end -}}
}{
	{{range $table := .Tables -}}
	{{titleCase $table.Name}}: "{{$table.Name}}",
	{{end -}}
}

var typeNameToTableName = map[string]string {
	{{range $table := .Tables -}}
	{{- $alias := $.Aliases.Table $table.Name -}}
	"{{$alias.UpSingular}}": TableNames.{{titleCase $table.Name}},
	{{end -}}
}

var typeNameToTableColumns = map[string][]string{
	{{range $table := .Tables -}}
	{{- $alias := $.Aliases.Table $table.Name -}}
	"{{$alias.UpSingular}}": {{$alias.DownSingular}}AllColumns,
	{{end -}}
}