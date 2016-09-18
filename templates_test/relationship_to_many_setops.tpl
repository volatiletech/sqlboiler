{{- if .Table.IsJoinTable -}}
{{- else -}}
	{{- $dot := . -}}
	{{- $table := .Table -}}
	{{- range .Table.ToManyRelationships -}}
	{{- $varNameSingular := .Table | singular | camelCase -}}
	{{- $foreignVarNameSingular := .ForeignTable | singular | camelCase -}}
	{{- $rel := textsFromRelationship $dot.Tables $table .}}
func test{{$rel.LocalTable.NameGo}}ToManyAddOp{{$rel.Function.Name}}(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a {{$rel.LocalTable.NameGo}}
	var b, c, d, e {{$rel.ForeignTable.NameGo}}

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, {{$varNameSingular}}DBTypes, false, strmangle.SetComplement({{$varNameSingular}}PrimaryKeyColumns, {{$varNameSingular}}ColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*{{$rel.ForeignTable.NameGo}}{&b, &c, &d, &e}
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

	foreignersSplitByInsertion := [][]*{{$rel.ForeignTable.NameGo}}{
		{&b, &c},
		{&d, &e},
	}

	for i, x := range foreignersSplitByInsertion {
		err = a.Add{{$rel.Function.Name}}(tx, i != 0, x...)
		if err != nil {
			t.Fatal(err)
		}

		first := x[0]
		second := x[1]
		{{- if .ToJoinTable}}

		if first.R.{{$rel.Function.ForeignName}}[0] != &a {
			t.Error("relationship was not added properly to the slice")
		}
		if second.R.{{$rel.Function.ForeignName}}[0] != &a {
			t.Error("relationship was not added properly to the slice")
		}
		{{- else}}

		{{if $rel.Function.UsesBytes -}}
		if 0 != bytes.Compare(a.{{$rel.Function.LocalAssignment}}, first.{{$rel.Function.ForeignAssignment}}) {
			t.Error("foreign key was wrong value", a.{{$rel.Function.LocalAssignment}}, first.{{$rel.Function.ForeignAssignment}})
		}
		if 0 != bytes.Compare(a.{{$rel.Function.LocalAssignment}}, second.{{$rel.Function.ForeignAssignment}}) {
			t.Error("foreign key was wrong value", a.{{$rel.Function.LocalAssignment}}, second.{{$rel.Function.ForeignAssignment}})
		}
		{{else -}}
		if a.{{$rel.Function.LocalAssignment}} != first.{{$rel.Function.ForeignAssignment}} {
			t.Error("foreign key was wrong value", a.{{$rel.Function.LocalAssignment}}, first.{{$rel.Function.ForeignAssignment}})
		}
		if a.{{$rel.Function.LocalAssignment}} != second.{{$rel.Function.ForeignAssignment}} {
			t.Error("foreign key was wrong value", a.{{$rel.Function.LocalAssignment}}, second.{{$rel.Function.ForeignAssignment}})
		}
		{{- end}}

		if first.R.{{$rel.Function.ForeignName}} != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}
		if second.R.{{$rel.Function.ForeignName}} != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}
		{{- end}}

		if a.R.{{$rel.Function.Name}}[i*2] != first {
			t.Error("relationship struct slice not set to correct value")
		}
		if a.R.{{$rel.Function.Name}}[i*2+1] != second {
			t.Error("relationship struct slice not set to correct value")
		}

		count, err := a.{{$rel.Function.Name}}(tx).Count()
		if err != nil {
			t.Fatal(err)
		}
		if want := int64((i+1)*2); count != want {
			t.Error("want", want, "got", count)
		}
	}
}
{{- if (or .ForeignColumnNullable .ToJoinTable)}}

func test{{$rel.LocalTable.NameGo}}ToManySetOp{{$rel.Function.Name}}(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a {{$rel.LocalTable.NameGo}}
	var b, c, d, e {{$rel.ForeignTable.NameGo}}

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, {{$varNameSingular}}DBTypes, false, strmangle.SetComplement({{$varNameSingular}}PrimaryKeyColumns, {{$varNameSingular}}ColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*{{$rel.ForeignTable.NameGo}}{&b, &c, &d, &e}
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

	err = a.Set{{$rel.Function.Name}}(tx, false, &b, &c)
	if err != nil {
		t.Fatal(err)
	}

	count, err := a.{{$rel.Function.Name}}(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	err = a.Set{{$rel.Function.Name}}(tx, true, &d, &e)
	if err != nil {
		t.Fatal(err)
	}

	count, err = a.{{$rel.Function.Name}}(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	{{- if .ToJoinTable}}

	if len(b.R.{{$rel.Function.ForeignName}}) != 0 {
		t.Error("relationship was not removed properly from the slice")
	}
	if len(c.R.{{$rel.Function.ForeignName}}) != 0 {
		t.Error("relationship was not removed properly from the slice")
	}
	if d.R.{{$rel.Function.ForeignName}}[0] != &a {
		t.Error("relationship was not added properly to the slice")
	}
	if e.R.{{$rel.Function.ForeignName}}[0] != &a {
		t.Error("relationship was not added properly to the slice")
	}
	{{- else}}

	if b.{{$rel.ForeignTable.ColumnNameGo}}.Valid {
		t.Error("want b's foreign key value to be nil")
	}
	if c.{{$rel.ForeignTable.ColumnNameGo}}.Valid {
		t.Error("want c's foreign key value to be nil")
	}
	{{if $rel.Function.UsesBytes -}}
	if 0 != bytes.Compare(a.{{$rel.Function.LocalAssignment}}, d.{{$rel.Function.ForeignAssignment}}) {
		t.Error("foreign key was wrong value", a.{{$rel.Function.LocalAssignment}}, d.{{$rel.Function.ForeignAssignment}})
	}
	if 0 != bytes.Compare(a.{{$rel.Function.LocalAssignment}}, e.{{$rel.Function.ForeignAssignment}}) {
		t.Error("foreign key was wrong value", a.{{$rel.Function.LocalAssignment}}, e.{{$rel.Function.ForeignAssignment}})
	}
	{{else -}}
	if a.{{$rel.Function.LocalAssignment}} != d.{{$rel.Function.ForeignAssignment}} {
		t.Error("foreign key was wrong value", a.{{$rel.Function.LocalAssignment}}, d.{{$rel.Function.ForeignAssignment}})
	}
	if a.{{$rel.Function.LocalAssignment}} != e.{{$rel.Function.ForeignAssignment}} {
		t.Error("foreign key was wrong value", a.{{$rel.Function.LocalAssignment}}, e.{{$rel.Function.ForeignAssignment}})
	}
	{{- end}}

	if b.R.{{$rel.Function.ForeignName}} != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if c.R.{{$rel.Function.ForeignName}} != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if d.R.{{$rel.Function.ForeignName}} != &a {
		t.Error("relationship was not added properly to the foreign struct")
	}
	if e.R.{{$rel.Function.ForeignName}} != &a {
		t.Error("relationship was not added properly to the foreign struct")
	}
	{{- end}}

	if a.R.{{$rel.Function.Name}}[0] != &d {
		t.Error("relationship struct slice not set to correct value")
	}
	if a.R.{{$rel.Function.Name}}[1] != &e {
		t.Error("relationship struct slice not set to correct value")
	}
}

func test{{$rel.LocalTable.NameGo}}ToManyRemoveOp{{$rel.Function.Name}}(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a {{$rel.LocalTable.NameGo}}
	var b, c, d, e {{$rel.ForeignTable.NameGo}}

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, {{$varNameSingular}}DBTypes, false, strmangle.SetComplement({{$varNameSingular}}PrimaryKeyColumns, {{$varNameSingular}}ColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*{{$rel.ForeignTable.NameGo}}{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, {{$foreignVarNameSingular}}DBTypes, false, strmangle.SetComplement({{$foreignVarNameSingular}}PrimaryKeyColumns, {{$foreignVarNameSingular}}ColumnsWithoutDefault)...); err != nil {
			t.Fatal(err)
		}
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	err = a.Add{{$rel.Function.Name}}(tx, true, foreigners...)
	if err != nil {
		t.Fatal(err)
	}

	count, err := a.{{$rel.Function.Name}}(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 4 {
		t.Error("count was wrong:", count)
	}

	err = a.Remove{{$rel.Function.Name}}(tx, foreigners[:2]...)
	if err != nil {
		t.Fatal(err)
	}

	count, err = a.{{$rel.Function.Name}}(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	{{- if .ToJoinTable}}

	if len(b.R.{{$rel.Function.ForeignName}}) != 0 {
		t.Error("relationship was not removed properly from the slice")
	}
	if len(c.R.{{$rel.Function.ForeignName}}) != 0 {
		t.Error("relationship was not removed properly from the slice")
	}
	if d.R.{{$rel.Function.ForeignName}}[0] != &a {
		t.Error("relationship was not added properly to the foreign struct")
	}
	if e.R.{{$rel.Function.ForeignName}}[0] != &a {
		t.Error("relationship was not added properly to the foreign struct")
	}
	{{- else}}

	if b.{{$rel.ForeignTable.ColumnNameGo}}.Valid {
		t.Error("want b's foreign key value to be nil")
	}
	if c.{{$rel.ForeignTable.ColumnNameGo}}.Valid {
		t.Error("want c's foreign key value to be nil")
	}

	if b.R.{{$rel.Function.ForeignName}} != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if c.R.{{$rel.Function.ForeignName}} != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if d.R.{{$rel.Function.ForeignName}} != &a {
		t.Error("relationship to a should have been preserved")
	}
	if e.R.{{$rel.Function.ForeignName}} != &a {
		t.Error("relationship to a should have been preserved")
	}
	{{- end}}

	if len(a.R.{{$rel.Function.Name}}) != 2 {
		t.Error("should have preserved two relationships")
	}

	// Removal doesn't do a stable deletion for performance so we have to flip the order
	if a.R.{{$rel.Function.Name}}[1] != &d {
		t.Error("relationship to d should have been preserved")
	}
	if a.R.{{$rel.Function.Name}}[0] != &e {
		t.Error("relationship to e should have been preserved")
	}
}
{{end -}}
{{- end -}}{{- /* range relationships */ -}}
{{- end -}}{{- /* outer if join table */ -}}
