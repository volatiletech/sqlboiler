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
	{{$varNameSingular}} := &{{$tableNameSingular}}{}
	if err = randomize.Struct(seed, {{$varNameSingular}}, {{$varNameSingular}}DBTypes, true, {{$varNameSingular}}ColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize {{$tableNameSingular}} struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = {{$varNameSingular}}.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := {{$tableNamePlural}}(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, {{$varNameSingular}}, {{$varNameSingular}}DBTypes, true, {{$varNameSingular}}PrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize {{$tableNameSingular}} struct: %s", err)
	}

	{{if .NoRowsAffected -}}
	if err = {{$varNameSingular}}.Update(tx); err != nil {
		t.Error(err)
	}
	{{else -}}
	if rowsAff, err := {{$varNameSingular}}.Update(tx); err != nil {
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
	{{$varNameSingular}} := &{{$tableNameSingular}}{}
	if err = randomize.Struct(seed, {{$varNameSingular}}, {{$varNameSingular}}DBTypes, true, {{$varNameSingular}}ColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize {{$tableNameSingular}} struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = {{$varNameSingular}}.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := {{$tableNamePlural}}(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, {{$varNameSingular}}, {{$varNameSingular}}DBTypes, true, {{$varNameSingular}}PrimaryKeyColumns...); err != nil {
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

	value := reflect.Indirect(reflect.ValueOf({{$varNameSingular}}))
	updateMap := M{}
	for _, col := range fields {
		updateMap[col] = value.FieldByName(strmangle.TitleCase(col)).Interface()
	}

	slice := {{$tableNameSingular}}Slice{{"{"}}{{$varNameSingular}}{{"}"}}
	{{if .NoRowsAffected -}}
	if err = slice.UpdateAll(tx, updateMap); err != nil {
		t.Error(err)
	}
	{{else -}}
	if rowsAff, err := slice.UpdateAll(tx, updateMap); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("wanted one record updated but got", rowsAff)
	}
	{{end -}}
}
