{{- $tableNameSingular := .Table.Name | singular | titleCase -}}
{{- $tableNamePlural := .Table.Name | plural | titleCase -}}
{{- $varNamePlural := .Table.Name | plural | camelCase -}}
{{- $varNameSingular := .Table.Name | singular | camelCase -}}
func Test{{$tableNamePlural}}Find(t *testing.T) {
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

  {{$varNameSingular}}Found, err := {{$tableNameSingular}}Find(tx, {{.Table.PKey.Columns | stringMap .StringFuncs.titleCase | prefixStringSlice (printf "%s." $varNameSingular) | join ", "}})
  if err != nil {
    t.Error(err)
  }

  if {{$varNameSingular}}Found == nil {
    t.Error("want a record, got nil")
  }
}
