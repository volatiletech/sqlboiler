{{- $tableName := .TableName -}}
// {{titleCase $tableName}} is an object representing the database table.
type {{titleCase $tableName}} struct {
  {{range $key, $value := .TableData -}}
  {{titleCase $value.Name}} {{$value.Type}} `db:"{{makeDBName $tableName $value.Name}}" json:"{{$value.Name}}"`
  {{end -}}
}
