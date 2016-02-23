package cmds

import "testing"

func testVerifyStructArgs(t *testing.T) {
	err := verifyStructArgs([]string{})
	if err == nil {
		t.Error("Expected an error")
	}

	err = verifyStructArgs([]string{"hello"})
	if err != nil {
		t.Errorf("Expected error nil, got: %s", err)
	}
}
