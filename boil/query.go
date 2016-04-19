package boil

import "database/sql"

type where struct {
	clause string
	args   []interface{}
}

type Query struct {
	executor   Executor
	delete     bool
	update     map[string]interface{}
	selectCols []string
	from       string
	joins      []string
	where      []where
	groupBy    []string
	orderBy    []string
	having     []string
	limit      int
}

func SetDelete(q *Query, flag bool) {
	q.delete = flag
}

func SetUpdate(q *Query, cols map[string]interface{}) {
	q.update = cols
}

func SetExecutor(q *Query, exec Executor) {
	q.executor = exec
}

func SetSelect() {

}

func SetFrom() {

}

func SetJoins() {

}

func SetWhere() {

}

func SetGroupBy() {

}

func SetOrderBy() {

}

func SetHaving() {

}

func SetLimit() {

}

func ExecQuery(q *Query) error {
	return nil
}

func ExecQueryOne(q *Query) (*sql.Row, error) {
	return nil, nil
}

func ExecQueryAll(q *Query) (*sql.Rows, error) {
	return nil, nil
}

func buildQuery(q *Query) string {
	return ""
}
