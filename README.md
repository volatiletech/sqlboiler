# SQLBoiler

[![License](https://img.shields.io/badge/license-BSD-blue.svg)](https://github.com/vattle/sqlboiler/blob/master/LICENSE)
[![GoDoc](https://godoc.org/github.com/vattle/sqlboiler?status.svg)](https://godoc.org/github.com/vattle/sqlboiler)
[![CircleCI](https://circleci.com/gh/vattle/sqlboiler.svg?style=shield)](https://circleci.com/gh/vattle/sqlboiler)
[![Go Report Card](https://goreportcard.com/badge/vattle/sqlboiler)](http://goreportcard.com/report/vattle/sqlboiler)

SQLBoiler is a tool to generate a Go data model tailored to your database schema.

It is a "database-first" ORM as opposed to "code-first" (like gorm/gorp).
That means you must first create your database schema. Please use something
like goose or some other migration tool to manage this part of the database's
lifecycle.

## About SQL Boiler

#### Features

- Full model generation
- High performance through generation
- Uses boil.Executor (simple interface, sql.DB, sqlx.DB etc. compatible)
- Easy workflow (models can always be regenerated, full auto-complete)
- Strongly typed querying (usually no converting or binding to pointers)
- Hooks (Before/After Create/Update)
- Automatic CreatedAt/UpdatedAt
- Relationships/Associations
- Eager loading
- Transactions
- Raw SQL fallbacks
- Compatibility tests (Run against your own DB schema)
- Debug logging

#### Supported Databases

- PostgreSQL

Note: Seeking contributors for other database engines.

#### Automatic CreatedAt/UpdatedAt

If your generated SQLBoiler models package can find columns with the
names `created_at` or `updated_at` it will automatically set them
to `time.Now()` in your database, and update your object appropriately.

Note: You can set the timezone for this feature by calling `boil.SetLocation()`

#### Example Queries

```go
import (
  // Import this so we don't have to use qm.Limit etc.
  . "github.com/vattle/sqlboiler/boil/qm"
)

// Open handle to database like normal
db, err := sql.Open("postgres", "dbname=fun user=abc")
if err != nil {
  return err
}

// Query all users
users, err := models.Users(db).All()

// Panic-able if you like to code that way
users := models.Users(db).AllP()

// More complex query
users, err := models.Users(db, Where("age > ?", 30), Limit(5), Offset(6)).All()

// Ultra complex query
users, err := models.Users(db,
  Select("id", "name"),
  InnerJoin("credit_cards c on c.user_id = users.id"),
  Where("age > ?", 30),
  AndIn("c.kind in ?", "visa", "mastercard"),
  Or("email like ?", "%aol.com%"),
  GroupBy("id", "name"),
  Having("count(c.id) > ?", 2),
  Limit(5),
  Offset(6),
).All()

// Use any "boil.Executor" implementation (*sql.DB, *sql.Tx, data-dog mock db)
// for any query.
tx, err := db.Begin()
if err != nil {
  return err
}
users, err := models.Users(tx).All()

// Relationships
user, err := models.Users(db).One()
if err != nil {
  return err
}
movies, err := user.FavoriteMovies(db).All()

// Eager loading
users, err := models.Users(db, Load("FavoriteMovies")).All()
if err != nil {
  return err
}
fmt.Println(len(users.Loaded.FavoriteMovies))
```

## How to boil your database

#### Download

```shell
go get -u -t github.com/vattle/sqlboiler
```

#### Configuration

Create a configuration file. Because the project uses [viper](github.com/spf13/viper), TOML, JSON and YAML
are all supported. Environment variables are also able to be used.
We will assume TOML for the rest of the documentation.

The configuration file is searched for in the following directories in this
order:

- `./`
- `$XDG_CONFIG_HOME/sqlboiler/`
- `$HOME/.config/sqlboiler/`

```shell
vim ./sqlboiler.toml
```

Currently the only section in the configuration is `postgres`, and it takes
configuration parameters that will be passed mostly directly to the
[pq](github.com/lib/pq) driver. Here is a rundown of all the different
values that can go in that section:

| Name | Required | Default |
| --- | --- | --- |
| dbname  | yes       | none      |
| host    | yes       | none      |
| port    | no        | 5432      |
| user    | yes       | none      |
| pass    | no        | none      |
| sslmode | no        | 'require' |

Example:

```toml
[postgres]
dbname="dbname"
host="localhost"
port=5432
user="dbusername"
pass="dbpassword"
sslmode="require"
```

## Usage

#### Initial Generation

After creating a configuration file that points at the database we want to
generate models for, we can invoke the sqlboiler command line utility.

```text
SQL Boiler generates a Go ORM from template files, tailored to your database schema.
Complete documentation is available at http://github.com/vattle/sqlboiler

Usage:
  sqlboiler [flags] <driver>

Examples:
sqlboiler postgres

Flags:
  -d, --debug                 Debug mode prints stack traces on error
  -x, --exclude stringSlice   Tables to be excluded from the generated package
  -o, --output string         The name of the folder to output to (default "models")
  -p, --pkgname string        The name you wish to assign to your generated package (default "models")
```

Follow the steps below to do some basic model generation. Once we've generated
our models, we can run the compatibility tests which will exercise the entirety
of the generated code. This way we can ensure that our database is compatible
with sqlboiler. If you find there are some failing tests, please check the
[faq](#FAQ) section.

```shell
# Generate our models and exclude the migrations table
sqlboiler -x goose_migrations postgres

# Run the generated tests
go test ./models # This requires an administrator postgres user because of some
                 # voodoo we do to disable triggers for the generated test db
```

## FAQ

Work in Progress
