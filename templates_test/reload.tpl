{{- $tableNameSingular := .Table.Name | singular | titleCase -}}
{{- $tableNamePlural := .Table.Name | plural | titleCase -}}
{{- $varNamePlural := .Table.Name | plural | camelCase -}}
{{- $varNameSingular := .Table.Name | singular | camelCase -}}
func Test{{$tableNamePlural}}Reload(t *testing.T) {
  var err error

  o := {{$tableNameSingular}}{}
  if err = boil.RandomizeStruct(&o, {{$varNameSingular}}DBTypes, true); err != nil {
    t.Errorf("Unable to randomize {{$tableNameSingular}} struct: %s", err)
  }

  if err = o.InsertG(); err != nil {
    t.Errorf("Unable to insert {{$tableNameSingular}}:\n%#v\nErr: %s", o, err)
  }

  // Create another copy of the object
  o1, err := {{$tableNameSingular}}FindG({{.Table.PKey.Columns | stringMap .StringFuncs.titleCase | prefixStringSlice "o." | join ", "}})
  if err != nil {
    t.Errorf("Unable to find {{$tableNameSingular}} row.")
  }

  // Randomize the struct values again, except for the primary key values, so we can call update.
  err = boil.RandomizeStruct(&o, {{$varNameSingular}}DBTypes, true, {{$varNameSingular}}PrimaryKeyColumns...)
  if err != nil {
    t.Errorf("Unable to randomize {{$tableNameSingular}} struct members excluding primary keys: %s", err)
  }

  colsWithoutPrimKeys := boil.SetComplement({{$varNameSingular}}Columns, {{$varNameSingular}}PrimaryKeyColumns)

  if err = o.UpdateG(colsWithoutPrimKeys...); err != nil {
    t.Errorf("Unable to update the {{$tableNameSingular}} row: %s", err)
  }

  if err = o1.ReloadG(); err != nil {
    t.Errorf("Unable to reload {{$tableNameSingular}} object: %s", err)
  }

  {{$varNameSingular}}CompareVals(&o, o1, t)

  {{$varNamePlural}}DeleteAllRows(t)
}

func Test{{$tableNamePlural}}ReloadAll(t *testing.T) {
  var err error

  o1 := make({{$tableNameSingular}}Slice, 3)
  o2 := make({{$tableNameSingular}}Slice, 3)
  if err = boil.RandomizeSlice(&o1, {{$varNameSingular}}DBTypes, false); err != nil {
    t.Errorf("Unable to randomize {{$tableNameSingular}} slice: %s", err)
  }

  for i := 0; i < len(o1); i++ {
    if err = o1[i].InsertG(); err != nil {
      t.Errorf("Unable to insert {{$tableNameSingular}}:\n%#v\nErr: %s", o1[i], err)
    }
  }

  for i := 0; i < len(o1); i++ {
    o2[i], err = {{$tableNameSingular}}FindG({{.Table.PKey.Columns | stringMap .StringFuncs.titleCase | prefixStringSlice "o1[i]." | join ", "}})
    if err != nil {
      t.Errorf("Unable to find {{$tableNameSingular}} row.")
    }

    {{$varNameSingular}}CompareVals(o1[i], o2[i], t)
  }

  // Randomize the struct values again, except for the primary key values, so we can call update.
  err = boil.RandomizeSlice(&o1, {{$varNameSingular}}DBTypes, false, {{$varNameSingular}}PrimaryKeyColumns...)
  if err != nil {
    t.Errorf("Unable to randomize {{$tableNameSingular}} slice excluding primary keys: %s", err)
  }

  colsWithoutPrimKeys := boil.SetComplement({{$varNameSingular}}Columns, {{$varNameSingular}}PrimaryKeyColumns)

  for i := 0; i < len(o1); i++ {
    if err = o1[i].UpdateG(colsWithoutPrimKeys...); err != nil {
      t.Errorf("Unable to update the {{$tableNameSingular}} row: %s", err)
    }
  }

  if err = o2.ReloadAllG(); err != nil {
    t.Errorf("Unable to reload {{$tableNameSingular}} object: %s", err)
  }

  for i := 0; i < len(o1); i++  {
    {{$varNameSingular}}CompareVals(o2[i], o1[i], t)
  }

  {{$varNamePlural}}DeleteAllRows(t)
}
