package boil

import (
	"database/sql"
	"testing"
)

func TestGetSetDB(t *testing.T) {
	t.Parallel()

	SetContextDB(&sql.DB{})

	if GetContextDB() == nil {
		t.Errorf("Expected GetDB to return a database handle, got nil")
	}
}
