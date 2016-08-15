// This test suite runs each operation test in parallel.
// Example, if your database has 3 tables, the suite will run:
// table1, table2 and table3 Delete in parallel
// table1, table2 and table3 Insert in parallel, and so forth.
// It does NOT run each operation group in parallel.
// Separating the tests thusly grants avoidance of Postgres deadlocks.
func TestAll(t *testing.T) {
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
  t.Run("{{$tableName}}", test{{$tableName}}QueryDeleteAll)
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

func TestFinishers(t *testing.T) {
  {{- range $index, $table := .Tables}}
  {{- if $table.IsJoinTable -}}
  {{- else -}}
  {{- $tableName := $table.Name | plural | titleCase -}}
  t.Run("{{$tableName}}", test{{$tableName}}Bind)
  t.Run("{{$tableName}}", test{{$tableName}}One)
  t.Run("{{$tableName}}", test{{$tableName}}All)
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

func TestRelationships(t *testing.T) {
  {{- $dot := .}}
  {{- range $index, $table := .Tables}}
    {{- $tableName := $table.Name | plural | titleCase -}}
    {{- if $table.IsJoinTable -}}
    {{- else -}}
      {{- range $table.ToManyRelationships -}}
        {{- $rel := textsFromRelationship $dot.Tables $table . -}}
        {{- if (and .ForeignColumnUnique (not .ToJoinTable)) -}}
          {{- $funcName := $rel.LocalTable.NameGo -}}
  t.Run("{{$rel.ForeignTable.NameGo}}ToOne", test{{$rel.ForeignTable.NameGo}}ToOne{{$rel.LocalTable.NameGo}}_{{$funcName}})
        {{else -}}
  t.Run("{{$rel.LocalTable.NameGo}}ToMany", test{{$rel.LocalTable.NameGo}}ToMany{{$rel.Function.Name}})
        {{end -}}{{- /* if unique */ -}}
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
