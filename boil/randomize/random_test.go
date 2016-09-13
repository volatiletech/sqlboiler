package randomize

import "testing"

func TestStableDBName(t *testing.T) {
	t.Parallel()

	db := "awesomedb"

	one, two := StableDBName(db), StableDBName(db)

	if len(one) != 40 {
		t.Error("want 40 characters:", len(one), one)
	}

	if one != two {
		t.Error("it should always produce the same value")
	}
}
