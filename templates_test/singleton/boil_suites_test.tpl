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

func TestHelpers(t *testing.T) {
  {{- range $index, $table := .Tables}}
  {{- if $table.IsJoinTable -}}
  {{- else -}}
  {{- $tableName := $table.Name | plural | titleCase -}}
  t.Run("{{$tableName}}", test{{$tableName}}InPrimaryKeyArgs)
  t.Run("{{$tableName}}", test{{$tableName}}SliceInPrimaryKeyArgs)
  {{end -}}
  {{- end -}}
}

func TestHooks(t *testing.T) {
  {{- range $index, $table := .Tables}}
  {{- if $table.IsJoinTable -}}
  {{- else -}}
  {{- $tableName := $table.Name | plural | titleCase -}}
  t.Run("{{$tableName}}", test{{$tableName}}Hooks)
  {{end -}}
  {{- end -}}
}

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

// TestToMany tests cannot be run in parallel
// or postgres deadlocks will occur.
func TestToMany(t *testing.T) {
  {{- $dot := .}}
  {{- range $index, $table := .Tables}}
    {{- $tableName := $table.Name | plural | titleCase -}}
    {{- if $table.IsJoinTable -}}
    {{- else -}}
      {{- range $table.ToManyRelationships -}}
        {{- $rel := textsFromRelationship $dot.Tables $table . -}}
        {{- if (and .ForeignColumnUnique (not .ToJoinTable)) -}}
          {{- $oneToOne := textsFromOneToOneRelationship $dot.PkgName $dot.Tables $table . -}}
  t.Run("{{$oneToOne.LocalTable.NameGo}}OneToOne{{$oneToOne.ForeignTable.NameGo}}_{{$oneToOne.Function.Name}}", test{{$oneToOne.LocalTable.NameGo}}ToOne{{$oneToOne.ForeignTable.NameGo}}_{{$oneToOne.Function.Name}})
        {{else -}}
  t.Run("{{$rel.LocalTable.NameGo}}ToMany{{$rel.Function.Name}}", test{{$rel.LocalTable.NameGo}}ToMany{{$rel.Function.Name}})
        {{end -}}{{- /* if unique */ -}}
      {{- end -}}{{- /* range */ -}}
    {{- end -}}{{- /* outer if join table */ -}}
  {{- end -}}{{- /* outer tables range */ -}}
}

// TestToOne tests cannot be run in parallel
// or postgres deadlocks will occur.
func TestToOne(t *testing.T) {
  {{- $dot := . -}}
{{- range $index, $table := .Tables}}
  {{- if $table.IsJoinTable -}}
  {{- else -}}
    {{- range $table.FKeys -}}
      {{- $rel := textsFromForeignKey $dot.PkgName $dot.Tables $table . -}}
  t.Run("{{$rel.LocalTable.NameGo}}To{{$rel.ForeignTable.NameGo}}_{{$rel.Function.Name}}", test{{$rel.LocalTable.NameGo}}ToOne{{$rel.ForeignTable.NameGo}}_{{$rel.Function.Name}})
    {{end -}}{{- /* fkey range */ -}}
  {{- end -}}{{- /* if join table */ -}}
{{- end -}}{{- /* tables range */ -}}
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
