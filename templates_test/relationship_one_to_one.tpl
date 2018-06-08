{{- if .Table.IsJoinTable -}}
{{- else -}}
	{{- range .Table.ToOneRelationships -}}
		{{- $txt := txtsFromOneToOne $.Tables $.Table . -}}
		{{- $varNameSingular := .Table | singular | camelCase -}}
		{{- $foreignVarNameSingular := .ForeignTable | singular | camelCase}}
func test{{$txt.LocalTable.NameGo}}OneToOne{{$txt.ForeignTable.NameGo}}Using{{$txt.Function.Name}}(t *testing.T) {
	{{if not $.NoContext}}ctx := context.Background(){{end}}
	tx := MustTx({{if $.NoContext}}boil.Begin(){{else}}boil.BeginTx(ctx, nil){{end}})
	defer tx.Rollback()

	var foreign {{$txt.ForeignTable.NameGo}}
	var local {{$txt.LocalTable.NameGo}}

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &foreign, {{$foreignVarNameSingular}}DBTypes, true, {{$foreignVarNameSingular}}ColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize {{$txt.ForeignTable.NameGo}} struct: %s", err)
	}
	if err := randomize.Struct(seed, &local, {{$varNameSingular}}DBTypes, true, {{$varNameSingular}}ColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize {{$txt.LocalTable.NameGo}} struct: %s", err)
	}

	if err := local.Insert({{if not $.NoContext}}ctx, {{end -}} tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	{{if $txt.Function.UsesPrimitives -}}
	foreign.{{$txt.ForeignTable.ColumnNameGo}} = local.{{$txt.LocalTable.ColumnNameGo}}
	{{else -}}
	queries.Assign(&foreign.{{$txt.ForeignTable.ColumnNameGo}}, local.{{$txt.LocalTable.ColumnNameGo}})
	{{end -}}
	if err := foreign.Insert({{if not $.NoContext}}ctx, {{end -}} tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	check, err := local.{{$txt.Function.Name}}().One({{if not $.NoContext}}ctx, {{end -}} tx)
	if err != nil {
		t.Fatal(err)
	}

	{{if $txt.Function.UsesPrimitives -}}
	if check.{{$txt.ForeignTable.ColumnNameGo}} != foreign.{{$txt.ForeignTable.ColumnNameGo}} {
	{{else -}}
	if !queries.Equal(check.{{$txt.ForeignTable.ColumnNameGo}}, foreign.{{$txt.ForeignTable.ColumnNameGo}}) {
	{{end -}}
		t.Errorf("want: %v, got %v", foreign.{{$txt.ForeignTable.ColumnNameGo}}, check.{{$txt.ForeignTable.ColumnNameGo}})
	}

	slice := {{$txt.LocalTable.NameGo}}Slice{&local}
	if err = local.L.Load{{$txt.Function.Name}}({{if not $.NoContext}}ctx, {{end -}} tx, false, (*[]*{{$txt.LocalTable.NameGo}})(&slice), nil); err != nil {
		t.Fatal(err)
	}
	if local.R.{{$txt.Function.Name}} == nil {
		t.Error("struct should have been eager loaded")
	}

	local.R.{{$txt.Function.Name}} = nil
	if err = local.L.Load{{$txt.Function.Name}}({{if not $.NoContext}}ctx, {{end -}} tx, true, &local, nil); err != nil {
		t.Fatal(err)
	}
	if local.R.{{$txt.Function.Name}} == nil {
		t.Error("struct should have been eager loaded")
	}
}

{{end -}}{{/* range */}}
{{- end -}}{{/* join table */}}
