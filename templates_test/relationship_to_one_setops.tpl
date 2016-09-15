{{- define "relationship_to_one_setops_test_helper" -}}
{{- $varNameSingular := .ForeignKey.Table | singular | camelCase -}}
{{- $foreignVarNameSingular := .ForeignKey.ForeignTable | singular | camelCase -}}
func test{{.LocalTable.NameGo}}ToOneSetOp{{.ForeignTable.NameGo}}_{{.Function.Name}}(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a {{.LocalTable.NameGo}}
	var b, c {{.ForeignTable.NameGo}}

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, {{$varNameSingular}}DBTypes, false, {{$varNameSingular}}PrimaryKeyColumns...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, {{$foreignVarNameSingular}}DBTypes, false, {{$foreignVarNameSingular}}PrimaryKeyColumns...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, {{$foreignVarNameSingular}}DBTypes, false, {{$foreignVarNameSingular}}PrimaryKeyColumns...); err != nil {
		t.Fatal(err)
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}

	for i, x := range []*{{.ForeignTable.NameGo}}{&b, &c} {
		err = a.Set{{.Function.Name}}(tx, i != 0, x)
		if err != nil {
			t.Fatal(err)
		}

		if a.{{.Function.LocalAssignment}} != x.{{.Function.ForeignAssignment}} {
			t.Error("foreign key was wrong value", a.{{.Function.LocalAssignment}})
		}
		if a.R.{{.Function.Name}} != x {
			t.Error("relationship struct not set to correct value")
		}

		zero := reflect.Zero(reflect.TypeOf(a.{{.Function.LocalAssignment}}))
		reflect.Indirect(reflect.ValueOf(&a.{{.Function.LocalAssignment}})).Set(zero)

		if err = a.Reload(tx); err != nil {
			t.Fatal("failed to reload", err)
		}

		if a.{{.Function.LocalAssignment}} != x.{{.Function.ForeignAssignment}} {
			t.Error("foreign key was wrong value", a.{{.Function.LocalAssignment}}, x.{{.Function.ForeignAssignment}})
		}

		{{if .ForeignKey.Unique -}}
		if x.R.{{.Function.ForeignName}} != &a {
			t.Error("failed to append to foreign relationship struct")
		}
		{{else -}}
		if x.R.{{.Function.ForeignName}}[0] != &a {
			t.Error("failed to append to foreign relationship struct")
		}
		{{end -}}
	}
}
{{- if .ForeignKey.Nullable}}

func test{{.LocalTable.NameGo}}ToOneRemoveOp{{.ForeignTable.NameGo}}_{{.Function.Name}}(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a {{.LocalTable.NameGo}}
	var b {{.ForeignTable.NameGo}}

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, {{$varNameSingular}}DBTypes, false, {{$varNameSingular}}PrimaryKeyColumns...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, {{$foreignVarNameSingular}}DBTypes, false, {{$foreignVarNameSingular}}PrimaryKeyColumns...); err != nil {
		t.Fatal(err)
	}

	if err = a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	if err = a.Set{{.Function.Name}}(tx, true, &b); err != nil {
		t.Fatal(err)
	}

	if err = a.Remove{{.Function.Name}}(tx, &b); err != nil {
		t.Error("failed to remove relationship")
	}

	count, err := a.{{.Function.Name}}(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 0 {
		t.Error("want no relationships remaining")
	}

	if a.R.{{.Function.Name}} != nil {
		t.Error("R struct entry should be nil")
	}

	if a.{{.LocalTable.ColumnNameGo}}.Valid {
		t.Error("R struct entry should be nil")
	}

	{{if .ForeignKey.Unique -}}
	if b.R.{{.Function.ForeignName}} != nil {
		t.Error("failed to remove a from b's relationships")
	}
	{{else -}}
	if len(b.R.{{.Function.ForeignName}}) != 0 {
		t.Error("failed to remove a from b's relationships")
	}
	{{end -}}
}
{{end -}}
{{- end -}}
{{- if .Table.IsJoinTable -}}
{{- else -}}
	{{- $dot := . -}}
	{{- range .Table.FKeys -}}
		{{- $rel := textsFromForeignKey $dot.PkgName $dot.Tables $dot.Table .}}

{{template "relationship_to_one_setops_test_helper" $rel -}}
{{- end -}}
{{- end -}}
