{{- if .Table.IsJoinTable -}}
{{- else -}}
	{{- range $rel := .Table.ToOneRelationships -}}
		{{- $ltable := $.Aliases.Table $rel.Table -}}
		{{- $ftable := $.Aliases.Table $rel.ForeignTable -}}
		{{- $relAlias := $ftable.Relationship $rel.Name -}}
		{{- $usesPrimitives := usesPrimitives $.Tables $rel.Table $rel.Column $rel.ForeignTable $rel.ForeignColumn -}}
		{{- $colField := $ltable.Column $rel.Column -}}
		{{- $fcolField := $ftable.Column $rel.ForeignColumn -}}
		{{- $foreignPKeyCols := (getTable $.Tables .ForeignTable).PKey.Columns }}
		{{- $canSoftDelete := (getTable $.Tables .ForeignTable).CanSoftDelete $.AutoColumns.Deleted }}
func test{{$ltable.UpSingular}}OneToOneSetOp{{$ftable.UpSingular}}Using{{$relAlias.Local}}(t *testing.T) {
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
		err = a.Set{{$relAlias.Local}}({{if not $.NoContext}}ctx, {{end -}} tx, i != 0, x)
		if err != nil {
			t.Fatal(err)
		}

		if a.R.{{$relAlias.Local}} != x {
			t.Error("relationship struct not set to correct value")
		}
		if x.R.{{$relAlias.Foreign}} != &a {
			t.Error("failed to append to foreign relationship struct")
		}

		{{if $usesPrimitives -}}
		if a.{{$colField}} != x.{{$fcolField}} {
		{{else -}}
		if !queries.Equal(a.{{$colField}}, x.{{$fcolField}}) {
		{{end -}}
			t.Error("foreign key was wrong value", a.{{$colField}})
		}

		{{if setInclude .ForeignColumn $foreignPKeyCols -}}
		if exists, err := {{$ftable.UpSingular}}Exists({{if not $.NoContext}}ctx, {{end -}} tx, x.{{$foreignPKeyCols | stringMap $.StringFuncs.titleCase | join ", x."}}); err != nil {
			t.Fatal(err)
		} else if !exists {
			t.Error("want 'x' to exist")
		}
		{{else -}}
		zero := reflect.Zero(reflect.TypeOf(x.{{$fcolField}}))
		reflect.Indirect(reflect.ValueOf(&x.{{$fcolField}})).Set(zero)

		if err = x.Reload({{if not $.NoContext}}ctx, {{end -}} tx); err != nil {
			t.Fatal("failed to reload", err)
		}
		{{- end}}

		{{if $usesPrimitives -}}
		if a.{{$colField}} != x.{{$fcolField}} {
		{{else -}}
		if !queries.Equal(a.{{$colField}}, x.{{$fcolField}}) {
		{{end -}}
			t.Error("foreign key was wrong value", a.{{$colField}}, x.{{$fcolField}})
		}

		if {{if not $.NoRowsAffected}}_, {{end -}} err = x.Delete({{if not $.NoContext}}ctx, {{end -}} tx {{- if and $.AddSoftDeletes $canSoftDelete}}, true{{end}}); err != nil {
			t.Fatal("failed to delete x", err)
		}
	}
}
{{- if $rel.ForeignColumnNullable}}

func test{{$ltable.UpSingular}}OneToOneRemoveOp{{$ftable.UpSingular}}Using{{$relAlias.Local}}(t *testing.T) {
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

	if err = a.Set{{$relAlias.Local}}({{if not $.NoContext}}ctx, {{end -}} tx, true, &b); err != nil {
		t.Fatal(err)
	}

	if err = a.Remove{{$relAlias.Local}}({{if not $.NoContext}}ctx, {{end -}} tx, &b); err != nil {
		t.Error("failed to remove relationship")
	}

	count, err := a.{{$relAlias.Local}}().Count({{if not $.NoContext}}ctx, {{end -}} tx)
	if err != nil {
		t.Error(err)
	}
	if count != 0 {
		t.Error("want no relationships remaining")
	}

	if a.R.{{$relAlias.Local}} != nil {
		t.Error("R struct entry should be nil")
	}

	if !queries.IsValuerNil(b.{{$fcolField}}) {
		t.Error("foreign key column should be nil")
	}

	if b.R.{{$relAlias.Foreign}} != nil {
		t.Error("failed to remove a from b's relationships")
	}
}
{{end -}}{{/* end if foreign key nullable */}}
{{- end -}}{{/* range */}}
{{- end -}}{{/* join table */}}
