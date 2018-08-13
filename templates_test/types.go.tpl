{{- $alias := .Aliases.Table .Table.Name}}
var (
	{{$alias.DownSingular}}DBTypes = map[string]string{{"{"}}{{.Table.Columns | columnDBTypes | makeStringMap}}{{"}"}}
	_ = bytes.MinRead
)
