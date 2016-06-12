{{- $tableNameSingular := titleCaseSingular .Table.Name -}}
{{- $dbName := singular .Table.Name -}}
{{- $tableNamePlural := titleCasePlural .Table.Name -}}
{{- $varNamePlural := camelCasePlural .Table.Name -}}
{{- $varNameSingular := camelCaseSingular .Table.Name -}}
func Test{{$tableNamePlural}}Bind(t *testing.T) {

}

func Test{{$tableNamePlural}}One(t *testing.T) {
  var err error

  {{$varNamePlural}}DeleteAllRows(t)

  o := {{$tableNameSingular}}{}
  if err = boil.RandomizeStruct(&o); err != nil {
    t.Errorf("Unable to randomize {{$tableNameSingular}} struct: %s", err)
  }

  if err = o.Insert(); err != nil {
    t.Errorf("Unable to insert {{$tableNameSingular}}:\n%#v\nErr: %s", o, err)
  }

  j, err := {{$tableNamePlural}}().One()
  if err != nil {
    t.Errorf("Unable to fetch One {{$tableNameSingular}} result:\n#%v\nErr: %s", j, err)
  }

  {{range $key, $value := .Table.Columns}}
  {{if eq $value.Type "null.Time"}}
  if o.{{titleCase $value.Name}}.Time.Format("02/01/2006") != j.{{titleCase $value.Name}}.Time.Format("02/01/2006") {
    t.Errorf("Expected NullTime {{$value.Name}} column string values to match, got:\nStruct: %#v\nResponse: %#v\n\n", o.{{titleCase $value.Name}}.Time.Format("02/01/2006"), j.{{titleCase $value.Name}}.Time.Format("02/01/2006"))
  }
  {{else if eq $value.Type "time.Time"}}
  if o.{{titleCase $value.Name}}.Format("02/01/2006") != j.{{titleCase $value.Name}}.Format("02/01/2006") {
    t.Errorf("Expected Time {{$value.Name}} column string values to match, got:\nStruct: %#v\nResponse: %#v\n\n", o.{{titleCase $value.Name}}.Format("02/01/2006"), j.{{titleCase $value.Name}}.Format("02/01/2006"))
  }
  {{else if eq $value.Type "[]byte"}}
  if !byteSliceEqual(o.{{titleCase $value.Name}}, j.{{titleCase $value.Name}}) {
    t.Errorf("%d) Expected {{$value.Name}} columns to match, got:\nStruct: %#v\nResponse: %#v\n\n", o.{{titleCase $value.Name}}, j.{{titleCase $value.Name}})

  }
  {{else}}
  if j.{{titleCase $value.Name}} != o.{{titleCase $value.Name}} {
    t.Errorf("Expected {{$value.Name}} columns to match, got:\nStruct: %#v\nResponse: %#v", o.{{titleCase $value.Name}}, j.{{titleCase $value.Name}})
  }
  {{end}}
  {{end}}
}

func Test{{$tableNamePlural}}All(t *testing.T) {
  var err error

  {{$varNamePlural}}DeleteAllRows(t)

  o := make({{$varNameSingular}}Slice, 3)
  if err = boil.RandomizeSlice(&o); err != nil {
    t.Errorf("Unable to randomize {{$tableNameSingular}} slice: %s", err)
  }

  // insert random columns to test DeleteAll
  for i := 0; i < len(o); i++ {
    err = o[i].Insert()
    if err != nil {
      t.Errorf("Unable to insert {{$tableNameSingular}}:\n%#v\nErr: %s", o[i], err)
    }
  }

  j, err := {{$tableNamePlural}}().All()
  if err != nil {
    t.Errorf("Unable to fetch All {{$tableNameSingular}} results: %s", err)
  }

  if len(j) != 3 {
    t.Errorf("Expected 3 results, got %d", len(j))
  }

  for i := 0; i < len(o); i++ {
    {{range $key, $value := .Table.Columns}}
    {{if eq $value.Type "null.Time"}}
    if o[i].{{titleCase $value.Name}}.Time.Format("02/01/2006") != j[i].{{titleCase $value.Name}}.Time.Format("02/01/2006") {
      t.Errorf("%d) Expected NullTime {{$value.Name}} column string values to match, got:\nStruct: %#v\nResponse: %#v\n\n", i, o[i].{{titleCase $value.Name}}.Time.Format("02/01/2006"), j[i].{{titleCase $value.Name}}.Time.Format("02/01/2006"))
    }
    {{else if eq $value.Type "time.Time"}}
    if o[i].{{titleCase $value.Name}}.Format("02/01/2006") != j[i].{{titleCase $value.Name}}.Format("02/01/2006") {
      t.Errorf("%d) Expected Time {{$value.Name}} column string values to match, got:\nStruct: %#v\nResponse: %#v\n\n", i, o[i].{{titleCase $value.Name}}.Format("02/01/2006"), j[i].{{titleCase $value.Name}}.Format("02/01/2006"))
    }
    {{else if eq $value.Type "[]byte"}}
    if !byteSliceEqual(o[i].{{titleCase $value.Name}}, j[i].{{titleCase $value.Name}}) {
      t.Errorf("%d) Expected {{$value.Name}} columns to match, got:\nStruct: %#v\nResponse: %#v\n\n", i, o[i].{{titleCase $value.Name}}, j[i].{{titleCase $value.Name}})

    }
    {{else}}
    if j[i].{{titleCase $value.Name}} != o[i].{{titleCase $value.Name}} {
      t.Errorf("%d) Expected {{$value.Name}} columns to match, got:\nStruct: %#v\nResponse: %#v\n\n", i, o[i].{{titleCase $value.Name}}, j[i].{{titleCase $value.Name}})
    }
    {{end}}
    {{end}}
  }
}

func Test{{$tableNamePlural}}Count(t *testing.T) {
  var err error

  {{$varNamePlural}}DeleteAllRows(t)

  o := make({{$varNameSingular}}Slice, 3)
  if err = boil.RandomizeSlice(&o); err != nil {
    t.Errorf("Unable to randomize {{$tableNameSingular}} slice: %s", err)
  }

  // insert random columns to test DeleteAll
  for i := 0; i < len(o); i++ {
    err = o[i].Insert()
    if err != nil {
      t.Errorf("Unable to insert {{$tableNameSingular}}:\n%#v\nErr: %s", o[i], err)
    }
  }
}
