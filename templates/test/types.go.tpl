{{- $alias := .Aliases.Table .Table.Name}}
var (
	{{$alias.DownSingular}}DBTypes = map[string]string{{"{"}}{{range $i, $col := .Table.Columns -}}{{- if ne $i 0}},{{end}}`{{$alias.Column $col.Name}}`: `{{$col.DBType}}`{{end}}{{"}"}}
	_ = bytes.MinRead
)
