var dbNameRand *rand.Rand

func MustTx(transactor boil.Transactor, err error) boil.Transactor {
	if err != nil {
		panic(fmt.Sprintf("Cannot create a transactor: %s", err))
	}
	return transactor
}

func initDBNameRand(input string) {
	sum := md5.Sum([]byte(input))

	var sumInt string
	for _, v := range sum {
		sumInt = sumInt + strconv.Itoa(int(v))
	}

	// Cut integer to 18 digits to ensure no int64 overflow.
	sumInt = sumInt[:18]

	sumTmp := sumInt
	for i, v := range sumInt {
		if v == '0' {
			sumTmp = sumInt[i+1:]
			continue
		}
		break
	}

	sumInt = sumTmp

	randSeed, err := strconv.ParseInt(sumInt, 0, 64)
	if err != nil {
		fmt.Printf("Unable to parse sumInt: %s", err)
		os.Exit(-1)
	}

	dbNameRand = rand.New(rand.NewSource(randSeed))
}

var alphabetChars = "abcdefghijklmnopqrstuvwxyz"
func randStr(length int) string {
	c := len(alphabetChars)

	output := make([]rune, length)
	for i := 0; i < length; i++ {
		output[i] = rune(alphabetChars[dbNameRand.Intn(c)])
	}

	return string(output)
}

// getDBNameHash takes a database name in, and generates
// a random string using the database name as the rand Seed.
// getDBNameHash is used to generate unique test database names.
func getDBNameHash(input string) string {
	initDBNameRand(input)
	return randStr(40)
}

// byteSliceEqual calls bytes.Equal to check that two
// byte slices are equal. bytes.Equal is not used directly
// to avoid an unecessary conditional type import.
func byteSliceEqual(a []byte, b []byte) bool {
	return bytes.Equal(a, b)
}
