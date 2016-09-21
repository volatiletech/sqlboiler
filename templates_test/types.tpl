{{- $varNameSingular := .Table.Name | singular | camelCase -}}
var (
	{{$varNameSingular}}DBTypes = map[string]string{{"{"}}{{.Table.Columns | columnDBTypes | makeStringMap}}{{"}"}}
	_ = bytes.MinRead
)
