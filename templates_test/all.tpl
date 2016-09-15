{{- $tableNameSingular := .Table.Name | singular | titleCase -}}
{{- $tableNamePlural := .Table.Name | plural | titleCase -}}
{{- $varNamePlural := .Table.Name | plural | camelCase -}}
{{- $varNameSingular := .Table.Name | singular | camelCase -}}
func test{{$tableNamePlural}}(t *testing.T) {
	t.Parallel()

	query := {{$tableNamePlural}}(nil)

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}
