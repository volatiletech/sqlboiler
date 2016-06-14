{{- $tableNameSingular := titleCaseSingular .Table.Name -}}
{{- $dbName := singular .Table.Name -}}
{{- $tableNamePlural := titleCasePlural .Table.Name -}}
{{- $varNamePlural := camelCasePlural .Table.Name -}}
{{- $varNameSingular := camelCaseSingular .Table.Name -}}
func Test{{$tableNamePlural}}Bind(t *testing.T) {
  var err error

  {{$varNamePlural}}DeleteAllRows(t)

  o := {{$tableNameSingular}}{}
  if err = boil.RandomizeStruct(&o); err != nil {
    t.Errorf("Unable to randomize {{$tableNameSingular}} struct: %s", err)
  }

  if err = o.Insert(); err != nil {
    t.Errorf("Unable to insert {{$tableNameSingular}}:\n%#v\nErr: %s", o, err)
  }

  j := {{$tableNameSingular}}{}

  err = {{$tableNamePlural}}(qm.Where("{{wherePrimaryKey .Table.PKey.Columns 1}}", {{titleCaseCommaList "o." .Table.PKey.Columns}})).Bind(&j)
  if err != nil {
    t.Errorf("Unable to call Bind on {{$tableNameSingular}} single object: %s", err)
  }

  {{$varNameSingular}}CompareVals(&o, &j, t)

  // insert 3 rows, attempt to bind into slice
  {{$varNamePlural}}DeleteAllRows(t)

  y := make({{$varNameSingular}}Slice, 3)
  if err = boil.RandomizeSlice(&y); err != nil {
    t.Errorf("Unable to randomize {{$tableNameSingular}} slice: %s", err)
  }

  // insert random columns to test DeleteAll
  for i := 0; i < len(y); i++ {
    err = y[i].Insert()
    if err != nil {
      t.Errorf("Unable to insert {{$tableNameSingular}}:\n%#v\nErr: %s", y[i], err)
    }
  }

  k := {{$varNameSingular}}Slice{}
  err = {{$tableNamePlural}}().Bind(&k)
  if err != nil {
    t.Errorf("Unable to call Bind on {{$tableNameSingular}} slice of objects: %s", err)
  }
  
  if len(k) != 3 {
    t.Errorf("Expected 3 results, got %d", len(k))
  }

  for i := 0; i < len(y); i++ {
    {{$varNameSingular}}CompareVals(y[i], k[i], t)
  }
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

  {{$varNameSingular}}CompareVals(&o, j, t)
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
    {{$varNameSingular}}CompareVals(o[i], j[i], t)
  }
}

func Test{{$tableNamePlural}}Count(t *testing.T) {
  var err error

  {{$varNamePlural}}DeleteAllRows(t)

  o := make({{$varNameSingular}}Slice, 3)
  if err = boil.RandomizeSlice(&o); err != nil {
    t.Errorf("Unable to randomize {{$tableNameSingular}} slice: %s", err)
  }

  // insert random columns to test Count
  for i := 0; i < len(o); i++ {
    err = o[i].Insert()
    if err != nil {
      t.Errorf("Unable to insert {{$tableNameSingular}}:\n%#v\nErr: %s", o[i], err)
    }
  }

  c, err := {{$tableNamePlural}}().Count()
  if err != nil {
    t.Errorf("Unable to count query {{$tableNameSingular}}: %s", err)
  }

  if c != 3 {
    t.Errorf("Expected 3 results from count {{$tableNameSingular}}, got %d", c)
  }
}
