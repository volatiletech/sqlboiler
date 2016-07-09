{{- $tableNameSingular := .Table.Name | singular | titleCase -}}
{{- $dbName := singular .Table.Name -}}
{{- $tableNamePlural := .Table.Name | plural | titleCase -}}
{{- $varNamePlural := .Table.Name | plural | camelCase -}}
{{- $varNameSingular := .Table.Name | singular | camelCase -}}
{{- $parent := . -}}
func Test{{$tableNamePlural}}Insert(t *testing.T) {
  t.Skip("don't need this ruining everything")
  var err error
  var errs []error
  emptyTime := time.Time{}.String()

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
    j[i], err = {{$tableNameSingular}}Find({{.Table.PKey.Columns | stringMap .StringFuncs.titleCase | prefixStringSlice "o[i]." | join ", "}})
    {{$varNameSingular}}CompareVals(o[i], j[i], t)
  }

  /**
   * Edge case test for:
   * No includes specified, all zero values.
   *
   * Expected result:
   * Columns with default values set to their default values.
   * Columns without default values set to their zero value.
   */

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


  /*{{with .Table.Columns | filterColumnsBySimpleDefault}}
	// Ensure the default value columns are returned in the object
    {{range .}}
      {{$tc := titleCase .Name}}
      {{$zv := zeroValue .}}
      {{$dv := "false"}}
      {{$ty := trimPrefix "null." .Type}}
      {{if and (ne $ty "[]byte") .Nullable}}
  if item.{{$tc}}.Valid == false {
    t.Errorf("Expected the nullable default value column {{$tc}} of {{$tableNameSingular}} to be valid.")
  }
        {{if eq .Type "null.Time"}}
  if (item.{{$tc}}.{{$ty}}.String() == emptyTime && !isZeroTime({{$dv}})) ||
    (item.{{$tc}}.{{$ty}}.String() > emptyTime && isZeroTime({{$dv}})) {
        {{else}}
  if item.{{$tc}}.{{$ty}} != {{$dv}} {
        {{- end -}}
    t.Errorf("Expected the nullable default value column {{$tc}} of {{$tableNameSingular}} to match the database default value:\n%#v\n%v\n\n", item.{{$tc}}.{{$ty}}, {{$dv}})
  }
      {{else}}
        {{if eq .Type "[]byte"}}
  if string(item.{{$tc}}) != string({{$dv}}) {
        {{else if eq .Type "time.Time"}}
  if (item.{{$tc}}.String() == emptyTime && !isZeroTime({{$dv}})) ||
    (item.{{$tc}}.String() > emptyTime && isZeroTime({{$dv}})) {
        {{else}}
  if item.{{$tc}} != {{$dv}} {
        {{- end -}}
    t.Errorf("Expected the default value column {{$tc}} of {{$tableNameSingular}} to match the database default value:\n%#v\n%v\n\n", item.{{$tc}}, {{$dv}})
  }
      {{end}}
    {{end}}
  {{end}}*/

  {{with .Table.Columns | filterColumnsByAutoIncrement false | filterColumnsByDefault false}}
  // Ensure the non-defaultvalue columns and non-autoincrement columns are stored correctly as zero or null values.
    {{range .}}
      {{- $tc := titleCase .Name -}}
      {{- $zv := zeroValue . -}}
      {{$ty := trimPrefix "null." .Type}}
      {{if and (ne $ty "[]byte") .Nullable}}
  if item.{{$tc}}.Valid == true {
    t.Errorf("Expected the nullable column {{$tc}} of {{$tableNameSingular}} to be invalid (null).")
  }
        {{if eq .Type "null.Time"}}
  if item.{{$tc}}.{{$ty}}.String() != emptyTime {
        {{else}}
  if item.{{$tc}}.{{$ty}} != {{$zv}} {
        {{- end -}}
    t.Errorf("Expected the nullable column {{$tc}} of {{$tableNameSingular}} to be a zero-value (null):\n%#v\n%v\n\n", item.{{$tc}}.{{$ty}}, {{$zv}})
  }
      {{else}}
        {{if eq .Type "[]byte"}}
  if string(item.{{$tc}}) != string({{$zv}}) {
        {{else if eq .Type "time.Time"}}
  if item.{{$tc}}.String() != emptyTime {
        {{else}}
  if item.{{$tc}} != {{$zv}} {
        {{- end -}}
    t.Errorf("Expected the column {{$tc}} of {{$tableNameSingular}} to be a zero-value (null):\n%#v\n%v\n\n", item.{{$tc}}, {{$zv}})
  }
      {{- end}}
    {{end}}
  {{end}}

  /**
   * Edge case test for:
   * No includes specified, all non-zero values.
   *
   * Expected result:
   * Non-zero auto-increment column values ignored by insert helper.
   * Object updated with correct auto-increment values.
   * All other column values in object remain.
   */

  {{$varNamePlural}}DeleteAllRows(t)

  item = &{{$tableNameSingular}}{}
  if err = item.Insert(); err != nil {
    t.Errorf("Unable to insert zero-value item {{$tableNameSingular}}:\n%#v\nErr: %s", item, err)
  }

  /**
   * Edge case test for:
   * Auto-increment columns and nullable columns includes specified, all zero values.
   *
   * Expected result:
   * Auto-increment value inserted and nullable column values inserted.
   * Default values for nullable columns should NOT be present in returned object.
   */

  /**
   * Edge case test for:
   * Auto-increment columns and nullable columns includes specified, all non-zero values.
   *
   * Expected result:
   * Auto-increment value inserted and nullable column values inserted.
   * Default values for nullable columns should NOT be present in returned object.
   * Should be no zero values anywhere.
   */

  /**
   * Edge case test for:
   * Non-nullable columns includes specified, all zero values.
   *
   * Expected result:
   * Auto-increment values ignored by insert helper.
   * Object updated with correct auto-increment values.
   * All non-nullable columns should be returned as zero values, regardless of default values.
   */

  /**
   * Edge case test for:
   * Non-nullable columns includes specified, all non-zero values.
   *
   * Expected result:
   * Auto-increment values ignored by insert helper.
   * Object updated with correct auto-increment values.
   */
}
