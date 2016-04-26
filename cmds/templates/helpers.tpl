{{if hasPrimaryKey .Table.PKey -}}
{{- $tableNameSingular := titleCaseSingular .Table.Name -}}
{{- $varNameSingular := camelCaseSingular .Table.Name -}}
func (o {{$varNameSingular}}Slice) inPrimaryKeyArgs() []interface{} {
  var args []interface{}

  for i := 0; i < len(o); i++ {
    {{- range $key, $value := .Table.PKey.Columns }}
    args = append(args, o[i].{{titleCase $value}})
    {{ end -}}
  }

  return args
}
{{- end}}
