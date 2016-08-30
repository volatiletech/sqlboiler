[![Build Status](https://travis-ci.org/DATA-DOG/go-sqlmock.png)](https://travis-ci.org/DATA-DOG/go-sqlmock)
[![GoDoc](https://godoc.org/github.com/DATA-DOG/go-sqlmock?status.png)](https://godoc.org/github.com/DATA-DOG/go-sqlmock)

# Sql driver mock for Golang

**sqlmock** is a mock library implementing [sql/driver](https://godoc.org/database/sql/driver). Which has one and only
purpose - to simulate any **sql** driver behavior in tests, without needing a real database connection. It helps to
maintain correct **TDD** workflow.

- this library is now complete and stable. (you may not find new changes for this reason)
- supports concurrency and multiple connections.
- does not require any modifications to your source code.
- the driver allows to mock any sql driver method behavior.
- has strict by default expectation order matching.
- has no vendor dependencies.

## Install

    go get gopkg.in/DATA-DOG/go-sqlmock.v1

If you need an old version, checkout **go-sqlmock** at gopkg.in:

    go get gopkg.in/DATA-DOG/go-sqlmock.v0

## Documentation and Examples

Visit [godoc](http://godoc.org/github.com/DATA-DOG/go-sqlmock) for general examples and public api reference.
See **.travis.yml** for supported **go** versions.
Different use case, is to functionally test with a real database - [go-txdb](https://github.com/DATA-DOG/go-txdb)
all database related actions are isolated within a single transaction so the database can remain in the same state.

See implementation examples:

- [blog API server](https://github.com/DATA-DOG/go-sqlmock/tree/master/examples/blog)
- [the same orders example](https://github.com/DATA-DOG/go-sqlmock/tree/master/examples/orders)

### Something you may want to test

``` go
package main

import "database/sql"

func recordStats(db *sql.DB, userID, productID int64) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return
	}

	defer func() {
		switch err {
		case nil:
			err = tx.Commit()
		default:
			tx.Rollback()
		}
	}()

	if _, err = tx.Exec("UPDATE products SET views = views + 1"); err != nil {
		return
	}
	if _, err = tx.Exec("INSERT INTO product_viewers (user_id, product_id) VALUES (?, ?)", userID, productID); err != nil {
		return
	}
	return
}

func main() {
	// @NOTE: the real connection is not required for tests
	db, err := sql.Open("mysql", "root@/blog")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if err = recordStats(db, 1 /*some user id*/, 5 /*some product id*/); err != nil {
		panic(err)
	}
}
```

### Tests with sqlmock

``` go
package main

import (
	"fmt"
	"testing"

	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

// a successful case
func TestShouldUpdateStats(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE products").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO product_viewers").WithArgs(2, 3).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// now we execute our method
	if err = recordStats(db, 2, 3); err != nil {
		t.Errorf("error was not expected while updating stats: %s", err)
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
}

// a failing test case
func TestShouldRollbackStatUpdatesOnFailure(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE products").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO product_viewers").
		WithArgs(2, 3).
		WillReturnError(fmt.Errorf("some error"))
	mock.ExpectRollback()

	// now we execute our method
	if err = recordStats(db, 2, 3); err == nil {
		t.Errorf("was expecting an error, but there was none")
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
}
```

## Matching arguments like time.Time

There may be arguments which are of `struct` type and cannot be compared easily by value like `time.Time`. In this case
**sqlmock** provides an [Argument](https://godoc.org/github.com/DATA-DOG/go-sqlmock#Argument) interface which
can be used in more sophisticated matching. Here is a simple example of time argument matching:

``` go
type AnyTime struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

func TestAnyTimeArgument(t *testing.T) {
	t.Parallel()
	db, mock, err := New()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectExec("INSERT INTO users").
		WithArgs("john", AnyTime{}).
		WillReturnResult(NewResult(1, 1))

	_, err = db.Exec("INSERT INTO users(name, created_at) VALUES (?, ?)", "john", time.Now())
	if err != nil {
		t.Errorf("error '%s' was not expected, while inserting a row", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
}
```

It only asserts that argument is of `time.Time` type.

## Run tests

    go test -race

## Changes

- **2016-02-23** - added **sqlmock.AnyArg()** function to provide any kind
  of argument matcher.
- **2016-02-23** - convert expected arguments to driver.Value as natural
  driver does, the change may affect time.Time comparison and will be
  stricter. See [issue](https://github.com/DATA-DOG/go-sqlmock/issues/31).
- **2015-08-27** - **v1** api change, concurrency support, all known issues fixed.
- **2014-08-16** instead of **panic** during reflect type mismatch when comparing query arguments - now return error
- **2014-08-14** added **sqlmock.NewErrorResult** which gives an option to return driver.Result with errors for
interface methods, see [issue](https://github.com/DATA-DOG/go-sqlmock/issues/5)
- **2014-05-29** allow to match arguments in more sophisticated ways, by providing an **sqlmock.Argument** interface
- **2014-04-21** introduce **sqlmock.New()** to open a mock database connection for tests. This method
calls sql.DB.Ping to ensure that connection is open, see [issue](https://github.com/DATA-DOG/go-sqlmock/issues/4).
This way on Close it will surely assert if all expectations are met, even if database was not triggered at all.
The old way is still available, but it is advisable to call db.Ping manually before asserting with db.Close.
- **2014-02-14** RowsFromCSVString is now a part of Rows interface named as FromCSVString.
It has changed to allow more ways to construct rows and to easily extend this API in future.
See [issue 1](https://github.com/DATA-DOG/go-sqlmock/issues/1)
**RowsFromCSVString** is deprecated and will be removed in future

## Contributions

Feel free to open a pull request. Note, if you wish to contribute an extension to public (exported methods or types) -
please open an issue before, to discuss whether these changes can be accepted. All backward incompatible changes are
and will be treated cautiously

## License

The [three clause BSD license](http://en.wikipedia.org/wiki/BSD_licenses)

