{{- $tableNameSingular := titleCaseSingular .Table.Name -}}
{{- $varNameSingular := camelCaseSingular .Table.Name -}}
type {{$tableNameSingular}}Slice []*{{$tableNameSingular}}

func ({{$varNameSingular}}Query) One() (*{{$tableNameSingular}}, error) {
return nil, nil
}

func ({{$varNameSingular}}Query) All() ({{$tableNameSingular}}Slice, error) {
return nil, nil
}

func ({{$varNameSingular}}Query) Count() (int64, error) {
return 0, nil
}
