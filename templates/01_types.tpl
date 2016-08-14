{{- $varNameSingular := .Table.Name | singular | camelCase -}}
{{- $tableNameSingular := .Table.Name | singular | titleCase -}}
var (
  {{$varNameSingular}}Columns                   = []string{{"{"}}{{.Table.Columns | columnNames | stringMap .StringFuncs.quoteWrap | join ", "}}{{"}"}}
  {{$varNameSingular}}ColumnsWithoutDefault     = []string{{"{"}}{{.Table.Columns | filterColumnsByDefault false | columnNames | stringMap .StringFuncs.quoteWrap | join ","}}{{"}"}}
  {{$varNameSingular}}ColumnsWithDefault        = []string{{"{"}}{{.Table.Columns | filterColumnsByDefault true | columnNames | stringMap .StringFuncs.quoteWrap | join ","}}{{"}"}}
  {{$varNameSingular}}PrimaryKeyColumns         = []string{{"{"}}{{.Table.PKey.Columns | stringMap .StringFuncs.quoteWrap | join ", "}}{{"}"}}
)

type (
  {{$tableNameSingular}}Slice []*{{$tableNameSingular}}
  {{$tableNameSingular}}Hook func(*{{$tableNameSingular}}) error

  {{$varNameSingular}}Query struct {
    *boil.Query
  }
)
