{{- $alias := .Aliases.Table .Table.Name}}
func test{{$alias.UpPlural}}Find(t *testing.T) {
	t.Parallel()

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

	{{$alias.DownSingular}}Found, err := Find{{$alias.UpSingular}}(ctx,  tx, {{.Table.PKey.Columns | stringMap (aliasCols $alias) | prefixStringSlice (printf "%s." "o") | join ", "}})
	if err != nil {
		t.Error(err)
	}

	if {{$alias.DownSingular}}Found == nil {
		t.Error("want a record, got nil")
	}
}
