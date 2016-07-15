{{- $tableNameSingular := .Table.Name | singular | titleCase -}}
{{- $dbName := singular .Table.Name -}}
{{- $tableNamePlural := .Table.Name | plural | titleCase -}}
{{- $varNamePlural := .Table.Name | plural | camelCase -}}
{{- $varNameSingular := .Table.Name | singular | camelCase -}}
func Test{{$tableNamePlural}}(t *testing.T) {
  var err error

  o := make({{$varNameSingular}}Slice, 2)
  if err = boil.RandomizeSlice(&o, {{$varNameSingular}}DBTypes, true); err != nil {
    t.Errorf("Unable to randomize {{$tableNameSingular}} slice: %s", err)
  }

  // insert two random objects to test DeleteAll
  for i := 0; i < len(o); i++ {
    err = o[i].Insert()
    if err != nil {
      t.Errorf("Unable to insert {{$tableNameSingular}}:\n%#v\nErr: %s", o[i], err)
    }
  }

  // Delete all rows to give a clean slate
  err = {{$tableNamePlural}}().DeleteAll()
  if err != nil {
    t.Errorf("Unable to delete all from {{$tableNamePlural}}: %s", err)
  }

  // Check number of rows in table to ensure DeleteAll was successful
  var c int64
  c, err = {{$tableNamePlural}}().Count()

  if c != 0 {
    t.Errorf("Expected {{.Table.Name}} table to be empty, but got %d rows", c)
  }

  o = make({{$varNameSingular}}Slice, 3)
  if err = boil.RandomizeSlice(&o, {{$varNameSingular}}DBTypes, true); err != nil {
    t.Errorf("Unable to randomize {{$tableNameSingular}} slice: %s", err)
  }

  for i := 0; i < len(o); i++ {
    err = o[i].Insert()
    if err != nil {
      t.Errorf("Unable to insert {{$tableNameSingular}}:\n%#v\nErr: %s", o[i], err)
    }
  }

  // Ensure Count is valid
  c, err = {{$tableNamePlural}}().Count()
  if c != 3 {
    t.Errorf("Expected {{.Table.Name}} table to have 3 rows, but got %d", c)
  }

  // Attempt to retrieve all objects
  res, err := {{$tableNamePlural}}().All()
  if err != nil {
    t.Errorf("Unable to retrieve all {{$tableNamePlural}}, err: %s", err)
  }

  if len(res) != 3 {
    t.Errorf("Expected 3 {{$tableNameSingular}} rows, got %d", len(res))
  }

  {{$varNamePlural}}DeleteAllRows(t)
}
