{{- $varNameSingular := camelCaseSingular .Table.Name -}}
{{- $tableNameSingular := titleCaseSingular .Table.Name -}}
var (
  {{$varNameSingular}}Columns = []string{{"{"}}{{columnsToStrings .Table.Columns | commaList}}{{"}"}}
  {{$varNameSingular}}ColumnsWithoutDefault = []string{{"{"}}{{filterColumnsByDefault .Table.Columns false}}{{"}"}}
  {{$varNameSingular}}ColumnsWithDefault = []string{{"{"}}{{filterColumnsByDefault .Table.Columns true}}{{"}"}}
  {{$varNameSingular}}PrimaryKeyColumns = []string{{"{"}}{{commaList .Table.PKey.Columns}}{{"}"}}
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
