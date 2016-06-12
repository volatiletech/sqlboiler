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

    // Compare saved objects from earlier to the found objects
    if !reflect.DeepEqual(j[i], o[i]) {
      t.Errorf("Expected j[%d] to match o[%d], got:\n\nj: #%v\n\no:#%v\n\n", i, i, j[i], o[i])
    }
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
