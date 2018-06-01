{{- $tableNameSingular := .Table.Name | singular | titleCase -}}
{{- $tableNamePlural := .Table.Name | plural | titleCase -}}
{{- $varNamePlural := .Table.Name | plural | camelCase -}}
{{- $varNameSingular := .Table.Name | singular | camelCase -}}
func test{{$tableNamePlural}}Upsert(t *testing.T) {
	t.Parallel()

	if len({{$varNameSingular}}Columns) == len({{$varNameSingular}}PrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	o := {{$tableNameSingular}}{}
	if err = randomize.Struct(seed, &o, {{$varNameSingular}}DBTypes, true); err != nil {
		t.Errorf("Unable to randomize {{$tableNameSingular}} struct: %s", err)
	}

	{{if not .NoContext}}ctx := context.Background(){{end}}
	tx := MustTx({{if .NoContext}}{{if .NoContext}}boil.Begin(){{else}}boil.BeginTx(ctx, nil){{end}}{{else}}boil.BeginTx(ctx, nil){{end}})
	defer tx.Rollback()
	if err = o.Upsert({{if not .NoContext}}ctx, {{end -}} tx, nil); err != nil {
		t.Errorf("Unable to upsert {{$tableNameSingular}}: %s", err)
	}

	count, err := {{$tableNamePlural}}().Count({{if not .NoContext}}ctx, {{end -}} tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &o, {{$varNameSingular}}DBTypes, false, {{$varNameSingular}}PrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize {{$tableNameSingular}} struct: %s", err)
	}

	if err = o.Upsert({{if not .NoContext}}ctx, {{end -}} tx, nil); err != nil {
		t.Errorf("Unable to upsert {{$tableNameSingular}}: %s", err)
	}

	count, err = {{$tableNamePlural}}().Count({{if not .NoContext}}ctx, {{end -}} tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}
