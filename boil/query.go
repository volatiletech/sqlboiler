package boil

type where struct {
	clause string
	args   []interface{}
}

type Query struct {
	limit    int
	where    []where
	executor Executor
	groupBy  []string
	orderBy  []string
	having   []string
	from     string
}

func (q *Query) buildQuery() string {
	return ""
}

// NewQuery initializes a new Query using the passed in QueryMods
func NewQuery(mods ...QueryMod) *Query {
	return NewQueryX(currentDB, mods...)
}

// NewQueryX initializes a new Query using the passed in QueryMods
func NewQueryX(executor Executor, mods ...QueryMod) *Query {
	q := &Query{executor: executor}
	q.Apply(mods...)

	return q
}
