{{- $alias := .Aliases.Table .Table.Name}}
func test{{$alias.UpPlural}}Insert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &{{$alias.UpSingular}}{}
	if err = randomize.Struct(seed, o, {{$alias.DownSingular}}DBTypes, true, {{$alias.DownSingular}}ColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize {{$alias.UpSingular}} struct: %s", err)
	}

	{{if not .NoContext}}ctx := context.Background(){{end}}
	tx := MustTx({{if .NoContext}}boil.Begin(){{else}}boil.BeginTx(ctx, nil){{end}})
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert({{if not .NoContext}}ctx, {{end -}} tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := {{$alias.UpPlural}}().Count({{if not .NoContext}}ctx, {{end -}} tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func test{{$alias.UpPlural}}InsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &{{$alias.UpSingular}}{}
	if err = randomize.Struct(seed, o, {{$alias.DownSingular}}DBTypes, true); err != nil {
		t.Errorf("Unable to randomize {{$alias.UpSingular}} struct: %s", err)
	}

	{{if not .NoContext}}ctx := context.Background(){{end}}
	tx := MustTx({{if .NoContext}}boil.Begin(){{else}}boil.BeginTx(ctx, nil){{end}})
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert({{if not .NoContext}}ctx, {{end -}} tx, boil.Whitelist({{$alias.DownSingular}}ColumnsWithoutDefault...)); err != nil {
		t.Error(err)
	}

	count, err := {{$alias.UpPlural}}().Count({{if not .NoContext}}ctx, {{end -}} tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func test{{$alias.UpPlural}}SliceInsertAll(t *testing.T) {
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
 	o := {{$alias.UpSingular}}Slice{ {{$alias.DownSingular}}One, {{$alias.DownSingular}}Two }
 	{{if not .NoContext}}ctx := context.Background(){{end}}
	tx := MustTx({{if .NoContext}}boil.Begin(){{else}}boil.BeginTx(ctx, nil){{end}})
	defer func() {_ = tx.Rollback()} ()
	if err = o.InsertAll({{if not .NoContext}}ctx, {{end -}}tx, boil.Infer()); err != nil {
		t.Error(err)
	}
 	count, err := {{$alias.UpPlural}}().Count({{if not .NoContext}}ctx, {{end -}} tx)
	if err != nil {
		t.Error(err)
	}
 	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}
