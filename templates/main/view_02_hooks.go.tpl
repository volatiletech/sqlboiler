{{- if not .NoHooks -}}
{{- $alias := .Aliases.View .View.Name}}

var {{$alias.DownSingular}}AfterSelectHooks []{{$alias.UpSingular}}Hook

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

// Add{{$alias.UpSingular}}Hook registers your hook function for all future operations.
func Add{{$alias.UpSingular}}Hook(hookPoint boil.HookPoint, {{$alias.DownSingular}}Hook {{$alias.UpSingular}}Hook) {
	switch hookPoint {
		case boil.AfterSelectHook:
			{{$alias.DownSingular}}AfterSelectHooks = append({{$alias.DownSingular}}AfterSelectHooks, {{$alias.DownSingular}}Hook)
	}
}
{{- end}}
