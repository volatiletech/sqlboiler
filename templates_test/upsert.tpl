{{- $tableNameSingular := .Table.Name | singular | titleCase -}}
{{- $tableNamePlural := .Table.Name | plural | titleCase -}}
{{- $varNamePlural := .Table.Name | plural | camelCase -}}
{{- $varNameSingular := .Table.Name | singular | camelCase -}}
func Test{{$tableNamePlural}}Upsert(t *testing.T) {
  var err error

  o := {{$tableNameSingular}}{}

  columns := o.generateUpsertColumns([]string{"one", "two"}, []string{"three", "four"}, []string{"five", "six"})
  if columns.conflict[0] != "one" || columns.conflict[1] != "two" {
    t.Errorf("Expected conflict to be %v, got %v", []string{"one", "two"}, columns.conflict)
  }

  if columns.update[0] != "three" || columns.update[1] != "four" {
    t.Errorf("Expected update to be %v, got %v", []string{"three", "four"}, columns.update)
  }

  if columns.whitelist[0] != "five" || columns.whitelist[1] != "six" {
    t.Errorf("Expected whitelist to be %v, got %v", []string{"five", "six"}, columns.whitelist)
  }

  columns = o.generateUpsertColumns(nil, nil, nil)
  if len(columns.whitelist) == 0 {
    t.Errorf("Expected whitelist to contain columns, but got len 0")
  }

  if len(columns.conflict) == 0 {
    t.Errorf("Expected conflict to contain columns, but got len 0")
  }

  if len(columns.update) == 0 {
    t.Errorf("expected update to contain columns, but got len 0")
  }

  upsertCols := upsertData{
    conflict: []string{"key1", `"key2"`},
    update: []string{"aaa", `"bbb"`},
    whitelist: []string{"thing", `"stuff"`},
    returning: []string{},
  }

  query := o.generateUpsertQuery(false, upsertCols)
  expectedQuery := `INSERT INTO {{.Table.Name}} ("thing", "stuff") VALUES ($1, $2) ON CONFLICT DO NOTHING`

  if query != expectedQuery {
    t.Errorf("Expected query mismatch:\n\n%s\n%s\n", query, expectedQuery)
  }

  query = o.generateUpsertQuery(true, upsertCols)
  expectedQuery = `INSERT INTO {{.Table.Name}} ("thing", "stuff") VALUES ($1, $2) ON CONFLICT ("key1", "key2") DO UPDATE SET "aaa" = EXCLUDED."aaa", "bbb" = EXCLUDED."bbb"`

  if query != expectedQuery {
    t.Errorf("Expected query mismatch:\n\n%s\n%s\n", query, expectedQuery)
  }

  upsertCols.returning = []string{"stuff"}
  query = o.generateUpsertQuery(true, upsertCols)
  expectedQuery = expectedQuery + ` RETURNING "stuff"`

  if query != expectedQuery {
    t.Errorf("Expected query mismatch:\n\n%s\n%s\n", query, expectedQuery)
  }

  // Attempt the INSERT side of an UPSERT
  if err = boil.RandomizeStruct(&o, {{$varNameSingular}}DBTypes, true); err != nil {
    t.Errorf("Unable to randomize {{$tableNameSingular}} struct: %s", err)
  }

  if err = o.UpsertG(false, nil, nil); err != nil {
    t.Errorf("Unable to upsert {{$tableNameSingular}}: %s", err)
  }

  compare, err := {{$tableNameSingular}}FindG({{.Table.PKey.Columns | stringMap .StringFuncs.titleCase | prefixStringSlice "o." | join ", "}})
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

  if err = o.UpsertG(true, nil, nil); err != nil {
    t.Errorf("Unable to upsert {{$tableNameSingular}}: %s", err)
  }

  compare, err = {{$tableNameSingular}}FindG({{.Table.PKey.Columns | stringMap .StringFuncs.titleCase | prefixStringSlice "o." | join ", "}})
  if err != nil {
    t.Errorf("Unable to find {{$tableNameSingular}}: %s", err)
  }
  err = {{$varNameSingular}}CompareVals(&o, compare, true); if err != nil {
    t.Error(err)
  }

  {{$varNamePlural}}DeleteAllRows(t)
}
