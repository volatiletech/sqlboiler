package cmds

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	"github.com/pobri19/sqlboiler/dbdrivers"
)

type ImportSorter []string

func (i importList) Len() int {
	return len(i)
}

func (i importList) Swap(k, j int) {
	i[k], i[j] = i[j], i[k]
}

func (i importList) Less(k, j int) bool {
	res := strings.Compare(strings.TrimLeft(i[k], "_ "), strings.TrimLeft(i[j], "_ "))
	if res <= 0 {
		return true
	}

	return false
}

func combineImports(a, b imports) imports {
	var c imports

	c.standard = removeDuplicates(combineStringSlices(a.standard, b.standard))
	c.thirdparty = removeDuplicates(combineStringSlices(a.thirdparty, b.thirdparty))

	sort.Sort(c.standard)
	sort.Sort(c.thirdparty)

	return c
}

func combineConditionalTypeImports(a imports, b map[string]imports, columns []dbdrivers.DBColumn) imports {
	tmpImp := imports{
		standard:   make(importList, len(a.standard)),
		thirdparty: make(importList, len(a.thirdparty)),
	}

	copy(tmpImp.standard, a.standard)
	copy(tmpImp.thirdparty, a.thirdparty)

	for _, col := range columns {
		for key, imp := range b {
			if col.Type == key {
				tmpImp.standard = append(tmpImp.standard, imp.standard...)
				tmpImp.thirdparty = append(tmpImp.thirdparty, imp.thirdparty...)
			}
		}
	}

	tmpImp.standard = removeDuplicates(tmpImp.standard)
	tmpImp.thirdparty = removeDuplicates(tmpImp.thirdparty)

	sort.Sort(tmpImp.standard)
	sort.Sort(tmpImp.thirdparty)

	return tmpImp
}

func buildImportString(imps *imports) []byte {
	stdlen, thirdlen := len(imps.standard), len(imps.thirdparty)
	if stdlen+thirdlen < 1 {
		return []byte{}
	}

	if stdlen+thirdlen == 1 {
		var imp string
		if stdlen == 1 {
			imp = imps.standard[0]
		} else {
			imp = imps.thirdparty[0]
		}
		return []byte(fmt.Sprintf(`import %s`, imp))
	}

	buf := &bytes.Buffer{}
	buf.WriteString("import (")
	for _, std := range imps.standard {
		fmt.Fprintf(buf, "\n\t%s", std)
	}
	if stdlen != 0 && thirdlen != 0 {
		buf.WriteString("\n")
	}
	for _, third := range imps.thirdparty {
		fmt.Fprintf(buf, "\n\t%s", third)
	}
	buf.WriteString("\n)\n")

	return buf.Bytes()
}
