{{- $tableNameSingular := .Table.Name | singular | titleCase -}}
{{- $dbName := singular .Table.Name -}}
// {{$tableNameSingular}} is an object representing the database table.
type {{$tableNameSingular}} struct {
  {{range $key, $value := .Table.Columns -}}
  {{titleCase $value.Name}} {{$value.Type}} `db:"{{makeDBName $dbName $value.Name}}" json:"{{$value.Name}}"`
  {{end -}}
}
