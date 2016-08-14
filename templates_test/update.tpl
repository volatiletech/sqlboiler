{{- $tableNameSingular := .Table.Name | singular | titleCase -}}
{{- $tableNamePlural := .Table.Name | plural | titleCase -}}
{{- $varNamePlural := .Table.Name | plural | camelCase -}}
{{- $varNameSingular := .Table.Name | singular | camelCase -}}
func Test{{$tableNamePlural}}Update(t *testing.T) {
  t.Parallel()

  {{$varNameSingular}} := &{{$tableNameSingular}}{}
  if err := boil.RandomizeStruct({{$varNameSingular}}, {{$varNameSingular}}DBTypes, true); err != nil {
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

  if err = boil.RandomizeStruct({{$varNameSingular}}, {{$varNameSingular}}DBTypes, true, {{$varNameSingular}}PrimaryKeyColumns...); err != nil {
    t.Errorf("Unable to randomize {{$tableNameSingular}} slice: %s", err)
  }

  if err = {{$varNameSingular}}.Update(tx); err != nil {
    t.Error(err)
  }
}

func Test{{$tableNamePlural}}SliceUpdateAll(t *testing.T) {
  t.Parallel()

  {{$varNameSingular}} := &{{$tableNameSingular}}{}
  if err := boil.RandomizeStruct({{$varNameSingular}}, {{$varNameSingular}}DBTypes, true); err != nil {
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

  if err = boil.RandomizeStruct({{$varNameSingular}}, {{$varNameSingular}}DBTypes, true, {{$varNameSingular}}PrimaryKeyColumns...); err != nil {
    t.Errorf("Unable to randomize {{$tableNameSingular}} slice: %s", err)
  }

  // Remove Primary keys and unique columns from what we plan to update
  fields := strmangle.SetComplement(
    {{$varNameSingular}}Columns,
    {{$varNameSingular}}PrimaryKeyColumns,
  )

	value := reflect.Indirect(reflect.ValueOf({{$varNameSingular}}))
  updateMap := M{}
  for _, col := range fields {
    updateMap[col] = value.FieldByName(strmangle.TitleCase(col)).Interface()
  }

  slice := {{$tableNameSingular}}Slice{{"{"}}{{$varNameSingular}}{{"}"}}
  if err = slice.UpdateAll(tx, updateMap); err != nil {
    t.Error(err)
  }
}
