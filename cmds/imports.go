package cmds

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
)

type ImportSorter []string

func (i ImportSorter) Len() int {
	return len(i)
}

func (i ImportSorter) Swap(k, j int) {
	i[k], i[j] = i[j], i[k]
}

func (i ImportSorter) Less(k, j int) bool {
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

	sort.Sort(ImportSorter(c.standard))
	sort.Sort(ImportSorter(c.thirdparty))

	return c
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
