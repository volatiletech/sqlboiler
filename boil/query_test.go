package boil

import (
	"database/sql"
	"reflect"
	"testing"
)

func TestSetLimit(t *testing.T) {
	t.Parallel()

	q := &Query{}
	SetLimit(q, 10)

	expect := 10
	if q.limit != expect {
		t.Errorf("Expected %d, got %d", expect, q.limit)
	}
}

func TestSetOffset(t *testing.T) {
	t.Parallel()

	q := &Query{}
	SetOffset(q, 10)

	expect := 10
	if q.offset != expect {
		t.Errorf("Expected %d, got %d", expect, q.offset)
	}
}

func TestSetSQL(t *testing.T) {
	t.Parallel()

	q := &Query{}
	SetSQL(q, "select * from thing", 5, 3)

	if len(q.plainSQL.args) != 2 {
		t.Errorf("Expected len 2, got %d", len(q.plainSQL.args))
	}

	if q.plainSQL.sql != "select * from thing" {
		t.Errorf("Was not expected string, got %s", q.plainSQL.sql)
	}
}

func TestWhere(t *testing.T) {
	t.Parallel()

	q := &Query{}
	expect := "x > $1 AND y > $2"
	AppendWhere(q, expect, 5, 3)
	AppendWhere(q, expect, 5, 3)

	if len(q.where) != 2 {
		t.Errorf("%#v", q.where)
	}

	if q.where[0].clause != expect || q.where[1].clause != expect {
		t.Errorf("Expected %s, got %#v", expect, q.where)
	}

	if len(q.where[0].args) != 2 || len(q.where[0].args) != 2 {
		t.Errorf("arg length wrong: %#v", q.where)
	}

	if q.where[0].args[0].(int) != 5 || q.where[0].args[1].(int) != 3 {
		t.Errorf("args wrong: %#v", q.where)
	}

	SetWhere(q, expect, 5, 3)
	if q.where[0].clause != expect {
		t.Errorf("Expected %s, got %v", expect, q.where)
	}

	if len(q.where[0].args) != 2 {
		t.Errorf("Expected %d args, got %d", 2, len(q.where[0].args))
	}

	if q.where[0].args[0].(int) != 5 || q.where[0].args[1].(int) != 3 {
		t.Errorf("Args not set correctly, expected 5 & 3, got: %#v", q.where[0].args)
	}

	if len(q.where) != 1 {
		t.Errorf("%#v", q.where)
	}
}

func TestSetLastWhereAsOr(t *testing.T) {
	t.Parallel()
	q := &Query{}

	AppendWhere(q, "")

	if q.where[0].orSeparator {
		t.Errorf("Do not want or separator")
	}

	SetLastWhereAsOr(q)

	if len(q.where) != 1 {
		t.Errorf("Want len 1")
	}
	if !q.where[0].orSeparator {
		t.Errorf("Want or separator")
	}

	AppendWhere(q, "")
	SetLastWhereAsOr(q)

	if len(q.where) != 2 {
		t.Errorf("Want len 2")
	}
	if q.where[0].orSeparator != true {
		t.Errorf("Expected true")
	}
	if q.where[1].orSeparator != true {
		t.Errorf("Expected true")
	}
}

func TestGroupBy(t *testing.T) {
	t.Parallel()

	q := &Query{}
	expect := "col1, col2"
	ApplyGroupBy(q, expect)
	ApplyGroupBy(q, expect)

	if len(q.groupBy) != 2 && (q.groupBy[0] != expect || q.groupBy[1] != expect) {
		t.Errorf("Expected %s, got %s %s", expect, q.groupBy[0], q.groupBy[1])
	}

	SetGroupBy(q, expect)
	if len(q.groupBy) != 1 && q.groupBy[0] != expect {
		t.Errorf("Expected %s, got %s", expect, q.groupBy[0])
	}
}

func TestOrderBy(t *testing.T) {
	t.Parallel()

	q := &Query{}
	expect := "col1 desc, col2 asc"
	ApplyOrderBy(q, expect)
	ApplyOrderBy(q, expect)

	if len(q.orderBy) != 2 && (q.orderBy[0] != expect || q.orderBy[1] != expect) {
		t.Errorf("Expected %s, got %s %s", expect, q.orderBy[0], q.orderBy[1])
	}

	SetOrderBy(q, "col1 desc, col2 asc")
	if len(q.orderBy) != 1 && q.orderBy[0] != expect {
		t.Errorf("Expected %s, got %s", expect, q.orderBy[0])
	}
}

func TestHaving(t *testing.T) {
	t.Parallel()

	q := &Query{}
	expect := "count(orders.order_id) > ?"
	ApplyHaving(q, expect, 10)
	ApplyHaving(q, expect, 10)

	if len(q.having) != 2 {
		t.Errorf("Expected 2, got %d", len(q.having))
	}

	if q.having[0].clause != expect || q.having[1].clause != expect {
		t.Errorf("Expected %s, got %s %s", expect, q.having[0].clause, q.having[1].clause)
	}

	if q.having[0].args[0] != 10 || q.having[1].args[0] != 10 {
		t.Errorf("Expected %v, got %v %v", 10, q.having[0].args[0], q.having[1].args[0])
	}

	SetHaving(q, expect, 10)
	if len(q.having) != 1 && (q.having[0].clause != expect || q.having[0].args[0] != 10) {
		t.Errorf("Expected %s, got %s %v", expect, q.having[0], q.having[0].args[0])
	}
}

func TestFrom(t *testing.T) {
	t.Parallel()

	q := &Query{}
	AppendFrom(q, "videos a", "orders b")
	AppendFrom(q, "videos a", "orders b")

	expect := []string{"videos a", "orders b", "videos a", "orders b"}
	if !reflect.DeepEqual(q.from, expect) {
		t.Errorf("Expected %s, got %s", expect, q.from)
	}

	SetFrom(q, "videos a", "orders b")
	if !reflect.DeepEqual(q.from, expect[:2]) {
		t.Errorf("Expected %s, got %s", expect, q.from)
	}
}

func TestSetCount(t *testing.T) {
	t.Parallel()

	q := &Query{}
	SetCount(q)

	if q.modFunction != "COUNT" {
		t.Errorf("Wrong modFunction, got %s", q.modFunction)
	}
}

func TestSetAvg(t *testing.T) {
	t.Parallel()

	q := &Query{}
	SetAvg(q)

	if q.modFunction != "AVG" {
		t.Errorf("Wrong modFunction, got %s", q.modFunction)
	}
}

func TestSetMax(t *testing.T) {
	t.Parallel()

	q := &Query{}
	SetMax(q)

	if q.modFunction != "MAX" {
		t.Errorf("Wrong modFunction, got %s", q.modFunction)
	}
}

func TestSetMin(t *testing.T) {
	t.Parallel()

	q := &Query{}
	SetMin(q)

	if q.modFunction != "MIN" {
		t.Errorf("Wrong modFunction, got %s", q.modFunction)
	}
}

func TestSetSum(t *testing.T) {
	t.Parallel()

	q := &Query{}
	SetSum(q)

	if q.modFunction != "SUM" {
		t.Errorf("Wrong modFunction, got %s", q.modFunction)
	}
}

func TestSetUpdate(t *testing.T) {
	t.Parallel()

	q := &Query{}
	SetUpdate(q, map[string]interface{}{"test": 5})

	if q.update["test"] != 5 {
		t.Errorf("Wrong update, got %v", q.update)
	}
}

func TestSetDelete(t *testing.T) {
	t.Parallel()

	q := &Query{}
	SetDelete(q)

	if q.delete != true {
		t.Errorf("Expected %t, got %t", true, q.delete)
	}
}

func TestSetExecutor(t *testing.T) {
	t.Parallel()

	q := &Query{}
	d := &sql.DB{}
	SetExecutor(q, d)

	if q.executor != d {
		t.Errorf("Expected executor to get set to d, but was: %#v", q.executor)
	}
}

func TestAppendSelect(t *testing.T) {
	t.Parallel()

	q := &Query{}
	AppendSelect(q, "col1", "col2")
	AppendSelect(q, "col1", "col2")

	if len(q.selectCols) != 4 {
		t.Errorf("Expected selectCols len 4, got %d", len(q.selectCols))
	}

	if q.selectCols[0] != `col1` && q.selectCols[1] != `col2` {
		t.Errorf("select cols value mismatch: %#v", q.selectCols)
	}
	if q.selectCols[2] != `col1` && q.selectCols[3] != `col2` {
		t.Errorf("select cols value mismatch: %#v", q.selectCols)
	}

	SetSelect(q, "col1", "col2")
	if q.selectCols[0] != `col1` && q.selectCols[1] != `col2` {
		t.Errorf("select cols value mismatch: %#v", q.selectCols)
	}
}

func TestSelect(t *testing.T) {
	t.Parallel()

	q := &Query{}
	q.selectCols = []string{"one"}

	ret := Select(q)
	if ret[0] != "one" {
		t.Errorf("Expected %q, got %s", "one", ret[0])
	}
}

func TestSQL(t *testing.T) {
	t.Parallel()

	q := SQL("thing", 5)
	if q.plainSQL.sql != "thing" {
		t.Errorf("Expected %q, got %s", "thing", q.plainSQL.sql)
	}
	if q.plainSQL.args[0].(int) != 5 {
		t.Errorf("Expected 5, got %v", q.plainSQL.args[0])
	}
}

func TestInnerJoin(t *testing.T) {
	t.Parallel()

	q := &Query{}
	AppendInnerJoin(q, "thing=$1 AND stuff=$2", 2, 5)
	AppendInnerJoin(q, "thing=$1 AND stuff=$2", 2, 5)

	if len(q.joins) != 2 {
		t.Errorf("Expected len 1, got %d", len(q.joins))
	}

	if q.joins[0].clause != "thing=$1 AND stuff=$2" {
		t.Errorf("Got invalid innerJoin on string: %#v", q.joins)
	}
	if q.joins[1].clause != "thing=$1 AND stuff=$2" {
		t.Errorf("Got invalid innerJoin on string: %#v", q.joins)
	}

	if len(q.joins[0].args) != 2 {
		t.Errorf("Expected len 2, got %d", len(q.joins[0].args))
	}
	if len(q.joins[1].args) != 2 {
		t.Errorf("Expected len 2, got %d", len(q.joins[1].args))
	}

	if q.joins[0].args[0] != 2 && q.joins[0].args[1] != 5 {
		t.Errorf("Invalid args values, got %#v", q.joins[0].args)
	}

	SetInnerJoin(q, "thing=$1 AND stuff=$2", 2, 5)
	if len(q.joins) != 1 {
		t.Errorf("Expected len 1, got %d", len(q.joins))
	}

	if q.joins[0].clause != "thing=$1 AND stuff=$2" {
		t.Errorf("Got invalid innerJoin on string: %#v", q.joins)
	}
}
