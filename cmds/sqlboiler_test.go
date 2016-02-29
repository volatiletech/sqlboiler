package cmds

import "testing"

func TestInitTemplates(t *testing.T) {
	templates, err := initTemplates()
	if err != nil {
		t.Errorf("Unable to init templates: %s", err)
	}

	if len(templates) < 2 {
		t.Errorf("Expected > 2 templates to be loaded from templates folder, only loaded: %d\n\n%#v", len(templates), templates)
	}
}
