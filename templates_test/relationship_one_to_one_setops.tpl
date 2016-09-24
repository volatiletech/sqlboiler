{{- if .Table.IsJoinTable -}}
{{- else -}}
	{{- $dot := . -}}
	{{- range .Table.ToOneRelationships -}}
		{{- $txt := txtsFromOneToOne $dot.Tables $dot.Table .}}
{{- $varNameSingular := .Table | singular | camelCase -}}
{{- $foreignVarNameSingular := .ForeignTable | singular | camelCase -}}
{{- $foreignPKeyCols := (getTable $dot.Tables .ForeignTable).PKey.Columns}}
func test{{$txt.LocalTable.NameGo}}OneToOneSetOp{{$txt.ForeignTable.NameGo}}Using{{$txt.Function.Name}}(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a {{$txt.LocalTable.NameGo}}
	var b, c {{$txt.ForeignTable.NameGo}}

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, {{$varNameSingular}}DBTypes, false, strmangle.SetComplement({{$varNameSingular}}PrimaryKeyColumns, {{$varNameSingular}}ColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, {{$foreignVarNameSingular}}DBTypes, false, strmangle.SetComplement({{$foreignVarNameSingular}}PrimaryKeyColumns, {{$foreignVarNameSingular}}ColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, {{$foreignVarNameSingular}}DBTypes, false, strmangle.SetComplement({{$foreignVarNameSingular}}PrimaryKeyColumns, {{$foreignVarNameSingular}}ColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}

	for i, x := range []*{{$txt.ForeignTable.NameGo}}{&b, &c} {
		err = a.Set{{$txt.Function.Name}}(tx, i != 0, x)
		if err != nil {
			t.Fatal(err)
		}

		if a.R.{{$txt.Function.Name}} != x {
			t.Error("relationship struct not set to correct value")
		}
		if x.R.{{$txt.Function.ForeignName}} != &a {
			t.Error("failed to append to foreign relationship struct")
		}

		{{if $txt.Function.UsesBytes -}}
		if 0 != bytes.Compare(a.{{$txt.Function.LocalAssignment}}, x.{{$txt.Function.ForeignAssignment}}) {
		{{else -}}
		if a.{{$txt.Function.LocalAssignment}} != x.{{$txt.Function.ForeignAssignment}} {
		{{end -}}
			t.Error("foreign key was wrong value", a.{{$txt.Function.LocalAssignment}})
		}

		{{if setInclude .ForeignColumn $foreignPKeyCols -}}
		if exists, err := {{$txt.ForeignTable.NameGo}}Exists(tx, x.{{$foreignPKeyCols | stringMap $dot.StringFuncs.titleCase | join ", x."}}); err != nil {
			t.Fatal(err)
		} else if !exists {
			t.Error("want 'x' to exist")
		}
		{{else -}}
		zero := reflect.Zero(reflect.TypeOf(x.{{$txt.Function.ForeignAssignment}}))
		reflect.Indirect(reflect.ValueOf(&x.{{$txt.Function.ForeignAssignment}})).Set(zero)

		if err = x.Reload(tx); err != nil {
			t.Fatal("failed to reload", err)
		}
		{{- end}}

		{{if $txt.Function.UsesBytes -}}
		if 0 != bytes.Compare(a.{{$txt.Function.LocalAssignment}}, x.{{$txt.Function.ForeignAssignment}}) {
		{{else -}}
		if a.{{$txt.Function.LocalAssignment}} != x.{{$txt.Function.ForeignAssignment}} {
		{{end -}}
			t.Error("foreign key was wrong value", a.{{$txt.Function.LocalAssignment}}, x.{{$txt.Function.ForeignAssignment}})
		}

		if err = x.Delete(tx); err != nil {
			t.Fatal("failed to delete x", err)
		}
	}
}
{{- if .ForeignColumnNullable}}

func test{{$txt.LocalTable.NameGo}}OneToOneRemoveOp{{$txt.ForeignTable.NameGo}}Using{{$txt.Function.Name}}(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a {{$txt.LocalTable.NameGo}}
	var b {{$txt.ForeignTable.NameGo}}

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, {{$varNameSingular}}DBTypes, false, strmangle.SetComplement({{$varNameSingular}}PrimaryKeyColumns, {{$varNameSingular}}ColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, {{$foreignVarNameSingular}}DBTypes, false, strmangle.SetComplement({{$foreignVarNameSingular}}PrimaryKeyColumns, {{$foreignVarNameSingular}}ColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err = a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	if err = a.Set{{$txt.Function.Name}}(tx, true, &b); err != nil {
		t.Fatal(err)
	}

	if err = a.Remove{{$txt.Function.Name}}(tx, &b); err != nil {
		t.Error("failed to remove relationship")
	}

	count, err := a.{{$txt.Function.Name}}(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 0 {
		t.Error("want no relationships remaining")
	}

	if a.R.{{$txt.Function.Name}} != nil {
		t.Error("R struct entry should be nil")
	}

	if b.{{$txt.ForeignTable.ColumnNameGo}}.Valid {
		t.Error("foreign key column should be nil")
	}

	if b.R.{{$txt.Function.ForeignName}} != nil {
		t.Error("failed to remove a from b's relationships")
	}
}
{{end -}}{{/* end if foreign key nullable */}}
{{- end -}}{{/* range */}}
{{- end -}}{{/* join table */}}
