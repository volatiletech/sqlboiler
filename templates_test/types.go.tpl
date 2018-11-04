{{- $alias := .Aliases.Table .Table.Name}}
var (
	{{$alias.DownSingular}}DBTypes = map[string]string{{"{"}}{{.Table.Columns | columnDBTypes $alias.UpSingular $alias.Columns | makeStringMap}}{{"}"}}
	_ = bytes.MinRead
)
