{{- $tableNameSingular := .Table.Name | singular | titleCase -}}
{{- $dbName := singular .Table.Name -}}
{{- $tableNamePlural := .Table.Name | plural | titleCase -}}
{{- $varNamePlural := .Table.Name | plural | camelCase -}}
{{- $varNameSingular := .Table.Name | singular | camelCase -}}
{{- $parent := . -}}
func Test{{$tableNamePlural}}Insert(t *testing.T) {
  var err error

  var errs []error
  _ = errs

  emptyTime := time.Time{}.String()
  _ = emptyTime

  nullTime := null.NewTime(time.Time{}, true)
  _ = nullTime

  {{$varNamePlural}}DeleteAllRows(t)

  o := make({{$varNameSingular}}Slice, 3)
  if err = boil.RandomizeSlice(&o, {{$varNameSingular}}DBTypes); err != nil {
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
    j[i], err = {{$tableNameSingular}}Find({{.Table.PKey.Columns | stringMap .StringFuncs.titleCase | prefixStringSlice "o[i]." | join ", "}})
    {{$varNameSingular}}CompareVals(o[i], j[i], t)
  }

  {{$varNamePlural}}DeleteAllRows(t)

  item := &{{$tableNameSingular}}{}
  if err = item.Insert(); err != nil {
    t.Errorf("Unable to insert zero-value item {{$tableNameSingular}}:\n%#v\nErr: %s", item, err)
  }

  {{with .Table.Columns | filterColumnsByAutoIncrement true | columnNames | stringMap $parent.StringFuncs.quoteWrap | join ", "}}
  // Ensure the auto increment columns are returned in the object
  if errs = boil.IsZeroValue(item, false, {{.}}); errs != nil {
    for _, e := range errs {
      t.Errorf("Expected auto-increment columns to be greater than 0, err: %s\n", e)
    }
  }
  {{end}}

  {{with .Table.Columns | filterColumnsBySimpleDefault}}
  simpleDefaults := []string{{"{"}}{{. | columnNames | stringMap $parent.StringFuncs.quoteWrap | join ", "}}{{"}"}}
  defaultValues := []interface{}{{"{"}}{{. | defaultValues | join ", "}}{{"}"}}

  if len(simpleDefaults) != len(defaultValues) {
    t.Fatalf("Mismatch between slice lengths: %d, %d", len(simpleDefaults), len(defaultValues))
  }

  if errs = boil.IsValueMatch(item, simpleDefaults, defaultValues); errs != nil {
    for _, e := range errs {
      t.Errorf("Expected default value to match column value, err: %s\n", e);
    }
  }
  {{end}}

  {{with .Table.Columns | filterColumnsByAutoIncrement false | filterColumnsByDefault false}}
  // Ensure the non-defaultvalue columns and non-autoincrement columns are stored correctly as zero or null values.
  regularCols := []string{{"{"}}{{. | columnNames | stringMap $parent.StringFuncs.quoteWrap | join ", "}}{{"}"}}

  for _, c := range regularCols {
    rv := reflect.Indirect(reflect.ValueOf(item))
    field := rv.FieldByName(strmangle.TitleCase(c))

    zv := reflect.Zero(field.Type()).Interface()
    fv := field.Interface()

    if !reflect.DeepEqual(zv, fv) {
      t.Errorf("Expected column %s to be zero value, got: %v, wanted: %v", c, fv, zv)
    }
  }
  {{end}}
}
