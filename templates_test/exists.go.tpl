{{- $alias := .Aliases.Table .Table.Name}}
func test{{$alias.UpPlural}}Exists(t *testing.T) {
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

	{{$pkeyArgs := .Table.PKey.Columns | stringMap (aliasCols $alias) | prefixStringSlice (printf "%s." "o") | join ", " -}}
	e, err := {{$alias.UpSingular}}Exists({{if not .NoContext}}ctx, {{end -}} tx, {{$pkeyArgs}})
	if err != nil {
		t.Errorf("Unable to check if {{$alias.UpSingular}} exists: %s", err)
	}
	if !e {
		t.Errorf("Expected {{$alias.UpSingular}}Exists to return true, but got false.")
	}
}
