{{- $tableNameSingular := .Table.Name | singular | titleCase -}}
{{- $tableNamePlural := .Table.Name | plural | titleCase -}}
{{- $varNamePlural := .Table.Name | plural | camelCase -}}
{{- $varNameSingular := .Table.Name | singular | camelCase -}}
func test{{$tableNamePlural}}Exists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &{{$tableNameSingular}}{}
	if err = randomize.Struct(seed, o, {{$varNameSingular}}DBTypes, true, {{$varNameSingular}}ColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize {{$tableNameSingular}} struct: %s", err)
	}

	{{if not .NoContext}}ctx := context.Background(){{end}}
	tx := MustTx({{if .NoContext}}boil.Begin(){{else}}boil.BeginTx(ctx, nil){{end}})
	defer tx.Rollback()
	if err = o.Insert({{if not .NoContext}}ctx, {{end -}} tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	{{$pkeyArgs := .Table.PKey.Columns | stringMap .StringFuncs.titleCase | prefixStringSlice (printf "%s." "o") | join ", " -}}
	e, err := {{$tableNameSingular}}Exists({{if not .NoContext}}ctx, {{end -}} tx, {{$pkeyArgs}})
	if err != nil {
		t.Errorf("Unable to check if {{$tableNameSingular}} exists: %s", err)
	}
	if !e {
		t.Errorf("Expected {{$tableNameSingular}}ExistsG to return true, but got false.")
	}
}
