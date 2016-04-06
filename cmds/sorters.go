package cmds

import "strings"

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

func (t templater) Len() int {
	return len(t)
}

func (t templater) Swap(k, j int) {
	t[k], t[j] = t[j], t[k]
}

func (t templater) Less(k, j int) bool {
	// Make sure "struct" goes to the front
	if t[k].Name() == "struct.tpl" {
		return true
	}

	res := strings.Compare(t[k].Name(), t[j].Name())
	if res <= 0 {
		return true
	}

	return false
}
