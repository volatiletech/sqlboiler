/*
Package sqlmock is a mock library implementing sql driver. Which has one and only
purpose - to simulate any sql driver behavior in tests, without needing a real
database connection. It helps to maintain correct **TDD** workflow.

It does not require any modifications to your source code in order to test
and mock database operations. Supports concurrency and multiple database mocking.

The driver allows to mock any sql driver method behavior.
*/
package sqlmock

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"regexp"
)

// Sqlmock interface serves to create expectations
// for any kind of database action in order to mock
// and test real database behavior.
type Sqlmock interface {

	// ExpectClose queues an expectation for this database
	// action to be triggered. the *ExpectedClose allows
	// to mock database response
	ExpectClose() *ExpectedClose

	// ExpectationsWereMet checks whether all queued expectations
	// were met in order. If any of them was not met - an error is returned.
	ExpectationsWereMet() error

	// ExpectPrepare expects Prepare() to be called with sql query
	// which match sqlRegexStr given regexp.
	// the *ExpectedPrepare allows to mock database response.
	// Note that you may expect Query() or Exec() on the *ExpectedPrepare
	// statement to prevent repeating sqlRegexStr
	ExpectPrepare(sqlRegexStr string) *ExpectedPrepare

	// ExpectQuery expects Query() or QueryRow() to be called with sql query
	// which match sqlRegexStr given regexp.
	// the *ExpectedQuery allows to mock database response.
	ExpectQuery(sqlRegexStr string) *ExpectedQuery

	// ExpectExec expects Exec() to be called with sql query
	// which match sqlRegexStr given regexp.
	// the *ExpectedExec allows to mock database response
	ExpectExec(sqlRegexStr string) *ExpectedExec

	// ExpectBegin expects *sql.DB.Begin to be called.
	// the *ExpectedBegin allows to mock database response
	ExpectBegin() *ExpectedBegin

	// ExpectCommit expects *sql.Tx.Commit to be called.
	// the *ExpectedCommit allows to mock database response
	ExpectCommit() *ExpectedCommit

	// ExpectRollback expects *sql.Tx.Rollback to be called.
	// the *ExpectedRollback allows to mock database response
	ExpectRollback() *ExpectedRollback

	// MatchExpectationsInOrder gives an option whether to match all
	// expectations in the order they were set or not.
	//
	// By default it is set to - true. But if you use goroutines
	// to parallelize your query executation, that option may
	// be handy.
	MatchExpectationsInOrder(bool)
}

type sqlmock struct {
	ordered bool
	dsn     string
	opened  int
	drv     *mockDriver

	expected []expectation
}

func (c *sqlmock) open() (*sql.DB, Sqlmock, error) {
	db, err := sql.Open("sqlmock", c.dsn)
	if err != nil {
		return db, c, err
	}
	return db, c, db.Ping()
}

func (c *sqlmock) ExpectClose() *ExpectedClose {
	e := &ExpectedClose{}
	c.expected = append(c.expected, e)
	return e
}

func (c *sqlmock) MatchExpectationsInOrder(b bool) {
	c.ordered = b
}

// Close a mock database driver connection. It may or may not
// be called depending on the sircumstances, but if it is called
// there must be an *ExpectedClose expectation satisfied.
// meets http://golang.org/pkg/database/sql/driver/#Conn interface
func (c *sqlmock) Close() error {
	c.drv.Lock()
	defer c.drv.Unlock()

	c.opened--
	if c.opened == 0 {
		delete(c.drv.conns, c.dsn)
	}

	var expected *ExpectedClose
	var fulfilled int
	var ok bool
	for _, next := range c.expected {
		next.Lock()
		if next.fulfilled() {
			next.Unlock()
			fulfilled++
			continue
		}

		if expected, ok = next.(*ExpectedClose); ok {
			break
		}

		next.Unlock()
		if c.ordered {
			return fmt.Errorf("call to database Close, was not expected, next expectation is: %s", next)
		}
	}

	if expected == nil {
		msg := "call to database Close was not expected"
		if fulfilled == len(c.expected) {
			msg = "all expectations were already fulfilled, " + msg
		}
		return fmt.Errorf(msg)
	}

	expected.triggered = true
	expected.Unlock()
	return expected.err
}

func (c *sqlmock) ExpectationsWereMet() error {
	for _, e := range c.expected {
		if !e.fulfilled() {
			return fmt.Errorf("there is a remaining expectation which was not matched: %s", e)
		}
	}
	return nil
}

// Begin meets http://golang.org/pkg/database/sql/driver/#Conn interface
func (c *sqlmock) Begin() (driver.Tx, error) {
	var expected *ExpectedBegin
	var ok bool
	var fulfilled int
	for _, next := range c.expected {
		next.Lock()
		if next.fulfilled() {
			next.Unlock()
			fulfilled++
			continue
		}

		if expected, ok = next.(*ExpectedBegin); ok {
			break
		}

		next.Unlock()
		if c.ordered {
			return nil, fmt.Errorf("call to database transaction Begin, was not expected, next expectation is: %s", next)
		}
	}
	if expected == nil {
		msg := "call to database transaction Begin was not expected"
		if fulfilled == len(c.expected) {
			msg = "all expectations were already fulfilled, " + msg
		}
		return nil, fmt.Errorf(msg)
	}

	expected.triggered = true
	expected.Unlock()
	return c, expected.err
}

func (c *sqlmock) ExpectBegin() *ExpectedBegin {
	e := &ExpectedBegin{}
	c.expected = append(c.expected, e)
	return e
}

// Exec meets http://golang.org/pkg/database/sql/driver/#Execer
func (c *sqlmock) Exec(query string, args []driver.Value) (res driver.Result, err error) {
	query = stripQuery(query)
	var expected *ExpectedExec
	var fulfilled int
	var ok bool
	for _, next := range c.expected {
		next.Lock()
		if next.fulfilled() {
			next.Unlock()
			fulfilled++
			continue
		}

		if c.ordered {
			if expected, ok = next.(*ExpectedExec); ok {
				break
			}
			next.Unlock()
			return nil, fmt.Errorf("call to exec query '%s' with args %+v, was not expected, next expectation is: %s", query, args, next)
		}
		if exec, ok := next.(*ExpectedExec); ok {
			if err := exec.attemptMatch(query, args); err == nil {
				expected = exec
				break
			}
		}
		next.Unlock()
	}
	if expected == nil {
		msg := "call to exec '%s' query with args %+v was not expected"
		if fulfilled == len(c.expected) {
			msg = "all expectations were already fulfilled, " + msg
		}
		return nil, fmt.Errorf(msg, query, args)
	}

	defer expected.Unlock()

	if !expected.queryMatches(query) {
		return nil, fmt.Errorf("exec query '%s', does not match regex '%s'", query, expected.sqlRegex.String())
	}

	if err := expected.argsMatches(args); err != nil {
		return nil, fmt.Errorf("exec query '%s', arguments do not match: %s", query, err)
	}

	expected.triggered = true

	if expected.err != nil {
		return nil, expected.err // mocked to return error
	}

	if expected.result == nil {
		return nil, fmt.Errorf("exec query '%s' with args %+v, must return a database/sql/driver.result, but it was not set for expectation %T as %+v", query, args, expected, expected)
	}

	return expected.result, err
}

func (c *sqlmock) ExpectExec(sqlRegexStr string) *ExpectedExec {
	e := &ExpectedExec{}
	e.sqlRegex = regexp.MustCompile(sqlRegexStr)
	c.expected = append(c.expected, e)
	return e
}

// Prepare meets http://golang.org/pkg/database/sql/driver/#Conn interface
func (c *sqlmock) Prepare(query string) (driver.Stmt, error) {
	var expected *ExpectedPrepare
	var fulfilled int
	var ok bool
	for _, next := range c.expected {
		next.Lock()
		if next.fulfilled() {
			next.Unlock()
			fulfilled++
			continue
		}

		if expected, ok = next.(*ExpectedPrepare); ok {
			break
		}

		next.Unlock()
		if c.ordered {
			return nil, fmt.Errorf("call to Prepare stetement with query '%s', was not expected, next expectation is: %s", query, next)
		}
	}

	query = stripQuery(query)
	if expected == nil {
		msg := "call to Prepare '%s' query was not expected"
		if fulfilled == len(c.expected) {
			msg = "all expectations were already fulfilled, " + msg
		}
		return nil, fmt.Errorf(msg, query)
	}

	expected.triggered = true
	expected.Unlock()
	return &statement{c, query, expected.closeErr}, expected.err
}

func (c *sqlmock) ExpectPrepare(sqlRegexStr string) *ExpectedPrepare {
	e := &ExpectedPrepare{sqlRegex: regexp.MustCompile(sqlRegexStr), mock: c}
	c.expected = append(c.expected, e)
	return e
}

// Query meets http://golang.org/pkg/database/sql/driver/#Queryer
func (c *sqlmock) Query(query string, args []driver.Value) (rw driver.Rows, err error) {
	query = stripQuery(query)
	var expected *ExpectedQuery
	var fulfilled int
	var ok bool
	for _, next := range c.expected {
		next.Lock()
		if next.fulfilled() {
			next.Unlock()
			fulfilled++
			continue
		}

		if c.ordered {
			if expected, ok = next.(*ExpectedQuery); ok {
				break
			}
			next.Unlock()
			return nil, fmt.Errorf("call to query '%s' with args %+v, was not expected, next expectation is: %s", query, args, next)
		}
		if qr, ok := next.(*ExpectedQuery); ok {
			if err := qr.attemptMatch(query, args); err == nil {
				expected = qr
				break
			}
		}
		next.Unlock()
	}

	if expected == nil {
		msg := "call to query '%s' with args %+v was not expected"
		if fulfilled == len(c.expected) {
			msg = "all expectations were already fulfilled, " + msg
		}
		return nil, fmt.Errorf(msg, query, args)
	}

	defer expected.Unlock()

	if !expected.queryMatches(query) {
		return nil, fmt.Errorf("query '%s', does not match regex [%s]", query, expected.sqlRegex.String())
	}

	if err := expected.argsMatches(args); err != nil {
		return nil, fmt.Errorf("exec query '%s', arguments do not match: %s", query, err)
	}

	expected.triggered = true

	if expected.err != nil {
		return nil, expected.err // mocked to return error
	}

	if expected.rows == nil {
		return nil, fmt.Errorf("query '%s' with args %+v, must return a database/sql/driver.rows, but it was not set for expectation %T as %+v", query, args, expected, expected)
	}

	return expected.rows, err
}

func (c *sqlmock) ExpectQuery(sqlRegexStr string) *ExpectedQuery {
	e := &ExpectedQuery{}
	e.sqlRegex = regexp.MustCompile(sqlRegexStr)
	c.expected = append(c.expected, e)
	return e
}

func (c *sqlmock) ExpectCommit() *ExpectedCommit {
	e := &ExpectedCommit{}
	c.expected = append(c.expected, e)
	return e
}

func (c *sqlmock) ExpectRollback() *ExpectedRollback {
	e := &ExpectedRollback{}
	c.expected = append(c.expected, e)
	return e
}

// Commit meets http://golang.org/pkg/database/sql/driver/#Tx
func (c *sqlmock) Commit() error {
	var expected *ExpectedCommit
	var fulfilled int
	var ok bool
	for _, next := range c.expected {
		next.Lock()
		if next.fulfilled() {
			next.Unlock()
			fulfilled++
			continue
		}

		if expected, ok = next.(*ExpectedCommit); ok {
			break
		}

		next.Unlock()
		if c.ordered {
			return fmt.Errorf("call to commit transaction, was not expected, next expectation is: %s", next)
		}
	}
	if expected == nil {
		msg := "call to commit transaction was not expected"
		if fulfilled == len(c.expected) {
			msg = "all expectations were already fulfilled, " + msg
		}
		return fmt.Errorf(msg)
	}

	expected.triggered = true
	expected.Unlock()
	return expected.err
}

// Rollback meets http://golang.org/pkg/database/sql/driver/#Tx
func (c *sqlmock) Rollback() error {
	var expected *ExpectedRollback
	var fulfilled int
	var ok bool
	for _, next := range c.expected {
		next.Lock()
		if next.fulfilled() {
			next.Unlock()
			fulfilled++
			continue
		}

		if expected, ok = next.(*ExpectedRollback); ok {
			break
		}

		next.Unlock()
		if c.ordered {
			return fmt.Errorf("call to rollback transaction, was not expected, next expectation is: %s", next)
		}
	}
	if expected == nil {
		msg := "call to rollback transaction was not expected"
		if fulfilled == len(c.expected) {
			msg = "all expectations were already fulfilled, " + msg
		}
		return fmt.Errorf(msg)
	}

	expected.triggered = true
	expected.Unlock()
	return expected.err
}
