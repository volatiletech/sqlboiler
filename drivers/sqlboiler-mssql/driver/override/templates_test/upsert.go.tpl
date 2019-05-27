{{- $alias := .Aliases.Table .Table.Name}}
func test{{$alias.UpPlural}}Upsert(t *testing.T) {
	t.Parallel()

	if len({{$alias.DownSingular}}AllColumns) == len({{$alias.DownSingular}}PrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	o := {{$alias.UpSingular}}{}
	if err = randomize.Struct(seed, &o, {{$alias.DownSingular}}DBTypes, true); err != nil {
		t.Errorf("Unable to randomize {{$alias.UpSingular}} struct: %s", err)
	}

	{{if not .NoContext}}ctx := context.Background(){{end}}
	tx := MustTx({{if .NoContext}}{{if .NoContext}}boil.Begin(){{else}}boil.BeginTx(ctx, nil){{end}}{{else}}boil.BeginTx(ctx, nil){{end}})
	defer func() { _ = tx.Rollback() }()
	if err = o.Upsert({{if not .NoContext}}ctx, {{end -}} tx, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert {{$alias.UpSingular}}: %s", err)
	}

	count, err := {{$alias.UpPlural}}().Count({{if not .NoContext}}ctx, {{end -}} tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &o, {{$alias.DownSingular}}DBTypes, false, {{$alias.DownSingular}}PrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize {{$alias.UpSingular}} struct: %s", err)
	}

	if err = o.Upsert({{if not .NoContext}}ctx, {{end -}} tx, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert {{$alias.UpSingular}}: %s", err)
	}

	count, err = {{$alias.UpPlural}}().Count({{if not .NoContext}}ctx, {{end -}} tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}
