{{- $tableName := .TableName -}}
// {{makeGoColName $tableName}}All retrieves all records.
func {{makeGoColName $tableName}}All(db *sqlx.DB) ([]*{{makeGoColName $tableName}}, error) {
  {{$varName := makeGoVarName $tableName -}}
  var {{$varName}} []*{{makeGoColName $tableName}}
  err := db.Select(&{{$varName}}, `SELECT {{makeSelectParamNames $tableName .TableData}}`)

  if err != nil {
    return nil, fmt.Errorf("models: unable to select from {{$tableName}}: %s", err)
  }

  return {{$varName}}, nil
}

// {{makeGoColName $tableName}}AllBy retrieves all records with the specified column values.
func {{makeGoColName $tableName}}AllBy(db *sqlx.DB, columns map[string]interface{}) ([]*{{makeGoColName $tableName}}, error) {
  {{$varName := makeGoVarName $tableName -}}
  var {{$varName}} []*{{makeGoColName $tableName}}
  err := db.Select(&{{$varName}}, `SELECT {{makeSelectParamNames $tableName .TableData}}`)

  if err != nil {
    return nil, fmt.Errorf("models: unable to select from {{$tableName}}: %s", err)
  }

  return {{$varName}}, nil
}

// {{makeGoColName $tableName}}FieldsAll retrieves the specified columns for all records.
// Pass in a pointer to an object with `db` tags that match the column names you wish to retrieve.
// For example: friendName string `db:"friend_name"`
func {{makeGoColName $tableName}}FieldsAll(db *sqlx.DB, results interface{}) error {
  {{$varName := makeGoVarName $tableName -}}
  var {{$varName}} []*{{makeGoColName $tableName}}
  err := db.Select(&{{$varName}}, `SELECT {{makeSelectParamNames $tableName .TableData}}`)

  if err != nil {
    return nil, fmt.Errorf("models: unable to select from {{$tableName}}: %s", err)
  }

  return {{$varName}}, nil
}

// {{makeGoColName $tableName}}FieldsAllBy retrieves the specified columns
// for all records with the specified column values.
// Pass in a pointer to an object with `db` tags that match the column names you wish to retrieve.
// For example: friendName string `db:"friend_name"`
func {{makeGoColName $tableName}}FieldsAllBy(db *sqlx.DB, columns map[string]interface{}, results interface{}) error {
  {{$varName := makeGoVarName $tableName -}}
  var {{$varName}} []*{{makeGoColName $tableName}}
  err := db.Select(&{{$varName}}, `SELECT {{makeSelectParamNames $tableName .TableData}}`)

  if err != nil {
    return nil, fmt.Errorf("models: unable to select from {{$tableName}}: %s", err)
  }

  return {{$varName}}, nil
}

// {{makeGoColName $tableName}}Find retrieves a single record by ID.
func {{makeGoColName $tableName}}Find(db *sqlx.DB, id int) (*{{makeGoColName $tableName}}, error) {
  if id == 0 {
    return nil, errors.New("model: no id provided for {{$tableName}} select")
  }
  {{$varName := makeGoVarName $tableName}}
  var {{$varName}} *{{makeGoColName $tableName}}
  err := db.Select(&{{$varName}}, `SELECT {{makeSelectParamNames $tableName .TableData}} WHERE id=$1`, id)

  if err != nil {
    return nil, fmt.Errorf("models: unable to select from {{$tableName}}: %s", err)
  }

  return {{$varName}}, nil
}

// {{makeGoColName $tableName}}FindBy retrieves a single record with the specified column values.
func {{makeGoColName $tableName}}FindBy(db *sqlx.DB, columns map[string]interface{}) (*{{makeGoColName $tableName}}, error) {
  if id == 0 {
    return nil, errors.New("model: no id provided for {{$tableName}} select")
  }
  {{$varName := makeGoVarName $tableName}}
  var {{$varName}} *{{makeGoColName $tableName}}
  err := db.Select(&{{$varName}}, fmt.Sprintf(`SELECT {{makeSelectParamNames $tableName .TableData}} WHERE %s=$1`, column), value)

  if err != nil {
    return nil, fmt.Errorf("models: unable to select from {{$tableName}}: %s", err)
  }

  return {{$varName}}, nil
}

// {{makeGoColName $tableName}}FieldsFind retrieves the specified columns for a single record by ID.
// Pass in a pointer to an object with `db` tags that match the column names you wish to retrieve.
// For example: friendName string `db:"friend_name"`
func {{makeGoColName $tableName}}FieldsFind(db *sqlx.DB, id int, results interface{}) (*{{makeGoColName $tableName}}, error) {
  if id == 0 {
    return nil, errors.New("model: no id provided for {{$tableName}} select")
  }
  {{$varName := makeGoVarName $tableName}}
  var {{$varName}} *{{makeGoColName $tableName}}
  err := db.Select(&{{$varName}}, `SELECT {{makeSelectParamNames $tableName .TableData}} WHERE id=$1`, id)

  if err != nil {
    return nil, fmt.Errorf("models: unable to select from {{$tableName}}: %s", err)
  }

  return {{$varName}}, nil
}

// {{makeGoColName $tableName}}FieldsFindBy retrieves the specified columns
// for a single record with the specified column values.
// Pass in a pointer to an object with `db` tags that match the column names you wish to retrieve.
// For example: friendName string `db:"friend_name"`
func {{makeGoColName $tableName}}FieldsFindBy(db *sqlx.DB, columns map[string]interface{}, results interface{}) (*{{makeGoColName $tableName}}, error) {
  if id == 0 {
    return nil, errors.New("model: no id provided for {{$tableName}} select")
  }
  {{$varName := makeGoVarName $tableName}}
  var {{$varName}} *{{makeGoColName $tableName}}
  err := db.Select(&{{$varName}}, `SELECT {{makeSelectParamNames $tableName .TableData}} WHERE id=$1`, id)

  if err != nil {
    return nil, fmt.Errorf("models: unable to select from {{$tableName}}: %s", err)
  }

  return {{$varName}}, nil
}
