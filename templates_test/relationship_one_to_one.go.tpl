{{- if .Table.IsJoinTable -}}
{{- else -}}
	{{- range $rel := .Table.ToOneRelationships -}}
		{{- $ltable := $.Aliases.Table $rel.Table -}}
		{{- $ftable := $.Aliases.Table $rel.ForeignTable -}}
		{{- $relAlias := $ftable.Relationship $rel.Name -}}
		{{- $usesPrimitives := usesPrimitives $.Tables $rel.Table $rel.Column $rel.ForeignTable $rel.ForeignColumn -}}
		{{- $colField := $ltable.Column $rel.Column -}}
		{{- $fcolField := $ftable.Column $rel.ForeignColumn }}
func test{{$ltable.UpSingular}}OneToOne{{$ftable.UpSingular}}Using{{$relAlias.Local}}(t *testing.T) {
	{{if not $.NoContext}}ctx := context.Background(){{end}}
	tx := MustTx({{if $.NoContext}}boil.Begin(){{else}}boil.BeginTx(ctx, nil){{end}})
	defer func() { _ = tx.Rollback() }()

	var foreign {{$ftable.UpSingular}}
	var local {{$ltable.UpSingular}}

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &foreign, {{$ftable.DownSingular}}DBTypes, true, {{$ftable.DownSingular}}ColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize {{$ftable.UpSingular}} struct: %s", err)
	}
	if err := randomize.Struct(seed, &local, {{$ltable.DownSingular}}DBTypes, true, {{$ltable.DownSingular}}ColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize {{$ltable.UpSingular}} struct: %s", err)
	}

	if err := local.Insert({{if not $.NoContext}}ctx, {{end -}} tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	{{if $usesPrimitives -}}
	foreign.{{$fcolField}} = local.{{$colField}}
	{{else -}}
	queries.Assign(&foreign.{{$fcolField}}, local.{{$colField}})
	{{end -}}
	if err := foreign.Insert({{if not $.NoContext}}ctx, {{end -}} tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	check, err := local.{{$relAlias.Local}}().One({{if not $.NoContext}}ctx, {{end -}} tx)
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
	if err = local.L.Load{{$relAlias.Local}}({{if not $.NoContext}}ctx, {{end -}} tx, false, (*[]*{{$ltable.UpSingular}})(&slice), nil); err != nil {
		t.Fatal(err)
	}
	if local.R.{{$relAlias.Local}} == nil {
		t.Error("struct should have been eager loaded")
	}

	local.R.{{$relAlias.Local}} = nil
	if err = local.L.Load{{$relAlias.Local}}({{if not $.NoContext}}ctx, {{end -}} tx, true, &local, nil); err != nil {
		t.Fatal(err)
	}
	if local.R.{{$relAlias.Local}} == nil {
		t.Error("struct should have been eager loaded")
	}
}

{{end -}}{{/* range */}}
{{- end -}}{{/* join table */}}
