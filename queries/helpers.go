package queries

import (
	"fmt"
	"reflect"
)

// NonZeroDefaultSet returns the fields included in the
// defaults slice that are non zero values
func NonZeroDefaultSet(defaults []string, obj interface{}) []string {
	c := make([]string, 0, len(defaults))

	val := reflect.Indirect(reflect.ValueOf(obj))
	typ := val.Type()
	nf := typ.NumField()

	for _, def := range defaults {
		found := false
		for i := 0; i < nf; i++ {
			field := typ.Field(i)
			name, _ := getBoilTag(field)

			if name != def {
				continue
			}

			found = true
			fieldVal := val.Field(i)

			zero := reflect.Zero(fieldVal.Type())
			if !reflect.DeepEqual(zero.Interface(), fieldVal.Interface()) {
				c = append(c, def)
			}
			break
		}

		if !found {
			panic(fmt.Sprintf("could not find field name %s in type %T", def, obj))
		}
	}

	return c
}
