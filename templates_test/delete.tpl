{{- $tableNameSingular := .Table.Name | singular | titleCase -}}
{{- $tableNamePlural := .Table.Name | plural | titleCase -}}
{{- $varNamePlural := .Table.Name | plural | camelCase -}}
{{- $varNameSingular := .Table.Name | singular | camelCase -}}
func Test{{$tableNamePlural}}Delete(t *testing.T) {
  t.Parallel()

  seed := new(boil.Seed)
  var err error
  {{$varNameSingular}} := &{{$tableNameSingular}}{}
  if err = seed.RandomizeStruct({{$varNameSingular}}, {{$varNameSingular}}DBTypes, true); err != nil {
    t.Errorf("Unable to randomize {{$tableNameSingular}} struct: %s", err)
  }

  tx := MustTx(boil.Begin())
  defer tx.Rollback()
  if err = {{$varNameSingular}}.Insert(tx); err != nil {
    t.Error(err)
  }

  if err = {{$varNameSingular}}.Delete(tx); err != nil {
    t.Error(err)
  }

  count, err := {{$tableNamePlural}}(tx).Count()
  if err != nil {
    t.Error(err)
  }

  if count != 0 {
    t.Error("want zero records, got:", count)
  }
}

func Test{{$tableNamePlural}}QueryDeleteAll(t *testing.T) {
  t.Parallel()

  seed := new(boil.Seed)
  var err error
  {{$varNameSingular}} := &{{$tableNameSingular}}{}
  if err = seed.RandomizeStruct({{$varNameSingular}}, {{$varNameSingular}}DBTypes, true); err != nil {
    t.Errorf("Unable to randomize {{$tableNameSingular}} struct: %s", err)
  }

  tx := MustTx(boil.Begin())
  defer tx.Rollback()
  if err = {{$varNameSingular}}.Insert(tx); err != nil {
    t.Error(err)
  }

  if err = {{$tableNamePlural}}(tx).DeleteAll(); err != nil {
    t.Error(err)
  }

  count, err := {{$tableNamePlural}}(tx).Count()
  if err != nil {
    t.Error(err)
  }

  if count != 0 {
    t.Error("want zero records, got:", count)
  }
}

func Test{{$tableNamePlural}}SliceDeleteAll(t *testing.T) {
  t.Parallel()

  seed := new(boil.Seed)
  var err error
  {{$varNameSingular}} := &{{$tableNameSingular}}{}
  if err = seed.RandomizeStruct({{$varNameSingular}}, {{$varNameSingular}}DBTypes, true); err != nil {
    t.Errorf("Unable to randomize {{$tableNameSingular}} struct: %s", err)
  }

  tx := MustTx(boil.Begin())
  defer tx.Rollback()
  if err = {{$varNameSingular}}.Insert(tx); err != nil {
    t.Error(err)
  }

  slice := {{$tableNameSingular}}Slice{{"{"}}{{$varNameSingular}}{{"}"}}

  if err = slice.DeleteAll(tx); err != nil {
    t.Error(err)
  }

  count, err := {{$tableNamePlural}}(tx).Count()
  if err != nil {
    t.Error(err)
  }

  if count != 0 {
    t.Error("want zero records, got:", count)
  }
}
