package qm

import (
	"reflect"
	"testing"
)

func TestWhereIn(t *testing.T) {
	tests := []struct {
		Input    func() whereInQueryMod
		Expected whereInQueryMod
	}{
		{
			// Test standard ints
			Input: func() whereInQueryMod {
				return WhereIn("id in ?", 1, 2).(whereInQueryMod)
			},
			Expected: whereInQueryMod{
				clause: "id in ?",
				args:   []interface{}{1, 2},
			},
		},
		{
			// Test a slice
			Input: func() whereInQueryMod {
				return WhereIn("id in ?", []int{1, 2}).(whereInQueryMod)
			},
			Expected: whereInQueryMod{
				clause: "id in ?",
				args:   []interface{}{1, 2},
			},
		},
	}

	for _, test := range tests {
		actual := test.Input()
		if !reflect.DeepEqual(actual, test.Expected) {
			t.Fatalf("actual %+v does not match expected %+v", actual, test.Expected)
		}
	}
}
