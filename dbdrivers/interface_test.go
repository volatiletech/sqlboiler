package dbdrivers

import "testing"

func TestTables(t *testing.T) {

}

func TestSetIsJoinTable(t *testing.T) {
	tests := []struct {
		Pkey   []string
		Fkey   []string
		Should bool
	}{
		{Pkey: []string{"one"}, Fkey: []string{"one"}, Should: true},
	}

	for i, test := range tests {
		var table Table

		table.PKey = &PrimaryKey{Columns: test.Pkey}
		for _, k := range test.Fkey {
			table.FKeys = append(table.FKeys, ForeignKey{Column: k})
		}

		setIsJoinTable(&table)
		if is := table.IsJoinTable; is != test.Should {
			t.Errorf("%d) want: %t, got: %t\nTest: %#v", i, test.Should, is, test)
		}
	}
}
