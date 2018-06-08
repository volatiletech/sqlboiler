{{- if .Table.IsJoinTable -}}
{{- else -}}
	{{- range .Table.ToOneRelationships -}}
		{{- $txt := txtsFromOneToOne $.Tables $.Table .}}
{{- $varNameSingular := .Table | singular | camelCase -}}
{{- $foreignVarNameSingular := .ForeignTable | singular | camelCase -}}
{{- $foreignPKeyCols := (getTable $.Tables .ForeignTable).PKey.Columns}}
func test{{$txt.LocalTable.NameGo}}OneToOneSetOp{{$txt.ForeignTable.NameGo}}Using{{$txt.Function.Name}}(t *testing.T) {
	var err error

	{{if not $.NoContext}}ctx := context.Background(){{end}}
	tx := MustTx({{if $.NoContext}}boil.Begin(){{else}}boil.BeginTx(ctx, nil){{end}})
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

	if err := a.Insert({{if not $.NoContext}}ctx, {{end -}} tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert({{if not $.NoContext}}ctx, {{end -}} tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	for i, x := range []*{{$txt.ForeignTable.NameGo}}{&b, &c} {
		err = a.Set{{$txt.Function.Name}}({{if not $.NoContext}}ctx, {{end -}} tx, i != 0, x)
		if err != nil {
			t.Fatal(err)
		}

		if a.R.{{$txt.Function.Name}} != x {
			t.Error("relationship struct not set to correct value")
		}
		if x.R.{{$txt.Function.ForeignName}} != &a {
			t.Error("failed to append to foreign relationship struct")
		}

		{{if $txt.Function.UsesPrimitives -}}
		if a.{{$txt.LocalTable.ColumnNameGo}} != x.{{$txt.ForeignTable.ColumnNameGo}} {
		{{else -}}
		if !queries.Equal(a.{{$txt.LocalTable.ColumnNameGo}}, x.{{$txt.ForeignTable.ColumnNameGo}}) {
		{{end -}}
			t.Error("foreign key was wrong value", a.{{$txt.LocalTable.ColumnNameGo}})
		}

		{{if setInclude .ForeignColumn $foreignPKeyCols -}}
		if exists, err := {{$txt.ForeignTable.NameGo}}Exists({{if not $.NoContext}}ctx, {{end -}} tx, x.{{$foreignPKeyCols | stringMap $.StringFuncs.titleCase | join ", x."}}); err != nil {
			t.Fatal(err)
		} else if !exists {
			t.Error("want 'x' to exist")
		}
		{{else -}}
		zero := reflect.Zero(reflect.TypeOf(x.{{$txt.ForeignTable.ColumnNameGo}}))
		reflect.Indirect(reflect.ValueOf(&x.{{$txt.ForeignTable.ColumnNameGo}})).Set(zero)

		if err = x.Reload({{if not $.NoContext}}ctx, {{end -}} tx); err != nil {
			t.Fatal("failed to reload", err)
		}
		{{- end}}

		{{if $txt.Function.UsesPrimitives -}}
		if a.{{$txt.LocalTable.ColumnNameGo}} != x.{{$txt.ForeignTable.ColumnNameGo}} {
		{{else -}}
		if !queries.Equal(a.{{$txt.LocalTable.ColumnNameGo}}, x.{{$txt.ForeignTable.ColumnNameGo}}) {
		{{end -}}
			t.Error("foreign key was wrong value", a.{{$txt.LocalTable.ColumnNameGo}}, x.{{$txt.ForeignTable.ColumnNameGo}})
		}

		if {{if not $.NoRowsAffected}}_, {{end -}} err = x.Delete({{if not $.NoContext}}ctx, {{end -}} tx); err != nil {
			t.Fatal("failed to delete x", err)
		}
	}
}
{{- if .ForeignColumnNullable}}

func test{{$txt.LocalTable.NameGo}}OneToOneRemoveOp{{$txt.ForeignTable.NameGo}}Using{{$txt.Function.Name}}(t *testing.T) {
	var err error

	{{if not $.NoContext}}ctx := context.Background(){{end}}
	tx := MustTx({{if $.NoContext}}boil.Begin(){{else}}boil.BeginTx(ctx, nil){{end}})
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

	if err = a.Insert({{if not $.NoContext}}ctx, {{end -}} tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	if err = a.Set{{$txt.Function.Name}}({{if not $.NoContext}}ctx, {{end -}} tx, true, &b); err != nil {
		t.Fatal(err)
	}

	if err = a.Remove{{$txt.Function.Name}}({{if not $.NoContext}}ctx, {{end -}} tx, &b); err != nil {
		t.Error("failed to remove relationship")
	}

	count, err := a.{{$txt.Function.Name}}().Count({{if not $.NoContext}}ctx, {{end -}} tx)
	if err != nil {
		t.Error(err)
	}
	if count != 0 {
		t.Error("want no relationships remaining")
	}

	if a.R.{{$txt.Function.Name}} != nil {
		t.Error("R struct entry should be nil")
	}

	if !queries.IsValuerNil(b.{{$txt.ForeignTable.ColumnNameGo}}) {
		t.Error("foreign key column should be nil")
	}

	if b.R.{{$txt.Function.ForeignName}} != nil {
		t.Error("failed to remove a from b's relationships")
	}
}
{{end -}}{{/* end if foreign key nullable */}}
{{- end -}}{{/* range */}}
{{- end -}}{{/* join table */}}
