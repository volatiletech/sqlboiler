package cmds

import (
	"reflect"
	"testing"

	"github.com/pobri19/sqlboiler/dbdrivers"
)

func TestCombineConditionalTypeImports(t *testing.T) {
	imports1 := imports{
		standard: importList{
			`"errors"`,
			`"fmt"`,
		},
		thirdparty: importList{
			`"github.com/pobri19/sqlboiler/boil"`,
		},
	}

	importsExpected := imports{
		standard: importList{
			`"errors"`,
			`"fmt"`,
			`"time"`,
		},
		thirdparty: importList{
			`"github.com/pobri19/sqlboiler/boil"`,
			`"gopkg.in/guregu/null.v3"`,
		},
	}

	cols := []dbdrivers.Column{
		dbdrivers.Column{
			Type: "null.Time",
		},
		dbdrivers.Column{
			Type: "null.Time",
		},
		dbdrivers.Column{
			Type: "time.Time",
		},
		dbdrivers.Column{
			Type: "null.Float",
		},
	}

	res1 := combineConditionalTypeImports(imports1, sqlBoilerConditionalTypeImports, cols)

	if !reflect.DeepEqual(res1, importsExpected) {
		t.Errorf("Expected res1 to match importsExpected, got:\n\n%#v\n", res1)
	}

	imports2 := imports{
		standard: importList{
			`"errors"`,
			`"fmt"`,
			`"time"`,
		},
		thirdparty: importList{
			`"github.com/pobri19/sqlboiler/boil"`,
			`"gopkg.in/guregu/null.v3"`,
		},
	}

	res2 := combineConditionalTypeImports(imports2, sqlBoilerConditionalTypeImports, cols)

	if !reflect.DeepEqual(res2, importsExpected) {
		t.Errorf("Expected res2 to match importsExpected, got:\n\n%#v\n", res1)
	}
}
