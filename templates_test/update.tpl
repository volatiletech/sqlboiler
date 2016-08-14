{{- $tableNameSingular := .Table.Name | singular | titleCase -}}
{{- $tableNamePlural := .Table.Name | plural | titleCase -}}
{{- $varNamePlural := .Table.Name | plural | camelCase -}}
{{- $varNameSingular := .Table.Name | singular | camelCase -}}
func Test{{$tableNamePlural}}Update(t *testing.T) {
  var err error

  item := {{$tableNameSingular}}{}
  boil.RandomizeValidatedStruct(&item, {{$varNameSingular}}ValidatedColumns, {{$varNameSingular}}DBTypes)
  if err = item.InsertG(); err != nil {
    t.Errorf("Unable to insert zero-value item {{$tableNameSingular}}:\n%#v\nErr: %s", item, err)
  }

  blacklistCols := strmangle.SetMerge({{$varNameSingular}}AutoIncrementColumns, {{$varNameSingular}}PrimaryKeyColumns)
  if err = boil.RandomizeStruct(&item, {{$varNameSingular}}DBTypes, false, blacklistCols...); err != nil {
    t.Errorf("Unable to randomize {{$tableNameSingular}} struct: %s", err)
  }

  whitelist := strmangle.SetComplement({{$varNameSingular}}Columns, {{$varNameSingular}}AutoIncrementColumns)
  if err = item.UpdateG(whitelist...); err != nil {
    t.Errorf("Unable to update {{$tableNameSingular}}: %s", err)
  }

  var j *{{$tableNameSingular}}
  j, err = {{$tableNameSingular}}FindG({{.Table.PKey.Columns | stringMap .StringFuncs.titleCase | prefixStringSlice "item." | join ", "}})
  if err != nil {
    t.Errorf("Unable to find {{$tableNameSingular}} row: %s", err)
  }

  err = {{$varNameSingular}}CompareVals(&item, j, true); if err != nil {
    t.Error(err)
  }

  {{$varNamePlural}}DeleteAllRows(t)
}

func Test{{$tableNamePlural}}SliceUpdateAll(t *testing.T) {
  var err error

  // insert random columns to test UpdateAll
  o := make({{$tableNameSingular}}Slice, 3)
  j := make({{$tableNameSingular}}Slice, 3)

  if err = boil.RandomizeSlice(&o, {{$varNameSingular}}DBTypes, false); err != nil {
    t.Errorf("Unable to randomize {{$tableNameSingular}} slice: %s", err)
  }

  for i := 0; i < len(o); i++ {
    if err = o[i].InsertG(); err != nil {
      t.Errorf("Unable to insert {{$tableNameSingular}}:\n%#v\nErr: %s", o[i], err)
    }
  }

  vals := M{}

  tmp := {{$tableNameSingular}}{}
  blacklist := strmangle.SetMerge({{$varNameSingular}}PrimaryKeyColumns, {{$varNameSingular}}UniqueColumns)
  if err = boil.RandomizeStruct(&tmp, {{$varNameSingular}}DBTypes, false, blacklist...); err != nil {
    t.Errorf("Unable to randomize struct {{$tableNameSingular}}: %s", err)
  }

  // Build the columns and column values from the randomized struct
	tmpVal := reflect.Indirect(reflect.ValueOf(tmp))
  nonBlacklist := strmangle.SetComplement({{$varNameSingular}}Columns, blacklist)
  for _, col := range nonBlacklist {
    vals[col] = tmpVal.FieldByName(strmangle.TitleCase(col)).Interface()
  }

  err = o.UpdateAllG(vals)
  if err != nil {
    t.Errorf("Failed to update all for {{$tableNameSingular}}: %s", err)
  }

  for i := 0; i < len(o); i++ {
    j[i], err = {{$tableNameSingular}}FindG({{.Table.PKey.Columns | stringMap .StringFuncs.titleCase | prefixStringSlice "o[i]." | join ", "}})
    if err != nil {
      t.Errorf("Unable to find {{$tableNameSingular}} row: %s", err)
    }

    err = {{$varNameSingular}}CompareVals(j[i], &tmp, true, blacklist...)
    if err != nil {
      t.Error(err)
    }
  }

	for i := 0; i < len(o); i++ {
		// Ensure Find found the correct primary key ID's
    orig := boil.GetStructValues(o[i], {{$varNameSingular}}PrimaryKeyColumns...)
		new := boil.GetStructValues(j[i], {{$varNameSingular}}PrimaryKeyColumns...)

    if !reflect.DeepEqual(orig, new) {
      t.Errorf("object %d): primary keys do not match:\n\n%#v\n%#v", i, orig, new)
    }
	}

  {{$varNamePlural}}DeleteAllRows(t)
}
