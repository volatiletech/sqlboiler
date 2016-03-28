{{- $tableNameSingular := titleCaseSingular .Table.Name -}}
{{- $dbName := singular .Table.Name -}}
{{- $tableNamePlural := titleCasePlural .Table.Name -}}
{{- $varNamePlural := camelCasePlural .Table.Name -}}
// {{$tableNamePlural}}All retrieves all records.
func Test{{$tableNamePlural}}All(t *testing.T) {

}
