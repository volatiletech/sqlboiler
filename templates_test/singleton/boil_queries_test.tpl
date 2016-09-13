var dbNameRand *rand.Rand

func MustTx(transactor boil.Transactor, err error) boil.Transactor {
	if err != nil {
		panic(fmt.Sprintf("Cannot create a transactor: %s", err))
	}
	return transactor
}

var rgxPGFkey = regexp.MustCompile(`(?m)(?s)^ALTER TABLE ONLY.*?ADD CONSTRAINT.*?FOREIGN KEY.*?;\n`)
var rgxMySQLkey = regexp.MustCompile(`(?m)((,\n)?\s+CONSTRAINT.*?FOREIGN KEY.*?\n)+`)

func newFKeyDestroyer(reader io.Reader) io.Reader {
	return &fKeyDestroyer{
		reader: reader,
	}
}

type fKeyDestroyer struct {
	reader io.Reader
	buf    *bytes.Buffer
}

func (f *fKeyDestroyer) Read(b []byte) (int, error) {
	if f.buf == nil {
		all, err := ioutil.ReadAll(f.reader)
		if err != nil {
			return 0, err
		}

		f.buf = bytes.NewBuffer(rgxMySQLkey.ReplaceAll(rgxPGFkey.ReplaceAll(all, []byte{}), []byte{}))
	}

	return f.buf.Read(b)
}
