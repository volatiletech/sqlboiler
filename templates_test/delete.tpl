{{- $tableNameSingular := .Table.Name | singular | titleCase -}}
{{- $dbName := singular .Table.Name -}}
{{- $tableNamePlural := .Table.Name | plural | titleCase -}}
{{- $varNamePlural := .Table.Name | plural | camelCase -}}
{{- $varNameSingular := .Table.Name | singular | camelCase -}}
func {{$varNamePlural}}DeleteAllRows(t *testing.T) {
  // Delete all rows to give a clean slate
  err := {{$tableNamePlural}}().DeleteAll()
  if err != nil {
    t.Errorf("Unable to delete all from {{$tableNamePlural}}: %s", err)
  }
}

func Test{{$tableNamePlural}}QueryDeleteAll(t *testing.T) {
  var err error
  var c int64

  // Start from a clean slate
  {{$varNamePlural}}DeleteAllRows(t)

  // Check number of rows in table to ensure DeleteAll was successful
  c, err = {{$tableNamePlural}}().Count()

  if c != 0 {
    t.Errorf("Expected 0 rows after ObjDeleteAllRows() call, but got %d rows", c)
  }

  o := make({{$varNameSingular}}Slice, 3)
  if err = boil.RandomizeSlice(&o); err != nil {
    t.Errorf("Unable to randomize {{$tableNameSingular}} slice: %s", err)
  }

  // insert random columns to test DeleteAll
  for i := 0; i < len(o); i++ {
    err = o[i].Insert()
    if err != nil {
      t.Errorf("Unable to insert {{$tableNameSingular}}:\n%#v\nErr: %s", o[i], err)
    }
  }

  // Test DeleteAll() query function
  err = {{$tableNamePlural}}().DeleteAll()
  if err != nil {
    t.Errorf("Unable to delete all from {{$tableNamePlural}}: %s", err)
  }

  // Check number of rows in table to ensure DeleteAll was successful
  c, err = {{$tableNamePlural}}().Count()

  if c != 0 {
    t.Errorf("Expected 0 rows after Obj().DeleteAll() call, but got %d rows", c)
  }
}

func Test{{$tableNamePlural}}SliceDeleteAll(t *testing.T) {
  var err error
  var c int64

  // insert random columns to test DeleteAll
  o := make({{$varNameSingular}}Slice, 3)
  if err = boil.RandomizeSlice(&o); err != nil {
    t.Errorf("Unable to randomize {{$tableNameSingular}} slice: %s", err)
  }

  for i := 0; i < len(o); i++ {
    err = o[i].Insert()
    if err != nil {
      t.Errorf("Unable to insert {{$tableNameSingular}}:\n%#v\nErr: %s", o[i], err)
    }
  }

  // test DeleteAll slice function
  if err = o.DeleteAll(); err != nil {
    t.Errorf("Unable to objSlice.DeleteAll(): %s", err)
  }

  // Check number of rows in table to ensure DeleteAll was successful
  c, err = {{$tableNamePlural}}().Count()

  if c != 0 {
    t.Errorf("Expected 0 rows after objSlice.DeleteAll() call, but got %d rows", c)
  }
}

func Test{{$tableNamePlural}}Delete(t *testing.T) {
  var err error
  var c int64

  // insert random columns to test Delete
  o := make({{$varNameSingular}}Slice, 3)
  if err = boil.RandomizeSlice(&o); err != nil {
    t.Errorf("Unable to randomize {{$tableNameSingular}} slice: %s", err)
  }

  for i := 0; i < len(o); i++ {
    err = o[i].Insert()
    if err != nil {
      t.Errorf("Unable to insert {{$tableNameSingular}}:\n%#v\nErr: %s", o[i], err)
    }
  }

  o[0].Delete()

  // Check number of rows in table to ensure DeleteAll was successful
  c, err = {{$tableNamePlural}}().Count()

  if c != 2 {
    t.Errorf("Expected 2 rows after obj.Delete() call, but got %d rows", c)
  }

  o[1].Delete()
  o[2].Delete()

  // Check number of rows in table to ensure Delete worked for all rows
  c, err = {{$tableNamePlural}}().Count()

  if c != 0 {
    t.Errorf("Expected 0 rows after all obj.Delete() calls, but got %d rows", c)
  }
}
