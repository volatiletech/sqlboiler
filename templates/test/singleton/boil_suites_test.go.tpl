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

{{if .AddSoftDeletes -}}
func TestSoftDelete(t *testing.T) {
  {{- range .Tables}}
  {{- if .IsJoinTable -}}
  {{- else -}}
  	{{- if .CanSoftDelete $.AutoColumns.Deleted -}}
      {{- $alias := $.Aliases.Table .Name -}}
      t.Run("{{$alias.UpPlural}}", test{{$alias.UpPlural}}SoftDelete)
  	{{end -}}
  {{end -}}
  {{- end -}}
}

func TestQuerySoftDeleteAll(t *testing.T) {
  {{- range .Tables}}
  {{- if .IsJoinTable -}}
  {{- else -}}
  	{{- if .CanSoftDelete $.AutoColumns.Deleted -}}
      {{- $alias := $.Aliases.Table .Name -}}
      t.Run("{{$alias.UpPlural}}", test{{$alias.UpPlural}}QuerySoftDeleteAll)
  	{{end -}}
  {{end -}}
  {{- end -}}
}

func TestSliceSoftDeleteAll(t *testing.T) {
  {{- range .Tables}}
  {{- if .IsJoinTable -}}
  {{- else -}}
  	{{- if .CanSoftDelete $.AutoColumns.Deleted -}}
      {{- $alias := $.Aliases.Table .Name -}}
      t.Run("{{$alias.UpPlural}}", test{{$alias.UpPlural}}SliceSoftDeleteAll)
  	{{end -}}
  {{end -}}
  {{- end -}}
}
{{- end}}

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
    {{- range $fkey := .FKeys -}}
      {{- $ltable := $.Aliases.Table $fkey.Table -}}
      {{- $ftable := $.Aliases.Table $fkey.ForeignTable -}}
      {{- $relAlias := $ltable.Relationship $fkey.Name -}}
  t.Run("{{$ltable.UpSingular}}To{{$ftable.UpSingular}}Using{{$relAlias.Foreign}}", test{{$ltable.UpSingular}}ToOne{{$ftable.UpSingular}}Using{{$relAlias.Foreign}})
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
	  {{- range $rel := .ToOneRelationships -}}
      {{- $ltable := $.Aliases.Table $rel.Table -}}
      {{- $ftable := $.Aliases.Table $rel.ForeignTable -}}
      {{- $relAlias := $ftable.Relationship $rel.Name -}}
	t.Run("{{$ltable.UpSingular}}To{{$ftable.UpSingular}}Using{{$relAlias.Local}}", test{{$ltable.UpSingular}}OneToOne{{$ftable.UpSingular}}Using{{$relAlias.Local}})
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
      {{- range $rel := .ToManyRelationships -}}
        {{- $ltable := $.Aliases.Table $rel.Table -}}
        {{- $relAlias := $.Aliases.ManyRelationship $rel.ForeignTable $rel.Name $rel.JoinTable $rel.JoinLocalFKeyName -}}
  t.Run("{{$ltable.UpSingular}}To{{$relAlias.Local}}", test{{$ltable.UpSingular}}ToMany{{$relAlias.Local}})
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
    {{- range $fkey := .FKeys -}}
      {{- $ltable := $.Aliases.Table $fkey.Table -}}
      {{- $ftable := $.Aliases.Table $fkey.ForeignTable -}}
      {{- $relAlias := $ltable.Relationship $fkey.Name -}}
  t.Run("{{$ltable.UpSingular}}To{{$ftable.UpSingular}}Using{{$relAlias.Local}}", test{{$ltable.UpSingular}}ToOneSetOp{{$ftable.UpSingular}}Using{{$relAlias.Foreign}})
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
    {{- range $fkey := .FKeys -}}
      {{- if $fkey.Nullable -}}
        {{- $ltable := $.Aliases.Table $fkey.Table -}}
        {{- $ftable := $.Aliases.Table $fkey.ForeignTable -}}
        {{- $relAlias := $ltable.Relationship $fkey.Name -}}
  t.Run("{{$ltable.UpSingular}}To{{$ftable.UpSingular}}Using{{$relAlias.Local}}", test{{$ltable.UpSingular}}ToOneRemoveOp{{$ftable.UpSingular}}Using{{$relAlias.Foreign}})
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
	  {{- range $rel := .ToOneRelationships -}}
      {{- $ltable := $.Aliases.Table $rel.Table -}}
      {{- $ftable := $.Aliases.Table $rel.ForeignTable -}}
      {{- $relAlias := $ftable.Relationship $rel.Name -}}
	t.Run("{{$ltable.UpSingular}}To{{$ftable.UpSingular}}Using{{$relAlias.Local}}", test{{$ltable.UpSingular}}OneToOneSetOp{{$ftable.UpSingular}}Using{{$relAlias.Local}})
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
	  {{- range $rel := .ToOneRelationships -}}
		{{- if $rel.ForeignColumnNullable -}}
      {{- $ltable := $.Aliases.Table $rel.Table -}}
      {{- $ftable := $.Aliases.Table $rel.ForeignTable -}}
      {{- $relAlias := $ftable.Relationship $rel.Name -}}
	t.Run("{{$ltable.UpSingular}}To{{$ftable.UpSingular}}Using{{$relAlias.Local}}", test{{$ltable.UpSingular}}OneToOneRemoveOp{{$ftable.UpSingular}}Using{{$relAlias.Local}})
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
      {{- range $rel := .ToManyRelationships -}}
        {{- $ltable := $.Aliases.Table $rel.Table -}}
        {{- $relAlias := $.Aliases.ManyRelationship $rel.ForeignTable $rel.Name $rel.JoinTable $rel.JoinLocalFKeyName -}}
  t.Run("{{$ltable.UpSingular}}To{{$relAlias.Local}}", test{{$ltable.UpSingular}}ToManyAddOp{{$relAlias.Local}})
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
      {{- range $rel := .ToManyRelationships -}}
        {{- if not (or $rel.ForeignColumnNullable $rel.ToJoinTable)}}
        {{- else -}}
          {{- $ltable := $.Aliases.Table $rel.Table -}}
          {{- $relAlias := $.Aliases.ManyRelationship $rel.ForeignTable $rel.Name $rel.JoinTable $rel.JoinLocalFKeyName -}}
  t.Run("{{$ltable.UpSingular}}To{{$relAlias.Local}}", test{{$ltable.UpSingular}}ToManySetOp{{$relAlias.Local}})
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
      {{- range $rel := .ToManyRelationships -}}
        {{- if not (or $rel.ForeignColumnNullable $rel.ToJoinTable)}}
        {{- else -}}
          {{- $ltable := $.Aliases.Table $rel.Table -}}
          {{- $relAlias := $.Aliases.ManyRelationship $rel.ForeignTable $rel.Name $rel.JoinTable $rel.JoinLocalFKeyName -}}
  t.Run("{{$ltable.UpSingular}}To{{$relAlias.Local}}", test{{$ltable.UpSingular}}ToManyRemoveOp{{$relAlias.Local}})
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
