{{- $tableNameSingular := titleCaseSingular .Table -}}
{{- $dbName := singular .Table -}}
// {{$tableNameSingular}} is an object representing the database table.
type {{$tableNameSingular}} struct {
  {{range $key, $value := .Columns -}}
  {{titleCase $value.Name}} {{$value.Type}} `db:"{{makeDBName $dbName $value.Name}}" json:"{{$value.Name}}"`
  {{end -}}
}
