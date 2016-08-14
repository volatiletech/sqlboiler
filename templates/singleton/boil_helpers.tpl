// NewQueryG initializes a new Query using the passed in QueryMods
func NewQueryG(mods ...qm.QueryMod) *boil.Query {
	return NewQuery(boil.GetDB(), mods...)
}

// NewQuery initializes a new Query using the passed in QueryMods
func NewQuery(exec boil.Executor, mods ...qm.QueryMod) *boil.Query {
	q := &boil.Query{}
	boil.SetExecutor(q, exec)
	qm.Apply(q, mods...)

	return q
}

// generateUpsertQuery builds a SQL statement string using the upsertData provided.
func generateUpsertQuery(tableName string, updateOnConflict bool, ret, update, conflict, whitelist []string) string {
  conflict = strmangle.IdentQuoteSlice(conflict)
  whitelist = strmangle.IdentQuoteSlice(whitelist)
  ret = strmangle.IdentQuoteSlice(ret)

	buf := strmangle.GetBuffer()
	defer strmangle.PutBuffer(buf)

  fmt.Fprintf(
		buf,
    "INSERT INTO %s (%s) VALUES (%s) ON CONFLICT ",
		tableName,
    strings.Join(whitelist, ", "),
    strmangle.Placeholders(len(whitelist), 1, 1),
  )

  if !updateOnConflict {
    buf.WriteString("DO NOTHING")
  } else {
		buf.WriteByte('(')
		buf.WriteString(strings.Join(conflict, ", "))
		buf.WriteString(") DO UPDATE SET")

	  for i, v := range update {
			if i != 0 {
				buf.WriteByte(',')
			}
	    quoted := strmangle.IdentQuote(v)
			buf.WriteString(quoted)
			buf.WriteString(" = EXCLUDED.")
			buf.WriteString(quoted)
	  }
  }

  if len(ret) != 0 {
		buf.WriteString(" RETURNING ")
		buf.WriteString(strings.Join(ret, ", "))
  }

  return buf.String()
}
