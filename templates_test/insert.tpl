{{- $tableNameSingular := .Table.Name | singular | titleCase -}}
{{- $tableNamePlural := .Table.Name | plural | titleCase -}}
{{- $varNamePlural := .Table.Name | plural | camelCase -}}
{{- $varNameSingular := .Table.Name | singular | camelCase -}}
{{- $parent := . -}}
func Test{{$tableNamePlural}}Insert(t *testing.T) {
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

  count, err := {{$tableNamePlural}}(tx).Count()
  if err != nil {
    t.Error(err)
  }

  if count != 1 {
    t.Error("want one record, got:", count)
  }
}

func Test{{$tableNamePlural}}InsertWhitelist(t *testing.T) {
  {{$varNameSingular}} := &{{$tableNameSingular}}{}
  if err := boil.RandomizeStruct({{$varNameSingular}}, {{$varNameSingular}}DBTypes, true); err != nil {
    t.Errorf("Unable to randomize {{$tableNameSingular}} slice: %s", err)
  }

  tx, err := boil.Begin()
  if err != nil {
    t.Fatal(err)
  }
  defer tx.Rollback()

  if err = {{$varNameSingular}}.Insert(tx, {{$varNameSingular}}Columns...); err != nil {
    t.Error(err)
  }

  count, err := {{$tableNamePlural}}(tx).Count()
  if err != nil {
    t.Error(err)
  }

  if count != 1 {
    t.Error("want one record, got:", count)
  }
}
