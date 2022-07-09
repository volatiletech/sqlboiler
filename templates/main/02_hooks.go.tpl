{{- if not .NoHooks -}}
{{- $alias := .Aliases.Table .Table.Name}}

var {{$alias.DownSingular}}AfterSelectHooks []{{$alias.UpSingular}}Hook

{{if or (not .Table.IsView) (.Table.ViewCapabilities.CanInsert) -}}
var {{$alias.DownSingular}}BeforeInsertHooks []{{$alias.UpSingular}}Hook
var {{$alias.DownSingular}}AfterInsertHooks []{{$alias.UpSingular}}Hook
{{- end}}

{{if not .Table.IsView -}}
var {{$alias.DownSingular}}BeforeUpdateHooks []{{$alias.UpSingular}}Hook
var {{$alias.DownSingular}}AfterUpdateHooks []{{$alias.UpSingular}}Hook

var {{$alias.DownSingular}}BeforeDeleteHooks []{{$alias.UpSingular}}Hook
var {{$alias.DownSingular}}AfterDeleteHooks []{{$alias.UpSingular}}Hook
{{- end}}

{{if or (not .Table.IsView) (.Table.ViewCapabilities.CanUpsert) -}}
var {{$alias.DownSingular}}BeforeUpsertHooks []{{$alias.UpSingular}}Hook
var {{$alias.DownSingular}}AfterUpsertHooks []{{$alias.UpSingular}}Hook
{{- end}}

// doAfterSelectHooks executes all "after Select" hooks.
func (o *{{$alias.UpSingular}}) doAfterSelectHooks({{if .NoContext}}exec boil.Executor{{else}}ctx context.Context, exec boil.ContextExecutor{{end}}) (err error) {
	{{if not .NoContext -}}
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	{{end -}}
	for _, hook := range {{$alias.DownSingular}}AfterSelectHooks {
		if err := hook({{if not .NoContext}}ctx, {{end -}} exec, o); err != nil {
			return err
		}
	}

	return nil
}

{{if or (not .Table.IsView) (.Table.ViewCapabilities.CanInsert) -}}
// doBeforeInsertHooks executes all "before insert" hooks.
func (o *{{$alias.UpSingular}}) doBeforeInsertHooks({{if .NoContext}}exec boil.Executor{{else}}ctx context.Context, exec boil.ContextExecutor{{end}}) (err error) {
	{{if not .NoContext -}}
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	{{end -}}
	for _, hook := range {{$alias.DownSingular}}BeforeInsertHooks {
		if err := hook({{if not .NoContext}}ctx, {{end -}} exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *{{$alias.UpSingular}}) doAfterInsertHooks({{if .NoContext}}exec boil.Executor{{else}}ctx context.Context, exec boil.ContextExecutor{{end}}) (err error) {
	{{if not .NoContext -}}
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	{{end -}}
	for _, hook := range {{$alias.DownSingular}}AfterInsertHooks {
		if err := hook({{if not .NoContext}}ctx, {{end -}} exec, o); err != nil {
			return err
		}
	}

	return nil
}
{{- end}}

{{if not .Table.IsView -}}
// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *{{$alias.UpSingular}}) doBeforeUpdateHooks({{if .NoContext}}exec boil.Executor{{else}}ctx context.Context, exec boil.ContextExecutor{{end}}) (err error) {
	{{if not .NoContext -}}
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	{{end -}}
	for _, hook := range {{$alias.DownSingular}}BeforeUpdateHooks {
		if err := hook({{if not .NoContext}}ctx, {{end -}} exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *{{$alias.UpSingular}}) doAfterUpdateHooks({{if .NoContext}}exec boil.Executor{{else}}ctx context.Context, exec boil.ContextExecutor{{end}}) (err error) {
	{{if not .NoContext -}}
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	{{end -}}
	for _, hook := range {{$alias.DownSingular}}AfterUpdateHooks {
		if err := hook({{if not .NoContext}}ctx, {{end -}} exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *{{$alias.UpSingular}}) doBeforeDeleteHooks({{if .NoContext}}exec boil.Executor{{else}}ctx context.Context, exec boil.ContextExecutor{{end}}) (err error) {
	{{if not .NoContext -}}
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	{{end -}}
	for _, hook := range {{$alias.DownSingular}}BeforeDeleteHooks {
		if err := hook({{if not .NoContext}}ctx, {{end -}} exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *{{$alias.UpSingular}}) doAfterDeleteHooks({{if .NoContext}}exec boil.Executor{{else}}ctx context.Context, exec boil.ContextExecutor{{end}}) (err error) {
	{{if not .NoContext -}}
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	{{end -}}
	for _, hook := range {{$alias.DownSingular}}AfterDeleteHooks {
		if err := hook({{if not .NoContext}}ctx, {{end -}} exec, o); err != nil {
			return err
		}
	}

	return nil
}
{{- end}}

{{if or (not .Table.IsView) (.Table.ViewCapabilities.CanUpsert) -}}
// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *{{$alias.UpSingular}}) doBeforeUpsertHooks({{if .NoContext}}exec boil.Executor{{else}}ctx context.Context, exec boil.ContextExecutor{{end}}) (err error) {
	{{if not .NoContext -}}
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	{{end -}}
	for _, hook := range {{$alias.DownSingular}}BeforeUpsertHooks {
		if err := hook({{if not .NoContext}}ctx, {{end -}} exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *{{$alias.UpSingular}}) doAfterUpsertHooks({{if .NoContext}}exec boil.Executor{{else}}ctx context.Context, exec boil.ContextExecutor{{end}}) (err error) {
	{{if not .NoContext -}}
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	{{end -}}
	for _, hook := range {{$alias.DownSingular}}AfterUpsertHooks {
		if err := hook({{if not .NoContext}}ctx, {{end -}} exec, o); err != nil {
			return err
		}
	}

	return nil
}
{{- end}}

// Add{{$alias.UpSingular}}Hook registers your hook function for all future operations.
func Add{{$alias.UpSingular}}Hook(hookPoint boil.HookPoint, {{$alias.DownSingular}}Hook {{$alias.UpSingular}}Hook) {
	switch hookPoint {
		case boil.AfterSelectHook:
			{{$alias.DownSingular}}AfterSelectHooks = append({{$alias.DownSingular}}AfterSelectHooks, {{$alias.DownSingular}}Hook)
		{{- if or (not .Table.IsView) (.Table.ViewCapabilities.CanInsert)}}
		case boil.BeforeInsertHook:
			{{$alias.DownSingular}}BeforeInsertHooks = append({{$alias.DownSingular}}BeforeInsertHooks, {{$alias.DownSingular}}Hook)
		case boil.AfterInsertHook:
			{{$alias.DownSingular}}AfterInsertHooks = append({{$alias.DownSingular}}AfterInsertHooks, {{$alias.DownSingular}}Hook)
		{{- end}}
		{{- if not .Table.IsView}}
		case boil.BeforeUpdateHook:
			{{$alias.DownSingular}}BeforeUpdateHooks = append({{$alias.DownSingular}}BeforeUpdateHooks, {{$alias.DownSingular}}Hook)
		case boil.AfterUpdateHook:
			{{$alias.DownSingular}}AfterUpdateHooks = append({{$alias.DownSingular}}AfterUpdateHooks, {{$alias.DownSingular}}Hook)
		case boil.BeforeDeleteHook:
			{{$alias.DownSingular}}BeforeDeleteHooks = append({{$alias.DownSingular}}BeforeDeleteHooks, {{$alias.DownSingular}}Hook)
		case boil.AfterDeleteHook:
			{{$alias.DownSingular}}AfterDeleteHooks = append({{$alias.DownSingular}}AfterDeleteHooks, {{$alias.DownSingular}}Hook)
		{{- end}}
		{{- if or (not .Table.IsView) (.Table.ViewCapabilities.CanInsert)}}
		case boil.BeforeUpsertHook:
			{{$alias.DownSingular}}BeforeUpsertHooks = append({{$alias.DownSingular}}BeforeUpsertHooks, {{$alias.DownSingular}}Hook)
		case boil.AfterUpsertHook:
			{{$alias.DownSingular}}AfterUpsertHooks = append({{$alias.DownSingular}}AfterUpsertHooks, {{$alias.DownSingular}}Hook)
		{{- end}}
	}
}
{{- end}}
