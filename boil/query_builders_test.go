package boil

import "testing"

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
			In:  Query{innerJoins: []join{{on: `inner join a on stuff = there`}}},
			Out: map[string]string{"a": "a"},
		},
		{
			In:  Query{innerJoins: []join{{on: `outer join "a" on stuff = there`}}},
			Out: map[string]string{"a": "a"},
		},
		{
			In:  Query{innerJoins: []join{{on: `natural join a as b on stuff = there`}}},
			Out: map[string]string{"b": "a"},
		},
		{
			In:  Query{innerJoins: []join{{on: `right outer join "a" as "b" on stuff = there`}}},
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
