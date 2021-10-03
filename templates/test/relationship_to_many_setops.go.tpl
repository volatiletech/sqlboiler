{{- if .Table.IsJoinTable -}}
{{- else -}}
	{{- $table := .Table -}}
	{{- range $rel := .Table.ToManyRelationships -}}
		{{- $ltable := $.Aliases.Table $rel.Table -}}
		{{- $ftable := $.Aliases.Table $rel.ForeignTable -}}
		{{- $relAlias := $.Aliases.ManyRelationship $rel.ForeignTable $rel.Name $rel.JoinTable $rel.JoinLocalFKeyName -}}
		{{- $usesPrimitives := usesPrimitives $.Tables $rel.Table $rel.Column $rel.ForeignTable $rel.ForeignColumn -}}
		{{- $colField := $ltable.Column $rel.Column -}}
		{{- $fcolField := $ftable.Column $rel.ForeignColumn }}
func test{{$ltable.UpSingular}}ToManyAddOp{{$relAlias.Local}}(t *testing.T) {
	var err error

	{{if not $.NoContext}}ctx := context.Background(){{end}}
	tx := MustTx({{if $.NoContext}}boil.Begin(){{else}}boil.BeginTx(ctx, nil){{end}})
	defer func() { _ = tx.Rollback() }()

	var a {{$ltable.UpSingular}}
	var b, c, d, e {{$ftable.UpSingular}}

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, {{$ltable.DownSingular}}DBTypes, false, strmangle.SetComplement({{$ltable.DownSingular}}PrimaryKeyColumns, {{$ltable.DownSingular}}ColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*{{$ftable.UpSingular}}{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, {{$ftable.DownSingular}}DBTypes, false, strmangle.SetComplement({{$ftable.DownSingular}}PrimaryKeyColumns, {{$ftable.DownSingular}}ColumnsWithoutDefault)...); err != nil {
			t.Fatal(err)
		}
	}

	if err := a.Insert({{if not $.NoContext}}ctx, {{end -}} tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert({{if not $.NoContext}}ctx, {{end -}} tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert({{if not $.NoContext}}ctx, {{end -}} tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	foreignersSplitByInsertion := [][]*{{$ftable.UpSingular}}{
		{&b, &c},
		{&d, &e},
	}

	for i, x := range foreignersSplitByInsertion {
		err = a.Add{{$relAlias.Local}}({{if not $.NoContext}}ctx, {{end -}} tx, i != 0, x...)
		if err != nil {
			t.Fatal(err)
		}

		first := x[0]
		second := x[1]
		{{- if .ToJoinTable}}

		if first.R.{{$relAlias.Foreign}}[0] != &a {
			t.Error("relationship was not added properly to the slice")
		}
		if second.R.{{$relAlias.Foreign}}[0] != &a {
			t.Error("relationship was not added properly to the slice")
		}
		{{- else}}

		{{if $usesPrimitives -}}
		if a.{{$colField}} != first.{{$fcolField}} {
			t.Error("foreign key was wrong value", a.{{$colField}}, first.{{$fcolField}})
		}
		if a.{{$colField}} != second.{{$fcolField}} {
			t.Error("foreign key was wrong value", a.{{$colField}}, second.{{$fcolField}})
		}
		{{else -}}
		if !queries.Equal(a.{{$colField}}, first.{{$fcolField}}) {
			t.Error("foreign key was wrong value", a.{{$colField}}, first.{{$fcolField}})
		}
		if !queries.Equal(a.{{$colField}}, second.{{$fcolField}}) {
			t.Error("foreign key was wrong value", a.{{$colField}}, second.{{$fcolField}})
		}
		{{- end}}

		if first.R.{{$relAlias.Foreign}} != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}
		if second.R.{{$relAlias.Foreign}} != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}
		{{- end}}

		if a.R.{{$relAlias.Local}}[i*2] != first {
			t.Error("relationship struct slice not set to correct value")
		}
		if a.R.{{$relAlias.Local}}[i*2+1] != second {
			t.Error("relationship struct slice not set to correct value")
		}

		count, err := a.{{$relAlias.Local}}().Count({{if not $.NoContext}}ctx, {{end -}} tx)
		if err != nil {
			t.Fatal(err)
		}
		if want := int64((i+1)*2); count != want {
			t.Error("want", want, "got", count)
		}
	}
}
{{- if (or $rel.ForeignColumnNullable $rel.ToJoinTable)}}

func test{{$ltable.UpSingular}}ToManySetOp{{$relAlias.Local}}(t *testing.T) {
	var err error

	{{if not $.NoContext}}ctx := context.Background(){{end}}
	tx := MustTx({{if $.NoContext}}boil.Begin(){{else}}boil.BeginTx(ctx, nil){{end}})
	defer func() { _ = tx.Rollback() }()

	var a {{$ltable.UpSingular}}
	var b, c, d, e {{$ftable.UpSingular}}

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, {{$ltable.DownSingular}}DBTypes, false, strmangle.SetComplement({{$ltable.DownSingular}}PrimaryKeyColumns, {{$ltable.DownSingular}}ColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*{{$ftable.UpSingular}}{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, {{$ftable.DownSingular}}DBTypes, false, strmangle.SetComplement({{$ftable.DownSingular}}PrimaryKeyColumns, {{$ftable.DownSingular}}ColumnsWithoutDefault)...); err != nil {
			t.Fatal(err)
		}
	}

	if err = a.Insert({{if not $.NoContext}}ctx, {{end -}} tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert({{if not $.NoContext}}ctx, {{end -}} tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert({{if not $.NoContext}}ctx, {{end -}} tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	err = a.Set{{$relAlias.Local}}({{if not $.NoContext}}ctx, {{end -}} tx, false, &b, &c)
	if err != nil {
		t.Fatal(err)
	}

	count, err := a.{{$relAlias.Local}}().Count({{if not $.NoContext}}ctx, {{end -}} tx)
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	err = a.Set{{$relAlias.Local}}({{if not $.NoContext}}ctx, {{end -}} tx, true, &d, &e)
	if err != nil {
		t.Fatal(err)
	}

	count, err = a.{{$relAlias.Local}}().Count({{if not $.NoContext}}ctx, {{end -}} tx)
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	{{- if .ToJoinTable}}

	// The following checks cannot be implemented since we have no handle
	// to these when we call Set(). Leaving them here as wishful thinking
	// and to let people know there's dragons.
	//
	// if len(b.R.{{$relAlias.Foreign}}) != 0 {
	// 	t.Error("relationship was not removed properly from the slice")
	// }
	// if len(c.R.{{$relAlias.Foreign}}) != 0 {
	// 	t.Error("relationship was not removed properly from the slice")
	// }
	if d.R.{{$relAlias.Foreign}}[0] != &a {
		t.Error("relationship was not added properly to the slice")
	}
	if e.R.{{$relAlias.Foreign}}[0] != &a {
		t.Error("relationship was not added properly to the slice")
	}
	{{- else}}

	if !queries.IsValuerNil(b.{{$fcolField}}) {
		t.Error("want b's foreign key value to be nil")
	}
	if !queries.IsValuerNil(c.{{$fcolField}}) {
		t.Error("want c's foreign key value to be nil")
	}
	{{if $usesPrimitives -}}
	if a.{{$colField}} != d.{{$fcolField}} {
		t.Error("foreign key was wrong value", a.{{$colField}}, d.{{$fcolField}})
	}
	if a.{{$colField}} != e.{{$fcolField}} {
		t.Error("foreign key was wrong value", a.{{$colField}}, e.{{$fcolField}})
	}
	{{else -}}
	if !queries.Equal(a.{{$colField}}, d.{{$fcolField}}) {
		t.Error("foreign key was wrong value", a.{{$colField}}, d.{{$fcolField}})
	}
	if !queries.Equal(a.{{$colField}}, e.{{$fcolField}}) {
		t.Error("foreign key was wrong value", a.{{$colField}}, e.{{$fcolField}})
	}
	{{- end}}

	if b.R.{{$relAlias.Foreign}} != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if c.R.{{$relAlias.Foreign}} != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if d.R.{{$relAlias.Foreign}} != &a {
		t.Error("relationship was not added properly to the foreign struct")
	}
	if e.R.{{$relAlias.Foreign}} != &a {
		t.Error("relationship was not added properly to the foreign struct")
	}
	{{- end}}

	if a.R.{{$relAlias.Local}}[0] != &d {
		t.Error("relationship struct slice not set to correct value")
	}
	if a.R.{{$relAlias.Local}}[1] != &e {
		t.Error("relationship struct slice not set to correct value")
	}
}

func test{{$ltable.UpSingular}}ToManyRemoveOp{{$relAlias.Local}}(t *testing.T) {
	var err error

	{{if not $.NoContext}}ctx := context.Background(){{end}}
	tx := MustTx({{if $.NoContext}}boil.Begin(){{else}}boil.BeginTx(ctx, nil){{end}})
	defer func() { _ = tx.Rollback() }()

	var a {{$ltable.UpSingular}}
	var b, c, d, e {{$ftable.UpSingular}}

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, {{$ltable.DownSingular}}DBTypes, false, strmangle.SetComplement({{$ltable.DownSingular}}PrimaryKeyColumns, {{$ltable.DownSingular}}ColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*{{$ftable.UpSingular}}{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, {{$ftable.DownSingular}}DBTypes, false, strmangle.SetComplement({{$ftable.DownSingular}}PrimaryKeyColumns, {{$ftable.DownSingular}}ColumnsWithoutDefault)...); err != nil {
			t.Fatal(err)
		}
	}

	if err := a.Insert({{if not $.NoContext}}ctx, {{end -}} tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	err = a.Add{{$relAlias.Local}}({{if not $.NoContext}}ctx, {{end -}} tx, true, foreigners...)
	if err != nil {
		t.Fatal(err)
	}

	count, err := a.{{$relAlias.Local}}().Count({{if not $.NoContext}}ctx, {{end -}} tx)
	if err != nil {
		t.Fatal(err)
	}
	if count != 4 {
		t.Error("count was wrong:", count)
	}

	err = a.Remove{{$relAlias.Local}}({{if not $.NoContext}}ctx, {{end -}} tx, foreigners[:2]...)
	if err != nil {
		t.Fatal(err)
	}

	count, err = a.{{$relAlias.Local}}().Count({{if not $.NoContext}}ctx, {{end -}} tx)
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	{{- if .ToJoinTable}}

	if len(b.R.{{$relAlias.Foreign}}) != 0 {
		t.Error("relationship was not removed properly from the slice")
	}
	if len(c.R.{{$relAlias.Foreign}}) != 0 {
		t.Error("relationship was not removed properly from the slice")
	}
	if d.R.{{$relAlias.Foreign}}[0] != &a {
		t.Error("relationship was not added properly to the foreign struct")
	}
	if e.R.{{$relAlias.Foreign}}[0] != &a {
		t.Error("relationship was not added properly to the foreign struct")
	}
	{{- else}}

	if !queries.IsValuerNil(b.{{$fcolField}}) {
		t.Error("want b's foreign key value to be nil")
	}
	if !queries.IsValuerNil(c.{{$fcolField}}) {
		t.Error("want c's foreign key value to be nil")
	}

	if b.R.{{$relAlias.Foreign}} != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if c.R.{{$relAlias.Foreign}} != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if d.R.{{$relAlias.Foreign}} != &a {
		t.Error("relationship to a should have been preserved")
	}
	if e.R.{{$relAlias.Foreign}} != &a {
		t.Error("relationship to a should have been preserved")
	}
	{{- end}}

	if len(a.R.{{$relAlias.Local}}) != 2 {
		t.Error("should have preserved two relationships")
	}

	// Removal doesn't do a stable deletion for performance so we have to flip the order
	if a.R.{{$relAlias.Local}}[1] != &d {
		t.Error("relationship to d should have been preserved")
	}
	if a.R.{{$relAlias.Local}}[0] != &e {
		t.Error("relationship to e should have been preserved")
	}
}
{{end -}}
{{- end -}}{{- /* range relationships */ -}}
{{- end -}}{{- /* outer if join table */ -}}
