package boil

import (
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
		{&Query{from: "t"}, []interface{}{}},
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

		if string(byt) != out {
			t.Errorf("[%02d] Test failed:\nWant:\n%s\nGot:\n%s", i, byt, out)
		}

		if !reflect.DeepEqual(args, test.args) {
			t.Errorf("[%02d] Test failed:\nWant:\n%s\nGot:\n%s", i, spew.Sdump(test.args), spew.Sdump(args))
		}
	}
}

// func TestApply(t *testing.T) {
// 	t.Parallel()
//
// 	q := &boil.Query{}
// 	qfn1 := Limit(10)
// 	qfn2 := Where("x > $1 AND y > $2", 5, 3)
//
// 	q.Apply(qfn1, qfn2)
//
// 	expect1 := 10
// 	if q.limit != expect1 {
// 		t.Errorf("Expected %d, got %d", expect1, q.limit)
// 	}
//
// 	expect2 := "x > $1 AND y > $2"
// 	if len(q.where) != 1 {
// 		t.Errorf("Expected %d where slices, got %d", len(q.where))
// 	}
//
// 	expect := "x > $1 AND y > $2"
// 	if q.where[0].clause != expect2 {
// 		t.Errorf("Expected %s, got %s", expect, q.where)
// 	}
//
// 	if len(q.where[0].args) != 2 {
// 		t.Errorf("Expected %d args, got %d", 2, len(q.where[0].args))
// 	}
//
// 	if q.where[0].args[0].(int) != 5 || q.where[0].args[1].(int) != 3 {
// 		t.Errorf("Args not set correctly, expected 5 & 3, got: %#v", q.where[0].args)
// 	}
// }
//
// func TestLimit(t *testing.T) {
// 	t.Parallel()
//
// 	q := &boil.Query{}
// 	qfn := Limit(10)
//
// 	qfn(q)
//
// 	expect := 10
// 	if q.limit != expect {
// 		t.Errorf("Expected %d, got %d", expect, q.limit)
// 	}
// }
//
// func TestWhere(t *testing.T) {
// 	t.Parallel()
//
// 	q := &boil.Query{}
// 	qfn := Where("x > $1 AND y > $2", 5, 3)
//
// 	qfn(q)
//
// 	if len(q.where) != 1 {
// 		t.Errorf("Expected %d where slices, got %d", len(q.where))
// 	}
//
// 	expect := "x > $1 AND y > $2"
// 	if q.where[0].clause != expect {
// 		t.Errorf("Expected %s, got %s", expect, q.where)
// 	}
//
// 	if len(q.where[0].args) != 2 {
// 		t.Errorf("Expected %d args, got %d", 2, len(q.where[0].args))
// 	}
//
// 	if q.where[0].args[0].(int) != 5 || q.where[0].args[1].(int) != 3 {
// 		t.Errorf("Args not set correctly, expected 5 & 3, got: %#v", q.where[0].args)
// 	}
// }
//
// func TestGroupBy(t *testing.T) {
// 	t.Parallel()
//
// 	q := &boil.Query{}
// 	qfn := GroupBy("col1, col2")
//
// 	qfn(q)
//
// 	expect := "col1, col2"
// 	if len(q.groupBy) != 1 && q.groupBy[0] != expect {
// 		t.Errorf("Expected %s, got %s", expect, q.groupBy[0])
// 	}
// }
//
// func TestOrderBy(t *testing.T) {
// 	t.Parallel()
//
// 	q := &boil.Query{}
// 	qfn := OrderBy("col1 desc, col2 asc")
//
// 	qfn(q)
//
// 	expect := "col1 desc, col2 asc"
// 	if len(q.orderBy) != 1 && q.orderBy[0] != expect {
// 		t.Errorf("Expected %s, got %s", expect, q.orderBy[0])
// 	}
// }
//
// func TestHaving(t *testing.T) {
// 	t.Parallel()
//
// 	q := &boil.Query{}
// 	qfn := Having("count(orders.order_id) > 10")
//
// 	qfn(q)
//
// 	expect := "count(orders.order_id) > 10"
// 	if len(q.having) != 1 && q.having[0] != expect {
// 		t.Errorf("Expected %s, got %s", expect, q.having[0])
// 	}
// }
//
// func TestFrom(t *testing.T) {
// 	t.Parallel()
//
// 	q := &boil.Query{}
// 	qfn := From("videos a, orders b")
//
// 	qfn(q)
//
// 	expect := "videos a, orders b"
// 	if q.from != expect {
// 		t.Errorf("Expected %s, got %s", expect, q.from)
// 	}
// }
// }
