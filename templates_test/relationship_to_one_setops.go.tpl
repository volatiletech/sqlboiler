{{- if .Table.IsJoinTable -}}
{{- else -}}
	{{- range $fkey := .Table.FKeys -}}
		{{- $ltable := $.Aliases.Table $fkey.Table -}}
		{{- $ftable := $.Aliases.Table $fkey.ForeignTable -}}
		{{- $rel := $ltable.Relationship $fkey.Name -}}
		{{- $colField := $ltable.Column $fkey.Column -}}
		{{- $fcolField := $ftable.Column $fkey.ForeignColumn -}}
		{{- $usesPrimitives := usesPrimitives $.Tables $fkey.Table $fkey.Column $fkey.ForeignTable $fkey.ForeignColumn }}
func test{{$ltable.UpSingular}}ToOneSetOp{{$ftable.UpSingular}}Using{{$rel.Foreign}}(t *testing.T) {
	var err error

	{{if not $.NoContext}}ctx := context.Background(){{end}}
	tx := MustTx({{if $.NoContext}}boil.Begin(){{else}}boil.BeginTx(ctx, nil){{end}})
	defer func() { _ = tx.Rollback() }()

	var a {{$ltable.UpSingular}}
	var b, c {{$ftable.UpSingular}}

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, {{$ltable.DownSingular}}DBTypes, false, strmangle.SetComplement({{$ltable.DownSingular}}PrimaryKeyColumns, {{$ltable.DownSingular}}ColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, {{$ftable.DownSingular}}DBTypes, false, strmangle.SetComplement({{$ftable.DownSingular}}PrimaryKeyColumns, {{$ftable.DownSingular}}ColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, {{$ftable.DownSingular}}DBTypes, false, strmangle.SetComplement({{$ftable.DownSingular}}PrimaryKeyColumns, {{$ftable.DownSingular}}ColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err := a.Insert({{if not $.NoContext}}ctx, {{end -}} tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert({{if not $.NoContext}}ctx, {{end -}} tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	for i, x := range []*{{$ftable.UpSingular}}{&b, &c} {
		err = a.Set{{$rel.Foreign}}({{if not $.NoContext}}ctx, {{end -}} tx, i != 0, x)
		if err != nil {
			t.Fatal(err)
		}

		if a.R.{{$rel.Foreign}} != x {
			t.Error("relationship struct not set to correct value")
		}

		{{if $fkey.Unique -}}
		if x.R.{{$rel.Local}} != &a {
			t.Error("failed to append to foreign relationship struct")
		}
		{{else -}}
		if x.R.{{$rel.Local}}[0] != &a {
			t.Error("failed to append to foreign relationship struct")
		}
		{{end -}}

		{{if $usesPrimitives -}}
		if a.{{$colField}} != x.{{$fcolField}} {
		{{else -}}
		if !queries.Equal(a.{{$colField}}, x.{{$fcolField}}) {
		{{end -}}
			t.Error("foreign key was wrong value", a.{{$colField}})
		}

		{{if setInclude $fkey.Column $.Table.PKey.Columns -}}
		if exists, err := {{$ltable.UpSingular}}Exists({{if not $.NoContext}}ctx, {{end -}} tx, a.{{$.Table.PKey.Columns | stringMap (aliasCols $ltable) | join ", a."}}); err != nil {
			t.Fatal(err)
		} else if !exists {
			t.Error("want 'a' to exist")
		}
		{{else -}}
		zero := reflect.Zero(reflect.TypeOf(a.{{$colField}}))
		reflect.Indirect(reflect.ValueOf(&a.{{$colField}})).Set(zero)

		if err = a.Reload({{if not $.NoContext}}ctx, {{end -}} tx); err != nil {
			t.Fatal("failed to reload", err)
		}

		{{if $usesPrimitives -}}
		if a.{{$colField}} != x.{{$fcolField}} {
		{{else -}}
		if !queries.Equal(a.{{$colField}}, x.{{$fcolField}}) {
		{{end -}}
			t.Error("foreign key was wrong value", a.{{$colField}}, x.{{$fcolField}})
		}
		{{- end}}
	}
}
{{- if $fkey.Nullable}}

func test{{$ltable.UpSingular}}ToOneRemoveOp{{$ftable.UpSingular}}Using{{$rel.Foreign}}(t *testing.T) {
	var err error

	{{if not $.NoContext}}ctx := context.Background(){{end}}
	tx := MustTx({{if $.NoContext}}boil.Begin(){{else}}boil.BeginTx(ctx, nil){{end}})
	defer func() { _ = tx.Rollback() }()

	var a {{$ltable.UpSingular}}
	var b {{$ftable.UpSingular}}

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, {{$ltable.DownSingular}}DBTypes, false, strmangle.SetComplement({{$ltable.DownSingular}}PrimaryKeyColumns, {{$ltable.DownSingular}}ColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, {{$ftable.DownSingular}}DBTypes, false, strmangle.SetComplement({{$ftable.DownSingular}}PrimaryKeyColumns, {{$ftable.DownSingular}}ColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err = a.Insert({{if not $.NoContext}}ctx, {{end -}} tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	if err = a.Set{{$rel.Foreign}}({{if not $.NoContext}}ctx, {{end -}} tx, true, &b); err != nil {
		t.Fatal(err)
	}

	if err = a.Remove{{$rel.Foreign}}({{if not $.NoContext}}ctx, {{end -}} tx, &b); err != nil {
		t.Error("failed to remove relationship")
	}

	count, err := a.{{$rel.Foreign}}().Count({{if not $.NoContext}}ctx, {{end -}} tx)
	if err != nil {
		t.Error(err)
	}
	if count != 0 {
		t.Error("want no relationships remaining")
	}

	if a.R.{{$rel.Foreign}} != nil {
		t.Error("R struct entry should be nil")
	}

	if !queries.IsValuerNil(a.{{$colField}}) {
		t.Error("foreign key value should be nil")
	}

	{{if $fkey.Unique -}}
	if b.R.{{$rel.Local}} != nil {
		t.Error("failed to remove a from b's relationships")
	}
	{{else -}}
	if len(b.R.{{$rel.Local}}) != 0 {
		t.Error("failed to remove a from b's relationships")
	}
	{{- end}}
}
{{end -}}{{/* end if foreign key nullable */}}
{{- end -}}{{/* range */}}
{{- end -}}{{/* join table */}}
