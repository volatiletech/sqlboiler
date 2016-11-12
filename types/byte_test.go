package types

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestByteString(t *testing.T) {
	t.Parallel()

	b := Byte('b')
	if b.String() != "b" {
		t.Errorf("Expected %q, got %s", "b", b.String())
	}
}

func TestByteUnmarshal(t *testing.T) {
	t.Parallel()

	var b Byte
	err := json.Unmarshal([]byte(`"b"`), &b)
	if err != nil {
		t.Error(err)
	}

	if b != 'b' {
		t.Errorf("Expected %q, got %s", "b", b)
	}
}

func TestByteMarshal(t *testing.T) {
	t.Parallel()

	b := Byte('b')
	res, err := json.Marshal(&b)
	if err != nil {
		t.Error(err)
	}

	if !bytes.Equal(res, []byte(`"b"`)) {
		t.Errorf("expected %s, got %s", `"b"`, b.String())
	}
}

func TestByteValue(t *testing.T) {
	t.Parallel()

	b := Byte('b')
	v, err := b.Value()
	if err != nil {
		t.Error(err)
	}

	if !bytes.Equal([]byte{byte(b)}, v.([]byte)) {
		t.Errorf("byte mismatch, %v %v", b, v)
	}
}

func TestByteScan(t *testing.T) {
	t.Parallel()

	var b Byte

	s := "b"
	err := b.Scan(s)
	if err != nil {
		t.Error(err)
	}

	if !bytes.Equal([]byte{byte(b)}, []byte{'b'}) {
		t.Errorf("bad []byte: %#v ≠ %#v\n", b, "b")
	}
}
