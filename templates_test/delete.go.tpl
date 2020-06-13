{{- $alias := .Aliases.Table .Table.Name -}}
{{- $canSoftDelete := .Table.CanSoftDelete -}}
{{- $soft := and .AddSoftDeletes $canSoftDelete }}
{{- $colNames := .Table.Columns | columnNames -}}
{{- $schemaTable := .Table.Name | .SchemaTable}}
{{if and $soft (containsAny $colNames "updated_at") -}}
func test{{$alias.UpPlural}}SoftDeleteWithUpdatedAt(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &{{$alias.UpSingular}}{}
	if err = randomize.Struct(seed, o, {{$alias.DownSingular}}DBTypes, true, {{$alias.DownSingular}}ColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize {{$alias.UpSingular}} struct: %s", err)
	}

	{{if not .NoContext}}ctx := context.Background(){{end}}
	tx := MustTx({{if .NoContext}}boil.Begin(){{else}}boil.BeginTx(ctx, nil){{end}})
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert({{if not .NoContext}}ctx, {{end -}} tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	{{if .NoRowsAffected -}}
	if err = o.Delete({{if not .NoContext}}ctx, {{end -}} tx, false); err != nil {
		t.Error(err)
	}

	{{else -}}
	if rowsAff, err := o.Delete({{if not .NoContext}}ctx, {{end -}} tx, false); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	{{end -}}

	count, err := {{$alias.UpPlural}}().Count({{if not .NoContext}}ctx, {{end -}} tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}

	oDeleted := &{{$alias.UpSingular}}{}
	{{$pkeyArgs := .Table.PKey.Columns | stringMap (aliasCols $alias) | prefixStringSlice (printf "%s." "o") | join ", " -}}
	err = {{$alias.DownSingular}}Query{NewQuery(
		qm.From("{{$schemaTable}}"),
		qm.Where("{{if .Dialect.UseIndexPlaceholders}}{{whereClause .LQ .RQ 1 .Table.PKey.Columns}}{{else}}{{whereClause .LQ .RQ 0 .Table.PKey.Columns}}{{end}}", {{$pkeyArgs}}),
	)}.Bind(ctx, tx, oDeleted)
	if err != nil {
		t.Error(err)
	}
	if oDeleted.DeletedAt.Time != oDeleted.UpdatedAt {
		t.Error("DeletedAt != UpdatedAt")
	}
}

func test{{$alias.UpPlural}}QuerySoftDeleteWithUpdatedAtAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &{{$alias.UpSingular}}{}
	if err = randomize.Struct(seed, o, {{$alias.DownSingular}}DBTypes, true, {{$alias.DownSingular}}ColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize {{$alias.UpSingular}} struct: %s", err)
	}

	{{if not .NoContext}}ctx := context.Background(){{end}}
	tx := MustTx({{if .NoContext}}boil.Begin(){{else}}boil.BeginTx(ctx, nil){{end}})
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert({{if not .NoContext}}ctx, {{end -}} tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	{{if .NoRowsAffected -}}
	if err = {{$alias.UpPlural}}().DeleteAll({{if not .NoContext}}ctx, {{end -}} tx, false); err != nil {
		t.Error(err)
	}

	{{else -}}
	if rowsAff, err := {{$alias.UpPlural}}().DeleteAll({{if not .NoContext}}ctx, {{end -}} tx, false); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	{{end -}}

	count, err := {{$alias.UpPlural}}().Count({{if not .NoContext}}ctx, {{end -}} tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func test{{$alias.UpPlural}}SliceSoftDeleteWithUpdatedAtAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &{{$alias.UpSingular}}{}
	if err = randomize.Struct(seed, o, {{$alias.DownSingular}}DBTypes, true, {{$alias.DownSingular}}ColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize {{$alias.UpSingular}} struct: %s", err)
	}

	{{if not .NoContext}}ctx := context.Background(){{end}}
	tx := MustTx({{if .NoContext}}boil.Begin(){{else}}boil.BeginTx(ctx, nil){{end}})
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert({{if not .NoContext}}ctx, {{end -}} tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := {{$alias.UpSingular}}Slice{{"{"}}o{{"}"}}

	{{if .NoRowsAffected -}}
	if err = slice.DeleteAll({{if not .NoContext}}ctx, {{end -}} tx, false); err != nil {
		t.Error(err)
	}

	{{else -}}
	if rowsAff, err := slice.DeleteAll({{if not .NoContext}}ctx, {{end -}} tx, false); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	{{end -}}

	count, err := {{$alias.UpPlural}}().Count({{if not .NoContext}}ctx, {{end -}} tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

{{ end -}}

{{if $soft -}}
func test{{$alias.UpPlural}}SoftDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &{{$alias.UpSingular}}{}
	if err = randomize.Struct(seed, o, {{$alias.DownSingular}}DBTypes, true, {{$alias.DownSingular}}ColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize {{$alias.UpSingular}} struct: %s", err)
	}

	{{if not .NoContext}}ctx := context.Background(){{end}}
	tx := MustTx({{if .NoContext}}boil.Begin(){{else}}boil.BeginTx(ctx, nil){{end}})
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert({{if not .NoContext}}ctx, {{end -}} tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	{{if .NoRowsAffected -}}
	if err = o.Delete({{if not .NoContext}}ctx, {{end -}} tx, false); err != nil {
		t.Error(err)
	}

	{{else -}}
	if rowsAff, err := o.Delete({{if not .NoContext}}ctx, {{end -}} tx, false); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	{{end -}}

	count, err := {{$alias.UpPlural}}().Count({{if not .NoContext}}ctx, {{end -}} tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func test{{$alias.UpPlural}}QuerySoftDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &{{$alias.UpSingular}}{}
	if err = randomize.Struct(seed, o, {{$alias.DownSingular}}DBTypes, true, {{$alias.DownSingular}}ColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize {{$alias.UpSingular}} struct: %s", err)
	}

	{{if not .NoContext}}ctx := context.Background(){{end}}
	tx := MustTx({{if .NoContext}}boil.Begin(){{else}}boil.BeginTx(ctx, nil){{end}})
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert({{if not .NoContext}}ctx, {{end -}} tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	{{if .NoRowsAffected -}}
	if err = {{$alias.UpPlural}}().DeleteAll({{if not .NoContext}}ctx, {{end -}} tx, false); err != nil {
		t.Error(err)
	}

	{{else -}}
	if rowsAff, err := {{$alias.UpPlural}}().DeleteAll({{if not .NoContext}}ctx, {{end -}} tx, false); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	{{end -}}

	count, err := {{$alias.UpPlural}}().Count({{if not .NoContext}}ctx, {{end -}} tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func test{{$alias.UpPlural}}SliceSoftDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &{{$alias.UpSingular}}{}
	if err = randomize.Struct(seed, o, {{$alias.DownSingular}}DBTypes, true, {{$alias.DownSingular}}ColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize {{$alias.UpSingular}} struct: %s", err)
	}

	{{if not .NoContext}}ctx := context.Background(){{end}}
	tx := MustTx({{if .NoContext}}boil.Begin(){{else}}boil.BeginTx(ctx, nil){{end}})
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert({{if not .NoContext}}ctx, {{end -}} tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := {{$alias.UpSingular}}Slice{{"{"}}o{{"}"}}

	{{if .NoRowsAffected -}}
	if err = slice.DeleteAll({{if not .NoContext}}ctx, {{end -}} tx, false); err != nil {
		t.Error(err)
	}

	{{else -}}
	if rowsAff, err := slice.DeleteAll({{if not .NoContext}}ctx, {{end -}} tx, false); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	{{end -}}

	count, err := {{$alias.UpPlural}}().Count({{if not .NoContext}}ctx, {{end -}} tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

{{end -}}

func test{{$alias.UpPlural}}Delete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &{{$alias.UpSingular}}{}
	if err = randomize.Struct(seed, o, {{$alias.DownSingular}}DBTypes, true, {{$alias.DownSingular}}ColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize {{$alias.UpSingular}} struct: %s", err)
	}

	{{if not .NoContext}}ctx := context.Background(){{end}}
	tx := MustTx({{if .NoContext}}boil.Begin(){{else}}boil.BeginTx(ctx, nil){{end}})
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert({{if not .NoContext}}ctx, {{end -}} tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	{{if .NoRowsAffected -}}
	if err = o.Delete({{if not .NoContext}}ctx, {{end -}} tx {{- if $soft}}, true{{end}}); err != nil {
		t.Error(err)
	}

	{{else -}}
	if rowsAff, err := o.Delete({{if not .NoContext}}ctx, {{end -}} tx {{- if $soft}}, true{{end}}); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	{{end -}}

	count, err := {{$alias.UpPlural}}().Count({{if not .NoContext}}ctx, {{end -}} tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func test{{$alias.UpPlural}}QueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &{{$alias.UpSingular}}{}
	if err = randomize.Struct(seed, o, {{$alias.DownSingular}}DBTypes, true, {{$alias.DownSingular}}ColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize {{$alias.UpSingular}} struct: %s", err)
	}

	{{if not .NoContext}}ctx := context.Background(){{end}}
	tx := MustTx({{if .NoContext}}boil.Begin(){{else}}boil.BeginTx(ctx, nil){{end}})
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert({{if not .NoContext}}ctx, {{end -}} tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	{{if .NoRowsAffected -}}
	if err = {{$alias.UpPlural}}().DeleteAll({{if not .NoContext}}ctx, {{end -}} tx {{- if $soft}}, true{{end}}); err != nil {
		t.Error(err)
	}

	{{else -}}
	if rowsAff, err := {{$alias.UpPlural}}().DeleteAll({{if not .NoContext}}ctx, {{end -}} tx {{- if $soft}}, true{{end}}); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	{{end -}}

	count, err := {{$alias.UpPlural}}().Count({{if not .NoContext}}ctx, {{end -}} tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func test{{$alias.UpPlural}}SliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &{{$alias.UpSingular}}{}
	if err = randomize.Struct(seed, o, {{$alias.DownSingular}}DBTypes, true, {{$alias.DownSingular}}ColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize {{$alias.UpSingular}} struct: %s", err)
	}

	{{if not .NoContext}}ctx := context.Background(){{end}}
	tx := MustTx({{if .NoContext}}boil.Begin(){{else}}boil.BeginTx(ctx, nil){{end}})
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert({{if not .NoContext}}ctx, {{end -}} tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := {{$alias.UpSingular}}Slice{{"{"}}o{{"}"}}

	{{if .NoRowsAffected -}}
	if err = slice.DeleteAll({{if not .NoContext}}ctx, {{end -}} tx {{- if $soft}}, true{{end}}); err != nil {
		t.Error(err)
	}

	{{else -}}
	if rowsAff, err := slice.DeleteAll({{if not .NoContext}}ctx, {{end -}} tx {{- if $soft}}, true{{end}}); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	{{end -}}

	count, err := {{$alias.UpPlural}}().Count({{if not .NoContext}}ctx, {{end -}} tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}
