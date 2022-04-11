{{- $alias := .Aliases.Table .Table.Name}}
func test{{$alias.UpPlural}}Bind(t *testing.T) {
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

	if err = {{$alias.UpPlural}}().Bind(ctx, tx, o); err != nil {
		t.Error(err)
	}
}

func test{{$alias.UpPlural}}One(t *testing.T) {
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

	if x, err := {{$alias.UpPlural}}().One(ctx,  tx); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func test{{$alias.UpPlural}}All(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	{{$alias.DownSingular}}One := &{{$alias.UpSingular}}{}
	{{$alias.DownSingular}}Two := &{{$alias.UpSingular}}{}
	if err = randomize.Struct(seed, {{$alias.DownSingular}}One, {{$alias.DownSingular}}DBTypes, false, {{$alias.DownSingular}}ColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize {{$alias.UpSingular}} struct: %s", err)
	}
	if err = randomize.Struct(seed, {{$alias.DownSingular}}Two, {{$alias.DownSingular}}DBTypes, false, {{$alias.DownSingular}}ColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize {{$alias.UpSingular}} struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = {{$alias.DownSingular}}One.Insert(ctx,  tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = {{$alias.DownSingular}}Two.Insert(ctx,  tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := {{$alias.UpPlural}}().All(ctx,  tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func test{{$alias.UpPlural}}Count(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	{{$alias.DownSingular}}One := &{{$alias.UpSingular}}{}
	{{$alias.DownSingular}}Two := &{{$alias.UpSingular}}{}
	if err = randomize.Struct(seed, {{$alias.DownSingular}}One, {{$alias.DownSingular}}DBTypes, false, {{$alias.DownSingular}}ColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize {{$alias.UpSingular}} struct: %s", err)
	}
	if err = randomize.Struct(seed, {{$alias.DownSingular}}Two, {{$alias.DownSingular}}DBTypes, false, {{$alias.DownSingular}}ColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize {{$alias.UpSingular}} struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = {{$alias.DownSingular}}One.Insert(ctx,  tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = {{$alias.DownSingular}}Two.Insert(ctx,  tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := {{$alias.UpPlural}}().Count(ctx,  tx)
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}
