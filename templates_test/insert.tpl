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

  o := make({{$tableNameSingular}}Slice, 3)
  if err = boil.RandomizeSlice(&o, {{$varNameSingular}}DBTypes, true); err != nil {
    t.Errorf("Unable to randomize {{$tableNameSingular}} slice: %s", err)
  }

  for i := 0; i < len(o); i++ {
    if err = o[i].InsertG(); err != nil {
      t.Errorf("Unable to insert {{$tableNameSingular}}:\n%#v\nErr: %s", o[i], err)
    }
  }

  j := make({{$tableNameSingular}}Slice, 3)
  // Perform all Find queries and assign result objects to slice for comparison
  for i := 0; i < len(j); i++ {
    j[i], err = {{$tableNameSingular}}FindG({{.Table.PKey.Columns | stringMap .StringFuncs.titleCase | prefixStringSlice "o[i]." | join ", "}})
    {{$varNameSingular}}CompareVals(o[i], j[i], t)
  }

  {{$varNamePlural}}DeleteAllRows(t)

  item := &{{$tableNameSingular}}{}
  boil.RandomizeEnforcedStruct(item, {{$varNameSingular}}ColumnsEnforced, {{$varNameSingular}}DBTypes)
  if err = item.InsertG(); err != nil {
    t.Errorf("Unable to insert zero-value item {{$tableNameSingular}}:\n%#v\nErr: %s", item, err)
  }

  for _, c := range {{$varNameSingular}}AutoIncrementColumns {
    // Ensure the auto increment columns are returned in the object.
    if errs = boil.IsZeroValue(item, false, c); errs != nil {
      for _, e := range errs {
        t.Errorf("Expected auto-increment columns to be greater than 0, err: %s\n", e)
      }
    }
  }

  defaultValues := []interface{}{{"{"}}{{.Table.Columns | filterColumnsBySimpleDefault | defaultValues | join ", "}}{{"}"}}

  // Ensure the simple default column values are returned correctly.
  if len({{$varNameSingular}}ColumnsWithSimpleDefault) > 0 && len(defaultValues) > 0 {
    if len({{$varNameSingular}}ColumnsWithSimpleDefault) != len(defaultValues) {
      t.Fatalf("Mismatch between slice lengths: %d, %d", len({{$varNameSingular}}ColumnsWithSimpleDefault), len(defaultValues))
    }

    if errs = boil.IsValueMatch(item, {{$varNameSingular}}ColumnsWithSimpleDefault, defaultValues); errs != nil {
      for _, e := range errs {
        t.Errorf("Expected default value to match column value, err: %s\n", e);
      }
    }
  }

  regularCols := []string{{"{"}}{{.Table.Columns | filterColumnsByAutoIncrement false | filterColumnsByDefault false | columnNames | stringMap $parent.StringFuncs.quoteWrap | join ", "}}{{"}"}}

  // Remove the enforced columns, they can never be zero values
  regularCols = boil.SetComplement(regularCols, {{$varNameSingular}}ColumnsEnforced)

  // Ensure the non-defaultvalue columns and non-autoincrement columns are stored correctly as zero or null values.
  for _, c := range regularCols {
    rv := reflect.Indirect(reflect.ValueOf(item))
    field := rv.FieldByName(strmangle.TitleCase(c))

    zv := reflect.Zero(field.Type()).Interface()
    fv := field.Interface()

    if !reflect.DeepEqual(zv, fv) {
      t.Errorf("Expected column %s to be zero value, got: %v, wanted: %v", c, fv, zv)
    }
  }

  item = &{{$tableNameSingular}}{}

  wl, rc := item.generateInsertColumns()
  if !reflect.DeepEqual(rc, {{$varNameSingular}}ColumnsWithDefault) {
    t.Errorf("Expected return columns to contain all columns with default values:\n\nGot: %v\nWanted: %v", rc, {{$varNameSingular}}ColumnsWithDefault)
  }

  if !reflect.DeepEqual(wl, {{$varNameSingular}}ColumnsWithoutDefault) {
    t.Errorf("Expected whitelist to contain all columns without default values:\n\nGot: %v\nWanted: %v", wl, {{$varNameSingular}}ColumnsWithoutDefault)
  }

  if err = boil.RandomizeStruct(item, {{$varNameSingular}}DBTypes, false); err != nil {
    t.Errorf("Unable to randomize item: %s", err)
  }

  wl, rc = item.generateInsertColumns()
  if len(rc) > 0 {
    t.Errorf("Expected return columns to contain no columns:\n\nGot: %v", rc)
  }

  if !reflect.DeepEqual(wl, {{$varNameSingular}}Columns) {
    t.Errorf("Expected whitelist to contain all columns values:\n\nGot: %v\nWanted: %v", wl, {{$varNameSingular}}Columns)
  }

  {{$varNamePlural}}DeleteAllRows(t)
}
