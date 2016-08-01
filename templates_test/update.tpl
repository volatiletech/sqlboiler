{{- $tableNameSingular := .Table.Name | singular | titleCase -}}
{{- $dbName := singular .Table.Name -}}
{{- $tableNamePlural := .Table.Name | plural | titleCase -}}
{{- $varNamePlural := .Table.Name | plural | camelCase -}}
{{- $varNameSingular := .Table.Name | singular | camelCase -}}
func Test{{$tableNamePlural}}Update(t *testing.T) {
  var err error

  item := {{$tableNameSingular}}{}
  if err = item.InsertG(); err != nil {
    t.Errorf("Unable to insert zero-value item {{$tableNameSingular}}:\n%#v\nErr: %s", item, err)
  }

  blacklistCols := boil.SetMerge({{$varNameSingular}}AutoIncrementColumns, {{$varNameSingular}}PrimaryKeyColumns)
  if err = boil.RandomizeStruct(&item, {{$varNameSingular}}DBTypes, false, blacklistCols...); err != nil {
    t.Errorf("Unable to randomize {{$tableNameSingular}} struct: %s", err)
  }

  whitelist := boil.SetComplement({{$varNameSingular}}Columns, {{$varNameSingular}}AutoIncrementColumns)
  if err = item.UpdateG(whitelist...); err != nil {
    t.Errorf("Unable to update {{$tableNameSingular}}: %s", err)
  }

  var j *{{$tableNameSingular}}
  j, err = {{$tableNameSingular}}FindG({{.Table.PKey.Columns | stringMap .StringFuncs.titleCase | prefixStringSlice "item." | join ", "}})
  if err != nil {
    t.Errorf("Unable to find {{$tableNameSingular}} row: %s", err)
  }

  {{$varNameSingular}}CompareVals(&item, j, t)

  wl := item.generateUpdateColumns("test")
  if len(wl) != 1 && wl[0] != "test" {
    t.Errorf("Expected generateUpdateColumns whitelist to match expected whitelist")
  }

  wl = item.generateUpdateColumns()
  if len(wl) == 0 && len({{$varNameSingular}}ColumnsWithoutDefault) > 0 {
    t.Errorf("Expected generateUpdateColumns to build a whitelist for {{$tableNameSingular}}, but got 0 results")
  }

  {{$varNamePlural}}DeleteAllRows(t)
}
