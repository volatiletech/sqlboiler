{{- $tableNameSingular := .Table.Name | singular | titleCase -}}
{{- $tableNamePlural := .Table.Name | plural | titleCase -}}
{{- $varNamePlural := .Table.Name | plural | camelCase -}}
{{- $varNameSingular := .Table.Name | singular | camelCase -}}
func test{{$tableNamePlural}}Upsert(t *testing.T) {
  t.Parallel()

  seed := boil.NewSeed()
  var err error
  // Attempt the INSERT side of an UPSERT
  {{$varNameSingular}} := {{$tableNameSingular}}{}
  if err = seed.RandomizeStruct(&{{$varNameSingular}}, {{$varNameSingular}}DBTypes, true); err != nil {
    t.Errorf("Unable to randomize {{$tableNameSingular}} struct: %s", err)
  }

  tx := MustTx(boil.Begin())
  defer tx.Rollback()
  if err = {{$varNameSingular}}.Upsert(tx, false, nil, nil); err != nil {
    t.Errorf("Unable to upsert {{$tableNameSingular}}: %s", err)
  }

  count, err := {{$tableNamePlural}}(tx).Count()
  if err != nil {
    t.Error(err)
  }
  if count != 1 {
    t.Error("want one record, got:", count)
  }

  // Attempt the UPDATE side of an UPSERT
  if err = seed.RandomizeStruct(&{{$varNameSingular}}, {{$varNameSingular}}DBTypes, false, {{$varNameSingular}}PrimaryKeyColumns...); err != nil {
    t.Errorf("Unable to randomize {{$tableNameSingular}} struct: %s", err)
  }

  if err = {{$varNameSingular}}.Upsert(tx, true, nil, nil); err != nil {
    t.Errorf("Unable to upsert {{$tableNameSingular}}: %s", err)
  }

  count, err = {{$tableNamePlural}}(tx).Count()
  if err != nil {
    t.Error(err)
  }
  if count != 1 {
    t.Error("want one record, got:", count)
  }
}
