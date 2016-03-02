package cmds

import "testing"

func TestBuildCommandList(t *testing.T) {
	list := buildCommandList()

	skips := []string{"struct", "boil"}

	for _, item := range list {
		for _, skipItem := range skips {
			if item == skipItem {
				t.Errorf("Did not expect to find: %s %#v", item, list)
			}
		}
	}

CommandNameLoop:
	for cmdName := range sqlBoilerCommands {
		for _, skipItem := range skips {
			if cmdName == skipItem {
				continue CommandNameLoop
			}
		}

		found := false
		for _, item := range list {
			if item == cmdName {
				found = true
				break
			}
		}

		if !found {
			t.Error("Expected to find command name:", cmdName)
		}
	}
}
