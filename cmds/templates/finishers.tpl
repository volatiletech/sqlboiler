{{- $tableNameSingular := titleCaseSingular .Table.Name -}}
{{- $varNameSingular := camelCaseSingular .Table.Name -}}
type {{$varNameSingular}}Slice []*{{$tableNameSingular}}

func (q {{$varNameSingular}}Query) One() (*{{$tableNameSingular}}, error) {
  //var o *{{$tableNameSingular}}

  //qs.Apply(q, qs.Limit(1))

  //res := boil.ExecQueryOne(q)

return nil, nil
}

func (q {{$varNameSingular}}Query) All() ({{$varNameSingular}}Slice, error) {
return nil, nil
}

func (q {{$varNameSingular}}Query) Count() (int64, error) {
return 0, nil
}
