package boil

import (
	"fmt"
	"reflect"

	"github.com/vattle/sqlboiler/strmangle"
)

// NonZeroDefaultSet returns the fields included in the
// defaults slice that are non zero values
func NonZeroDefaultSet(defaults []string, titleCases map[string]string, obj interface{}) []string {
	c := make([]string, 0, len(defaults))

	val := reflect.Indirect(reflect.ValueOf(obj))

	for _, d := range defaults {
		var fieldName string
		if titleCases == nil {
			fieldName = strmangle.TitleCase(d)
		} else {
			fieldName = titleCases[d]
		}
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
