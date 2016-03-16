{{- $tableName := titleCase .Table -}}
{{- $varName := camelCase $tableName -}}
// {{$tableName}}All retrieves all records.
func {{$tableName}}All(db boil.DB) ([]*{{$tableName}}, error) {
  var {{$varName}} []*{{$tableName}}

	rows, err := db.Query(`SELECT {{selectParamNames .Table .Columns}} FROM {{.Table}}`)
  if err != nil {
    return nil, fmt.Errorf("models: failed to query: %v", err)
  }

	for rows.Next() {
		{{$varName}}Tmp := {{$tableName}}{}

		if err := rows.Scan({{scanParamNames $varName .Columns}}); err != nil {
			return nil, fmt.Errorf("models: failed to scan row: %v", err)
		}

		{{$varName}} = append({{$varName}}, {{$varName}}Tmp)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("models: failed to read rows: %v", err)
	}

	return {{$varName}}, nil
}
