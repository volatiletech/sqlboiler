{{- $tableNameSingular := .Table.Name | singular | titleCase -}}
{{- $dbName := singular .Table.Name -}}
{{- $tableNamePlural := .Table.Name | plural | titleCase -}}
{{- $varNamePlural := .Table.Name | plural | camelCase -}}
{{- $varNameSingular := .Table.Name | singular | camelCase -}}
func Test{{$tableNamePlural}}Insert(t *testing.T) {
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
    j[i], err = {{$tableNameSingular}}Find({{.Table.PKey.Columns | stringMap .StringFuncs.titleCase | prefixStringSlice "o[i]." | join ", "}})
    {{$varNameSingular}}CompareVals(o[i], j[i], t)
  }

  // Ensure we can store zero-valued object successfully.
  item := &{{$tableNameSingular}}{}
  if err = item.Insert(); err != nil {
    t.Errorf("Unable to insert zero-value item {{$tableNameSingular}}:\n%#v\nErr: %s", item, err)
  }

  {{with .Table.Columns | filterColumnsByAutoIncrement true | columnNames}}
  // Ensure the auto increment columns are returned in the object
    {{range .}}
  if item.{{titleCase .}} <= 0 {
    t.Errorf("Expected the auto-increment primary key to be greater than 0, got: %d", item.{{titleCase .}})
  }
    {{end}}
  {{end}}

  emptyTime := time.Time{}.String()
  {{with .Table.Columns | filterColumnsBySimpleDefault}}
	// Ensure the default value columns are returned in the object
    {{range .}}
      {{$tc := titleCase .Name}}
      {{$zv := zeroValue .}}
      {{$dv := defaultValue .}}
      {{$ty := trimPrefix "null." .Type}}
      {{if and (ne $ty "[]byte") .IsNullable}}
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
  {{end}}

  {{with .Table.Columns | filterColumnsByAutoIncrement false | filterColumnsByDefault false}}
  // Ensure the non-defaultvalue columns and non-autoincrement columns are stored correctly as zero or null values.
    {{range .}}
      {{- $tc := titleCase .Name -}}
      {{- $zv := zeroValue . -}}
      {{$ty := trimPrefix "null." .Type}}
      {{if and (ne $ty "[]byte") .IsNullable}}
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
}
