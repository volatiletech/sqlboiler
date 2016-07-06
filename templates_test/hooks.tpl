{{- $tableNameSingular := .Table.Name | singular | titleCase -}}
{{- $dbName := singular .Table.Name -}}
{{- $tableNamePlural := .Table.Name | plural | titleCase -}}
{{- $varNamePlural := .Table.Name | plural | camelCase -}}
{{- $varNameSingular := .Table.Name | singular | camelCase -}}
func {{$varNameSingular}}BeforeCreateHook(o *{{$tableNameSingular}}) error {
  return nil
}

func {{$varNameSingular}}AfterCreateHook(o *{{$tableNameSingular}}) error {
  return nil
}

func {{$varNameSingular}}BeforeUpdateHook(o *{{$tableNameSingular}}) error {
  return nil
}

func {{$varNameSingular}}AfterUpdateHook(o *{{$tableNameSingular}}) error {
  return nil
}

func Test{{$tableNamePlural}}Hooks(t *testing.T) {
  // var err error

  {{$varNamePlural}}DeleteAllRows(t)
  t.Errorf("Hook tests not implemented")
}
