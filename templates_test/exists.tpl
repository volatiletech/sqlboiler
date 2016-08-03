{{- $tableNameSingular := .Table.Name | singular | titleCase -}}
{{- $dbName := singular .Table.Name -}}
{{- $tableNamePlural := .Table.Name | plural | titleCase -}}
{{- $varNamePlural := .Table.Name | plural | camelCase -}}
{{- $varNameSingular := .Table.Name | singular | camelCase -}}
{{- $pkeyArgs := .Table.PKey.Columns | stringMap .StringFuncs.titleCase | prefixStringSlice "o." | join ", " -}}
func Test{{$tableNamePlural}}Exists(t *testing.T) {
  var err error

  o := {{$tableNameSingular}}{}
  if err = boil.RandomizeStruct(&o, {{$varNameSingular}}DBTypes, true); err != nil {
    t.Errorf("Unable to randomize {{$tableNameSingular}} struct: %s", err)
  }

  if err = o.InsertG(); err != nil {
    t.Errorf("Unable to insert {{$tableNameSingular}}:\n%#v\nErr: %s", o, err)
  }

  // Check Exists finds existing rows
  e, err := {{$tableNameSingular}}ExistsG({{$pkeyArgs}})
  if err != nil {
    t.Errorf("Unable to check if {{$tableNameSingular}} exists: %s", err)
  }
  if e != true {
    t.Errorf("Expected {{$tableNameSingular}}ExistsG to return true, but got false.")
  }

  whereClause := strmangle.WhereClause({{$varNameSingular}}PrimaryKeyColumns, 1)
  e, err = {{$tableNamePlural}}G(qm.Where(whereClause, boil.GetStructValues(o, {{$varNameSingular}}PrimaryKeyColumns...)...)).Exists()
  if err != nil {
    t.Errorf("Unable to check if {{$tableNameSingular}} exists: %s", err)
  }
  if e != true {
    t.Errorf("Expected ExistsG to return true, but got false.")
  }

  // Check Exists does not find non-existing rows
  o = {{$tableNameSingular}}{}
  e, err = {{$tableNameSingular}}ExistsG({{$pkeyArgs}})
  if err != nil {
    t.Errorf("Unable to check if {{$tableNameSingular}} exists: %s", err)
  }
  if e != false {
    t.Errorf("Expected {{$tableNameSingular}}ExistsG to return false, but got true.")
  }

  e, err = {{$tableNamePlural}}G(qm.Where(whereClause, boil.GetStructValues(o, {{$varNameSingular}}PrimaryKeyColumns...)...)).Exists()
  if err != nil {
    t.Errorf("Unable to check if {{$tableNameSingular}} exists: %s", err)
  }
  if e != false {
    t.Errorf("Expected ExistsG to return false, but got true.")
  }

  {{$varNamePlural}}DeleteAllRows(t)
}
