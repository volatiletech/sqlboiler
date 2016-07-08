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
	"golden",
	false,
	"Write golden files.",
)

func TestBuildQuery(t *testing.T) {
	t.Parallel()

	tests := []struct {
		q    *Query
		args []interface{}
	}{
		{&Query{table: "t"}, []interface{}{}},
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

func TestSetLimit(t *testing.T) {
	t.Parallel()

	q := &Query{}
	SetLimit(q, 10)

	expect := 10
	if q.limit != expect {
		t.Errorf("Expected %d, got %d", expect, q.limit)
	}
}

func TestSetWhere(t *testing.T) {
	t.Parallel()

	q := &Query{}
	SetWhere(q, "x > $1 AND y > $2", 5, 3)

	if len(q.where) != 1 {
		t.Errorf("Expected %d where slices, got %d", 1, len(q.where))
	}

	expect := "x > $1 AND y > $2"
	if q.where[0].clause != expect {
		t.Errorf("Expected %s, got %v", expect, q.where)
	}

	if len(q.where[0].args) != 2 {
		t.Errorf("Expected %d args, got %d", 2, len(q.where[0].args))
	}

	if q.where[0].args[0].(int) != 5 || q.where[0].args[1].(int) != 3 {
		t.Errorf("Args not set correctly, expected 5 & 3, got: %#v", q.where[0].args)
	}
}

func TestSetGroupBy(t *testing.T) {
	t.Parallel()

	q := &Query{}
	SetGroupBy(q, "col1, col2")

	expect := "col1, col2"
	if len(q.groupBy) != 1 && q.groupBy[0] != expect {
		t.Errorf("Expected %s, got %s", expect, q.groupBy[0])
	}
}

func TestSetOrderBy(t *testing.T) {
	t.Parallel()

	q := &Query{}
	SetOrderBy(q, "col1 desc, col2 asc")

	expect := "col1 desc, col2 asc"
	if len(q.orderBy) != 1 && q.orderBy[0] != expect {
		t.Errorf("Expected %s, got %s", expect, q.orderBy[0])
	}
}

func TestSetHaving(t *testing.T) {
	t.Parallel()

	q := &Query{}
	SetHaving(q, "count(orders.order_id) > 10")

	expect := "count(orders.order_id) > 10"
	if len(q.having) != 1 && q.having[0] != expect {
		t.Errorf("Expected %s, got %s", expect, q.having[0])
	}
}

func TestSetTable(t *testing.T) {
	t.Parallel()

	q := &Query{}
	SetTable(q, "videos a, orders b")

	expect := "videos a, orders b"
	if q.table != expect {
		t.Errorf("Expected %s, got %s", expect, q.table)
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

func TestSetSelect(t *testing.T) {
	t.Parallel()

	q := &Query{}
	SetSelect(q, "col1", "col2")

	if len(q.selectCols) != 2 {
		t.Errorf("Expected selectCols len 2, got %d", len(q.selectCols))
	}

	if q.selectCols[0] != "col1" && q.selectCols[1] != "col2" {
		t.Errorf("select cols value mismatch: %#v", q.selectCols)
	}
}

func TestSetInnerJoin(t *testing.T) {
	t.Parallel()

	q := &Query{}
	SetInnerJoin(q, "thing=$1 AND stuff=$2", 2, 5)

	if len(q.innerJoins) != 1 {
		t.Errorf("Expected len 1, got %d", len(q.innerJoins))
	}

	if q.innerJoins[0].on != "thing=$1 AND stuff=$2" {
		t.Errorf("Got invalid innerJoin on string: %#v", q.innerJoins)
	}

	if len(q.innerJoins[0].args) != 2 {
		t.Errorf("Expected len 2, got %d", len(q.innerJoins[0].args))
	}

	if q.innerJoins[0].args[0] != 2 && q.innerJoins[0].args[1] != 5 {
		t.Errorf("Invalid args values, got %#v", q.innerJoins[0].args)
	}
}

func TestSetOuterJoin(t *testing.T) {
	t.Parallel()
	q := &Query{}
	SetOuterJoin(q, "thing=$1 AND stuff=$2", 2, 5)

	if len(q.outerJoins) != 1 {
		t.Errorf("Expected len 1, got %d", len(q.outerJoins))
	}

	if q.outerJoins[0].on != "thing=$1 AND stuff=$2" {
		t.Errorf("Got invalid innerJoin on string: %#v", q.outerJoins)
	}

	if len(q.outerJoins[0].args) != 2 {
		t.Errorf("Expected len 2, got %d", len(q.outerJoins[0].args))
	}

	if q.outerJoins[0].args[0] != 2 && q.outerJoins[0].args[1] != 5 {
		t.Errorf("Invalid args values, got %#v", q.outerJoins[0].args)
	}
}

func TestSetLeftOuterJoin(t *testing.T) {
	t.Parallel()
	q := &Query{}
	SetLeftOuterJoin(q, "thing=$1 AND stuff=$2", 2, 5)

	if len(q.leftOuterJoins) != 1 {
		t.Errorf("Expected len 1, got %d", len(q.leftOuterJoins))
	}

	if q.leftOuterJoins[0].on != "thing=$1 AND stuff=$2" {
		t.Errorf("Got invalid innerJoin on string: %#v", q.leftOuterJoins)
	}

	if len(q.leftOuterJoins[0].args) != 2 {
		t.Errorf("Expected len 2, got %d", len(q.leftOuterJoins[0].args))
	}

	if q.leftOuterJoins[0].args[0] != 2 && q.leftOuterJoins[0].args[1] != 5 {
		t.Errorf("Invalid args values, got %#v", q.leftOuterJoins[0].args)
	}
}

func TestSetRightOuterJoin(t *testing.T) {
	t.Parallel()
	q := &Query{}
	SetRightOuterJoin(q, "thing=$1 AND stuff=$2", 2, 5)

	if len(q.rightOuterJoins) != 1 {
		t.Errorf("Expected len 1, got %d", len(q.rightOuterJoins))
	}

	if q.rightOuterJoins[0].on != "thing=$1 AND stuff=$2" {
		t.Errorf("Got invalid innerJoin on string: %#v", q.rightOuterJoins)
	}

	if len(q.rightOuterJoins[0].args) != 2 {
		t.Errorf("Expected len 2, got %d", len(q.rightOuterJoins[0].args))
	}

	if q.rightOuterJoins[0].args[0] != 2 && q.rightOuterJoins[0].args[1] != 5 {
		t.Errorf("Invalid args values, got %#v", q.rightOuterJoins[0].args)
	}
}
