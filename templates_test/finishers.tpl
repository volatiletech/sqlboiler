{{- $tableNameSingular := .Table.Name | singular | titleCase -}}
{{- $tableNamePlural := .Table.Name | plural | titleCase -}}
{{- $varNamePlural := .Table.Name | plural | camelCase -}}
{{- $varNameSingular := .Table.Name | singular | camelCase -}}
func Test{{$tableNamePlural}}Bind(t *testing.T) {
  var err error

  o := {{$tableNameSingular}}{}
  if err = boil.RandomizeStruct(&o, {{$varNameSingular}}DBTypes, true); err != nil {
    t.Errorf("Unable to randomize {{$tableNameSingular}} struct: %s", err)
  }

  if err = o.InsertG(); err != nil {
    t.Errorf("Unable to insert {{$tableNameSingular}}:\n%#v\nErr: %s", o, err)
  }

  j := {{$tableNameSingular}}{}

  err = {{$tableNamePlural}}G(qm.Where(`{{whereClause 1 .Table.PKey.Columns}}`, {{.Table.PKey.Columns | stringMap .StringFuncs.titleCase | prefixStringSlice "o." | join ", "}})).Bind(&j)
  if err != nil {
    t.Errorf("Unable to call Bind on {{$tableNameSingular}} single object: %s", err)
  }

  {{$varNameSingular}}CompareVals(&o, &j, t)

  // insert 3 rows, attempt to bind into slice
  {{$varNamePlural}}DeleteAllRows(t)

  y := make({{$tableNameSingular}}Slice, 3)
  if err = boil.RandomizeSlice(&y, {{$varNameSingular}}DBTypes, true); err != nil {
    t.Errorf("Unable to randomize {{$tableNameSingular}} slice: %s", err)
  }

  // insert random columns to test DeleteAll
  for i := 0; i < len(y); i++ {
    err = y[i].InsertG()
    if err != nil {
      t.Errorf("Unable to insert {{$tableNameSingular}}:\n%#v\nErr: %s", y[i], err)
    }
  }

  k := {{$tableNameSingular}}Slice{}
  err = {{$tableNamePlural}}G().Bind(&k)
  if err != nil {
    t.Errorf("Unable to call Bind on {{$tableNameSingular}} slice of objects: %s", err)
  }

  if len(k) != 3 {
    t.Errorf("Expected 3 results, got %d", len(k))
  }

  for i := 0; i < len(y); i++ {
    {{$varNameSingular}}CompareVals(y[i], k[i], t)
  }

  {{$varNamePlural}}DeleteAllRows(t)
}

func Test{{$tableNamePlural}}One(t *testing.T) {
  var err error

  o := {{$tableNameSingular}}{}
  if err = boil.RandomizeStruct(&o, {{$varNameSingular}}DBTypes, true); err != nil {
    t.Errorf("Unable to randomize {{$tableNameSingular}} struct: %s", err)
  }

  if err = o.InsertG(); err != nil {
    t.Errorf("Unable to insert {{$tableNameSingular}}:\n%#v\nErr: %s", o, err)
  }

  j, err := {{$tableNamePlural}}G().One()
  if err != nil {
    t.Errorf("Unable to fetch One {{$tableNameSingular}} result:\n#%v\nErr: %s", j, err)
  }

  {{$varNameSingular}}CompareVals(&o, j, t)

  {{$varNamePlural}}DeleteAllRows(t)
}

func Test{{$tableNamePlural}}All(t *testing.T) {
  var err error

  o := make({{$tableNameSingular}}Slice, 3)
  if err = boil.RandomizeSlice(&o, {{$varNameSingular}}DBTypes, true); err != nil {
    t.Errorf("Unable to randomize {{$tableNameSingular}} slice: %s", err)
  }

  // insert random columns to test DeleteAll
  for i := 0; i < len(o); i++ {
    err = o[i].InsertG()
    if err != nil {
      t.Errorf("Unable to insert {{$tableNameSingular}}:\n%#v\nErr: %s", o[i], err)
    }
  }

  j, err := {{$tableNamePlural}}G().All()
  if err != nil {
    t.Errorf("Unable to fetch All {{$tableNameSingular}} results: %s", err)
  }

  if len(j) != 3 {
    t.Errorf("Expected 3 results, got %d", len(j))
  }

  for i := 0; i < len(o); i++ {
    {{$varNameSingular}}CompareVals(o[i], j[i], t)
  }

  {{$varNamePlural}}DeleteAllRows(t)
}

func Test{{$tableNamePlural}}Count(t *testing.T) {
  var err error

  o := make({{$tableNameSingular}}Slice, 3)
  if err = boil.RandomizeSlice(&o, {{$varNameSingular}}DBTypes, true); err != nil {
    t.Errorf("Unable to randomize {{$tableNameSingular}} slice: %s", err)
  }

  // insert random columns to test Count
  for i := 0; i < len(o); i++ {
    err = o[i].InsertG()
    if err != nil {
      t.Errorf("Unable to insert {{$tableNameSingular}}:\n%#v\nErr: %s", o[i], err)
    }
  }

  c, err := {{$tableNamePlural}}G().Count()
  if err != nil {
    t.Errorf("Unable to count query {{$tableNameSingular}}: %s", err)
  }

  if c != 3 {
    t.Errorf("Expected 3 results from count {{$tableNameSingular}}, got %d", c)
  }

  {{$varNamePlural}}DeleteAllRows(t)
}
