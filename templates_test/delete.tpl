{{- $tableNameSingular := .Table.Name | singular | titleCase -}}
{{- $tableNamePlural := .Table.Name | plural | titleCase -}}
{{- $varNamePlural := .Table.Name | plural | camelCase -}}
{{- $varNameSingular := .Table.Name | singular | camelCase -}}
func Test{{$tableNamePlural}}Delete(t *testing.T) {
  t.Parallel()

  tx, err := boil.Begin()
  if err != nil {
    t.Fatal(err)
  }
  defer tx.Rollback()

  {{$varNameSingular}} := &{{$tableNameSingular}}{}
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

  tx, err := boil.Begin()
  if err != nil {
    t.Fatal(err)
  }
  defer tx.Rollback()

  {{$varNameSingular}} := &{{$tableNameSingular}}{}
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

  tx, err := boil.Begin()
  if err != nil {
    t.Fatal(err)
  }
  defer tx.Rollback()

  {{$varNameSingular}} := &{{$tableNameSingular}}{}
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
