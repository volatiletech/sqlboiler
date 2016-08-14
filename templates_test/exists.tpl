{{- $tableNameSingular := .Table.Name | singular | titleCase -}}
{{- $tableNamePlural := .Table.Name | plural | titleCase -}}
{{- $varNamePlural := .Table.Name | plural | camelCase -}}
{{- $varNameSingular := .Table.Name | singular | camelCase -}}
func Test{{$tableNamePlural}}Exists(t *testing.T) {
  t.Parallel()

  {{$varNameSingular}} := &{{$tableNameSingular}}{}
  if err := boil.RandomizeStruct({{$varNameSingular}}, {{$varNameSingular}}DBTypes, true, {{$varNameSingular}}ColumnsWithDefault...); err != nil {
    t.Errorf("Unable to randomize {{$tableNameSingular}} slice: %s", err)
  }

  tx, err := boil.Begin()
  if err != nil {
    t.Fatal(err)
  }
  defer tx.Rollback()

  if err = {{$varNameSingular}}.Insert(tx); err != nil {
    t.Error(err)
  }

  {{$pkeyArgs := .Table.PKey.Columns | stringMap .StringFuncs.titleCase | prefixStringSlice (printf "%s." $varNameSingular) | join ", " -}}
  e, err := {{$tableNameSingular}}Exists(tx, {{$pkeyArgs}})
  if err != nil {
    t.Errorf("Unable to check if {{$tableNameSingular}} exists: %s", err)
  }
  if e != true {
    t.Errorf("Expected {{$tableNameSingular}}ExistsG to return true, but got false.")
  }
}
