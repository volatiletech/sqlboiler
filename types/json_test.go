package types

import (
	"bytes"
	"strings"
	"testing"
)

func TestJSONString(t *testing.T) {
	t.Parallel()

	j := JSON("hello")
	if j.String() != "hello" {
		t.Errorf("Expected %q, got %s", "hello", j.String())
	}
}

func TestJSONUnmarshal(t *testing.T) {
	t.Parallel()

	type JSONTest struct {
		Name string
		Age  int
	}
	var jt JSONTest

	j := JSON(`{"Name":"hi","Age":15}`)
	err := j.Unmarshal(&jt)
	if err != nil {
		t.Error(err)
	}

	if jt.Name != "hi" {
		t.Errorf("Expected %q, got %s", "hi", jt.Name)
	}
	if jt.Age != 15 {
		t.Errorf("Expected %v, got %v", 15, jt.Age)
	}
}

func TestJSONMarshal(t *testing.T) {
	t.Parallel()

	type JSONTest struct {
		Name string
		Age  int
	}
	jt := JSONTest{
		Name: "hi",
		Age:  15,
	}

	var j JSON
	err := j.Marshal(jt)
	if err != nil {
		t.Error(err)
	}

	if j.String() != `{"Name":"hi","Age":15}` {
		t.Errorf("expected %s, got %s", `{"Name":"hi","Age":15}`, j.String())
	}
}

func TestJSONUnmarshalJSON(t *testing.T) {
	t.Parallel()

	j := JSON(nil)

	err := j.UnmarshalJSON(JSON(`"hi"`))
	if err != nil {
		t.Error(err)
	}

	if j.String() != `"hi"` {
		t.Errorf("Expected %q, got %s", "hi", j.String())
	}
}

func TestJSONMarshalJSON_Null(t *testing.T) {
	t.Parallel()

	var j JSON
	res, err := j.MarshalJSON()
	if err != nil {
		t.Error(err)
	}

	if !bytes.Equal(res, []byte(`null`)) {
		t.Errorf("Expected %q, got %v", `null`, res)
	}
}

func TestJSONMarshalJSON(t *testing.T) {
	t.Parallel()

	j := JSON(`"hi"`)
	res, err := j.MarshalJSON()
	if err != nil {
		t.Error(err)
	}

	if !bytes.Equal(res, []byte(`"hi"`)) {
		t.Errorf("Expected %q, got %v", `"hi"`, res)
	}
}

func TestJSONValue(t *testing.T) {
	t.Parallel()

	j := JSON(`{"Name":"hi","Age":15}`)
	v, err := j.Value()
	if err != nil {
		t.Error(err)
	}

	if !bytes.Equal(j, v.([]byte)) {
		t.Errorf("byte mismatch, %v %v", j, v)
	}
}

func TestJSONScan(t *testing.T) {
	t.Parallel()

	j := JSON{}

	err := j.Scan(`"hello"`)
	if err != nil {
		t.Error(err)
	}

	if !bytes.Equal(j, []byte(`"hello"`)) {
		t.Errorf("bad []byte: %#v ≠ %#v\n", j, string([]byte(`"hello"`)))
	}
}

func BenchmarkJSON_Scan(b *testing.B) {
	data := `"` + strings.Repeat("A", 1024) + `"`
	for i := 0; i < b.N; i++ {
		var j JSON
		err := j.Scan(data)
		if err != nil {
			b.Error(err)
		}
	}
}
