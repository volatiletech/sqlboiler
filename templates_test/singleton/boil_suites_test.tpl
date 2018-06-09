// This test suite runs each operation test in parallel.
// Example, if your database has 3 tables, the suite will run:
// table1, table2 and table3 Delete in parallel
// table1, table2 and table3 Insert in parallel, and so forth.
// It does NOT run each operation group in parallel.
// Separating the tests thusly grants avoidance of Postgres deadlocks.
func TestParent(t *testing.T) {
  {{- range .Tables}}
  {{- if .IsJoinTable -}}
  {{- else -}}
  {{- $alias := $.Aliases.Table .Name -}}
  t.Run("{{$alias.UpPlural}}", test{{$alias.UpPlural}})
  {{end -}}
  {{- end -}}
}

func TestDelete(t *testing.T) {
  {{- range .Tables}}
  {{- if .IsJoinTable -}}
  {{- else -}}
  {{- $alias := $.Aliases.Table .Name -}}
  t.Run("{{$alias.UpPlural}}", test{{$alias.UpPlural}}Delete)
  {{end -}}
  {{- end -}}
}

func TestQueryDeleteAll(t *testing.T) {
  {{- range .Tables}}
  {{- if .IsJoinTable -}}
  {{- else -}}
  {{- $alias := $.Aliases.Table .Name -}}
  t.Run("{{$alias.UpPlural}}", test{{$alias.UpPlural}}QueryDeleteAll)
  {{end -}}
  {{- end -}}
}

func TestSliceDeleteAll(t *testing.T) {
  {{- range .Tables}}
  {{- if .IsJoinTable -}}
  {{- else -}}
  {{- $alias := $.Aliases.Table .Name -}}
  t.Run("{{$alias.UpPlural}}", test{{$alias.UpPlural}}SliceDeleteAll)
  {{end -}}
  {{- end -}}
}

func TestExists(t *testing.T) {
  {{- range .Tables}}
  {{- if .IsJoinTable -}}
  {{- else -}}
  {{- $alias := $.Aliases.Table .Name -}}
  t.Run("{{$alias.UpPlural}}", test{{$alias.UpPlural}}Exists)
  {{end -}}
  {{- end -}}
}

func TestFind(t *testing.T) {
  {{- range .Tables}}
  {{- if .IsJoinTable -}}
  {{- else -}}
  {{- $alias := $.Aliases.Table .Name -}}
  t.Run("{{$alias.UpPlural}}", test{{$alias.UpPlural}}Find)
  {{end -}}
  {{- end -}}
}

func TestBind(t *testing.T) {
  {{- range .Tables}}
  {{- if .IsJoinTable -}}
  {{- else -}}
  {{- $alias := $.Aliases.Table .Name -}}
  t.Run("{{$alias.UpPlural}}", test{{$alias.UpPlural}}Bind)
  {{end -}}
  {{- end -}}
}

func TestOne(t *testing.T) {
  {{- range .Tables}}
  {{- if .IsJoinTable -}}
  {{- else -}}
  {{- $alias := $.Aliases.Table .Name -}}
  t.Run("{{$alias.UpPlural}}", test{{$alias.UpPlural}}One)
  {{end -}}
  {{- end -}}
}

func TestAll(t *testing.T) {
  {{- range .Tables}}
  {{- if .IsJoinTable -}}
  {{- else -}}
  {{- $alias := $.Aliases.Table .Name -}}
  t.Run("{{$alias.UpPlural}}", test{{$alias.UpPlural}}All)
  {{end -}}
  {{- end -}}
}

func TestCount(t *testing.T) {
  {{- range .Tables}}
  {{- if .IsJoinTable -}}
  {{- else -}}
  {{- $alias := $.Aliases.Table .Name -}}
  t.Run("{{$alias.UpPlural}}", test{{$alias.UpPlural}}Count)
  {{end -}}
  {{- end -}}
}

{{if not .NoHooks -}}
func TestHooks(t *testing.T) {
  {{- range .Tables}}
  {{- if .IsJoinTable -}}
  {{- else -}}
  {{- $alias := $.Aliases.Table .Name -}}
  t.Run("{{$alias.UpPlural}}", test{{$alias.UpPlural}}Hooks)
  {{end -}}
  {{- end -}}
}
{{- end}}

func TestInsert(t *testing.T) {
  {{- range .Tables}}
  {{- if .IsJoinTable -}}
  {{- else -}}
  {{- $alias := $.Aliases.Table .Name -}}
  t.Run("{{$alias.UpPlural}}", test{{$alias.UpPlural}}Insert)
  t.Run("{{$alias.UpPlural}}", test{{$alias.UpPlural}}InsertWhitelist)
  {{end -}}
  {{- end -}}
}

// TestToOne tests cannot be run in parallel
// or deadlocks can occur.
func TestToOne(t *testing.T) {
{{- range .Tables}}
  {{- if .IsJoinTable -}}
  {{- else -}}
    {{- range $i, $rel := .FKeys -}}
      {{- $txt := txtsFromFKey $.Tables $rel -}}
  t.Run("{{$txt.LocalTable.NameGo}}To{{$txt.ForeignTable.NameGo}}Using{{$txt.Function.Name}}", test{{$txt.LocalTable.NameGo}}ToOne{{$txt.ForeignTable.NameGo}}Using{{$txt.Function.Name}})
    {{end -}}{{- /* fkey range */ -}}
  {{- end -}}{{- /* if join table */ -}}
{{- end -}}{{- /* tables range */ -}}
}

// TestOneToOne tests cannot be run in parallel
// or deadlocks can occur.
func TestOneToOne(t *testing.T) {
  {{- range .Tables}}
	{{- if .IsJoinTable -}}
	{{- else -}}
	  {{- range $i, $rel := .ToOneRelationships -}}
		{{- $txt := txtsFromOneToOne $.Tables . $rel -}}
  t.Run("{{$txt.LocalTable.NameGo}}To{{$txt.ForeignTable.NameGo}}Using{{$txt.Function.Name}}", test{{$txt.LocalTable.NameGo}}OneToOne{{$txt.ForeignTable.NameGo}}Using{{$txt.Function.Name}})
	  {{end -}}{{- /* range */ -}}
	{{- end -}}{{- /* outer if join table */ -}}
  {{- end -}}{{- /* outer tables range */ -}}
}

// TestToMany tests cannot be run in parallel
// or deadlocks can occur.
func TestToMany(t *testing.T) {
  {{- range .Tables}}
    {{- if .IsJoinTable -}}
    {{- else -}}
      {{- range $i, $rel := .ToManyRelationships -}}
        {{- $txt := txtsFromToMany $.Tables . $rel -}}
  t.Run("{{$txt.LocalTable.NameGo}}To{{$txt.Function.Name}}", test{{$txt.LocalTable.NameGo}}ToMany{{$txt.Function.Name}})
      {{end -}}{{- /* range */ -}}
    {{- end -}}{{- /* outer if join table */ -}}
  {{- end -}}{{- /* outer tables range */ -}}
}

// TestToOneSet tests cannot be run in parallel
// or deadlocks can occur.
func TestToOneSet(t *testing.T) {
{{- range .Tables}}
  {{- if .IsJoinTable -}}
  {{- else -}}
    {{- range $i, $fkey := .FKeys -}}
      {{- $txt := txtsFromFKey $.Tables . $fkey -}}
  t.Run("{{$txt.LocalTable.NameGo}}To{{$txt.ForeignTable.NameGo}}Using{{$txt.Function.Name}}", test{{$txt.LocalTable.NameGo}}ToOneSetOp{{$txt.ForeignTable.NameGo}}Using{{$txt.Function.Name}})
    {{end -}}{{- /* fkey range */ -}}
  {{- end -}}{{- /* if join table */ -}}
{{- end -}}{{- /* tables range */ -}}
}

// TestToOneRemove tests cannot be run in parallel
// or deadlocks can occur.
func TestToOneRemove(t *testing.T) {
{{- range .Tables}}
  {{- if .IsJoinTable -}}
  {{- else -}}
    {{- range $i, $fkey := .FKeys -}}
      {{- $txt := txtsFromFKey $.Tables . $fkey -}}
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
  {{- range .Tables}}
	{{- if .IsJoinTable -}}
	{{- else -}}
	  {{- range $i, $rel := .ToOneRelationships -}}
		  {{- $txt := txtsFromOneToOne $.Tables . $rel -}}
	t.Run("{{$txt.LocalTable.NameGo}}To{{$txt.ForeignTable.NameGo}}Using{{$txt.Function.Name}}", test{{$txt.LocalTable.NameGo}}OneToOneSetOp{{$txt.ForeignTable.NameGo}}Using{{$txt.Function.Name}})
	  {{end -}}{{- /* range to one relationships */ -}}
	{{- end -}}{{- /* outer if join table */ -}}
  {{- end -}}{{- /* outer tables range */ -}}
}

// TestOneToOneRemove tests cannot be run in parallel
// or deadlocks can occur.
func TestOneToOneRemove(t *testing.T) {
  {{- range .Tables}}
	{{- if .IsJoinTable -}}
	{{- else -}}
	  {{- range $i, $rel := .ToOneRelationships -}}
		{{- if $rel.ForeignColumnNullable -}}
		  {{- $txt := txtsFromOneToOne $.Tables . $rel -}}
	t.Run("{{$txt.LocalTable.NameGo}}To{{$txt.ForeignTable.NameGo}}Using{{$txt.Function.Name}}", test{{$txt.LocalTable.NameGo}}OneToOneRemoveOp{{$txt.ForeignTable.NameGo}}Using{{$txt.Function.Name}})
		{{end -}}{{- /* if foreign column nullable */ -}}
	  {{- end -}}{{- /* range */ -}}
	{{- end -}}{{- /* outer if join table */ -}}
  {{- end -}}{{- /* outer tables range */ -}}
}

// TestToManyAdd tests cannot be run in parallel
// or deadlocks can occur.
func TestToManyAdd(t *testing.T) {
  {{- range .Tables}}
    {{- if .IsJoinTable -}}
    {{- else -}}
      {{- range $i, $rel := .ToManyRelationships -}}
        {{- $txt := txtsFromToMany $.Tables . $rel -}}
  t.Run("{{$txt.LocalTable.NameGo}}To{{$txt.Function.Name}}", test{{$txt.LocalTable.NameGo}}ToManyAddOp{{$txt.Function.Name}})
      {{end -}}{{- /* range */ -}}
    {{- end -}}{{- /* outer if join table */ -}}
  {{- end -}}{{- /* outer tables range */ -}}
}

// TestToManySet tests cannot be run in parallel
// or deadlocks can occur.
func TestToManySet(t *testing.T) {
  {{- range .Tables}}
    {{- if .IsJoinTable -}}
    {{- else -}}
      {{- range $i, $rel := .ToManyRelationships -}}
        {{- if not (or $rel.ForeignColumnNullable $rel.ToJoinTable)}}
        {{- else -}}
          {{- $txt := txtsFromToMany $.Tables . $rel -}}
    t.Run("{{$txt.LocalTable.NameGo}}To{{$txt.Function.Name}}", test{{$txt.LocalTable.NameGo}}ToManySetOp{{$txt.Function.Name}})
        {{end -}}{{- /* if foreign column nullable */ -}}
      {{- end -}}{{- /* range */ -}}
    {{- end -}}{{- /* outer if join table */ -}}
  {{- end -}}{{- /* outer tables range */ -}}
}

// TestToManyRemove tests cannot be run in parallel
// or deadlocks can occur.
func TestToManyRemove(t *testing.T) {
  {{- range .Tables}}
    {{- if .IsJoinTable -}}
    {{- else -}}
      {{- range $i, $rel := .ToManyRelationships -}}
        {{- if not (or $rel.ForeignColumnNullable $rel.ToJoinTable)}}
        {{- else -}}
          {{- $txt := txtsFromToMany $.Tables . $rel -}}
    t.Run("{{$txt.LocalTable.NameGo}}To{{$txt.Function.Name}}", test{{$txt.LocalTable.NameGo}}ToManyRemoveOp{{$txt.Function.Name}})
        {{end -}}{{- /* if foreign column nullable */ -}}
      {{- end -}}{{- /* range */ -}}
    {{- end -}}{{- /* outer if join table */ -}}
  {{- end -}}{{- /* outer tables range */ -}}
}

func TestReload(t *testing.T) {
  {{- range .Tables}}
  {{- if .IsJoinTable -}}
  {{- else -}}
  {{- $alias := $.Aliases.Table .Name -}}
  t.Run("{{$alias.UpPlural}}", test{{$alias.UpPlural}}Reload)
  {{end -}}
  {{- end -}}
}

func TestReloadAll(t *testing.T) {
  {{- range .Tables}}
  {{- if .IsJoinTable -}}
  {{- else -}}
  {{- $alias := $.Aliases.Table .Name -}}
  t.Run("{{$alias.UpPlural}}", test{{$alias.UpPlural}}ReloadAll)
  {{end -}}
  {{- end -}}
}

func TestSelect(t *testing.T) {
  {{- range .Tables}}
  {{- if .IsJoinTable -}}
  {{- else -}}
  {{- $alias := $.Aliases.Table .Name -}}
  t.Run("{{$alias.UpPlural}}", test{{$alias.UpPlural}}Select)
  {{end -}}
  {{- end -}}
}

func TestUpdate(t *testing.T) {
  {{- range .Tables}}
  {{- if .IsJoinTable -}}
  {{- else -}}
  {{- $alias := $.Aliases.Table .Name -}}
  t.Run("{{$alias.UpPlural}}", test{{$alias.UpPlural}}Update)
  {{end -}}
  {{- end -}}
}

func TestSliceUpdateAll(t *testing.T) {
  {{- range .Tables}}
  {{- if .IsJoinTable -}}
  {{- else -}}
  {{- $alias := $.Aliases.Table .Name -}}
  t.Run("{{$alias.UpPlural}}", test{{$alias.UpPlural}}SliceUpdateAll)
  {{end -}}
  {{- end -}}
}
