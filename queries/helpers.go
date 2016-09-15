package queries

import (
	"fmt"
	"reflect"

	"github.com/vattle/sqlboiler/strmangle"
)

// NonZeroDefaultSet returns the fields included in the
// defaults slice that are non zero values
func NonZeroDefaultSet(defaults []string, obj interface{}) []string {
	c := make([]string, 0, len(defaults))

	val := reflect.Indirect(reflect.ValueOf(obj))

	for _, d := range defaults {
		fieldName := strmangle.TitleCase(d)
		field := val.FieldByName(fieldName)
		if !field.IsValid() {
			panic(fmt.Sprintf("Could not find field name %s in type %T", fieldName, obj))
		}

		zero := reflect.Zero(field.Type())
		if !reflect.DeepEqual(zero.Interface(), field.Interface()) {
			c = append(c, d)
		}
	}

	return c
}
