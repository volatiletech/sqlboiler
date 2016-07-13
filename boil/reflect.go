package boil

import (
	"database/sql"
	"fmt"
	"reflect"

	"github.com/nullbio/sqlboiler/strmangle"
)

// Bind executes the query and inserts the
// result into the passed in object pointer
func (q *Query) Bind(obj interface{}) error {
	typ := reflect.TypeOf(obj)
	kind := typ.Kind()

	if kind != reflect.Ptr {
		return fmt.Errorf("Bind not given a pointer to a slice or struct: %s", typ.String())
	}

	typ = typ.Elem()
	kind = typ.Kind()

	if kind == reflect.Struct {
		row := ExecQueryOne(q)
		err := BindOne(row, q.selectCols, obj)
		if err != nil {
			return fmt.Errorf("Failed to execute Bind query for %s: %s", q.table, err)
		}
	} else if kind == reflect.Slice {
		rows, err := ExecQueryAll(q)
		if err != nil {
			return fmt.Errorf("Failed to execute Bind query for %s: %s", q.table, err)
		}
		err = BindAll(rows, q.selectCols, obj)
		if err != nil {
			return fmt.Errorf("Failed to Bind results to object provided for %s: %s", q.table, err)
		}
	} else {
		return fmt.Errorf("Bind given a pointer to a non-slice or non-struct: %s", typ.String())
	}

	return nil
}

// BindOne inserts the returned row columns into the
// passed in object pointer
func BindOne(row *sql.Row, selectCols []string, obj interface{}) error {
	kind := reflect.ValueOf(obj).Kind()
	if kind != reflect.Ptr {
		return fmt.Errorf("BindOne given a non-pointer type")
	}

	pointers := GetStructPointers(obj, selectCols...)
	if err := row.Scan(pointers...); err != nil {
		return fmt.Errorf("Unable to scan into pointers: %s", err)
	}

	return nil
}

// BindAll inserts the returned rows columns into the
// passed in slice of object pointers
func BindAll(rows *sql.Rows, selectCols []string, obj interface{}) error {
	ptrSlice := reflect.ValueOf(obj)
	typ := ptrSlice.Type()
	ptrSlice = ptrSlice.Elem()
	kind := typ.Kind()

	var structTyp reflect.Type

	for i, exp := range []reflect.Kind{reflect.Ptr, reflect.Slice, reflect.Ptr, reflect.Struct} {
		if i != 0 {
			typ = typ.Elem()
			kind = typ.Kind()
		}

		if kind != exp {
			return fmt.Errorf("[%d] BindAll object type should be *[]*Type but was: %s", i, ptrSlice.Type().String())
		}

		if kind == reflect.Struct {
			structTyp = typ
		}
	}

	for rows.Next() {
		newStruct := reflect.New(structTyp)
		pointers := GetStructPointers(newStruct.Interface(), selectCols...)
		if err := rows.Scan(pointers...); err != nil {
			return fmt.Errorf("Unable to scan into pointers: %s", err)
		}

		ptrSlice.Set(reflect.Append(ptrSlice, newStruct))
	}

	return nil
}

// GetStructValues returns the values (as interface) of the matching columns in obj
func GetStructValues(obj interface{}, columns ...string) []interface{} {
	ret := make([]interface{}, len(columns))
	val := reflect.Indirect(reflect.ValueOf(obj))

	for i, c := range columns {
		field := val.FieldByName(strmangle.TitleCase(c))
		ret[i] = field.Interface()
	}

	return ret
}

// GetStructPointers returns a slice of pointers to the matching columns in obj
func GetStructPointers(obj interface{}, columns ...string) []interface{} {
	val := reflect.ValueOf(obj).Elem()
	var ret []interface{}

	if len(columns) == 0 {
		fieldsLen := val.NumField()
		ret = make([]interface{}, fieldsLen)
		for i := 0; i < fieldsLen; i++ {
			ret[i] = val.Field(i).Addr().Interface()
		}
		return ret
	}

	ret = make([]interface{}, len(columns))

	for i, c := range columns {
		field := val.FieldByName(strmangle.TitleCase(c))
		if !field.IsValid() {
			panic(fmt.Sprintf("Could not find field on struct %T for field %s", obj, strmangle.TitleCase(c)))
		}

		field = field.Addr()
		ret[i] = field.Interface()
	}

	return ret
}
