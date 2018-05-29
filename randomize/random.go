package randomize

import (
	"crypto/md5"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/volatiletech/sqlboiler/strmangle"
)

const alphabetAll = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const alphabetLowerAlpha = "abcdefghijklmnopqrstuvwxyz"

// Str creates a randomized string from printable characters in the alphabet
func Str(s *Seed, ln int) string {
	str := make([]byte, ln)
	for i := 0; i < ln; i++ {
		str[i] = byte(alphabetAll[s.NextInt()%len(alphabetAll)])
	}

	return string(str)
}

// FormattedString checks a field type to see if it's in a range of special
// values and if so returns a randomized string for it.
func FormattedString(s *Seed, fieldType string) (string, bool) {
	if strings.HasPrefix(fieldType, "enum") {
		enum, err := EnumValue(s, fieldType)
		if err != nil {
			panic(err)
		}

		return enum, true
	}

	switch fieldType {
	case "json", "jsonb":
		return `"` + Str(s, 1) + `"`, true
	case "interval":
		return strconv.Itoa((s.NextInt()%26)+2) + " days", true
	case "uuid":
		randomUUID, err := uuid.NewV4()
		if err != nil {
			panic(err)
		}
		return randomUUID.String(), true
	case "cidr", "inet":
		return randNetAddr(s), true
	case "macaddr":
		return randMacAddr(s), true
	case "pg_lsn":
		return randLsn(s), true
	case "txid_snapshot":
		return randTxID(s), true
	case "money":
		return randMoney(s), true
	}

	return "", false
}

// MediumInt is a special case in mysql (thanks for that -_-)
// this function checks if the fieldtype matches and if so returns
// a random value in the proper range.
func MediumInt(s *Seed, fieldType string) (int32, bool) {
	if fieldType == "mediumint" {
		return int32(s.NextInt()) % 8388607, true
	}

	return 0, false
}

// Date generates a random time.Time between 1850 and 2050.
// Only the Day/Month/Year columns are set so that Dates and DateTimes do
// not cause mismatches in the test data comparisons.
func Date(s *Seed) time.Time {
	t := time.Date(
		1972+s.NextInt()%60,
		time.Month(1+(s.NextInt()%12)),
		1+(s.NextInt()%25),
		0,
		0,
		0,
		0,
		time.UTC,
	)

	return t
}

// EnumValue takes an enum field type, parses it's definition
// to figure out valid values, and selects a random one from within them.
func EnumValue(s *Seed, enum string) (string, error) {
	vals := strmangle.ParseEnumVals(enum)
	if vals == nil || len(vals) == 0 {
		return "", fmt.Errorf("unable to parse enum string: %s", enum)
	}

	return vals[s.NextInt()%len(vals)], nil
}

// ByteSlice creates a random set of bytes (non-printables included)
func ByteSlice(s *Seed, ln int) []byte {
	str := make([]byte, ln)
	for i := 0; i < ln; i++ {
		str[i] = byte(s.NextInt() % 256)
	}

	return str
}

func randNetAddr(s *Seed) string {
	return fmt.Sprintf(
		"%d.%d.%d.%d",
		s.NextInt()%254+1,
		s.NextInt()%254+1,
		s.NextInt()%254+1,
		s.NextInt()%254+1,
	)
}

func randMacAddr(s *Seed) string {
	buf := make([]byte, 6)
	for i := range buf {
		buf[i] = byte(s.NextInt())
	}

	// Set the local bit
	buf[0] |= 2
	return fmt.Sprintf(
		"%02x:%02x:%02x:%02x:%02x:%02x",
		buf[0], buf[1], buf[2], buf[3], buf[4], buf[5],
	)
}

func randLsn(s *Seed) string {
	a := s.NextInt() % 9000000
	b := s.NextInt() % 9000000
	return fmt.Sprintf("%d/%d", a, b)
}

func randTxID(s *Seed) string {
	// Order of integers is relevant
	a := s.NextInt()%200 + 100
	b := a + 100
	c := a
	d := a + 50
	return fmt.Sprintf("%d:%d:%d,%d", a, b, c, d)
}

func randMoney(s *Seed) string {
	return fmt.Sprintf("%d.00", s.NextInt())
}

// StableDBName takes a database name in, and generates
// a random string using the database name as the rand Seed.
// getDBNameHash is used to generate unique test database names.
func StableDBName(input string) string {
	return randStrFromSource(stableSource(input), 40)
}

// stableSource takes an input value, and produces a random
// seed from it that will produce very few collisions in
// a 40 character random string made from a different alphabet.
func stableSource(input string) *rand.Rand {
	sum := md5.Sum([]byte(input))
	var seed int64
	for i, byt := range sum {
		seed ^= int64(byt) << uint((i*4)%64)
	}
	return rand.New(rand.NewSource(seed))
}

func randStrFromSource(r *rand.Rand, length int) string {
	ln := len(alphabetLowerAlpha)

	output := make([]rune, length)
	for i := 0; i < length; i++ {
		output[i] = rune(alphabetLowerAlpha[r.Intn(ln)])
	}

	return string(output)
}
