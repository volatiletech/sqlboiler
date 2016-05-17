{{- $tableNameSingular := titleCaseSingular .Table.Name -}}
{{- $dbName := singular .Table.Name -}}
{{- $tableNamePlural := titleCasePlural .Table.Name -}}
{{- $varNamePlural := camelCasePlural .Table.Name -}}
// {{$tableNamePlural}}All retrieves all records.
func Test{{$tableNamePlural}}All(t *testing.T) {
  var err error

  r := make([]{{$tableNameSingular}}, 2)

  // insert two random columns to test DeleteAll
  for i, v := range r {
    err = boil.RandomizeStruct(&v)
    if err != nil {
      t.Errorf("%d: Unable to randomize {{$tableNameSingular}} struct: %s", i, err)
    }

    err = v.Insert()
    if err != nil {
      t.Errorf("Unable to insert {{$tableNameSingular}}:\n%#v\nErr: %s", v, err)
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

  o := make([]{{$tableNameSingular}}, 3)

  for i, v := range o {
    err = boil.RandomizeStruct(&v)
    if err != nil {
      t.Errorf("%d: Unable to randomize {{$tableNameSingular}} struct: %s", i, err)
    }

    err = v.Insert()
    if err != nil {
      t.Errorf("Unable to insert {{$tableNameSingular}}:\n%#v\nErr: %s", v, err)
    }
  }

  // Attempt to retrieve all objects
  res, err := {{$tableNamePlural}}().All()
  if err != nil {
    t.Errorf("Unable to retrieve all {{$tableNamePlural}}, err: %s", err)
  }

  if len(res) != 3 {
    t.Errorf("Expected 3 {{$tableNameSingular}} rows, got %d", len(res))
  }
}
