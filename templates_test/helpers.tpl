{{- $tableNameSingular := .Table.Name | singular | titleCase -}}
{{- $tableNamePlural := .Table.Name | plural | titleCase -}}
{{- $varNamePlural := .Table.Name | plural | camelCase -}}
{{- $varNameSingular := .Table.Name | singular | camelCase -}}
var {{$varNameSingular}}DBTypes = map[string]string{{"{"}}{{.Table.Columns | columnDBTypes | makeStringMap}}{{"}"}}

func test{{$tableNamePlural}}InPrimaryKeyArgs(t *testing.T) {
  t.Parallel()

  var err error
  var o {{$tableNameSingular}}
  o = {{$tableNameSingular}}{}

  seed := boil.NewSeed()
  if err = seed.RandomizeStruct(&o, {{$varNameSingular}}DBTypes, true); err != nil {
    t.Errorf("Could not randomize struct: %s", err)
  }

  args := o.inPrimaryKeyArgs()

  if len(args) != len({{$varNameSingular}}PrimaryKeyColumns) {
    t.Errorf("Expected args to be len %d, but got %d", len({{$varNameSingular}}PrimaryKeyColumns), len(args))
  }

  {{range $key, $value := .Table.PKey.Columns}}
  if o.{{titleCase $value}} != args[{{$key}}] {
    t.Errorf("Expected args[{{$key}}] to be value of o.{{titleCase $value}}, but got %#v", args[{{$key}}])
  }
  {{- end}}
}

func test{{$tableNamePlural}}SliceInPrimaryKeyArgs(t *testing.T) {
  t.Parallel()

  var err error
  o := make({{$tableNameSingular}}Slice, 3)

  seed := boil.NewSeed()
  for i := range o {
    o[i] = &{{$tableNameSingular}}{}
    if err = seed.RandomizeStruct(o[i], {{$varNameSingular}}DBTypes, true); err != nil {
      t.Errorf("Could not randomize struct: %s", err)
    }
  }

  args := o.inPrimaryKeyArgs()

  if len(args) != len({{$varNameSingular}}PrimaryKeyColumns) * 3 {
    t.Errorf("Expected args to be len %d, but got %d", len({{$varNameSingular}}PrimaryKeyColumns) * 3, len(args))
  }

  for i := 0; i < len({{$varNameSingular}}PrimaryKeyColumns) * 3; i++ {
    {{range $key, $value := .Table.PKey.Columns}}
    if o[i].{{titleCase $value}} != args[i] {
      t.Errorf("Expected args[%d] to be value of o.{{titleCase $value}}, but got %#v", i, args[i])
    }
    {{- end}}
  }
}
