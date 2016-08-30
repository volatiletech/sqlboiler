package sqlmock

import (
	"database/sql/driver"
	"encoding/csv"
	"io"
	"strings"
)

// CSVColumnParser is a function which converts trimmed csv
// column string to a []byte representation. currently
// transforms NULL to nil
var CSVColumnParser = func(s string) []byte {
	switch {
	case strings.ToLower(s) == "null":
		return nil
	}
	return []byte(s)
}

// Rows interface allows to construct rows
// which also satisfies database/sql/driver.Rows interface
type Rows interface {
	// composed interface, supports sql driver.Rows
	driver.Rows

	// AddRow composed from database driver.Value slice
	// return the same instance to perform subsequent actions.
	// Note that the number of values must match the number
	// of columns
	AddRow(columns ...driver.Value) Rows

	// FromCSVString build rows from csv string.
	// return the same instance to perform subsequent actions.
	// Note that the number of values must match the number
	// of columns
	FromCSVString(s string) Rows

	// RowError allows to set an error
	// which will be returned when a given
	// row number is read
	RowError(row int, err error) Rows

	// CloseError allows to set an error
	// which will be returned by rows.Close
	// function.
	//
	// The close error will be triggered only in cases
	// when rows.Next() EOF was not yet reached, that is
	// a default sql library behavior
	CloseError(err error) Rows
}

type rows struct {
	cols     []string
	rows     [][]driver.Value
	pos      int
	nextErr  map[int]error
	closeErr error
}

func (r *rows) Columns() []string {
	return r.cols
}

func (r *rows) Close() error {
	return r.closeErr
}

// advances to next row
func (r *rows) Next(dest []driver.Value) error {
	r.pos++
	if r.pos > len(r.rows) {
		return io.EOF // per interface spec
	}

	for i, col := range r.rows[r.pos-1] {
		dest[i] = col
	}

	return r.nextErr[r.pos-1]
}

// NewRows allows Rows to be created from a
// sql driver.Value slice or from the CSV string and
// to be used as sql driver.Rows
func NewRows(columns []string) Rows {
	return &rows{cols: columns, nextErr: make(map[int]error)}
}

func (r *rows) CloseError(err error) Rows {
	r.closeErr = err
	return r
}

func (r *rows) RowError(row int, err error) Rows {
	r.nextErr[row] = err
	return r
}

func (r *rows) AddRow(values ...driver.Value) Rows {
	if len(values) != len(r.cols) {
		panic("Expected number of values to match number of columns")
	}

	row := make([]driver.Value, len(r.cols))
	for i, v := range values {
		row[i] = v
	}

	r.rows = append(r.rows, row)
	return r
}

func (r *rows) FromCSVString(s string) Rows {
	res := strings.NewReader(strings.TrimSpace(s))
	csvReader := csv.NewReader(res)

	for {
		res, err := csvReader.Read()
		if err != nil || res == nil {
			break
		}

		row := make([]driver.Value, len(r.cols))
		for i, v := range res {
			row[i] = CSVColumnParser(strings.TrimSpace(v))
		}
		r.rows = append(r.rows, row)
	}
	return r
}
