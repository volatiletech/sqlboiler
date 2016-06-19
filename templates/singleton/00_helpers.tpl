// M type is for providing where filters to Where helpers.
type M map[string]interface{}

// NewQuery initializes a new Query using the passed in QueryMods
func NewQuery(mods ...qm.QueryMod) *boil.Query {
	return NewQueryX(boil.GetDB(), mods...)
}

// NewQueryX initializes a new Query using the passed in QueryMods
func NewQueryX(exec boil.Executor, mods ...qm.QueryMod) *boil.Query {
	q := &boil.Query{}
	boil.SetExecutor(q, exec)
	qm.Apply(q, mods...)

	return q
}
