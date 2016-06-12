package sqlboiler

import (
	"sort"
	"testing"
	"text/template"
)

func TestTemplateListSort(t *testing.T) {
	templs := templateList{
		template.New("bob.tpl"),
		template.New("all.tpl"),
		template.New("struct.tpl"),
		template.New("ttt.tpl"),
	}

	expected := []string{"bob.tpl", "all.tpl", "struct.tpl", "ttt.tpl"}

	for i, v := range templs {
		if v.Name() != expected[i] {
			t.Errorf("Order mismatch, expected: %s, got: %s", expected[i], v.Name())
		}
	}

	expected = []string{"struct.tpl", "all.tpl", "bob.tpl", "ttt.tpl"}

	sort.Sort(templs)

	for i, v := range templs {
		if v.Name() != expected[i] {
			t.Errorf("Order mismatch, expected: %s, got: %s", expected[i], v.Name())
		}
	}
}
