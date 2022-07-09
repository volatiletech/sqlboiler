package drivers

import (
	"reflect"
	"testing"
)

func TestConfigMustString(t *testing.T) {
	t.Parallel()

	key := "string"
	tests := []struct {
		Config map[string]interface{}
		Value  string
	}{
		{
			Config: map[string]interface{}{key: "str"},
			Value:  "str",
		},
		{
			Config: map[string]interface{}{key: ""},
			Value:  "",
		},
		{
			Config: map[string]interface{}{key: 5},
			Value:  "",
		},
		{
			Config: map[string]interface{}{},
			Value:  "",
		},
	}

	for i, test := range tests {
		var value string
		var paniced interface{}

		func() {
			defer func() {
				if r := recover(); r != nil {
					paniced = r
				}
			}()
			value = Config(test.Config).MustString(key)
		}()

		if len(test.Value) != 0 {
			if paniced != nil {
				t.Error(i, "wanted a value, but panic'd:", paniced)
			} else if value != test.Value {
				t.Error(i, "want:", test.Value, "got:", value)
			}
		} else {
			if paniced == nil {
				t.Error(i, "expected it to panic")
			}
		}
	}
}

func TestConfigMustInt(t *testing.T) {
	t.Parallel()

	key := "integer"
	tests := []struct {
		Config map[string]interface{}
		Value  int
	}{
		{
			Config: map[string]interface{}{key: 5},
			Value:  5,
		},
		{
			Config: map[string]interface{}{key: 5.0},
			Value:  5,
		},
		{
			Config: map[string]interface{}{key: 0},
			Value:  0,
		},
		{
			Config: map[string]interface{}{},
			Value:  0,
		},
	}

	for i, test := range tests {
		var value int
		var paniced interface{}

		func() {
			defer func() {
				if r := recover(); r != nil {
					paniced = r
				}
			}()
			value = Config(test.Config).MustInt(key)
		}()

		if test.Value != 0 {
			if paniced != nil {
				t.Error(i, "wanted a value, but panic'd:", paniced)
			} else if value != test.Value {
				t.Error(i, "want:", test.Value, "got:", value)
			}
		} else {
			if paniced == nil {
				t.Error(i, "expected it to panic")
			}
		}
	}
}

func TestConfigString(t *testing.T) {
	t.Parallel()

	key := "string"
	tests := []struct {
		Config map[string]interface{}
		Value  string
		Ok     bool
	}{
		{
			Config: map[string]interface{}{key: "str"},
			Value:  "str",
			Ok:     true,
		},
		{
			Config: map[string]interface{}{key: ""},
			Value:  "",
			Ok:     false,
		},
		{
			Config: map[string]interface{}{key: 5},
			Value:  "",
			Ok:     false,
		},
		{
			Config: map[string]interface{}{},
			Value:  "",
			Ok:     false,
		},
	}

	for i, test := range tests {
		value, ok := Config(test.Config).String(key)

		if ok != test.Ok {
			t.Error(i, "ok =", ok)
		}
		if value != test.Value {
			t.Error(i, "want:", test.Value, "got:", value)
		}
	}
}

func TestConfigInt(t *testing.T) {
	t.Parallel()

	key := "integer"
	tests := []struct {
		Config map[string]interface{}
		Value  int
		Ok     bool
	}{
		{
			Config: map[string]interface{}{key: 5},
			Value:  5,
			Ok:     true,
		},
		{
			Config: map[string]interface{}{key: 5.0},
			Value:  5,
			Ok:     true,
		},
		{
			Config: map[string]interface{}{key: 0},
			Value:  0,
			Ok:     false,
		},
		{
			Config: map[string]interface{}{},
			Value:  0,
			Ok:     false,
		},
	}

	for i, test := range tests {
		value, ok := Config(test.Config).Int(key)

		if ok != test.Ok {
			t.Error(i, "ok =", ok)
		}
		if value != test.Value {
			t.Error(i, "want:", test.Value, "got:", value)
		}
	}
}

func TestConfigStringSlice(t *testing.T) {
	t.Parallel()

	key := "slice"
	tests := []struct {
		Config map[string]interface{}
		Value  []string
		Ok     bool
	}{
		{
			Config: map[string]interface{}{key: []string{"str"}},
			Value:  []string{"str"},
			Ok:     true,
		},
		{
			Config: map[string]interface{}{key: []interface{}{"str"}},
			Value:  []string{"str"},
			Ok:     true,
		},
		{
			Config: map[string]interface{}{key: []string{}},
			Value:  nil,
			Ok:     false,
		},
		{
			Config: map[string]interface{}{key: 5},
			Value:  nil,
			Ok:     false,
		},
		{
			Config: map[string]interface{}{},
			Value:  nil,
			Ok:     false,
		},
	}

	for i, test := range tests {
		value, ok := Config(test.Config).StringSlice(key)

		if ok != test.Ok {
			t.Error(i, "ok =", ok)
		}
		if !reflect.DeepEqual(value, test.Value) {
			t.Error(i, "want:", test.Value, "got:", value)
		}
	}
}

func TestTablesFromList(t *testing.T) {
	t.Parallel()

	if TablesFromList(nil) != nil {
		t.Error("expected a shortcut to getting nil back")
	}

	if got := TablesFromList([]string{"a.b", "b", "c.d"}); !reflect.DeepEqual(got, []string{"b"}) {
		t.Error("list was wrong:", got)
	}
}

func TestColumnsFromList(t *testing.T) {
	t.Parallel()

	if ColumnsFromList(nil, "table") != nil {
		t.Error("expected a shortcut to getting nil back")
	}

	if got := ColumnsFromList([]string{"a.b", "b", "c.d", "c.a"}, "c"); !reflect.DeepEqual(got, []string{"d", "a"}) {
		t.Error("list was wrong:", got)
	}
	if got := ColumnsFromList([]string{"a.b", "b", "c.d", "c.a"}, "b"); len(got) != 0 {
		t.Error("list was wrong:", got)
	}
	if got := ColumnsFromList([]string{"*.b", "b", "c.d"}, "c"); !reflect.DeepEqual(got, []string{"b", "d"}) {
		t.Error("list was wrong:", got)
	}
}
