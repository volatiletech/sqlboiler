{{- $tableNameSingular := .Table.Name | singular | titleCase -}}
{{- $dbName := singular .Table.Name -}}
{{- $tableNamePlural := .Table.Name | plural | titleCase -}}
{{- $varNamePlural := .Table.Name | plural | camelCase -}}
{{- $varNameSingular := .Table.Name | singular | camelCase -}}
func {{$varNameSingular}}CompareVals(o *{{$tableNameSingular}}, j *{{$tableNameSingular}}, t *testing.T) {
  {{range $key, $value := .Table.Columns}}
  {{if eq $value.Type "null.Time"}}
  if o.{{titleCase $value.Name}}.Time.Format("02/01/2006") != j.{{titleCase $value.Name}}.Time.Format("02/01/2006") {
    t.Errorf("Expected NullTime {{$value.Name}} column string values to match, got:\nStruct: %#v\nResponse: %#v\n\n", o.{{titleCase $value.Name}}.Time.Format("02/01/2006"), j.{{titleCase $value.Name}}.Time.Format("02/01/2006"))
  }
  {{else if eq $value.Type "time.Time"}}
  if o.{{titleCase $value.Name}}.Format("02/01/2006") != j.{{titleCase $value.Name}}.Format("02/01/2006") {
    t.Errorf("Expected Time {{$value.Name}} column string values to match, got:\nStruct: %#v\nResponse: %#v\n\n", o.{{titleCase $value.Name}}.Format("02/01/2006"), j.{{titleCase $value.Name}}.Format("02/01/2006"))
  }
  {{else if eq $value.Type "[]byte"}}
  if !byteSliceEqual(o.{{titleCase $value.Name}}, j.{{titleCase $value.Name}}) {
    t.Errorf("Expected {{$value.Name}} columns to match, got:\nStruct: %#v\nResponse: %#v\n\n", o.{{titleCase $value.Name}}, j.{{titleCase $value.Name}})
  }
  {{else}}
  if j.{{titleCase $value.Name}} != o.{{titleCase $value.Name}} {
    t.Errorf("Expected {{$value.Name}} columns to match, got:\nStruct: %#v\nResponse: %#v\n\n", o.{{titleCase $value.Name}}, j.{{titleCase $value.Name}})
  }
  {{end}}
  {{end}}
}

func Test{{$tableNamePlural}}(t *testing.T) {
  var err error

  // Start from a clean slate
  {{$varNamePlural}}DeleteAllRows(t)

  o := make({{$varNameSingular}}Slice, 2)
  if err = boil.RandomizeSlice(&o); err != nil {
    t.Errorf("Unable to randomize {{$tableNameSingular}} slice: %s", err)
  }

  // insert two random columns to test DeleteAll
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
  if err = boil.RandomizeSlice(&o); err != nil {
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
}
