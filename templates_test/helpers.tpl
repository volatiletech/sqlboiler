{{- $tableNameSingular := .Table.Name | singular | titleCase -}}
{{- $tableNamePlural := .Table.Name | plural | titleCase -}}
{{- $varNamePlural := .Table.Name | plural | camelCase -}}
{{- $varNameSingular := .Table.Name | singular | camelCase -}}
var {{$varNameSingular}}DBTypes = map[string]string{{"{"}}{{.Table.Columns | columnDBTypes | makeStringMap}}{{"}"}}

func {{$varNameSingular}}CompareVals(o *{{$tableNameSingular}}, j *{{$tableNameSingular}}, equal bool, blacklist ...string) error {
  {{- range $key, $value := .Table.Columns -}}
  {{if eq $value.Type "null.Time"}}
  if ((equal && o.{{titleCase $value.Name}}.Time.UTC().Format("02/01/2006") != j.{{titleCase $value.Name}}.Time.UTC().Format("02/01/2006")) ||
    (!equal && o.{{titleCase $value.Name}}.Time.UTC().Format("02/01/2006") == j.{{titleCase $value.Name}}.Time.UTC().Format("02/01/2006"))) &&
    !strmangle.SetInclude("{{$value.Name}}", blacklist) {
    return errors.New(fmt.Sprintf("NullTime {{$value.Name}} unexpected value, got:\nStruct: %#v\nResponse: %#v\n\n", o.{{titleCase $value.Name}}.Time.Format("02/01/2006"), j.{{titleCase $value.Name}}.Time.Format("02/01/2006")))
  }
  {{else if eq $value.Type "time.Time"}}
  if ((equal && o.{{titleCase $value.Name}}.UTC().Format("02/01/2006") != j.{{titleCase $value.Name}}.UTC().Format("02/01/2006")) ||
    (!equal && o.{{titleCase $value.Name}}.UTC().Format("02/01/2006") == j.{{titleCase $value.Name}}.UTC().Format("02/01/2006"))) &&
    !strmangle.SetInclude("{{$value.Name}}", blacklist) {
    return errors.New(fmt.Sprintf("Time {{$value.Name}} unexpected value, got:\nStruct: %#v\nResponse: %#v\n\n", o.{{titleCase $value.Name}}.Format("02/01/2006"), j.{{titleCase $value.Name}}.Format("02/01/2006")))
  }
  {{else if eq $value.Type "[]byte"}}
  if ((equal && !byteSliceEqual(o.{{titleCase $value.Name}}, j.{{titleCase $value.Name}})) ||
    (!equal && byteSliceEqual(o.{{titleCase $value.Name}}, j.{{titleCase $value.Name}}))) &&
    !strmangle.SetInclude("{{$value.Name}}", blacklist) {
    return errors.New(fmt.Sprintf("Expected {{$value.Name}} columns to match, got:\nStruct: %#v\nResponse: %#v\n\n", o.{{titleCase $value.Name}}, j.{{titleCase $value.Name}}))
  }
  {{else}}
  if ((equal && j.{{titleCase $value.Name}} != o.{{titleCase $value.Name}}) ||
    (!equal && j.{{titleCase $value.Name}} == o.{{titleCase $value.Name}})) &&
    !strmangle.SetInclude("{{$value.Name}}", blacklist) {
    return errors.New(fmt.Sprintf("Expected {{$value.Name}} columns to match, got:\nStruct: %#v\nResponse: %#v\n\n", o.{{titleCase $value.Name}}, j.{{titleCase $value.Name}}))
  }
  {{end}}
  {{- end -}}

  return nil
}

func Test{{$tableNamePlural}}InPrimaryKeyArgs(t *testing.T) {
  var err error
  var o {{$tableNameSingular}}
  o = {{$tableNameSingular}}{}

  if err = boil.RandomizeStruct(&o, {{$varNameSingular}}DBTypes, true); err != nil {
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

func Test{{$tableNamePlural}}SliceInPrimaryKeyArgs(t *testing.T) {
  var err error
  o := make({{$tableNameSingular}}Slice, 3)

  if err = boil.RandomizeSlice(&o, {{$varNameSingular}}DBTypes, true); err != nil {
    t.Errorf("Could not randomize slice: %s", err)
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
