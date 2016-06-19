{{- $varNameSingular := camelCaseSingular .Table.Name -}}
{{- if hasPrimaryKey .Table.PKey -}}
{{- $tableNameSingular := titleCaseSingular .Table.Name -}}
func (o {{$tableNameSingular}}) inPrimaryKeyArgs() []interface{} {
  var args []interface{}

  {{- range $key, $value := .Table.PKey.Columns }}
  args = append(args, o.{{titleCase $value}})
  {{ end -}}

  return args
}

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
