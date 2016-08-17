package boil

import (
	"bytes"
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
		{&Query{
			selectCols: []string{"a.happy", "r.fun", "q"},
			from:       []string{"happiness as a"},
			joins:      []join{{clause: "rainbows r on a.id = r.happy_id"}},
		}, nil},
		{&Query{
			from:  []string{"happiness as a"},
			joins: []join{{clause: "rainbows r on a.id = r.happy_id"}},
		}, nil},
		{&Query{
			from: []string{"videos"},
			joins: []join{{
				clause: "(select id from users where deleted = ?) u on u.id = videos.user_id",
				args:   []interface{}{true},
			}},
			where: []where{{clause: "videos.deleted = ?", args: []interface{}{false}}},
		}, []interface{}{true, false}},
		{&Query{
			from:    []string{"a"},
			groupBy: []string{"id", "name"},
			where: []where{
				{clause: "a=? or b=?", args: []interface{}{1, 2}},
				{clause: "c=?", args: []interface{}{3}},
			},
			having: []having{
				{clause: "id <> ?", args: []interface{}{1}},
				{clause: "length(name, ?) > ?", args: []interface{}{"utf8", 5}},
			},
		}, []interface{}{1, 2, 3, 1, "utf8", 5}},
		{&Query{
			delete: true,
			from:   []string{"thing happy", `upset as "sad"`, "fun", "thing as stuff", `"angry" as mad`},
			where: []where{
				{clause: "a=?", args: []interface{}{}},
				{clause: "b=?", args: []interface{}{}},
				{clause: "c=?", args: []interface{}{}},
			},
		}, nil},
		{&Query{
			delete: true,
			from:   []string{"thing happy", `upset as "sad"`, "fun", "thing as stuff", `"angry" as mad`},
			where: []where{
				{clause: "(id=? and thing=?) or stuff=?", args: []interface{}{}},
			},
			limit: 5,
		}, nil},
		{&Query{
			from: []string{"thing happy", `"fun"`, `stuff`},
			update: map[string]interface{}{
				"col1":       1,
				`"col2"`:     2,
				`"fun".col3`: 3,
			},
			where: []where{
				{clause: "aa=? or bb=? or cc=?", orSeparator: true, args: []interface{}{4, 5, 6}},
				{clause: "dd=? or ee=? or ff=? and gg=?", args: []interface{}{7, 8, 9, 10}},
			},
			limit: 5,
		}, []interface{}{2, 3, 1, 4, 5, 6, 7, 8, 9, 10}},
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

func TestIdentifierMapping(t *testing.T) {
	t.Parallel()

	tests := []struct {
		In  Query
		Out map[string]string
	}{
		{
			In:  Query{from: []string{`a`}},
			Out: map[string]string{"a": "a"},
		},
		{
			In:  Query{from: []string{`"a"`, `b`}},
			Out: map[string]string{"a": "a", "b": "b"},
		},
		{
			In:  Query{from: []string{`a as b`}},
			Out: map[string]string{"b": "a"},
		},
		{
			In:  Query{from: []string{`a AS "b"`, `"c" as d`}},
			Out: map[string]string{"b": "a", "d": "c"},
		},
		{
			In:  Query{joins: []join{{kind: JoinInner, clause: `a on stuff = there`}}},
			Out: map[string]string{"a": "a"},
		},
		{
			In:  Query{joins: []join{{kind: JoinNatural, clause: `"a" on stuff = there`}}},
			Out: map[string]string{"a": "a"},
		},
		{
			In:  Query{joins: []join{{kind: JoinNatural, clause: `a as b on stuff = there`}}},
			Out: map[string]string{"b": "a"},
		},
		{
			In:  Query{joins: []join{{kind: JoinOuterRight, clause: `"a" as "b" on stuff = there`}}},
			Out: map[string]string{"b": "a"},
		},
	}

	for i, test := range tests {
		m := identifierMapping(&test.In)

		for k, v := range test.Out {
			val, ok := m[k]
			if !ok {
				t.Errorf("%d) want: %s = %s, but was missing", i, k, v)
			}
			if val != v {
				t.Errorf("%d) want: %s = %s, got: %s", i, k, v, val)
			}
		}
	}
}

func TestWriteStars(t *testing.T) {
	t.Parallel()

	tests := []struct {
		In  Query
		Out []string
	}{
		{
			In:  Query{from: []string{`a`}},
			Out: []string{`"a".*`},
		},
		{
			In:  Query{from: []string{`a as b`}},
			Out: []string{`"b".*`},
		},
		{
			In:  Query{from: []string{`a as b`, `c`}},
			Out: []string{`"b".*`, `"c".*`},
		},
		{
			In:  Query{from: []string{`a as b`, `c as d`}},
			Out: []string{`"b".*`, `"d".*`},
		},
	}

	for i, test := range tests {
		selects := writeStars(&test.In)
		if !reflect.DeepEqual(selects, test.Out) {
			t.Errorf("writeStar test fail %d\nwant: %v\ngot:  %v", i, test.Out, selects)
		}
	}
}

func TestWhereClause(t *testing.T) {
	t.Parallel()

	tests := []struct {
		q      Query
		expect string
	}{
		// Or("a=?")
		{
			q: Query{
				where: []where{where{clause: "a=?", orSeparator: true}},
			},
			expect: " WHERE (a=$1)",
		},
		// Where("a=?")
		{
			q: Query{
				where: []where{where{clause: "a=?"}},
			},
			expect: " WHERE (a=$1)",
		},
		// Where("(a=?)")
		{
			q: Query{
				where: []where{where{clause: "(a=?)"}},
			},
			expect: " WHERE ((a=$1))",
		},
		// Where("((a=? OR b=?))")
		{
			q: Query{
				where: []where{where{clause: "((a=? OR b=?))"}},
			},
			expect: " WHERE (((a=$1 OR b=$2)))",
		},
		// Where("(a=?)", Or("(b=?)")
		{
			q: Query{
				where: []where{
					where{clause: "(a=?)", orSeparator: true},
					where{clause: "(b=?)"},
				},
			},
			expect: " WHERE ((a=$1)) OR ((b=$2))",
		},
		// Where("a=? OR b=?")
		{
			q: Query{
				where: []where{where{clause: "a=? OR b=?"}},
			},
			expect: " WHERE (a=$1 OR b=$2)",
		},
		// Where("a=?"), Where("b=?")
		{
			q: Query{
				where: []where{where{clause: "a=?"}, where{clause: "b=?"}},
			},
			expect: " WHERE (a=$1) AND (b=$2)",
		},
		// Where("(a=? AND b=?) OR c=?")
		{
			q: Query{
				where: []where{where{clause: "(a=? AND b=?) OR c=?"}},
			},
			expect: " WHERE ((a=$1 AND b=$2) OR c=$3)",
		},
		// Where("a=? OR b=?"), Where("c=? OR d=? OR e=?")
		{
			q: Query{
				where: []where{
					where{clause: "(a=? OR b=?)"},
					where{clause: "(c=? OR d=? OR e=?)"},
				},
			},
			expect: " WHERE ((a=$1 OR b=$2)) AND ((c=$3 OR d=$4 OR e=$5))",
		},
		// Where("(a=? AND b=?) OR (c=? AND d=? AND e=?) OR f=? OR f=?")
		{
			q: Query{
				where: []where{
					where{clause: "(a=? AND b=?) OR (c=? AND d=? AND e=?) OR f=? OR g=?"},
				},
			},
			expect: " WHERE ((a=$1 AND b=$2) OR (c=$3 AND d=$4 AND e=$5) OR f=$6 OR g=$7)",
		},
		// Where("(a=? AND b=?) OR (c=? AND d=? OR e=?) OR f=? OR g=?")
		{
			q: Query{
				where: []where{
					where{clause: "(a=? AND b=?) OR (c=? AND d=? OR e=?) OR f=? OR g=?"},
				},
			},
			expect: " WHERE ((a=$1 AND b=$2) OR (c=$3 AND d=$4 OR e=$5) OR f=$6 OR g=$7)",
		},
		// Where("a=? or b=?"), Or("c=? and d=?"), Or("e=? or f=?")
		{
			q: Query{
				where: []where{
					where{clause: "a=? or b=?", orSeparator: true},
					where{clause: "c=? and d=?", orSeparator: true},
					where{clause: "e=? or f=?", orSeparator: true},
				},
			},
			expect: " WHERE (a=$1 or b=$2) OR (c=$3 and d=$4) OR (e=$5 or f=$6)",
		},
		// Where("a=? or b=?"), Or("c=? and d=?"), Or("e=? or f=?")
		{
			q: Query{
				where: []where{
					where{clause: "a=? or b=?"},
					where{clause: "c=? and d=?"},
					where{clause: "e=? or f=?"},
				},
			},
			expect: " WHERE (a=$1 or b=$2) AND (c=$3 and d=$4) AND (e=$5 or f=$6)",
		},
	}

	for i, test := range tests {
		result, _ := whereClause(&test.q, 1)
		if result != test.expect {
			t.Errorf("%d) Mismatch between expect and result:\n%s\n%s\n", i, test.expect, result)
		}
	}
}

func TestInClause(t *testing.T) {
	t.Parallel()

	tests := []struct {
		q      Query
		expect string
	}{
		// Or("a=?")
		{
			q: Query{
				in: []in{{clause: "a in ?", args: []interface{}{1}, orSeparator: true}},
			},
			expect: " WHERE a IN ($1)",
		},
		{
			q: Query{
				in: []in{{clause: "a in ?", args: []interface{}{1, 2, 3}}},
			},
			expect: " WHERE a IN ($1,$2,$3)",
		},
		// // Where("a=?")
		// {
		// 	q: Query{
		// 		where: []where{where{clause: "a=?"}},
		// 	},
		// 	expect: " WHERE (a=$1)",
		// },
		// // Where("(a=?)")
		// {
		// 	q: Query{
		// 		where: []where{where{clause: "(a=?)"}},
		// 	},
		// 	expect: " WHERE ((a=$1))",
		// },
		// // Where("((a=? OR b=?))")
		// {
		// 	q: Query{
		// 		where: []where{where{clause: "((a=? OR b=?))"}},
		// 	},
		// 	expect: " WHERE (((a=$1 OR b=$2)))",
		// },
		// // Where("(a=?)", Or("(b=?)")
		// {
		// 	q: Query{
		// 		where: []where{
		// 			where{clause: "(a=?)", orSeparator: true},
		// 			where{clause: "(b=?)"},
		// 		},
		// 	},
		// 	expect: " WHERE ((a=$1)) OR ((b=$2))",
		// },
		// // Where("a=? OR b=?")
		// {
		// 	q: Query{
		// 		where: []where{where{clause: "a=? OR b=?"}},
		// 	},
		// 	expect: " WHERE (a=$1 OR b=$2)",
		// },
		// // Where("a=?"), Where("b=?")
		// {
		// 	q: Query{
		// 		where: []where{where{clause: "a=?"}, where{clause: "b=?"}},
		// 	},
		// 	expect: " WHERE (a=$1) AND (b=$2)",
		// },
		// // Where("(a=? AND b=?) OR c=?")
		// {
		// 	q: Query{
		// 		where: []where{where{clause: "(a=? AND b=?) OR c=?"}},
		// 	},
		// 	expect: " WHERE ((a=$1 AND b=$2) OR c=$3)",
		// },
		// // Where("a=? OR b=?"), Where("c=? OR d=? OR e=?")
		// {
		// 	q: Query{
		// 		where: []where{
		// 			where{clause: "(a=? OR b=?)"},
		// 			where{clause: "(c=? OR d=? OR e=?)"},
		// 		},
		// 	},
		// 	expect: " WHERE ((a=$1 OR b=$2)) AND ((c=$3 OR d=$4 OR e=$5))",
		// },
		// // Where("(a=? AND b=?) OR (c=? AND d=? AND e=?) OR f=? OR f=?")
		// {
		// 	q: Query{
		// 		where: []where{
		// 			where{clause: "(a=? AND b=?) OR (c=? AND d=? AND e=?) OR f=? OR g=?"},
		// 		},
		// 	},
		// 	expect: " WHERE ((a=$1 AND b=$2) OR (c=$3 AND d=$4 AND e=$5) OR f=$6 OR g=$7)",
		// },
		// // Where("(a=? AND b=?) OR (c=? AND d=? OR e=?) OR f=? OR g=?")
		// {
		// 	q: Query{
		// 		where: []where{
		// 			where{clause: "(a=? AND b=?) OR (c=? AND d=? OR e=?) OR f=? OR g=?"},
		// 		},
		// 	},
		// 	expect: " WHERE ((a=$1 AND b=$2) OR (c=$3 AND d=$4 OR e=$5) OR f=$6 OR g=$7)",
		// },
		// // Where("a=? or b=?"), Or("c=? and d=?"), Or("e=? or f=?")
		// {
		// 	q: Query{
		// 		where: []where{
		// 			where{clause: "a=? or b=?", orSeparator: true},
		// 			where{clause: "c=? and d=?", orSeparator: true},
		// 			where{clause: "e=? or f=?", orSeparator: true},
		// 		},
		// 	},
		// 	expect: " WHERE (a=$1 or b=$2) OR (c=$3 and d=$4) OR (e=$5 or f=$6)",
		// },
		// // Where("a=? or b=?"), Or("c=? and d=?"), Or("e=? or f=?")
		// {
		// 	q: Query{
		// 		where: []where{
		// 			where{clause: "a=? or b=?"},
		// 			where{clause: "c=? and d=?"},
		// 			where{clause: "e=? or f=?"},
		// 		},
		// 	},
		// 	expect: " WHERE (a=$1 or b=$2) AND (c=$3 and d=$4) AND (e=$5 or f=$6)",
		// },
	}

	for i, test := range tests {
		result, _ := inClause(&test.q, 1)
		if result != test.expect {
			t.Errorf("%d) Mismatch between expect and result:\n%s\n%s\n", i, test.expect, result)
		}
	}
}

func TestConvertQuestionMarks(t *testing.T) {
	t.Parallel()

	tests := []struct {
		clause string
		start  int
		expect string
	}{
		{clause: "hello friend", start: 1, expect: "hello friend"},
		{clause: "thing=?", start: 2, expect: "thing=$2"},
		{clause: "thing=? and stuff=? and happy=?", start: 2, expect: "thing=$2 and stuff=$3 and happy=$4"},
		{clause: `thing \? stuff`, start: 2, expect: `thing ? stuff`},
		{clause: `thing \? stuff and happy \? fun`, start: 2, expect: `thing ? stuff and happy ? fun`},
		{
			clause: `thing \? stuff ? happy \? and mad ? fun \? \? \?`,
			start:  2,
			expect: `thing ? stuff $2 happy ? and mad $3 fun ? ? ?`,
		},
		{
			clause: `thing ? stuff ? happy \? fun \? ? ?`,
			start:  1,
			expect: `thing $1 stuff $2 happy ? fun ? $3 $4`,
		},
		{clause: `?`, start: 1, expect: `$1`},
		{clause: `???`, start: 1, expect: `$1$2$3`},
		{clause: `\?`, start: 1, expect: `?`},
		{clause: `\?\?\?`, start: 1, expect: `???`},
		{clause: `\??\??\??`, start: 1, expect: `?$1?$2?$3`},
		{clause: `?\??\??\?`, start: 1, expect: `$1?$2?$3?`},
	}

	for i, test := range tests {
		res := convertQuestionMarks(test.clause, test.start)
		if res != test.expect {
			t.Errorf("%d) Mismatch between expect and result:\n%s\n%s\n", i, test.expect, res)
		}
	}
}

func TestConvertInQuestionMarks(t *testing.T) {
	t.Parallel()

	tests := []struct {
		clause string
		start  int
		expect string
	}{
		{clause: "hello friend", start: 1, expect: "hello friend"},
		{clause: "thing=?", start: 2, expect: "thing=$2"},
		{clause: "thing=? and stuff=? and happy=?", start: 2, expect: "thing=$2 and stuff=$3 and happy=$4"},
		{clause: `thing \? stuff`, start: 2, expect: `thing ? stuff`},
		{clause: `thing \? stuff and happy \? fun`, start: 2, expect: `thing ? stuff and happy ? fun`},
		{
			clause: `thing \? stuff ? happy \? and mad ? fun \? \? \?`,
			start:  2,
			expect: `thing ? stuff $2 happy ? and mad $3 fun ? ? ?`,
		},
		{
			clause: `thing ? stuff ? happy \? fun \? ? ?`,
			start:  1,
			expect: `thing $1 stuff $2 happy ? fun ? $3 $4`,
		},
		{clause: `?`, start: 1, expect: `$1`},
		{clause: `???`, start: 1, expect: `$1$2$3`},
		{clause: `\?`, start: 1, expect: `?`},
		{clause: `\?\?\?`, start: 1, expect: `???`},
		{clause: `\??\??\??`, start: 1, expect: `?$1?$2?$3`},
		{clause: `?\??\??\?`, start: 1, expect: `$1?$2?$3?`},
	}

	for i, test := range tests {
		res := convertQuestionMarks(test.clause, test.start)
		if res != test.expect {
			t.Errorf("%d) Mismatch between expect and result:\n%s\n%s\n", i, test.expect, res)
		}
	}
}

func TestWriteAsStatements(t *testing.T) {
	t.Parallel()

	query := Query{
		selectCols: []string{
			`a`,
			`a.fun`,
			`"b"."fun"`,
			`"b".fun`,
			`b."fun"`,
			`a.clown.run`,
			`COUNT(a)`,
		},
	}

	expect := []string{
		`"a"`,
		`"a"."fun" as "a.fun"`,
		`"b"."fun" as "b.fun"`,
		`"b"."fun" as "b.fun"`,
		`"b"."fun" as "b.fun"`,
		`"a"."clown"."run" as "a.clown.run"`,
		`COUNT(a)`,
	}

	gots := writeAsStatements(&query)

	for i, got := range gots {
		if expect[i] != got {
			t.Errorf(`%d) want: %s, got: %s`, i, expect[i], got)
		}
	}
}
