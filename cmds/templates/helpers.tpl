{{- $varNameSingular := camelCaseSingular .Table.Name -}}
var {{$varNameSingular}}Columns = []string{{"{"}}{{columnsToStrings .Table.Columns | commaList}}{{"}"}}
var {{$varNameSingular}}ColumnsWithoutDefault = []string{{"{"}}{{filterColumnsByDefault .Table.Columns false}}{{"}"}}
var {{$varNameSingular}}ColumnsWithDefault = []string{{"{"}}{{filterColumnsByDefault .Table.Columns true}}{{"}"}}
var {{$varNameSingular}}PrimaryKeyColumns = []string{{"{"}}{{commaList .Table.PKey.Columns}}{{"}"}}
var {{$varNameSingular}}AutoIncrementColumns = []string{{"{"}}{{filterColumnsByAutoIncrement .Table.Columns}}{{"}"}}
var {{$varNameSingular}}AutoIncPrimaryKey = "{{autoIncPrimaryKey .Table.Columns .Table.PKey}}"

{{if hasPrimaryKey .Table.PKey -}}
{{- $tableNameSingular := titleCaseSingular .Table.Name -}}
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
