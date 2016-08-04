package boil

import (
	"bytes"
	"database/sql"
	"flag"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

var writeGoldenFiles = flag.Bool(
	"test.golden",
	false,
	"Write golden files.",
)

func TestBuildQuery(t *testing.T) {
	t.Parallel()

	tests := []struct {
		q    *Query
		args []interface{}
	}{
		{&Query{from: []string{"t"}}, nil},
		{&Query{from: []string{"q"}, limit: 5, offset: 6}, nil},
		{&Query{from: []string{"q"}, orderBy: []string{"a ASC", "b DESC"}}, nil},
		{&Query{from: []string{"t"}, selectCols: []string{"count(*) as ab, thing as bd", `"stuff"`}}, nil},
		{&Query{from: []string{"a", "b"}, selectCols: []string{"count(*) as ab, thing as bd", `"stuff"`}}, nil},
	}

	for i, test := range tests {
		filename := filepath.Join("_fixtures", fmt.Sprintf("%02d.sql", i))
		out, args := buildQuery(test.q)

		if *writeGoldenFiles {
			err := ioutil.WriteFile(filename, []byte(out), 0664)
			if err != nil {
				t.Fatalf("Failed to write golden file %s: %s\n", filename, err)
			}
			t.Logf("wrote golden file: %s\n", filename)
			continue
		}

		byt, err := ioutil.ReadFile(filename)
		if err != nil {
			t.Fatalf("Failed to read golden file %q: %v", filename, err)
		}

		if string(bytes.TrimSpace(byt)) != out {
			t.Errorf("[%02d] Test failed:\nWant:\n%s\nGot:\n%s", i, byt, out)
		}

		if !reflect.DeepEqual(args, test.args) {
			t.Errorf("[%02d] Test failed:\nWant:\n%s\nGot:\n%s", i, spew.Sdump(test.args), spew.Sdump(args))
		}
	}
}

func TestSetLastWhereAsOr(t *testing.T) {
	t.Parallel()
	q := &Query{}

	AppendWhere(q, "")

	if q.where[0].orSeperator {
		t.Errorf("Do not want or seperator")
	}

	SetLastWhereAsOr(q)

	if len(q.where) != 1 {
		t.Errorf("Want len 1")
	}
	if !q.where[0].orSeperator {
		t.Errorf("Want or seperator")
	}

	AppendWhere(q, "")
	SetLastWhereAsOr(q)

	if len(q.where) != 2 {
		t.Errorf("Want len 2")
	}
	if q.where[0].orSeperator != true {
		t.Errorf("Expected true")
	}
	if q.where[1].orSeperator != true {
		t.Errorf("Expected true")
	}
}

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
	expect := "count(orders.order_id) > 10"
	ApplyHaving(q, expect)
	ApplyHaving(q, expect)

	if len(q.having) != 2 && (q.having[0] != expect || q.having[1] != expect) {
		t.Errorf("Expected %s, got %s %s", expect, q.having[0], q.having[1])
	}

	SetHaving(q, expect)
	if len(q.having) != 1 && q.having[0] != expect {
		t.Errorf("Expected %s, got %s", expect, q.having[0])
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

func TestSetDelete(t *testing.T) {
	t.Parallel()

	q := &Query{}
	SetDelete(q)

	if q.delete != true {
		t.Errorf("Expected %t, got %t", true, q.delete)
	}
}

func TestSetUpdate(t *testing.T) {
	t.Parallel()

	q := &Query{}
	SetUpdate(q, map[string]interface{}{"col1": 1, "col2": 2})

	if len(q.update) != 2 {
		t.Errorf("Expected len 2, got %d", len(q.update))
	}

	if q.update["col1"] != 1 && q.update["col2"] != 2 {
		t.Errorf("Value misatch: %#v", q.update)
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

func TestSelect(t *testing.T) {
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

func TestInnerJoin(t *testing.T) {
	t.Parallel()

	q := &Query{}
	AppendInnerJoin(q, "thing=$1 AND stuff=$2", 2, 5)
	AppendInnerJoin(q, "thing=$1 AND stuff=$2", 2, 5)

	if len(q.innerJoins) != 2 {
		t.Errorf("Expected len 1, got %d", len(q.innerJoins))
	}

	if q.innerJoins[0].on != "thing=$1 AND stuff=$2" {
		t.Errorf("Got invalid innerJoin on string: %#v", q.innerJoins)
	}
	if q.innerJoins[1].on != "thing=$1 AND stuff=$2" {
		t.Errorf("Got invalid innerJoin on string: %#v", q.innerJoins)
	}

	if len(q.innerJoins[0].args) != 2 {
		t.Errorf("Expected len 2, got %d", len(q.innerJoins[0].args))
	}
	if len(q.innerJoins[1].args) != 2 {
		t.Errorf("Expected len 2, got %d", len(q.innerJoins[1].args))
	}

	if q.innerJoins[0].args[0] != 2 && q.innerJoins[0].args[1] != 5 {
		t.Errorf("Invalid args values, got %#v", q.innerJoins[0].args)
	}

	SetInnerJoin(q, "thing=$1 AND stuff=$2", 2, 5)
	if len(q.innerJoins) != 1 {
		t.Errorf("Expected len 1, got %d", len(q.innerJoins))
	}

	if q.innerJoins[0].on != "thing=$1 AND stuff=$2" {
		t.Errorf("Got invalid innerJoin on string: %#v", q.innerJoins)
	}
}
