package boil

import (
	"errors"
	"testing"
)

func TestErrors(t *testing.T) {
	t.Parallel()

	err := errors.New("test error")
	if IsBoilErr(err) == true {
		t.Errorf("Expected false")
	}

	err = WrapErr(errors.New("test error"))
	if err.Error() != "test error" {
		t.Errorf(`Expected "test error", got %v`, err.Error())
	}

	if IsBoilErr(err) != true {
		t.Errorf("Expected true")
	}
}
