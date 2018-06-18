{{- $alias := .Aliases.Table .Table.Name}}
func test{{$alias.UpPlural}}Reload(t *testing.T) {
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

	if err = o.Reload({{if not .NoContext}}ctx, {{end -}} tx); err != nil {
		t.Error(err)
	}
}

func test{{$alias.UpPlural}}ReloadAll(t *testing.T) {
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

	slice := {{$alias.UpSingular}}Slice{{"{"}}o{{"}"}}

	if err = slice.ReloadAll({{if not .NoContext}}ctx, {{end -}} tx); err != nil {
		t.Error(err)
	}
}
