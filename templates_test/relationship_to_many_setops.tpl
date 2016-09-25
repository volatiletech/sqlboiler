{{- if .Table.IsJoinTable -}}
{{- else -}}
	{{- $dot := . -}}
	{{- $table := .Table -}}
	{{- range .Table.ToManyRelationships -}}
	{{- $varNameSingular := .Table | singular | camelCase -}}
	{{- $foreignVarNameSingular := .ForeignTable | singular | camelCase -}}
	{{- $txt := txtsFromToMany $dot.Tables $table .}}
func test{{$txt.LocalTable.NameGo}}ToManyAddOp{{$txt.Function.Name}}(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a {{$txt.LocalTable.NameGo}}
	var b, c, d, e {{$txt.ForeignTable.NameGo}}

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, {{$varNameSingular}}DBTypes, false, strmangle.SetComplement({{$varNameSingular}}PrimaryKeyColumns, {{$varNameSingular}}ColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*{{$txt.ForeignTable.NameGo}}{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, {{$foreignVarNameSingular}}DBTypes, false, strmangle.SetComplement({{$foreignVarNameSingular}}PrimaryKeyColumns, {{$foreignVarNameSingular}}ColumnsWithoutDefault)...); err != nil {
			t.Fatal(err)
		}
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(tx); err != nil {
		t.Fatal(err)
	}

	foreignersSplitByInsertion := [][]*{{$txt.ForeignTable.NameGo}}{
		{&b, &c},
		{&d, &e},
	}

	for i, x := range foreignersSplitByInsertion {
		err = a.Add{{$txt.Function.Name}}(tx, i != 0, x...)
		if err != nil {
			t.Fatal(err)
		}

		first := x[0]
		second := x[1]
		{{- if .ToJoinTable}}

		if first.R.{{$txt.Function.ForeignName}}[0] != &a {
			t.Error("relationship was not added properly to the slice")
		}
		if second.R.{{$txt.Function.ForeignName}}[0] != &a {
			t.Error("relationship was not added properly to the slice")
		}
		{{- else}}

		{{if $txt.Function.UsesBytes -}}
		if 0 != bytes.Compare(a.{{$txt.Function.LocalAssignment}}, first.{{$txt.Function.ForeignAssignment}}) {
			t.Error("foreign key was wrong value", a.{{$txt.Function.LocalAssignment}}, first.{{$txt.Function.ForeignAssignment}})
		}
		if 0 != bytes.Compare(a.{{$txt.Function.LocalAssignment}}, second.{{$txt.Function.ForeignAssignment}}) {
			t.Error("foreign key was wrong value", a.{{$txt.Function.LocalAssignment}}, second.{{$txt.Function.ForeignAssignment}})
		}
		{{else -}}
		if a.{{$txt.Function.LocalAssignment}} != first.{{$txt.Function.ForeignAssignment}} {
			t.Error("foreign key was wrong value", a.{{$txt.Function.LocalAssignment}}, first.{{$txt.Function.ForeignAssignment}})
		}
		if a.{{$txt.Function.LocalAssignment}} != second.{{$txt.Function.ForeignAssignment}} {
			t.Error("foreign key was wrong value", a.{{$txt.Function.LocalAssignment}}, second.{{$txt.Function.ForeignAssignment}})
		}
		{{- end}}

		if first.R.{{$txt.Function.ForeignName}} != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}
		if second.R.{{$txt.Function.ForeignName}} != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}
		{{- end}}

		if a.R.{{$txt.Function.Name}}[i*2] != first {
			t.Error("relationship struct slice not set to correct value")
		}
		if a.R.{{$txt.Function.Name}}[i*2+1] != second {
			t.Error("relationship struct slice not set to correct value")
		}

		count, err := a.{{$txt.Function.Name}}(tx).Count()
		if err != nil {
			t.Fatal(err)
		}
		if want := int64((i+1)*2); count != want {
			t.Error("want", want, "got", count)
		}
	}
}
{{- if (or .ForeignColumnNullable .ToJoinTable)}}

func test{{$txt.LocalTable.NameGo}}ToManySetOp{{$txt.Function.Name}}(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a {{$txt.LocalTable.NameGo}}
	var b, c, d, e {{$txt.ForeignTable.NameGo}}

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, {{$varNameSingular}}DBTypes, false, strmangle.SetComplement({{$varNameSingular}}PrimaryKeyColumns, {{$varNameSingular}}ColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*{{$txt.ForeignTable.NameGo}}{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, {{$foreignVarNameSingular}}DBTypes, false, strmangle.SetComplement({{$foreignVarNameSingular}}PrimaryKeyColumns, {{$foreignVarNameSingular}}ColumnsWithoutDefault)...); err != nil {
			t.Fatal(err)
		}
	}

	if err = a.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(tx); err != nil {
		t.Fatal(err)
	}

	err = a.Set{{$txt.Function.Name}}(tx, false, &b, &c)
	if err != nil {
		t.Fatal(err)
	}

	count, err := a.{{$txt.Function.Name}}(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	err = a.Set{{$txt.Function.Name}}(tx, true, &d, &e)
	if err != nil {
		t.Fatal(err)
	}

	count, err = a.{{$txt.Function.Name}}(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	{{- if .ToJoinTable}}

	if len(b.R.{{$txt.Function.ForeignName}}) != 0 {
		t.Error("relationship was not removed properly from the slice")
	}
	if len(c.R.{{$txt.Function.ForeignName}}) != 0 {
		t.Error("relationship was not removed properly from the slice")
	}
	if d.R.{{$txt.Function.ForeignName}}[0] != &a {
		t.Error("relationship was not added properly to the slice")
	}
	if e.R.{{$txt.Function.ForeignName}}[0] != &a {
		t.Error("relationship was not added properly to the slice")
	}
	{{- else}}

	if b.{{$txt.ForeignTable.ColumnNameGo}}.Valid {
		t.Error("want b's foreign key value to be nil")
	}
	if c.{{$txt.ForeignTable.ColumnNameGo}}.Valid {
		t.Error("want c's foreign key value to be nil")
	}
	{{if $txt.Function.UsesBytes -}}
	if 0 != bytes.Compare(a.{{$txt.Function.LocalAssignment}}, d.{{$txt.Function.ForeignAssignment}}) {
		t.Error("foreign key was wrong value", a.{{$txt.Function.LocalAssignment}}, d.{{$txt.Function.ForeignAssignment}})
	}
	if 0 != bytes.Compare(a.{{$txt.Function.LocalAssignment}}, e.{{$txt.Function.ForeignAssignment}}) {
		t.Error("foreign key was wrong value", a.{{$txt.Function.LocalAssignment}}, e.{{$txt.Function.ForeignAssignment}})
	}
	{{else -}}
	if a.{{$txt.Function.LocalAssignment}} != d.{{$txt.Function.ForeignAssignment}} {
		t.Error("foreign key was wrong value", a.{{$txt.Function.LocalAssignment}}, d.{{$txt.Function.ForeignAssignment}})
	}
	if a.{{$txt.Function.LocalAssignment}} != e.{{$txt.Function.ForeignAssignment}} {
		t.Error("foreign key was wrong value", a.{{$txt.Function.LocalAssignment}}, e.{{$txt.Function.ForeignAssignment}})
	}
	{{- end}}

	if b.R.{{$txt.Function.ForeignName}} != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if c.R.{{$txt.Function.ForeignName}} != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if d.R.{{$txt.Function.ForeignName}} != &a {
		t.Error("relationship was not added properly to the foreign struct")
	}
	if e.R.{{$txt.Function.ForeignName}} != &a {
		t.Error("relationship was not added properly to the foreign struct")
	}
	{{- end}}

	if a.R.{{$txt.Function.Name}}[0] != &d {
		t.Error("relationship struct slice not set to correct value")
	}
	if a.R.{{$txt.Function.Name}}[1] != &e {
		t.Error("relationship struct slice not set to correct value")
	}
}

func test{{$txt.LocalTable.NameGo}}ToManyRemoveOp{{$txt.Function.Name}}(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a {{$txt.LocalTable.NameGo}}
	var b, c, d, e {{$txt.ForeignTable.NameGo}}

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, {{$varNameSingular}}DBTypes, false, strmangle.SetComplement({{$varNameSingular}}PrimaryKeyColumns, {{$varNameSingular}}ColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*{{$txt.ForeignTable.NameGo}}{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, {{$foreignVarNameSingular}}DBTypes, false, strmangle.SetComplement({{$foreignVarNameSingular}}PrimaryKeyColumns, {{$foreignVarNameSingular}}ColumnsWithoutDefault)...); err != nil {
			t.Fatal(err)
		}
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	err = a.Add{{$txt.Function.Name}}(tx, true, foreigners...)
	if err != nil {
		t.Fatal(err)
	}

	count, err := a.{{$txt.Function.Name}}(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 4 {
		t.Error("count was wrong:", count)
	}

	err = a.Remove{{$txt.Function.Name}}(tx, foreigners[:2]...)
	if err != nil {
		t.Fatal(err)
	}

	count, err = a.{{$txt.Function.Name}}(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	{{- if .ToJoinTable}}

	if len(b.R.{{$txt.Function.ForeignName}}) != 0 {
		t.Error("relationship was not removed properly from the slice")
	}
	if len(c.R.{{$txt.Function.ForeignName}}) != 0 {
		t.Error("relationship was not removed properly from the slice")
	}
	if d.R.{{$txt.Function.ForeignName}}[0] != &a {
		t.Error("relationship was not added properly to the foreign struct")
	}
	if e.R.{{$txt.Function.ForeignName}}[0] != &a {
		t.Error("relationship was not added properly to the foreign struct")
	}
	{{- else}}

	if b.{{$txt.ForeignTable.ColumnNameGo}}.Valid {
		t.Error("want b's foreign key value to be nil")
	}
	if c.{{$txt.ForeignTable.ColumnNameGo}}.Valid {
		t.Error("want c's foreign key value to be nil")
	}

	if b.R.{{$txt.Function.ForeignName}} != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if c.R.{{$txt.Function.ForeignName}} != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if d.R.{{$txt.Function.ForeignName}} != &a {
		t.Error("relationship to a should have been preserved")
	}
	if e.R.{{$txt.Function.ForeignName}} != &a {
		t.Error("relationship to a should have been preserved")
	}
	{{- end}}

	if len(a.R.{{$txt.Function.Name}}) != 2 {
		t.Error("should have preserved two relationships")
	}

	// Removal doesn't do a stable deletion for performance so we have to flip the order
	if a.R.{{$txt.Function.Name}}[1] != &d {
		t.Error("relationship to d should have been preserved")
	}
	if a.R.{{$txt.Function.Name}}[0] != &e {
		t.Error("relationship to e should have been preserved")
	}
}
{{end -}}
{{- end -}}{{- /* range relationships */ -}}
{{- end -}}{{- /* outer if join table */ -}}
