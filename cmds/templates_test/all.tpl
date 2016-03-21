{{- $tableNameSingular := titleCaseSingular .Table -}}
{{- $dbName := singular .Table -}}
{{- $tableNamePlural := titleCasePlural .Table -}}
{{- $varNamePlural := camelCasePlural .Table -}}
// {{$tableNamePlural}}All retrieves all records.
func Test{{$tableNamePlural}}All(t *testing.T) {

}
