var TableNames = struct {
	{{range $table := .Tables -}}
	{{titleCase $table.Name}} string
	{{end -}}
}{
	{{range $table := .Tables -}}
	{{titleCase $table.Name}}: "{{$table.Name}}",
	{{end -}}
}

var TypeNameToTableName = map[string]string {
	{{range $table := .Tables -}}
	{{- $alias := $.Aliases.Table $table.Name -}}
	"{{$alias.UpSingular}}": TableNames.{{titleCase $table.Name}},
	{{end -}}
}

var TypeNameToTableColumns = map[string][]string{
	{{range $table := .Tables -}}
	{{- $alias := $.Aliases.Table $table.Name -}}
	"{{$alias.UpSingular}}": {{$alias.DownSingular}}AllColumns,
	{{end -}}
}

var TableNameToTableColumns = map[string][]string{
{{range $table := .Tables -}}
	{{- $alias := $.Aliases.Table $table.Name -}}
	TableNames.{{titleCase $table.Name}}: {{$alias.DownSingular}}AllColumns,
{{end -}}
}