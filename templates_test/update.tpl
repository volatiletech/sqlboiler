{{- $tableNameSingular := .Table.Name | singular | titleCase -}}
{{- $tableNamePlural := .Table.Name | plural | titleCase -}}
{{- $varNamePlural := .Table.Name | plural | camelCase -}}
{{- $varNameSingular := .Table.Name | singular | camelCase -}}
func test{{$tableNamePlural}}Update(t *testing.T) {
	t.Parallel()

	if 0 == len({{$varNameSingular}}PrimaryKeyColumns) {
		t.Skip("Skipping table with no primary key columns")
	}
	if len({{$varNameSingular}}Columns) == len({{$varNameSingular}}PrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &{{$tableNameSingular}}{}
	if err = randomize.Struct(seed, o, {{$varNameSingular}}DBTypes, true, {{$varNameSingular}}ColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize {{$tableNameSingular}} struct: %s", err)
	}

	{{if not .NoContext}}ctx := context.Background(){{end}}
	tx := MustTx({{if .NoContext}}boil.Begin(){{else}}boil.BeginTx(ctx, nil){{end}})
	defer tx.Rollback()
	if err = o.Insert({{if not .NoContext}}ctx, {{end -}} tx); err != nil {
		t.Error(err)
	}

	count, err := {{$tableNamePlural}}().Count({{if not .NoContext}}ctx, {{end -}} tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, {{$varNameSingular}}DBTypes, true, {{$varNameSingular}}PrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize {{$tableNameSingular}} struct: %s", err)
	}

	{{if .NoRowsAffected -}}
	if err = o.Update({{if not .NoContext}}ctx, {{end -}} tx); err != nil {
		t.Error(err)
	}
	{{else -}}
	if rowsAff, err := o.Update({{if not .NoContext}}ctx, {{end -}} tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only affect one row but affected", rowsAff)
	}
	{{end -}}
}

func test{{$tableNamePlural}}SliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len({{$varNameSingular}}Columns) == len({{$varNameSingular}}PrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &{{$tableNameSingular}}{}
	if err = randomize.Struct(seed, o, {{$varNameSingular}}DBTypes, true, {{$varNameSingular}}ColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize {{$tableNameSingular}} struct: %s", err)
	}

	{{if not .NoContext}}ctx := context.Background(){{end}}
	tx := MustTx({{if .NoContext}}boil.Begin(){{else}}boil.BeginTx(ctx, nil){{end}})
	defer tx.Rollback()
	if err = o.Insert({{if not .NoContext}}ctx, {{end -}} tx); err != nil {
		t.Error(err)
	}

	count, err := {{$tableNamePlural}}().Count({{if not .NoContext}}ctx, {{end -}} tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, {{$varNameSingular}}DBTypes, true, {{$varNameSingular}}PrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize {{$tableNameSingular}} struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch({{$varNameSingular}}Columns, {{$varNameSingular}}PrimaryKeyColumns) {
		fields = {{$varNameSingular}}Columns
	} else {
		fields = strmangle.SetComplement(
			{{$varNameSingular}}Columns,
			{{$varNameSingular}}PrimaryKeyColumns,
		)
		{{- if .Dialect.UseAutoColumns }}
		fields = strmangle.SetComplement(
			fields,
			{{$varNameSingular}}ColumnsWithAuto,
		)
		{{- end}}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	updateMap := M{}
	for _, col := range fields {
		updateMap[col] = value.FieldByName(strmangle.TitleCase(col)).Interface()
	}

	slice := {{$tableNameSingular}}Slice{{"{"}}o{{"}"}}
	{{if .NoRowsAffected -}}
	if err = slice.UpdateAll({{if not .NoContext}}ctx, {{end -}} tx, updateMap); err != nil {
		t.Error(err)
	}
	{{else -}}
	if rowsAff, err := slice.UpdateAll({{if not .NoContext}}ctx, {{end -}} tx, updateMap); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("wanted one record updated but got", rowsAff)
	}
	{{end -}}
}
