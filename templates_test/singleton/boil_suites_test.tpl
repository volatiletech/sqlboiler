{{- $dot := .}}
// This test suite runs each operation test in parallel.
// Example, if your database has 3 tables, the suite will run:
// table1, table2 and table3 Delete in parallel
// table1, table2 and table3 Insert in parallel, and so forth.
// It does NOT run each operation group in parallel.
// Separating the tests thusly grants avoidance of Postgres deadlocks.
func TestParent(t *testing.T) {
  {{- range $index, $table := .Tables}}
  {{- if $table.IsJoinTable -}}
  {{- else -}}
  {{- $tableName := $table.Name | plural | titleCase -}}
  t.Run("{{$tableName}}", test{{$tableName}})
  {{end -}}
  {{- end -}}
}

func TestDelete(t *testing.T) {
  {{- range $index, $table := .Tables}}
  {{- if $table.IsJoinTable -}}
  {{- else -}}
  {{- $tableName := $table.Name | plural | titleCase -}}
  t.Run("{{$tableName}}", test{{$tableName}}Delete)
  {{end -}}
  {{- end -}}
}

func TestQueryDeleteAll(t *testing.T) {
  {{- range $index, $table := .Tables}}
  {{- if $table.IsJoinTable -}}
  {{- else -}}
  {{- $tableName := $table.Name | plural | titleCase -}}
  t.Run("{{$tableName}}", test{{$tableName}}QueryDeleteAll)
  {{end -}}
  {{- end -}}
}

func TestSliceDeleteAll(t *testing.T) {
  {{- range $index, $table := .Tables}}
  {{- if $table.IsJoinTable -}}
  {{- else -}}
  {{- $tableName := $table.Name | plural | titleCase -}}
  t.Run("{{$tableName}}", test{{$tableName}}SliceDeleteAll)
  {{end -}}
  {{- end -}}
}

func TestExists(t *testing.T) {
  {{- range $index, $table := .Tables}}
  {{- if $table.IsJoinTable -}}
  {{- else -}}
  {{- $tableName := $table.Name | plural | titleCase -}}
  t.Run("{{$tableName}}", test{{$tableName}}Exists)
  {{end -}}
  {{- end -}}
}

func TestFind(t *testing.T) {
  {{- range $index, $table := .Tables}}
  {{- if $table.IsJoinTable -}}
  {{- else -}}
  {{- $tableName := $table.Name | plural | titleCase -}}
  t.Run("{{$tableName}}", test{{$tableName}}Find)
  {{end -}}
  {{- end -}}
}

func TestBind(t *testing.T) {
  {{- range $index, $table := .Tables}}
  {{- if $table.IsJoinTable -}}
  {{- else -}}
  {{- $tableName := $table.Name | plural | titleCase -}}
  t.Run("{{$tableName}}", test{{$tableName}}Bind)
  {{end -}}
  {{- end -}}
}

func TestOne(t *testing.T) {
  {{- range $index, $table := .Tables}}
  {{- if $table.IsJoinTable -}}
  {{- else -}}
  {{- $tableName := $table.Name | plural | titleCase -}}
  t.Run("{{$tableName}}", test{{$tableName}}One)
  {{end -}}
  {{- end -}}
}

func TestAll(t *testing.T) {
  {{- range $index, $table := .Tables}}
  {{- if $table.IsJoinTable -}}
  {{- else -}}
  {{- $tableName := $table.Name | plural | titleCase -}}
  t.Run("{{$tableName}}", test{{$tableName}}All)
  {{end -}}
  {{- end -}}
}

func TestCount(t *testing.T) {
  {{- range $index, $table := .Tables}}
  {{- if $table.IsJoinTable -}}
  {{- else -}}
  {{- $tableName := $table.Name | plural | titleCase -}}
  t.Run("{{$tableName}}", test{{$tableName}}Count)
  {{end -}}
  {{- end -}}
}

{{if not .NoHooks -}}
func TestHooks(t *testing.T) {
  {{- range $index, $table := .Tables}}
  {{- if $table.IsJoinTable -}}
  {{- else -}}
  {{- $tableName := $table.Name | plural | titleCase -}}
  t.Run("{{$tableName}}", test{{$tableName}}Hooks)
  {{end -}}
  {{- end -}}
}
{{- end}}

func TestInsert(t *testing.T) {
  {{- range $index, $table := .Tables}}
  {{- if $table.IsJoinTable -}}
  {{- else -}}
  {{- $tableName := $table.Name | plural | titleCase -}}
  t.Run("{{$tableName}}", test{{$tableName}}Insert)
  t.Run("{{$tableName}}", test{{$tableName}}InsertWhitelist)
  {{end -}}
  {{- end -}}
}

// TestToOne tests cannot be run in parallel
// or deadlocks can occur.
func TestToOne(t *testing.T) {
{{- range $index, $table := .Tables}}
  {{- if $table.IsJoinTable -}}
  {{- else -}}
    {{- range $table.FKeys -}}
      {{- $txt := txtsFromFKey $dot.Tables $table . -}}
  t.Run("{{$txt.LocalTable.NameGo}}To{{$txt.ForeignTable.NameGo}}Using{{$txt.Function.Name}}", test{{$txt.LocalTable.NameGo}}ToOne{{$txt.ForeignTable.NameGo}}Using{{$txt.Function.Name}})
    {{end -}}{{- /* fkey range */ -}}
  {{- end -}}{{- /* if join table */ -}}
{{- end -}}{{- /* tables range */ -}}
}

// TestOneToOne tests cannot be run in parallel
// or deadlocks can occur.
func TestOneToOne(t *testing.T) {
  {{- range $index, $table := .Tables}}
	{{- if $table.IsJoinTable -}}
	{{- else -}}
	  {{- range $table.ToOneRelationships -}}
		{{- $txt := txtsFromOneToOne $dot.Tables $table . -}}
  t.Run("{{$txt.LocalTable.NameGo}}To{{$txt.ForeignTable.NameGo}}Using{{$txt.Function.Name}}", test{{$txt.LocalTable.NameGo}}OneToOne{{$txt.ForeignTable.NameGo}}Using{{$txt.Function.Name}})
	  {{end -}}{{- /* range */ -}}
	{{- end -}}{{- /* outer if join table */ -}}
  {{- end -}}{{- /* outer tables range */ -}}
}

// TestToMany tests cannot be run in parallel
// or deadlocks can occur.
func TestToMany(t *testing.T) {
  {{- range $index, $table := .Tables}}
    {{- if $table.IsJoinTable -}}
    {{- else -}}
      {{- range $table.ToManyRelationships -}}
        {{- $txt := txtsFromToMany $dot.Tables $table . -}}
  t.Run("{{$txt.LocalTable.NameGo}}To{{$txt.Function.Name}}", test{{$txt.LocalTable.NameGo}}ToMany{{$txt.Function.Name}})
      {{end -}}{{- /* range */ -}}
    {{- end -}}{{- /* outer if join table */ -}}
  {{- end -}}{{- /* outer tables range */ -}}
}

// TestToOneSet tests cannot be run in parallel
// or deadlocks can occur.
func TestToOneSet(t *testing.T) {
{{- range $index, $table := .Tables}}
  {{- if $table.IsJoinTable -}}
  {{- else -}}
    {{- range $table.FKeys -}}
      {{- $txt := txtsFromFKey $dot.Tables $table . -}}
  t.Run("{{$txt.LocalTable.NameGo}}To{{$txt.ForeignTable.NameGo}}Using{{$txt.Function.Name}}", test{{$txt.LocalTable.NameGo}}ToOneSetOp{{$txt.ForeignTable.NameGo}}Using{{$txt.Function.Name}})
    {{end -}}{{- /* fkey range */ -}}
  {{- end -}}{{- /* if join table */ -}}
{{- end -}}{{- /* tables range */ -}}
}

// TestToOneRemove tests cannot be run in parallel
// or deadlocks can occur.
func TestToOneRemove(t *testing.T) {
{{- range $index, $table := .Tables}}
  {{- if $table.IsJoinTable -}}
  {{- else -}}
    {{- range $table.FKeys -}}
      {{- $txt := txtsFromFKey $dot.Tables $table . -}}
      {{- if $txt.ForeignKey.Nullable -}}
  t.Run("{{$txt.LocalTable.NameGo}}To{{$txt.ForeignTable.NameGo}}Using{{$txt.Function.Name}}", test{{$txt.LocalTable.NameGo}}ToOneRemoveOp{{$txt.ForeignTable.NameGo}}Using{{$txt.Function.Name}})
      {{end -}}{{- /* if foreign key nullable */ -}}
    {{- end -}}{{- /* fkey range */ -}}
  {{- end -}}{{- /* if join table */ -}}
{{- end -}}{{- /* tables range */ -}}
}

// TestOneToOneSet tests cannot be run in parallel
// or deadlocks can occur.
func TestOneToOneSet(t *testing.T) {
  {{- range $index, $table := .Tables}}
	{{- if $table.IsJoinTable -}}
	{{- else -}}
	  {{- range $table.ToOneRelationships -}}
		  {{- $txt := txtsFromOneToOne $dot.Tables $table . -}}
	t.Run("{{$txt.LocalTable.NameGo}}To{{$txt.ForeignTable.NameGo}}Using{{$txt.Function.Name}}", test{{$txt.LocalTable.NameGo}}OneToOneSetOp{{$txt.ForeignTable.NameGo}}Using{{$txt.Function.Name}})
	  {{end -}}{{- /* range to one relationships */ -}}
	{{- end -}}{{- /* outer if join table */ -}}
  {{- end -}}{{- /* outer tables range */ -}}
}

// TestOneToOneRemove tests cannot be run in parallel
// or deadlocks can occur.
func TestOneToOneRemove(t *testing.T) {
  {{- range $index, $table := .Tables}}
	{{- if $table.IsJoinTable -}}
	{{- else -}}
	  {{- range $table.ToOneRelationships -}}
		{{- if .ForeignColumnNullable -}}
		  {{- $txt := txtsFromOneToOne $dot.Tables $table . -}}
	t.Run("{{$txt.LocalTable.NameGo}}To{{$txt.ForeignTable.NameGo}}Using{{$txt.Function.Name}}", test{{$txt.LocalTable.NameGo}}OneToOneRemoveOp{{$txt.ForeignTable.NameGo}}Using{{$txt.Function.Name}})
		{{end -}}{{- /* if foreign column nullable */ -}}
	  {{- end -}}{{- /* range */ -}}
	{{- end -}}{{- /* outer if join table */ -}}
  {{- end -}}{{- /* outer tables range */ -}}
}

// TestToManyAdd tests cannot be run in parallel
// or deadlocks can occur.
func TestToManyAdd(t *testing.T) {
  {{- range $index, $table := .Tables}}
    {{- if $table.IsJoinTable -}}
    {{- else -}}
      {{- range $table.ToManyRelationships -}}
        {{- $txt := txtsFromToMany $dot.Tables $table . -}}
  t.Run("{{$txt.LocalTable.NameGo}}To{{$txt.Function.Name}}", test{{$txt.LocalTable.NameGo}}ToManyAddOp{{$txt.Function.Name}})
      {{end -}}{{- /* range */ -}}
    {{- end -}}{{- /* outer if join table */ -}}
  {{- end -}}{{- /* outer tables range */ -}}
}

// TestToManySet tests cannot be run in parallel
// or deadlocks can occur.
func TestToManySet(t *testing.T) {
  {{- range $index, $table := .Tables}}
    {{- if $table.IsJoinTable -}}
    {{- else -}}
      {{- range $table.ToManyRelationships -}}
        {{- if not .ForeignColumnNullable -}}
        {{- else -}}
          {{- $txt := txtsFromToMany $dot.Tables $table . -}}
    t.Run("{{$txt.LocalTable.NameGo}}To{{$txt.Function.Name}}", test{{$txt.LocalTable.NameGo}}ToManySetOp{{$txt.Function.Name}})
        {{end -}}{{- /* if foreign column nullable */ -}}
      {{- end -}}{{- /* range */ -}}
    {{- end -}}{{- /* outer if join table */ -}}
  {{- end -}}{{- /* outer tables range */ -}}
}

// TestToManyRemove tests cannot be run in parallel
// or deadlocks can occur.
func TestToManyRemove(t *testing.T) {
  {{- range $index, $table := .Tables}}
    {{- if $table.IsJoinTable -}}
    {{- else -}}
      {{- range $table.ToManyRelationships -}}
        {{- if not .ForeignColumnNullable -}}
        {{- else -}}
          {{- $txt := txtsFromToMany $dot.Tables $table . -}}
    t.Run("{{$txt.LocalTable.NameGo}}To{{$txt.Function.Name}}", test{{$txt.LocalTable.NameGo}}ToManyRemoveOp{{$txt.Function.Name}})
        {{end -}}{{- /* if foreign column nullable */ -}}
      {{- end -}}{{- /* range */ -}}
    {{- end -}}{{- /* outer if join table */ -}}
  {{- end -}}{{- /* outer tables range */ -}}
}

func TestReload(t *testing.T) {
  {{- range $index, $table := .Tables}}
  {{- if $table.IsJoinTable -}}
  {{- else -}}
  {{- $tableName := $table.Name | plural | titleCase -}}
  t.Run("{{$tableName}}", test{{$tableName}}Reload)
  {{end -}}
  {{- end -}}
}

func TestReloadAll(t *testing.T) {
  {{- range $index, $table := .Tables}}
  {{- if $table.IsJoinTable -}}
  {{- else -}}
  {{- $tableName := $table.Name | plural | titleCase -}}
  t.Run("{{$tableName}}", test{{$tableName}}ReloadAll)
  {{end -}}
  {{- end -}}
}

func TestSelect(t *testing.T) {
  {{- range $index, $table := .Tables}}
  {{- if $table.IsJoinTable -}}
  {{- else -}}
  {{- $tableName := $table.Name | plural | titleCase -}}
  t.Run("{{$tableName}}", test{{$tableName}}Select)
  {{end -}}
  {{- end -}}
}

func TestUpdate(t *testing.T) {
  {{- range $index, $table := .Tables}}
  {{- if $table.IsJoinTable -}}
  {{- else -}}
  {{- $tableName := $table.Name | plural | titleCase -}}
  t.Run("{{$tableName}}", test{{$tableName}}Update)
  {{end -}}
  {{- end -}}
}

func TestSliceUpdateAll(t *testing.T) {
  {{- range $index, $table := .Tables}}
  {{- if $table.IsJoinTable -}}
  {{- else -}}
  {{- $tableName := $table.Name | plural | titleCase -}}
  t.Run("{{$tableName}}", test{{$tableName}}SliceUpdateAll)
  {{end -}}
  {{- end -}}
}

func TestUpsert(t *testing.T) {
  {{- range $index, $table := .Tables}}
  {{- if $table.IsJoinTable -}}
  {{- else -}}
  {{- $tableName := $table.Name | plural | titleCase -}}
  t.Run("{{$tableName}}", test{{$tableName}}Upsert)
  {{end -}}
  {{- end -}}
}
