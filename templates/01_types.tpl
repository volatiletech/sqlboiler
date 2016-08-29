{{- $varNameSingular := .Table.Name | singular | camelCase -}}
{{- $tableNameSingular := .Table.Name | singular | titleCase -}}
var (
  {{$varNameSingular}}Columns                   = []string{{"{"}}{{.Table.Columns | columnNames | stringMap .StringFuncs.quoteWrap | join ", "}}{{"}"}}
  {{$varNameSingular}}ColumnsWithoutDefault     = []string{{"{"}}{{.Table.Columns | filterColumnsByDefault false | columnNames | stringMap .StringFuncs.quoteWrap | join ","}}{{"}"}}
  {{$varNameSingular}}ColumnsWithDefault        = []string{{"{"}}{{.Table.Columns | filterColumnsByDefault true | columnNames | stringMap .StringFuncs.quoteWrap | join ","}}{{"}"}}
  {{$varNameSingular}}PrimaryKeyColumns         = []string{{"{"}}{{.Table.PKey.Columns | stringMap .StringFuncs.quoteWrap | join ", "}}{{"}"}}
  {{$varNameSingular}}TitleCases                = map[string]string{
    {{range $col := .Table.Columns | columnNames -}}
    "{{$col}}": "{{titleCase $col}}",
    {{end -}}
  }
)

type (
  {{$tableNameSingular}}Slice []*{{$tableNameSingular}}
  {{if eq .NoHooks false -}}
  {{$tableNameSingular}}Hook func(boil.Executor, *{{$tableNameSingular}}) error
  {{- end}}

  {{$varNameSingular}}Query struct {
    *boil.Query
  }
)

// Force time package dependency for automated UpdatedAt/CreatedAt.
var _ = time.Second
