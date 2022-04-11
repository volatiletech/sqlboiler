{{- if .Table.IsJoinTable -}}
{{- else -}}
	{{- range $fkey := .Table.FKeys -}}
		{{- $ltable := $.Aliases.Table $fkey.Table -}}
		{{- $ftable := $.Aliases.Table $fkey.ForeignTable -}}
		{{- $rel := $ltable.Relationship $fkey.Name -}}
		{{- $colField := $ltable.Column $fkey.Column -}}
		{{- $fcolField := $ftable.Column $fkey.ForeignColumn -}}
		{{- $usesPrimitives := usesPrimitives $.Tables $fkey.Table $fkey.Column $fkey.ForeignTable $fkey.ForeignColumn }}
func test{{$ltable.UpSingular}}ToOne{{$ftable.UpSingular}}Using{{$rel.Foreign}}(t *testing.T) {
	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var local {{$ltable.UpSingular}}
	var foreign {{$ftable.UpSingular}}

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, {{$ltable.DownSingular}}DBTypes, {{if $fkey.Nullable}}true{{else}}false{{end}}, {{$ltable.DownSingular}}ColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize {{$ltable.UpSingular}} struct: %s", err)
	}
	if err := randomize.Struct(seed, &foreign, {{$ftable.DownSingular}}DBTypes, {{if $fkey.ForeignColumnNullable}}true{{else}}false{{end}}, {{$ftable.DownSingular}}ColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize {{$ftable.UpSingular}} struct: %s", err)
	}

	if err := foreign.Insert(ctx,  tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	{{if $usesPrimitives -}}
	local.{{$colField}} = foreign.{{$fcolField}}
	{{else -}}
	queries.Assign(&local.{{$colField}}, foreign.{{$fcolField}})
	{{end -}}
	if err := local.Insert(ctx,  tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	check, err := local.{{$rel.Foreign}}().One(ctx,  tx)
	if err != nil {
		t.Fatal(err)
	}

	{{if $usesPrimitives -}}
	if check.{{$fcolField}} != foreign.{{$fcolField}} {
	{{else -}}
	if !queries.Equal(check.{{$fcolField}}, foreign.{{$fcolField}}) {
	{{end -}}
		t.Errorf("want: %v, got %v", foreign.{{$fcolField}}, check.{{$fcolField}})
	}

	slice := {{$ltable.UpSingular}}Slice{&local}
	if err = local.L.Load{{$rel.Foreign}}(ctx,  tx, false, (*[]*{{$ltable.UpSingular}})(&slice), nil); err != nil {
		t.Fatal(err)
	}
	if local.R.{{$rel.Foreign}} == nil {
		t.Error("struct should have been eager loaded")
	}

	local.R.{{$rel.Foreign}} = nil
	if err = local.L.Load{{$rel.Foreign}}(ctx,  tx, true, &local, nil); err != nil {
		t.Fatal(err)
	}
	if local.R.{{$rel.Foreign}} == nil {
		t.Error("struct should have been eager loaded")
	}
}

{{end -}}{{/* range */}}
{{- end -}}{{/* join table */}}
