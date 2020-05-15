package queries

// QueryReset holds the state for the built up query
type QueryReset struct {
	q     *Query
	undo  []undoFunc
	saved Query
}

type undoFunc = func(*QueryReset)

// Save removes the effect of a finisher
func (q *QueryReset) Save() {
	q.saved = *q.q
}

// Reset removes the effect of a finisher
func (q *QueryReset) Reset() {
	*q.q = q.saved
}

// AddReset adds makes a QueryReset type from a *Query
func (q *Query) AddReset() *QueryReset {
	return &QueryReset{q: q}
}
