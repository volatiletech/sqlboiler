{{- $alias := .Aliases.View .View.Name -}}
var (
	{{$alias.DownSingular}}AllColumns               = []string{{"{"}}{{.View.Columns | columnNames | stringMap .StringFuncs.quoteWrap | join ", "}}{{"}"}}
)

type (
	// {{$alias.UpSingular}}Slice is an alias for a slice of pointers to {{$alias.UpSingular}}.
	// This should almost always be used instead of []{{$alias.UpSingular}}.
	{{$alias.UpSingular}}Slice []*{{$alias.UpSingular}}
	{{if not .NoHooks -}}
	// {{$alias.UpSingular}}Hook is the signature for custom {{$alias.UpSingular}} hook methods
	{{$alias.UpSingular}}Hook func({{if .NoContext}}boil.Executor{{else}}context.Context, boil.ContextExecutor{{end}}, *{{$alias.UpSingular}}) error
	{{- end}}

	{{$alias.DownSingular}}Query struct {
		*queries.Query
	}
)

var (
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)
