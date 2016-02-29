{{- $tableName := .TableName -}}
// {{makeGoColName $tableName}} is an object representing the database table.
type {{makeGoColName $tableName}} struct {
  {{range $key, $value := .TableData -}}
  {{makeGoColName $value.ColName}} {{$value.ColType}} `db:"{{makeDBColName $tableName $value.ColName}}" json:"{{$value.ColName}}"`
  {{end -}}
}
