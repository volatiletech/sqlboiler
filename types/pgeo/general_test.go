package pgeo

import (
	"testing"
)

func BenchmarkParsePoint(b *testing.B) {
	pString := "(-13.735219157895635,-72.7159127785469)"
	for i := 0; i < 100000; i++ {
		_, err := parsePoint(pString)
		if err != nil {
			b.Error("parsePoint failed", err)
		}
	}
}
