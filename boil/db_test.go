package boil

import (
	"database/sql"
	"testing"
)

func TestGetSetDB(t *testing.T) {
	t.Parallel()

	SetDB(&sql.DB{})

	if GetDB() == nil {
		t.Errorf("Expected GetDB to return a database handle, got nil")
	}
}
