{{- $tableNameSingular := titleCaseSingular .Table.Name -}}
{{- $dbName := singular .Table.Name -}}
{{- $tableNamePlural := titleCasePlural .Table.Name -}}
{{- $varNamePlural := camelCasePlural .Table.Name -}}
{{- $varNameSingular := camelCaseSingular .Table.Name -}}
func {{$varNamePlural}}DeleteAllRows(t *testing.T) {
  // Delete all rows to give a clean slate
  err := {{$tableNamePlural}}().DeleteAll()
  if err != nil {
    t.Errorf("Unable to delete all from {{$tableNamePlural}}: %s", err)
  }
}

func Test{{$tableNamePlural}}Delete(t *testing.T) {
  var err error

  // Start from a clean slate
  {{$varNamePlural}}DeleteAllRows(t)

  r := make({{$varNameSingular}}Slice, 3)

  // insert random columns to test DeleteAll
  for i := 0; i < len(r); i++ {
    err = boil.RandomizeStruct(r[i])
    if err != nil {
      t.Errorf("%d: Unable to randomize {{$tableNameSingular}} struct: %s", i, err)
    }

    err = r[i].Insert()
    if err != nil {
      t.Errorf("Unable to insert {{$tableNameSingular}}:\n%#v\nErr: %s", r[i], err)
    }
  }

  // Test DeleteAll() query function
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

  // insert random columns to test DeleteAll
  o := make({{$varNameSingular}}Slice, 3)
  for i := 0; i < len(o); i++ {
    err = boil.RandomizeStruct(o[i])
    if err != nil {
      t.Errorf("%d: Unable to randomize {{$tableNameSingular}} struct: %s", i, err)
    }

    err = o[i].Insert()
    if err != nil {
      t.Errorf("Unable to insert {{$tableNameSingular}}:\n%#v\nErr: %s", o[i], err)
    }
  }

  // test DeleteAll slice function
  err = o.DeleteAll()

  // Check number of rows in table to ensure DeleteAll was successful
  c, err = {{$tableNamePlural}}().Count()

  if c != 0 {
    t.Errorf("Expected {{.Table.Name}} table to be empty, but got %d rows", c)
  }

  // insert random columns to test Delete
  o = make({{$varNameSingular}}Slice, 3)
  for i := 0; i < len(o); i++ {
    err = boil.RandomizeStruct(o[i])
    if err != nil {
      t.Errorf("%d: Unable to randomize {{$tableNameSingular}} struct: %s", i, err)
    }

    err = o[i].Insert()
    if err != nil {
      t.Errorf("Unable to insert {{$tableNameSingular}}:\n%#v\nErr: %s", o[i], err)
    }
  }

  //o[0].Delete()

  // Check number of rows in table to ensure DeleteAll was successful
  c, err = {{$tableNamePlural}}().Count()

  if c != 2 {
    t.Errorf("Expected {{.Table.Name}} table to have 2 rows, but got %d rows", c)
  }
}
