{{- $alias := .Aliases.Table .Table.Name -}}

var {{$alias.UpSingular}}Order = struct {
	{{range $column := .Table.Columns -}}
	{{- $colAlias := $alias.Column $column.Name -}}
	{{$colAlias}}Asc string
	{{$colAlias}}Desc string
	{{end -}}
}{
	{{range $column := .Table.Columns -}}
	{{- $colAlias := $alias.Column $column.Name -}}
	{{$colAlias}}Asc: "{{$.Table.Name | $.SchemaTable}}.{{$column.Name | $.Quotes}} ASC",
	{{$colAlias}}Desc: "{{$.Table.Name | $.SchemaTable}}.{{$column.Name | $.Quotes}} DESC",
	{{end -}}
}
