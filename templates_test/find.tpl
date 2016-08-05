{{- $tableNameSingular := .Table.Name | singular | titleCase -}}
{{- $tableNamePlural := .Table.Name | plural | titleCase -}}
{{- $varNamePlural := .Table.Name | plural | camelCase -}}
{{- $varNameSingular := .Table.Name | singular | camelCase -}}
func Test{{$tableNamePlural}}Find(t *testing.T) {
  var err error

  o := make({{$tableNameSingular}}Slice, 3)
  if err = boil.RandomizeSlice(&o, {{$varNameSingular}}DBTypes, true); err != nil {
    t.Errorf("Unable to randomize {{$tableNameSingular}} slice: %s", err)
  }

  for i := 0; i < len(o); i++ {
    if err = o[i].InsertG(); err != nil {
      t.Errorf("Unable to insert {{$tableNameSingular}}:\n%#v\nErr: %s", o[i], err)
    }
  }

  j := make({{$tableNameSingular}}Slice, 3)
  // Perform all Find queries and assign result objects to slice for comparison
  for i := 0; i < len(j); i++ {
    j[i], err = {{$tableNameSingular}}FindG({{.Table.PKey.Columns | stringMap .StringFuncs.titleCase | prefixStringSlice "o[i]." | join ", "}})
    {{$varNameSingular}}CompareVals(o[i], j[i], t)
  }

  f, err := {{$tableNameSingular}}FindG({{.Table.PKey.Columns | stringMap .StringFuncs.titleCase | prefixStringSlice "o[0]." | join ", "}}, {{$varNameSingular}}PrimaryKeyColumns...)
  {{range $key, $value := .Table.PKey.Columns}}
  if o[0].{{titleCase $value}} != f.{{titleCase $value}} {
    t.Errorf("Expected primary key values to match, {{titleCase $value}} did not match")
  }
  {{end}}

  colsWithoutPrimKeys := boil.SetComplement({{$varNameSingular}}Columns, {{$varNameSingular}}PrimaryKeyColumns)
  fRef := reflect.ValueOf(f).Elem()
  for _, v := range colsWithoutPrimKeys {
    val := fRef.FieldByName(v)
    if val.IsValid() {
      t.Errorf("Expected all other columns to be zero value, but column %s was %#v", v, val.Interface())
    }
  }

  {{$varNamePlural}}DeleteAllRows(t)
}
