{{- $tableNamePlural := .Table.Name | plural | titleCase -}}
{{- $tableNameSingular := .Table.Name | singular | titleCase -}}
{{- if .Table.IsJoinTable -}}
{{- else -}}
	{{- $dot := . }}
// Merge combines two {{$tableNamePlural}} into one. The primary record will be kept, and the secondary will be deleted.
func Merge{{$tableNamePlural}}(exec boil.Executor, primaryID uint64, secondaryID uint64) (err error) {
	tx, ok := exec.(boil.Transactor)
	if !ok {
		txdb, ok := exec.(boil.Beginner)
		if !ok {
			return errors.New("database does not support transactions")
		}

		tx, err = txdb.Begin()
		if err != nil {
			return err
		}

		defer func() {
			if p := recover(); p != nil {
				tx.Rollback()
				panic(p) // Rollback, then propagate panic
			} else if err != nil {
				tx.Rollback()
			} else {
				err = tx.Commit()
			}
		}()
	}

  primary, err := Find{{$tableNameSingular}}(tx, primaryID)
  if err != nil {
    return err
  } else if primary == nil {
		return errors.New("Primary {{$tableNameSingular}} not found")
	}

  secondary, err := Find{{$tableNameSingular}}(tx, secondaryID)
  if err != nil {
    return err
  } else if secondary == nil {
		return errors.New("Secondary {{$tableNameSingular}} not found")
	}

  foreignKeys := []foreignKey{
	{{- range .Tables -}}
	  {{- range .FKeys -}}
	    {{- if eq $dot.Table.Name .ForeignTable }}
		  {foreignTable: "{{.Table}}", foreignColumn: "{{.Column}}"},
      {{- end -}}
    {{- end -}}
  {{- end }}
  }

  conflictingKeys := []conflictingUniqueKey{
    {{- range .Tables -}}
      {{- $table := . -}}
      {{- range .FKeys -}}
        {{- $fk := . -}}
        {{- if eq $dot.Table.Name .ForeignTable -}}
          {{- range $table.UKeys -}}
            {{- if setInclude $fk.Column .Columns }}
              {table: "{{$fk.Table}}", objectIdColumn: "{{$fk.Column}}", columns: []string{`{{ .Columns | join "`,`" }}`}},
            {{- end -}}
          {{- end -}}
        {{- end -}}
      {{- end -}}
    {{- end }}
  }

  sqlTx, ok := tx.(*sql.Tx)
  if !ok {
    return errors.New("tx must be an sql.Tx")
  }

  err = mergeModels(sqlTx, primaryID, secondaryID, foreignKeys, conflictingKeys)
  if err != nil {
    return err
  }

	pr := reflect.ValueOf(primary)
	sr := reflect.ValueOf(secondary)
	// for any column thats null on the primary and not null on the secondary, copy from secondary to primary
	for i := 0; i < sr.Elem().NumField(); i++ {
		pf := pr.Elem().Field(i)
		sf := sr.Elem().Field(i)
		if sf.IsValid() {
			if nullable, ok := sf.Interface().(null.Nullable); ok && !nullable.IsNull() && pf.Interface().(null.Nullable).IsNull() {
				pf.Set(sf)
			}
		}
	}

	err = primary.Update(tx)
	if err != nil {
		return errors.WithStack(err)
	}

	err = secondary.Delete(tx)
	if err != nil {
		return errors.WithStack(err)
	}

  return nil
}

// Merge combines two {{$tableNamePlural}} into one. The primary record will be kept, and the secondary will be deleted.
func Merge{{$tableNamePlural}}G(primaryID uint64, secondaryID uint64) error {
  return Merge{{$tableNamePlural}}(boil.GetDB(), primaryID, secondaryID)
}
{{- end -}}{{/* join table */}}