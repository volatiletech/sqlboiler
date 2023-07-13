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

func BenchmarkParsePointScientificNotation(b *testing.B) {
	pString := "(-1.73521E-5,-2.7159127785469e-7)"
	for i := 0; i < 100000; i++ {
		_, err := parsePoint(pString)
		if err != nil {
			b.Error("parsePoint failed", err)
		}
	}
}
