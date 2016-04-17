{{- $tableNameSingular := titleCaseSingular .Table.Name -}}
{{- $varNameSingular := camelCaseSingular .Table.Name -}}
type {{$tableNameSingular}}Slice []*{{$tableNameSingular}}

func ({{$varNameSingular}}Query) One() (*{{$tableNameSingular}}, error) {

}

func ({{$varNameSingular}}Query) All() ({{$tableNameSingular}}Slice, error) {

}

func ({{$varNameSingular}}Query) Count() (int64, error) {

}
