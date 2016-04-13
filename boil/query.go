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

// makes a new empty query ?????
func New() {

}
