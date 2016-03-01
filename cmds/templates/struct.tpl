{{- $tableName := .TableName -}}
// {{makeGoName $tableName}} is an object representing the database table.
type {{makeGoName $tableName}} struct {
  {{range $key, $value := .TableData -}}
  {{makeGoName $value.Name}} {{$value.Type}} `db:"{{makeDBName $tableName $value.Name}}" json:"{{$value.Name}}"`
  {{end -}}
}
