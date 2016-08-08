{{- $tableNameSingular := .Table.Name | singular -}}
{{- $modelName := $tableNameSingular | titleCase -}}
// {{$modelName}} is an object representing the database table.
type {{$modelName}} struct {
  {{range $column := .Table.Columns -}}
  {{titleCase $column.Name}} {{$column.Type}} `boil:"{{$column.Name}}" json:"{{$column.Name}}" toml:"{{$column.Name}}" yaml:"{{$column.Name}}"`
  {{end -}}
}
