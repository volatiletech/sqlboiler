{{- $tableNameSingular := titleCaseSingular .Table.Name -}}
{{- $dbName := singular .Table.Name -}}
{{- $tableNamePlural := titleCasePlural .Table.Name -}}
{{- $varNamePlural := camelCasePlural .Table.Name -}}
{{- $varNameSingular := camelCaseSingular .Table.Name -}}
func Test{{$tableNamePlural}}Find(t *testing.T) {
  var err error

  {{$varNamePlural}}DeleteAllRows(t)

  o := make({{$varNameSingular}}Slice, 3)
  if err = boil.RandomizeSlice(&o); err != nil {
    t.Errorf("Unable to randomize {{$tableNameSingular}} slice: %s", err)
  }

  for i := 0; i < len(o); i++ {
    if err = o[i].Insert(); err != nil {
      t.Errorf("Unable to insert {{$tableNameSingular}}:\n%#v\nErr: %s", o[i], err)
    }
  }

  j := make({{$varNameSingular}}Slice, 3)
  // Perform all Find queries and assign result objects to slice for comparison
  for i := 0; i < len(j); i++ {
    j[i], err = {{$tableNameSingular}}Find({{titleCaseCommaList "o[i]." .Table.PKey.Columns}})

    {{range $key, $value := .Table.Columns}}
    {{if eq $value.Type "null.Time"}}
    if o[i].{{titleCase $value.Name}}.Time.Format("02/01/2006") != j[i].{{titleCase $value.Name}}.Time.Format("02/01/2006") {
      t.Errorf("%d) Expected NullTime {{$value.Name}} column string values to match, got:\nStruct: %#v\nResponse: %#v\n\n", i, o[i].{{titleCase $value.Name}}.Time.Format("02/01/2006"), j[i].{{titleCase $value.Name}}.Time.Format("02/01/2006"))
    }
    {{else if eq $value.Type "time.Time"}}
    if o[i].{{titleCase $value.Name}}.Format("02/01/2006") != j[i].{{titleCase $value.Name}}.Format("02/01/2006") {
      t.Errorf("%d) Expected Time {{$value.Name}} column string values to match, got:\nStruct: %#v\nResponse: %#v\n\n", i, o[i].{{titleCase $value.Name}}.Format("02/01/2006"), j[i].{{titleCase $value.Name}}.Format("02/01/2006"))
    }
    {{else if eq $value.Type "[]byte"}}
    if !byteSliceEqual(o[i].{{titleCase $value.Name}}, j[i].{{titleCase $value.Name}}) {
      t.Errorf("%d) Expected {{$value.Name}} columns to match, got:\nStruct: %#v\nResponse: %#v\n\n", i, o[i].{{titleCase $value.Name}}, j[i].{{titleCase $value.Name}})
    }
    {{else}}
    if j[i].{{titleCase $value.Name}} != o[i].{{titleCase $value.Name}} {
      t.Errorf("%d) Expected {{$value.Name}} columns to match, got:\nStruct: %#v\nResponse: %#v\n\n", i, o[i].{{titleCase $value.Name}}, j[i].{{titleCase $value.Name}})
    }
    {{end}}
    {{end}}
  }

  {{if hasPrimaryKey .Table.PKey}}
  f, err := {{$tableNameSingular}}Find({{titleCaseCommaList "o[0]." .Table.PKey.Columns}}, {{$varNameSingular}}PrimaryKeyColumns...)
  {{range $key, $value := .Table.PKey.Columns}}
  if o[0].{{titleCase $value}} != f.{{titleCase $value}} {
    t.Errorf("Expected primary key values to match, {{titleCase $value}} did not match")
  }
  {{end}}

  colsWithoutPrimKeys := boil.SetComplement({{$varNameSingular}}Columns, {{$varNameSingular}}PrimaryKeyColumns)
  fRef := reflect.ValueOf(f).Elem()
  for _, v := range colsWithoutPrimKeys {
    val := fRef.FieldByName(v)
    if val.IsValid() {
      t.Errorf("Expected all other columns to be zero value, but column %s was %#v", v, val.Interface())
    }
  }
  {{end}}
}
