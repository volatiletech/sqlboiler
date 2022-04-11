{{- $alias := .Aliases.Table .Table.Name}}
func test{{$alias.UpPlural}}Update(t *testing.T) {
	t.Parallel()

	if 0 == len({{$alias.DownSingular}}PrimaryKeyColumns) {
		t.Skip("Skipping table with no primary key columns")
	}
	if len({{$alias.DownSingular}}AllColumns) == len({{$alias.DownSingular}}PrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &{{$alias.UpSingular}}{}
	if err = randomize.Struct(seed, o, {{$alias.DownSingular}}DBTypes, true, {{$alias.DownSingular}}ColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize {{$alias.UpSingular}} struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx,  tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := {{$alias.UpPlural}}().Count(ctx,  tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, {{$alias.DownSingular}}DBTypes, true, {{$alias.DownSingular}}PrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize {{$alias.UpSingular}} struct: %s", err)
	}

	if rowsAff, err := o.Update(ctx,  tx, boil.Infer()); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only affect one row but affected", rowsAff)
	}
}

func test{{$alias.UpPlural}}SliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len({{$alias.DownSingular}}AllColumns) == len({{$alias.DownSingular}}PrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &{{$alias.UpSingular}}{}
	if err = randomize.Struct(seed, o, {{$alias.DownSingular}}DBTypes, true, {{$alias.DownSingular}}ColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize {{$alias.UpSingular}} struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx,  tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := {{$alias.UpPlural}}().Count(ctx,  tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, {{$alias.DownSingular}}DBTypes, true, {{$alias.DownSingular}}PrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize {{$alias.UpSingular}} struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch({{$alias.DownSingular}}AllColumns, {{$alias.DownSingular}}PrimaryKeyColumns) {
		fields = {{$alias.DownSingular}}AllColumns
	} else {
		fields = strmangle.SetComplement(
			{{$alias.DownSingular}}AllColumns,
			{{$alias.DownSingular}}PrimaryKeyColumns,
		)
		{{- if filterColumnsByAuto true .Table.Columns }}
		fields = strmangle.SetComplement(fields, {{$alias.DownSingular}}GeneratedColumns)
		{{- end}}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	typ := reflect.TypeOf(o).Elem()
	n := typ.NumField()

	updateMap := M{}
	for _, col := range fields {
		for i := 0; i < n; i++ {
			f := typ.Field(i)
			if f.Tag.Get("boil") == col {
				updateMap[col] = value.Field(i).Interface()
			}
		}
	}

	slice := {{$alias.UpSingular}}Slice{{"{"}}o{{"}"}}
	if rowsAff, err := slice.UpdateAll(ctx,  tx, updateMap); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("wanted one record updated but got", rowsAff)
	}
}
