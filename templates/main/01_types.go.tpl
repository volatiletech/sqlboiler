{{if .Table.IsJoinTable -}}
{{else -}}
{{- $alias := .Aliases.Table .Table.Name -}}
var (
	{{$alias.DownSingular}}AllColumns               = []string{{"{"}}{{.Table.Columns | columnNames | stringMap .StringFuncs.quoteWrap | join ", "}}{{"}"}}
	{{$alias.DownSingular}}ColumnsWithoutDefault = []string{{"{"}}{{.Table.Columns | filterColumnsByDefault false | columnNames | stringMap .StringFuncs.quoteWrap | join ","}}{{"}"}}
	{{$alias.DownSingular}}ColumnsWithDefault    = []string{{"{"}}{{.Table.Columns | filterColumnsByDefault true | columnNames | stringMap .StringFuncs.quoteWrap | join ","}}{{"}"}}
	{{if .Table.IsView -}}
	{{$alias.DownSingular}}PrimaryKeyColumns     = []string{}
	{{else -}}
	{{$alias.DownSingular}}PrimaryKeyColumns     = []string{{"{"}}{{.Table.PKey.Columns | stringMap .StringFuncs.quoteWrap | join ", "}}{{"}"}}
	{{end -}}
	{{$alias.DownSingular}}GeneratedColumns = []string{{"{"}}{{.Table.Columns | filterColumnsByAuto true | columnNames | stringMap .StringFuncs.quoteWrap | join ","}}{{"}"}}
)

// {{$alias.DownSingular}}Table satisfies a constraint for helpers.Table[*{{$alias.UpSingular}}]
type {{$alias.DownSingular}}Table struct{}

// New returns a new model
func ({{$alias.DownSingular}}Table) New() *{{$alias.UpSingular}} {
	return &{{$alias.UpSingular}}{}
}

{{- $canSoftDelete := .Table.CanSoftDelete $.AutoColumns.Deleted -}}
{{- $soft := and .AddSoftDeletes $canSoftDelete }}
// TableName returns the name of the table
func ({{$alias.DownSingular}}Table) TableInfo() helpers.TableInfo {
  return helpers.TableInfo{
	  Name: "{{.Table.Name | .SchemaTable}}",
		Dialect: dialect,
	  Type: {{$alias.DownSingular}}Type,
	  Mapping: {{$alias.DownSingular}}Mapping,
		{{if not .Table.IsView -}}
	  PrimaryKeyMapping: {{$alias.DownSingular}}PrimaryKeyMapping,
		{{- end}}

    AllColumns: {{$alias.DownSingular}}AllColumns,
    ColumnsWithDefault: {{$alias.DownSingular}}ColumnsWithDefault,
    ColumnsWithoutDefault: {{$alias.DownSingular}}ColumnsWithoutDefault,
    PrimaryKeyColumns: {{$alias.DownSingular}}PrimaryKeyColumns,
    GeneratedColumns: {{$alias.DownSingular}}GeneratedColumns,

		{{- if $soft}}
			DeletionColumnName: "{{or $.AutoColumns.Deleted "deleted_at" | $.Quotes}}",
		{{- end}}
	}
}

// SetAsSoftDeleted set the deleted column when soft deleting
// does nothing when not supported
func ({{$alias.DownSingular}}Table) SetAsSoftDeleted(o *{{$alias.UpSingular}}, t time.Time) {
		{{- if $soft}}
			o.{{$alias.Column (or $.AutoColumns.Deleted "deleted_at")}} = null.TimeFrom(t)
		{{- end}}
}

{{if not .Table.IsView -}}
{{if eq (len .Table.PKey.Columns) 1 -}}
	{{$pk := .Table.GetColumn (index .Table.PKey.Columns 0)}}
	type {{$alias.UpSingular}}PK {{$pk.Type}}

	func (pk {{$alias.UpSingular}}PK) Values() []any {
		return []any{pk}
	}
{{- else -}}
	{{$pk := .Table.GetColumn (index .Table.PKey.Columns 0)}}
	type {{$alias.UpSingular}}PK struct{
		{{range .Table.PKey.Columns -}}
		{{- $column := $.Table.GetColumn . -}}
		{{$alias.Column $column.Name}} {{$column.Type}}
		{{end -}}
	}

	func (pk {{$alias.UpSingular}}PK) Values() []any {
		return []any{
			{{- range .Table.PKey.Columns -}}
			{{- $column := $.Table.GetColumn . -}}
			pk.{{$alias.Column $column.Name}},
			{{- end -}}
		}
	}
{{- end}}
{{- end}}

type (
	// {{$alias.UpSingular}}Slice is an alias for a slice of pointers to {{$alias.UpSingular}}.
	// This should almost always be used instead of []{{$alias.UpSingular}}.
	{{if .Table.IsView -}}
	{{$alias.UpSingular}}Slice = helpers.ViewSlice[*{{$alias.UpSingular}}]
	{{- else -}}
	{{$alias.UpSingular}}Slice = helpers.TableSlice[*{{$alias.UpSingular}}, {{$alias.DownSingular}}Table, {{$alias.DownSingular}}Hooks]
	{{end -}}

	{{if not .NoHooks -}}
	// {{$alias.UpSingular}}Hook is the signature for custom {{$alias.UpSingular}} hook methods
	{{$alias.UpSingular}}Hook = helpers.TableHook[*{{$alias.UpSingular}}]
	// {{$alias.DownSingular}}Hooks is contains methods to retrieve registered hooks
	{{$alias.DownSingular}}Hooks struct{}
	{{- else -}}
	// {{$alias.DownSingular}}Hooks is a no-op hook. Returns nothing
	{{$alias.DownSingular}}Hooks = helpers.NoOpHooks[*{{$alias.UpSingular}}]
	{{end -}}

	{{- $canSoftDelete := .Table.CanSoftDelete $.AutoColumns.Deleted -}}
	{{- $soft := and .AddSoftDeletes $canSoftDelete }}
	// {{$alias.DownSingular}}Query to query this model
	{{$alias.DownSingular}}Query struct{
	  *queries.Query // so it can be used as a query itself
		helpers.BaseQuery[*{{$alias.UpSingular}}, {{$alias.DownSingular}}Table]
		helpers.SelectQuery[*{{$alias.UpSingular}}, {{$alias.DownSingular}}Table, {{$alias.DownSingular}}Hooks, {{$alias.UpSingular}}Slice]
		helpers.DeleteQuery[*{{$alias.UpSingular}}, {{$alias.DownSingular}}Table]
	}
)

// Cache for insert, update and upsert
var (
	{{$alias.DownSingular}}Type = reflect.TypeOf(&{{$alias.UpSingular}}{})
	{{$alias.DownSingular}}Mapping = queries.MakeStructMapping({{$alias.DownSingular}}Type)
	{{if not .Table.IsView -}}
	{{$alias.DownSingular}}PrimaryKeyMapping, _ = queries.BindMapping({{$alias.DownSingular}}Type, {{$alias.DownSingular}}Mapping, {{$alias.DownSingular}}PrimaryKeyColumns)
	{{end -}}
	{{$alias.DownSingular}}InsertCacheMut sync.RWMutex
	{{$alias.DownSingular}}InsertCache = make(map[string]insertCache)
	{{$alias.DownSingular}}UpdateCacheMut sync.RWMutex
	{{$alias.DownSingular}}UpdateCache = make(map[string]updateCache)
	{{$alias.DownSingular}}UpsertCacheMut sync.RWMutex
	{{$alias.DownSingular}}UpsertCache = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
	{{if .Table.IsView -}}
	// These are used in some views
	_ = fmt.Sprintln("")
	_ = reflect.Int
	_ = strings.Builder{}
	_ = sync.Mutex{}
	_ = strmangle.Plural("")
	_ = strconv.IntSize
	{{if not .Table.ViewCapabilities.CanUpsert -}}
	_ = sql.ErrNoRows
	_ = errors.New("")
	{{- end}}
	{{- end}}
)
{{end -}}
