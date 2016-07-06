{{- $tableNameSingular := .Table.Name | singular | titleCase -}}
{{- $dbName := singular .Table.Name -}}
{{- $tableNamePlural := .Table.Name | plural | titleCase -}}
{{- $varNamePlural := .Table.Name | plural | camelCase -}}
{{- $varNameSingular := .Table.Name | singular | camelCase -}}
func Test{{$tableNamePlural}}Update(t *testing.T) {
  t.Errorf("test update not implemented")
}
