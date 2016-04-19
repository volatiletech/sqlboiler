// M type is for providing where filters to Where helpers.
type M map[string]interface{}

// NewQuery initializes a new Query using the passed in QueryMods
func NewQuery(mods ...qs.QueryMod) *boil.Query {
	return NewQueryX(boil.GetDB(), mods...)
}

// NewQueryX initializes a new Query using the passed in QueryMods
func NewQueryX(executor boil.Executor, mods ...qs.QueryMod) *boil.Query {
	q := &Query{executor: executor}
	q.Apply(mods...)

	return q
}
