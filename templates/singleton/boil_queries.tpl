var dialect boil.Dialect = boil.Dialect{
  LQ: 0x{{printf "%x" .Dialect.LQ}},
  RQ: 0x{{printf "%x" .Dialect.RQ}},
  IndexPlaceholders: {{.Dialect.IndexPlaceholders}},
}

// NewQueryG initializes a new Query using the passed in QueryMods
func NewQueryG(mods ...qm.QueryMod) *boil.Query {
  return NewQuery(boil.GetDB(), mods...)
}

// NewQuery initializes a new Query using the passed in QueryMods
func NewQuery(exec boil.Executor, mods ...qm.QueryMod) *boil.Query {
  q := &boil.Query{}
  boil.SetExecutor(q, exec)
  boil.SetDialect(q, &dialect)
  qm.Apply(q, mods...)

  return q
}
