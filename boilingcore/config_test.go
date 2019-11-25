package boilingcore

import (
	"testing"

	"github.com/razor-1/sqlboiler/drivers"
)

func TestConfig_OutputDirDepth(t *testing.T) {
	t.Parallel()

	tests := []struct {
		Output string
		Depth  int
	}{
		{
			Output: ".",
			Depth:  0,
		},
		{
			Output: "./",
			Depth:  0,
		},
		{
			Output: "foo",
			Depth:  1,
		},
		{
			Output: "foo/bar",
			Depth:  2,
		},
	}

	for i, test := range tests {
		cfg := Config{
			OutFolder: test.Output,
		}
		if want, got := test.Depth, cfg.OutputDirDepth(); got != want {
			t.Errorf("%d) wrong depth, want: %d, got: %d", i, want, got)
		}
	}
}

func TestConvertAliases(t *testing.T) {
	t.Parallel()

	var intf interface{} = map[string]interface{}{
		"tables": map[string]interface{}{
			"table_name": map[string]interface{}{
				"up_plural":     "a",
				"up_singular":   "b",
				"down_plural":   "c",
				"down_singular": "d",

				"columns": map[string]interface{}{
					"a": "b",
				},
				"relationships": map[string]interface{}{
					"ib_fk_1": map[string]interface{}{
						"local":   "a",
						"foreign": "b",
					},
				},
			},
		},
	}

	aliases := ConvertAliases(intf)

	if len(aliases.Tables) != 1 {
		t.Fatalf("should have one table alias: %#v", aliases.Tables)
	}

	table := aliases.Tables["table_name"]
	if table.UpPlural != "a" {
		t.Error("value was wrong:", table.UpPlural)
	}
	if table.UpSingular != "b" {
		t.Error("value was wrong:", table.UpSingular)
	}
	if table.DownPlural != "c" {
		t.Error("value was wrong:", table.DownPlural)
	}
	if table.DownSingular != "d" {
		t.Error("value was wrong:", table.DownSingular)
	}

	if len(table.Columns) != 1 {
		t.Error("should have one column")
	}

	if table.Columns["a"] != "b" {
		t.Error("column alias was wrong")
	}

	if len(aliases.Tables) != 1 {
		t.Fatal("should have one relationship alias")
	}

	rel := table.Relationships["ib_fk_1"]
	if rel.Local != "a" {
		t.Error("value was wrong:", rel.Local)
	}
	if rel.Foreign != "b" {
		t.Error("value was wrong:", rel.Foreign)
	}
}

func TestConvertAliasesAltSyntax(t *testing.T) {
	t.Parallel()

	var intf interface{} = map[string]interface{}{
		"tables": []interface{}{
			map[string]interface{}{
				"name":          "table_name",
				"up_plural":     "a",
				"up_singular":   "b",
				"down_plural":   "c",
				"down_singular": "d",

				"columns": []interface{}{
					map[string]interface{}{
						"name":  "a",
						"alias": "b",
					},
				},
				"relationships": []interface{}{
					map[string]interface{}{
						"name":    "ib_fk_1",
						"local":   "a",
						"foreign": "b",
					},
				},
			},
		},
	}

	aliases := ConvertAliases(intf)

	if len(aliases.Tables) != 1 {
		t.Fatalf("should have one table alias: %#v", aliases.Tables)
	}

	table := aliases.Tables["table_name"]
	if table.UpPlural != "a" {
		t.Error("value was wrong:", table.UpPlural)
	}
	if table.UpSingular != "b" {
		t.Error("value was wrong:", table.UpSingular)
	}
	if table.DownPlural != "c" {
		t.Error("value was wrong:", table.DownPlural)
	}
	if table.DownSingular != "d" {
		t.Error("value was wrong:", table.DownSingular)
	}

	if len(table.Columns) != 1 {
		t.Error("should have one column")
	}

	if table.Columns["a"] != "b" {
		t.Error("column alias was wrong")
	}

	if len(aliases.Tables) != 1 {
		t.Fatal("should have one relationship alias")
	}

	rel := table.Relationships["ib_fk_1"]
	if rel.Local != "a" {
		t.Error("value was wrong:", rel.Local)
	}
	if rel.Foreign != "b" {
		t.Error("value was wrong:", rel.Foreign)
	}
}

func TestConvertTypeReplace(t *testing.T) {
	t.Parallel()

	fullColumn := map[string]interface{}{
		"name":           "a",
		"type":           "b",
		"db_type":        "c",
		"udt_name":       "d",
		"full_db_type":   "e",
		"arr_type":       "f",
		"auto_generated": true,
		"nullable":       true,
	}

	var intf interface{} = []interface{}{
		map[string]interface{}{
			"match":   fullColumn,
			"replace": fullColumn,
			"imports": map[string]interface{}{
				"standard": []interface{}{
					"abc",
				},
				"third_party": []interface{}{
					"github.com/abc",
				},
			},
		},
	}

	typeReplace := ConvertTypeReplace(intf)
	if len(typeReplace) != 1 {
		t.Error("should have one entry")
	}

	checkColumn := func(t *testing.T, c drivers.Column) {
		t.Helper()
		if c.Name != "a" {
			t.Error("value was wrong:", c.Name)
		}
		if c.Type != "b" {
			t.Error("value was wrong:", c.Type)
		}
		if c.DBType != "c" {
			t.Error("value was wrong:", c.DBType)
		}
		if c.UDTName != "d" {
			t.Error("value was wrong:", c.UDTName)
		}
		if c.FullDBType != "e" {
			t.Error("value was wrong:", c.FullDBType)
		}
		if *c.ArrType != "f" {
			t.Error("value was wrong:", c.ArrType)
		}
		if c.AutoGenerated != true {
			t.Error("value was wrong:", c.AutoGenerated)
		}
		if c.Nullable != true {
			t.Error("value was wrong:", c.Nullable)
		}
	}

	r := typeReplace[0]
	checkColumn(t, r.Match)
	checkColumn(t, r.Replace)

	if got := r.Imports.Standard[0]; got != "abc" {
		t.Error("standard import wrong:", got)
	}
	if got := r.Imports.ThirdParty[0]; got != "github.com/abc" {
		t.Error("standard import wrong:", got)
	}
}
