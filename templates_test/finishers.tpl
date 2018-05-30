{{- $tableNameSingular := .Table.Name | singular | titleCase -}}
{{- $tableNamePlural := .Table.Name | plural | titleCase -}}
{{- $varNamePlural := .Table.Name | plural | camelCase -}}
{{- $varNameSingular := .Table.Name | singular | camelCase -}}
func test{{$tableNamePlural}}Bind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	{{$varNameSingular}} := &{{$tableNameSingular}}{}
	if err = randomize.Struct(seed, {{$varNameSingular}}, {{$varNameSingular}}DBTypes, true, {{$varNameSingular}}ColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize {{$tableNameSingular}} struct: %s", err)
	}

	{{if not .NoContext}}ctx := context.Background(){{end}}
	tx := MustTx({{if .NoContext}}boil.Begin(){{else}}boil.BeginTx(ctx, nil){{end}})
	defer tx.Rollback()
	if err = {{$varNameSingular}}.Insert({{if not .NoContext}}ctx, {{end -}} tx); err != nil {
		t.Error(err)
	}

	if err = {{$tableNamePlural}}().Bind({{if .NoContext}}nil{{else}}ctx{{end}}, tx, {{$varNameSingular}}); err != nil {
		t.Error(err)
	}
}

func test{{$tableNamePlural}}One(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	{{$varNameSingular}} := &{{$tableNameSingular}}{}
	if err = randomize.Struct(seed, {{$varNameSingular}}, {{$varNameSingular}}DBTypes, true, {{$varNameSingular}}ColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize {{$tableNameSingular}} struct: %s", err)
	}

	{{if not .NoContext}}ctx := context.Background(){{end}}
	tx := MustTx({{if .NoContext}}boil.Begin(){{else}}boil.BeginTx(ctx, nil){{end}})
	defer tx.Rollback()
	if err = {{$varNameSingular}}.Insert({{if not .NoContext}}ctx, {{end -}} tx); err != nil {
		t.Error(err)
	}

	if x, err := {{$tableNamePlural}}().One({{if not .NoContext}}ctx, {{end -}} tx); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func test{{$tableNamePlural}}All(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	{{$varNameSingular}}One := &{{$tableNameSingular}}{}
	{{$varNameSingular}}Two := &{{$tableNameSingular}}{}
	if err = randomize.Struct(seed, {{$varNameSingular}}One, {{$varNameSingular}}DBTypes, false, {{$varNameSingular}}ColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize {{$tableNameSingular}} struct: %s", err)
	}
	if err = randomize.Struct(seed, {{$varNameSingular}}Two, {{$varNameSingular}}DBTypes, false, {{$varNameSingular}}ColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize {{$tableNameSingular}} struct: %s", err)
	}

	{{if not .NoContext}}ctx := context.Background(){{end}}
	tx := MustTx({{if .NoContext}}boil.Begin(){{else}}boil.BeginTx(ctx, nil){{end}})
	defer tx.Rollback()
	if err = {{$varNameSingular}}One.Insert({{if not .NoContext}}ctx, {{end -}} tx); err != nil {
		t.Error(err)
	}
	if err = {{$varNameSingular}}Two.Insert({{if not .NoContext}}ctx, {{end -}} tx); err != nil {
		t.Error(err)
	}

	slice, err := {{$tableNamePlural}}().All({{if not .NoContext}}ctx, {{end -}} tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func test{{$tableNamePlural}}Count(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	{{$varNameSingular}}One := &{{$tableNameSingular}}{}
	{{$varNameSingular}}Two := &{{$tableNameSingular}}{}
	if err = randomize.Struct(seed, {{$varNameSingular}}One, {{$varNameSingular}}DBTypes, false, {{$varNameSingular}}ColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize {{$tableNameSingular}} struct: %s", err)
	}
	if err = randomize.Struct(seed, {{$varNameSingular}}Two, {{$varNameSingular}}DBTypes, false, {{$varNameSingular}}ColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize {{$tableNameSingular}} struct: %s", err)
	}

	{{if not .NoContext}}ctx := context.Background(){{end}}
	tx := MustTx({{if .NoContext}}boil.Begin(){{else}}boil.BeginTx(ctx, nil){{end}})
	defer tx.Rollback()
	if err = {{$varNameSingular}}One.Insert({{if not .NoContext}}ctx, {{end -}} tx); err != nil {
		t.Error(err)
	}
	if err = {{$varNameSingular}}Two.Insert({{if not .NoContext}}ctx, {{end -}} tx); err != nil {
		t.Error(err)
	}

	count, err := {{$tableNamePlural}}().Count({{if not .NoContext}}ctx, {{end -}} tx)
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}
