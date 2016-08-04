{{- $tableNameSingular := .Table.Name | singular | titleCase -}}
{{- $dbName := singular .Table.Name -}}
{{- $tableNamePlural := .Table.Name | plural | titleCase -}}
{{- $varNamePlural := .Table.Name | plural | camelCase -}}
{{- $varNameSingular := .Table.Name | singular | camelCase -}}
{{- $parent := . -}}
func Test{{$tableNamePlural}}Upsert(t *testing.T) {
  //var err error

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
    conflict: []string{},
    update: []string{},
    whitelist: []string{"thing"},
    returning: []string{},
  }

  query := o.generateUpsertQuery(false, upsertCols)
  expectedQuery := `INSERT INTO {{.Table.Name}} ("thing") VALUES ($1) ON CONFLICT DO NOTHING`

  if query != expectedQuery {
    t.Errorf("Expected query mismatch:\n\n%s\n%s\n", query, expectedQuery)
  }

  /*
  query = o.generateUpsertQuery(true, upsertCols)
  primKeys := strings.Join(strmangle.IdentQuote())
  expectedQuery = `INSERT INTO {{.Table.Name}} ("thing") VALUES ($1) ON CONFLICT DO UPDATE()`

  if query != expectedQuery {
    t.Errorf("Expected query mismatch:\n\n%s\n%s\n", query, expectedQuery)
  }
  */

  /*
  create empty row
  assign random values to it

  attempt to insert it using upsert
  make sure values come back appropriately

  attempt to upsert row again, make sure comes back as prim key error
  attempt upsert again, set update to false, ensure it ignores error

  attempt to randomize everything except primary keys on duplicate row
  attempt upsert again, set update to true, nil, nil
  perform a find on the the row
  check if the found row matches the upsert object to ensure returning cols worked appropriately and update worked appropriately




  */

  {{$varNamePlural}}DeleteAllRows(t)
}
