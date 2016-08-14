{{- $tableNameSingular := .Table.Name | singular | titleCase -}}
{{- $tableNamePlural := .Table.Name | plural | titleCase -}}
{{- $varNamePlural := .Table.Name | plural | camelCase -}}
{{- $varNameSingular := .Table.Name | singular | camelCase -}}
func Test{{$tableNamePlural}}Upsert(t *testing.T) {
  var err error

  o := {{$tableNameSingular}}{}

  // Attempt the INSERT side of an UPSERT
  if err = boil.RandomizeStruct(&o, {{$varNameSingular}}DBTypes, true); err != nil {
    t.Errorf("Unable to randomize {{$tableNameSingular}} struct: %s", err)
  }

  tx, err := boil.Begin()
  if err != nil {
    t.Fatal(err)
  }
  defer tx.Rollback()

  if err = o.Upsert(tx, false, nil, nil); err != nil {
    t.Errorf("Unable to upsert {{$tableNameSingular}}: %s", err)
  }

  compare, err := {{$tableNameSingular}}Find(tx, {{.Table.PKey.Columns | stringMap .StringFuncs.titleCase | prefixStringSlice "o." | join ", "}})
  if err != nil {
    t.Errorf("Unable to find {{$tableNameSingular}}: %s", err)
  }
  err = {{$varNameSingular}}CompareVals(&o, compare, true); if err != nil {
    t.Error(err)
  }

  // Attempt the UPDATE side of an UPSERT
  if err = boil.RandomizeStruct(&o, {{$varNameSingular}}DBTypes, false, {{$varNameSingular}}PrimaryKeyColumns...); err != nil {
    t.Errorf("Unable to randomize {{$tableNameSingular}} struct: %s", err)
  }

  if err = o.Upsert(tx, true, nil, nil); err != nil {
    t.Errorf("Unable to upsert {{$tableNameSingular}}: %s", err)
  }

  compare, err = {{$tableNameSingular}}Find(tx, {{.Table.PKey.Columns | stringMap .StringFuncs.titleCase | prefixStringSlice "o." | join ", "}})
  if err != nil {
    t.Errorf("Unable to find {{$tableNameSingular}}: %s", err)
  }
  err = {{$varNameSingular}}CompareVals(&o, compare, true); if err != nil {
    t.Error(err)
  }
}
