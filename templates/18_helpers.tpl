{{- $varNameSingular := .Table.Name | singular | camelCase -}}
{{- $tableNameSingular := .Table.Name | singular | titleCase -}}
func (o {{$tableNameSingular}}) inPrimaryKeyArgs() []interface{} {
	var args []interface{}

	{{- range $key, $value := .Table.PKey.Columns }}
	args = append(args, o.{{titleCase $value}})
	{{ end -}}

	return args
}

func (o {{$tableNameSingular}}Slice) inPrimaryKeyArgs() []interface{} {
	var args []interface{}

	for i := 0; i < len(o); i++ {
		{{- range $key, $value := .Table.PKey.Columns }}
		args = append(args, o[i].{{titleCase $value}})
		{{ end -}}
	}

	return args
}
