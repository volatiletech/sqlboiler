{{- $varNameSingular := .Table.Name | singular | camelCase -}}
{{- $tableNameSingular := .Table.Name | singular | titleCase -}}
var (
  {{$varNameSingular}}Columns = []string{{"{"}}{{.Table.Columns | columnNames | join ", "}}{{"}"}}
  {{$varNameSingular}}ColumnsWithoutDefault = []string{{"{"}}{{filterColumnsByDefault .Table.Columns false}}{{"}"}}
  {{$varNameSingular}}ColumnsWithDefault = []string{{"{"}}{{filterColumnsByDefault .Table.Columns true}}{{"}"}}
  {{$varNameSingular}}PrimaryKeyColumns = []string{{"{"}}{{.Table.PKey.Columns | join ", "}}{{"}"}}
  {{$varNameSingular}}AutoIncrementColumns = []string{{"{"}}{{filterColumnsByAutoIncrement .Table.Columns}}{{"}"}}
  {{$varNameSingular}}AutoIncPrimaryKey = "{{autoIncPrimaryKey .Table.Columns .Table.PKey}}"
)

type (
  {{$varNameSingular}}Slice []*{{$tableNameSingular}}
  {{$tableNameSingular}}Hook func(*{{$tableNameSingular}}) error
  {{$varNameSingular}}Query struct {
    *boil.Query
  }
)
