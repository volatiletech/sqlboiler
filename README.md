![sqlboiler logo](http://i.imgur.com/NJtCT7y.png)

[![License](https://img.shields.io/badge/license-BSD-blue.svg)](https://github.com/vattle/sqlboiler/blob/master/LICENSE)
[![GoDoc](https://godoc.org/github.com/vattle/sqlboiler?status.svg)](https://godoc.org/github.com/vattle/sqlboiler)
[![Slack](https://img.shields.io/badge/slack-%23general-lightgrey.svg)](https://sqlboiler.from-the.cloud)
[![CircleCI](https://circleci.com/gh/vattle/sqlboiler.svg?style=shield)](https://circleci.com/gh/vattle/sqlboiler)
[![Go Report Card](https://goreportcard.com/badge/vattle/sqlboiler)](http://goreportcard.com/report/vattle/sqlboiler)

SQLBoiler is a tool to generate a Go ORM tailored to your database schema.

It is a "database-first" ORM as opposed to "code-first" (like gorm/gorp).
That means you must first create your database schema. Please use something
like [goose](https://bitbucket.org/liamstask/goose), [sql-migrate](https://github.com/rubenv/sql-migrate)
or some other migration tool to manage this part of the database's life-cycle.

## Why another ORM

While attempting to migrate a legacy Rails database, we realized how much ActiveRecord benefitted us in terms of development velocity.
Coming over to the Go `database/sql` package after using ActiveRecord feels extremely repetitive, super long-winded and down-right boring.
Being Go veterans we knew the state of ORMs was shaky, and after a quick review we found what our fears confirmed. Most packages out
there are code-first, reflect-based and have a very weak story around relationships between models. So with that we set out with these goals:

* Work with existing databases: Don't be the tool to define the schema, that's better left to other tools.
* ActiveRecord-like productivity: Eliminate all sql boilerplate, have relationships as a first-class concept.
* Go-like feel: Work with normal structs, call functions, no hyper-magical struct tags, small interfaces.
* Go-like performance: [Benchmark](#benchmarks) and optimize the hot-paths, perform like hand-rolled `sql.DB` code.

We believe with SQLBoiler and our database-first code-generation approach we've been able to successfully meet all of these goals. On top
of that SQLBoiler also confers the following benefits:

* The models package is type safe. This means no chance of random panics due to passing in the wrong type. No need for interface{}.
* Our types closely correlate to your database column types. This is expanded by our extended null package which supports nearly all Go data types.
* A system that is easy to debug. Your ORM is tailored to your schema, the code paths should be easy to trace since it's not all buried in reflect.
* Auto-completion provides work-flow efficiency gains.

Table of Contents
=================

  * [SQLBoiler](#sqlboiler)
    * [Why another ORM](#why-another-orm)
    * [About SQL Boiler](#about-sql-boiler)
      * [Features](#features)
      * [Supported Databases](#supported-databases)
      * [A Small Taste](#a-small-taste)
    * [Requirements &amp; Pro Tips](#requirements--pro-tips)
      * [Requirements](#requirements)
      * [Pro Tips](#pro-tips)
    * [Getting started](#getting-started)
        * [Download](#download)
        * [Configuration](#configuration)
        * [Initial Generation](#initial-generation)
        * [Regeneration](#regeneration)
        * [Extending Generated Models](#extending-generated-models)
    * [Diagnosing Problems](#diagnosing-problems)
    * [Features &amp; Examples](#features--examples)
      * [Automatic CreatedAt/UpdatedAt](#automatic-createdatupdatedat)
        * [Overriding Automatic Timestamps](#overriding-automatic-timestamps)
      * [Query Building](#query-building)
      * [Query Mod System](#query-mod-system)
      * [Function Variations](#function-variations)
      * [Finishers](#finishers)
      * [Raw Query](#raw-query)
      * [Binding](#binding)
      * [Relationships](#relationships)
      * [Hooks](#hooks)
      * [Transactions](#transactions)
      * [Debug Logging](#debug-logging)
      * [Select](#select)
      * [Find](#find)
      * [Insert](#insert)
      * [Update](#update)
      * [Delete](#delete)
      * [Upsert](#upsert)
      * [Reload](#reload)
      * [Exists](#exists)
      * [Enums](#enums)
    * [FAQ](#faq)
        * [Won't compiling models for a huge database be very slow?](#wont-compiling-models-for-a-huge-database-be-very-slow)
        * [Missing imports for generated package](#missing-imports-for-generated-package)
  * [Benchmarks](#benchmarks)

## About SQL Boiler

### Features

- Full model generation
- Extremely fast code generation
- High performance through generation & intelligent caching
- Uses boil.Executor (simple interface, sql.DB, sqlx.DB etc. compatible)
- Easy workflow (models can always be regenerated, full auto-complete)
- Strongly typed querying (usually no converting or binding to pointers)
- Hooks (Before/After Create/Select/Update/Delete/Upsert)
- Automatic CreatedAt/UpdatedAt
- Table whitelist/blacklist
- Relationships/Associations
- Eager loading (recursive)
- Custom struct tags
- Transactions
- Raw SQL fallback
- Compatibility tests (Run against your own DB schema)
- Debug logging
- Schemas support
- 1d arrays, json, hstore & more
- Enum types

### Supported Databases

- PostgreSQL
- MySQL
- Microsoft SQL Server

*Note: Seeking contributors for other database engines.*

*Microsoft SQL Server: Limit with offset support only for SQL Server 2012 and above.*

### A Small Taste

For a comprehensive list of available operations and examples please see [Features & Examples](#features--examples).

```go
import (
  // Import this so we don't have to use qm.Limit etc.
  . "github.com/vattle/sqlboiler/queries/qm"
)

// Open handle to database like normal
db, err := sql.Open("postgres", "dbname=fun user=abc")
if err != nil {
  return err
}

// If you don't want to pass in db to all generated methods
// you can use boil.SetDB to set it globally, and then use
// the G variant methods like so:
boil.SetDB(db)
users, err := models.UsersG().All()

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
  Or("email like ?", `%aol.com%`),
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
fmt.Println(len(users.R.FavoriteMovies))
```

## Requirements & Pro Tips

### Requirements

* Go 1.6 minimum, and Go 1.7 for compatibility tests.
* Table names and column names should use `snake_case` format.
  * We require `snake_case` table names and column names. This is a recommended default in Postgres,
  and we agree that it's good form, so we're enforcing this format for all drivers for the time being.
* Join tables should use a *composite primary key*.
  * For join tables to be used transparently for relationships your join table must have
  a *composite primary key* that encompasses both foreign table foreign keys. For example, on a
  join table named `user_videos` you should have: `primary key(user_id, video_id)`, with both `user_id`
  and `video_id` being foreign key columns to the users and videos tables respectively.
* For MySQL if using the `github.com/go-sql-driver/mysql` driver, please activate
  [time.Time parsing](https://github.com/go-sql-driver/mysql#timetime-support) when making your
  MySQL database connection. SQLBoiler uses `time.Time` and `null.Time` to represent time in
  it's models and without this enabled any models with `DATE`/`DATETIME` columns will not work.

### Pro Tips
* Foreign key column names should end with `_id`.
  * Foreign key column names in the format `x_id` will generate clearer method names.
  It is advisable to use this naming convention whenever it makes sense for your database schema.
* If you never plan on using the hooks functionality you can disable generation of this
  feature using the `--no-hooks` flag. This will save you some binary size.

## Getting started

[SQLBoiler Screencast #1: How to get started](https://www.youtube.com/watch?v=fKmRemtmi0Y)

#### Download

```shell
go get -u -t github.com/vattle/sqlboiler
```

#### Configuration

Create a configuration file. Because the project uses [viper](github.com/spf13/viper), TOML, JSON and YAML
are all supported. Environment variables are also able to be used.
We will assume TOML for the rest of the documentation.

The configuration file should be named `sqlboiler.toml` and is searched for in the following directories in this
order:

- `./`
- `$XDG_CONFIG_HOME/sqlboiler/`
- `$HOME/.config/sqlboiler/`

We require you pass in your `postgres` and `mysql` database configuration via the configuration file rather than env vars.
There is no command line argument support for database configuration. Values given under the `postgres` and `mysql`
block are passed directly to the postgres and mysql drivers. Here is a rundown of all the different
values that can go in that section:

| Name | Required | Postgres Default | MySQL Default |
| --- | --- | --- | --- |
| dbname  | yes       | none      | none   |
| host    | yes       | none      | none   |
| port    | no        | 5432      | 3306   |
| user    | yes       | none      | none   |
| pass    | no        | none      | none   |
| sslmode | no        | "require" | "true" |

You can also pass in these top level configuration values if you would prefer
not to pass them through the command line or environment variables:

| Name               | Defaults  |
| ------------------ | --------- |
| basedir            | none      |
| schema             | "public" *(or dbname for mysql)* |
| pkgname            | "models"  |
| output             | "models"  |
| whitelist          | []        |
| blacklist          | []        |
| tag                | []        |
| debug              | false     |
| no-hooks           | false     |
| no-tests           | false     |
| no-auto-timestamps | false     |

Example:

```toml
blacklist=["migrations", "other"]
schema="myschema"
[postgres]
  dbname="dbname"
  host="localhost"
  port=5432
  user="dbusername"
  pass="dbpassword"
[mysql]
  dbname="dbname"
  host="localhost"
  port=3306
  user="dbusername"
  pass="dbpassword"
  sslmode="false"
[mssql]
  dbname="dbname"
  host="localhost"
  port=1433
  user="dbusername"
  pass="dbpassword"
  sslmode="disable"
```

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
      --basedir string          The base directory has the templates and templates_test folders
  -b, --blacklist stringSlice   Do not include these tables in your generated package
  -d, --debug                   Debug mode prints stack traces on error
      --no-auto-timestamps      Disable automatic timestamps for created_at/updated_at
      --no-hooks                Disable hooks feature for your models
      --no-tests                Disable generated go test files
  -o, --output string           The name of the folder to output to (default "models")
  -p, --pkgname string          The name you wish to assign to your generated package (default "models")
  -s, --schema string           The name of your database schema, for databases that support real schemas (default "public")
  -t, --tag stringSlice         Struct tags to be included on your models in addition to json, yaml, toml
      --version                 Print the version
  -w, --whitelist stringSlice   Only include these tables in your generated package
```

Follow the steps below to do some basic model generation. Once you've generated
your models, you can run the compatibility tests which will exercise the entirety
of the generated code. This way you can ensure that your database is compatible
with SQLBoiler. If you find there are some failing tests, please check the
[Diagnosing Problems](#diagnosing-problems) section.

```sh
# Generate our models and exclude the migrations table
sqlboiler -b goose_migrations postgres

# Run the generated tests
go test ./models
```

*Note: No `mysqldump` or `pg_dump` equivalent for Microsoft SQL Server, so generated tests must be supplemented by `tables_schema.sql` with `CREATE TABLE ...` queries*


You can use `go generate` for SQLBoiler if you want to to make it easy to
run the command.

It's important to not modify anything in the output folder, which brings us to
the next topic: regeneration.

#### Regeneration

When regenerating the models it's recommended that you completely delete the
generated directory in a build script or use the `--wipe` flag in SQLBoiler.
The reasons for this are that sqlboiler doesn't try to diff your files in any
smart way, it simply writes the files it's going to write whether they're there
or not and doesn't delete any files that were added by you or previous runs of
SQLBoiler. In the best case this can cause compilation errors, in the worst case
this may leave extraneous and unusable code that was generated against tables
that are no longer in the database.

The bottom line is that this tool should always produce the same result from
the same source. And the intention is to always regenerate from a pure state.
The only reason the `--wipe` flag isn't defaulted to on is because we don't
like programs that `rm -rf` things on the filesystem without being asked to.

#### Extending generated models

There will probably come a time when you want to extend the generated models
with some kinds of helper functions. A general guideline is to put your
extension functions into a separate package so that your functions aren't
accidentally deleted when regenerating. Past that there are 3 main ways to
extend the models, the first way is the most desirable:

**Method 1: Simple Functions**

```go
// Package modext is for SQLBoiler helper methods
package modext

// UserFirstTimeSetup is an extension of the user model.
func UserFirstTimeSetup(db *sql.DB, u *models.User) error { ... }
```

Code organization is accomplished by using multiple files, and everything
is passed as a parameter so these kinds of methods are very easy to test.

Calling code is also very straightforward:

```go
user, err := Users(db).One()
// elided error check

err = modext.UserFirstTimeSetup(db, user)
// elided error check
```

**Method 2: Empty struct methods**

The above is the best way to code extensions for SQLBoiler, however there may
be times when the number of methods grows too large and code completion is
not as helpful anymore. In these cases you may consider structuring the code
like this:

```go
// Package modext is for SQLBoiler helper methods
package modext

type users struct {}

var Users = users{}

// FirstTimeSetup is an extension of the user model.
func (u users) FirstTimeSetup(db *sql.DB, u *models.User) error { ... }
```

Calling code then looks a little bit different:

```go
user, err := Users(db).One()
// elided error check

err = modext.Users.FirstTimeSetup(db, user)
// elided error check
```

This is almost identical to the method above, but gives slight amounts more
organization at virtually no cost at runtime. It is however not as desirable
as the first method since it does have some runtime cost and doesn't offer that
much benefit over it.

**Method 3: Embedding**

This pattern is not for the faint of heart, what it provides in benefits it
more than makes up for in downsides. It's possible to embed the SQLBoiler
structs inside your own to enhance them. However it's subject to easy breakages
and a dependency on these additional objects. It can also introduce
inconsistencies as some objects may have no extended functionality and therefore
have no reason to be embedded so you either have to have a struct for each
generated struct even if it's empty, or have inconsistencies, some places where
you use the enhanced model, and some where you do not.

```go
user, err := Users(db).One()
// elided error check

enhUser := modext.User{user}
err = ehnUser.FirstTimeSetup(db)
// elided error check
```

I don't recommend this pattern, but included it so that people know it's an
option and also know the problems with it.

## Diagnosing Problems

The most common causes of problems and panics are:

- Forgetting to exclude tables you do not want included in your generation, like migration tables.
- Tables without a primary key. All tables require one.
- Forgetting to put foreign key constraints on your columns that reference other tables.
- The compatibility tests require privileges to create a database for testing purposes, ensure the user
  supplied in your `sqlboiler.toml` config has adequate privileges.
- A nil or closed database handle. Ensure your passed in `boil.Executor` is not nil.
  - If you decide to use the `G` variant of functions instead, make sure you've initialized your
    global database handle using `boil.SetDB()`.

For errors with other causes, it may be simple to debug yourself by looking at the generated code.
Setting `boil.DebugMode` to `true` can help with this. You can change the output using `boil.DebugWriter` (defaults to `os.Stdout`).

If you're still stuck and/or you think you've found a bug, feel free to leave an issue and we'll do our best to help you.

## Features & Examples

Most examples in this section will be demonstrated using the following Postgres schema, structs and variables:

```sql
CREATE TABLE pilots (
  id integer NOT NULL,
  name text NOT NULL
);

ALTER TABLE pilots ADD CONSTRAINT pilot_pkey PRIMARY KEY (id);

CREATE TABLE jets (
  id integer NOT NULL,
  pilot_id integer NOT NULL,
  age integer NOT NULL,
  name text NOT NULL,
  color text NOT NULL
);

ALTER TABLE jets ADD CONSTRAINT jet_pkey PRIMARY KEY (id);
ALTER TABLE jets ADD CONSTRAINT jet_pilots_fkey FOREIGN KEY (pilot_id) REFERENCES pilots(id);

CREATE TABLE languages (
  id integer NOT NULL,
  language text NOT NULL
);

ALTER TABLE languages ADD CONSTRAINT language_pkey PRIMARY KEY (id);

-- Join table
CREATE TABLE pilot_languages (
  pilot_id integer NOT NULL,
  language_id integer NOT NULL
);

-- Composite primary key
ALTER TABLE pilot_languages ADD CONSTRAINT pilot_language_pkey PRIMARY KEY (pilot_id, language_id);
ALTER TABLE pilot_languages ADD CONSTRAINT pilot_language_pilots_fkey FOREIGN KEY (pilot_id) REFERENCES pilots(id);
ALTER TABLE pilot_languages ADD CONSTRAINT pilot_language_languages_fkey FOREIGN KEY (language_id) REFERENCES languages(id);
```

The generated model structs for this schema look like the following. Note that we've included the relationship
structs as well so you can see how it all pieces together:

```go
type Pilot struct {
  ID   int    `boil:"id" json:"id" toml:"id" yaml:"id"`
  Name string `boil:"name" json:"name" toml:"name" yaml:"name"`

  R *pilotR `boil:"-" json:"-" toml:"-" yaml:"-"`
  L pilotR  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

type pilotR struct {
  Licenses  LicenseSlice
  Languages LanguageSlice
  Jets      JetSlice
}

type Jet struct {
  ID      int    `boil:"id" json:"id" toml:"id" yaml:"id"`
  PilotID int    `boil:"pilot_id" json:"pilot_id" toml:"pilot_id" yaml:"pilot_id"`
  Age     int    `boil:"age" json:"age" toml:"age" yaml:"age"`
  Name    string `boil:"name" json:"name" toml:"name" yaml:"name"`
  Color   string `boil:"color" json:"color" toml:"color" yaml:"color"`

  R *jetR `boil:"-" json:"-" toml:"-" yaml:"-"`
  L jetR  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

type jetR struct {
  Pilot *Pilot
}

type Language struct {
  ID       int    `boil:"id" json:"id" toml:"id" yaml:"id"`
  Language string `boil:"language" json:"language" toml:"language" yaml:"language"`

  R *languageR `boil:"-" json:"-" toml:"-" yaml:"-"`
  L languageR  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

type languageR struct {
  Pilots PilotSlice
}
```

```go
// Open handle to database like normal
db, err := sql.Open("postgres", "dbname=fun user=abc")
if err != nil {
  return err
}
```

### Automatic CreatedAt/UpdatedAt

If your generated SQLBoiler models package can find columns with the
names `created_at` or `updated_at` it will automatically set them
to `time.Now()` in your database, and update your object appropriately.
To disable this feature use `--no-auto-timestamps`.

Note: You can set the timezone for this feature by calling `boil.SetLocation()`

#### Overriding Automatic Timestamps

* **Insert**
  * Timestamps for both `updated_at` and `created_at` that are zero values will be set automatically.
  * To set the timestamp to null, set `Valid` to false and `Time` to a non-zero value.
  This is somewhat of a work around until we can devise a better solution in a later version.
* **Update**
  * The `updated_at` column will always be set to `time.Now()`. If you need to override
  this value you will need to fall back to another method in the meantime: `queries.Raw()`,
  overriding `updated_at` in all of your objects using a hook, or create your own wrapper.
* **Upsert**
  * `created_at` will be set automatically if it is a zero value, otherwise your supplied value
  will be used. To set `created_at` to `null`, set `Valid` to false and `Time` to a non-zero value.
  * The `updated_at` column will always be set to `time.Now()`.

### Query Building

We generate "Starter" methods for you. These methods are named as the plural versions of your model,
for example: `models.Jets()`. Starter methods are used to build queries using our
[Query Mod System](#query-mod-system). They take a slice of [Query Mods](#query-mod-system)
as parameters, and end with a call to a [Finisher](#finishers) method.

Here are a few examples:

```go
// SELECT COUNT(*) FROM pilots;
count, err := models.Pilots().Count()

// SELECT * FROM "pilots" LIMIT 5;
pilots, err := models.Pilots(qm.Limit(5)).All()

// DELETE FROM "pilots" WHERE "id"=$1;
err := models.Pilots(qm.Where("id=?", 1)).DeleteAll()
```

In the event that you would like to build a query and specify the table yourself, you
can do so using `models.NewQuery()`:

```go
// Select all rows from the pilots table by using the From query mod.
err := models.NewQuery(db, From("pilots")).All()
```

As you can see, [Query Mods](#query-mods) allow you to modify your queries, and [Finishers](#finishers)
allow you to execute the final action.

We also generate query building helper methods for your relationships as well. Take a look at our
[Relationships Query Building](#relationships) section for some additional query building information.


### Query Mod System

The query mod system allows you to modify queries created with [Starter](#query-building) methods
when performing query building. Here is a list of all of your generated query mods using examples:

```go
// Dot import so we can access query mods directly instead of prefixing with "qm."
import . "github.com/vattle/sqlboiler/queries/qm"

// Use a raw query against a generated struct (Pilot in this example)
// If this query mod exists in your call, it will override the others.
// "?" placeholders are not supported here, use "$1, $2" etc.
SQL("select * from pilots where id=$1", 10)
models.Pilots(SQL("select * from pilots where id=$1", 10)).All()

Select("id", "name") // Select specific columns.
From("pilots as p") // Specify the FROM table manually, can be useful for doing complex queries.

// WHERE clause building
Where("name=?", "John")
And("age=?", 24)
Or("height=?", 183)

// WHERE IN clause building
WhereIn("name, age in ?", "John" 24, "Tim", 33) // Generates: WHERE ("name","age") IN (($1,$2),($3,$4))
AndIn("weight in ?", 84)
OrIn("height in ?", 183, 177, 204)

InnerJoin("pilots p on jets.pilot_id=?", 10)

GroupBy("name")
OrderBy("age, height")

Having("count(jets) > 2")

Limit(15)
Offset(5)

// Explicit locking
For("update nowait")

// Eager Loading -- Load takes the relationship name, ie the struct field name of the
// Relationship struct field you want to load.
Load("Languages") // If it's a ToOne relationship it's in singular form, ToMany is plural.
```

Note: We don't force you to break queries apart like this if you don't want to, the following
is also valid and supported by query mods that take a clause:

```go
Where("(name=? OR age=?) AND height=?", "John", 24, 183)
```

### Function Variations

You will find that most functions have the following variations. We've used the
```Delete``` method to demonstrate:

```go
// Set the global db handle for G method variants.
boil.SetDB(db)

pilot, _ := models.FindPilot(db, 1)

err := pilot.Delete(db) // Regular variant, takes a db handle (boil.Executor interface).
pilot.DeleteP(db)       // Panic variant, takes a db handle and panics on error.
err := pilot.DeleteG()  // Global variant, uses the globally set db handle (boil.SetDB()).
pilot.DeleteGP()        // Global&Panic variant, combines the global db handle and panic on error.

db.Begin()              // Normal sql package way of creating a transaction
boil.Begin()            // Uses the global database handle set by boil.SetDB()
```

Note that it's slightly different for query building.

### Finishers

Here are a list of all of the finishers that can be used in combination with
[Query Building](#query-building).

Finishers all have `P` (panic) [method variations](#function-variations). To specify
your db handle use the `G` or regular variation of the [Starter](#query-building) method.

```go
// These are called like the following:
models.Pilots(db).All()

One() // Retrieve one row as object (same as LIMIT(1))
All() // Retrieve all rows as objects (same as SELECT * FROM)
Count() // Number of rows (same as COUNT(*))
UpdateAll(models.M{"name": "John", "age": 23}) // Update all rows matching the built query.
DeleteAll() // Delete all rows matching the built query.
Exists() // Returns a bool indicating whether the row(s) for the built query exists.
Bind(&myObj) // Bind the results of a query to your own struct object.
Exec() // Execute an SQL query that does not require any rows returned.
QueryRow() // Execute an SQL query expected to return only a single row.
Query() // Execute an SQL query expected to return multiple rows.
```

### Raw Query

We provide `queries.Raw()` for executing raw queries. Generally you will want to use `Bind()` with
this, like the following:

```go
err := queries.Raw(db, "select * from pilots where id=$1", 5).Bind(&obj)
```

You can use your own structs or a generated struct as a parameter to Bind. Bind supports both
a single object for single row queries and a slice of objects for multiple row queries.

`queries.Raw()` also has a method that can execute a query without binding to an object, if required.

You also have `models.NewQuery()` at your disposal if you would still like to use [Query Building](#query-building)
in combination with your own custom, non-generated model.

### Binding

For a comprehensive ruleset for `Bind()` you can refer to our [godoc](https://godoc.org/github.com/vattle/sqlboiler/queries#Bind).

The `Bind()` [Finisher](#finisher) allows the results of a query built with
the [Raw SQL](#raw-query) method or the [Query Builder](#query-building) methods to be bound
to your generated struct objects, or your own custom struct objects.

This can be useful for complex queries, queries that only require a small subset of data
and have no need for the rest of the object variables, or custom join struct objects like
the following:

```go
// Custom struct using two generated structs
type PilotAndJet struct {
  models.Pilot `boil:",bind"`
  models.Jet   `boil:",bind"`
}

var paj PilotAndJet
// Use a raw query
err := queries.Raw(`
  select pilots.id as "pilots.id", pilots.name as "pilots.name",
  jets.id as "jets.id", jets.pilot_id as "jets.pilot_id",
  jets.age as "jets.age", jets.name as "jets.name", jets.color as "jets.color"
  from pilots inner join jets on jets.pilot_id=?`, 23,
).Bind(&paj)

// Use query building
err := models.NewQuery(db,
  Select("pilots.id", "pilots.name", "jets.id", "jets.pilot_id", "jets.age", "jets.name", "jets.color"),
  From("pilots"),
  InnerJoin("jets on jets.pilot_id = pilots.id"),
).Bind(&paj)
```

```go
// Custom struct for selecting a subset of data
type JetInfo struct {
  AgeSum int `boil:"age_sum"`
  Count int `boil:"juicy_count"`
}

var info JetInfo

// Use query building
err := models.NewQuery(db, Select("sum(age) as age_sum", "count(*) as juicy_count", From("jets"))).Bind(&info)

// Use a raw query
err := queries.Raw(`select sum(age) as "age_sum", count(*) as "juicy_count" from jets`).Bind(&info)
```

We support the following struct tag modes for `Bind()` control:

```go
type CoolObject struct {
  // Don't specify a name, Bind will TitleCase the column
  // name, and try to match against this.
  Frog int

  // Specify an alternative name for the column, it will
  // be titlecased for matching, can be whatever you like.
  Cat int  `boil:"kitten"`

  // Ignore this struct field, do not attempt to bind it.
  Pig int  `boil:"-"`

  // Instead of binding to this as a regular struct field
  // (like other sql-able structs eg. time.Time)
  // Recursively search inside the Dog struct for field names from the query.
  Dog      `boil:",bind"`

  // Same as the above, except specify a different table name
  Mouse    `boil:"rodent,bind"`

  // Ignore this struct field, do not attempt to bind it.
  Bird     `boil:"-"`
}
```

### Relationships

Helper methods will be generated for every to one and to many relationship structure
you have defined in your database by using foreign keys.

We attach these helpers directly to your model struct, for example:

```go
jet, _ := models.FindJet(db, 1)

// "to one" relationship helper method.
// This will retrieve the pilot for the jet.
pilot, err := jet.Pilot(db).One()

// "to many" relationship helper method.
// This will retrieve all languages for the pilot.
languages, err := pilot.Languages(db).All()
```

If your relationship involves a join table SQLBoiler will figure it out for you transparently.

It is important to note that you should use `Eager Loading` if you plan
on loading large collections of rows, to avoid N+1 performance problems.

For example, take the following:

```go
// Avoid this loop query pattern, it is slow.
jets, _ := models.Jets(db).All()
pilots := make([]models.Pilot, len(jets))
for i := 0; i < len(jets); i++ {
  pilots[i] = jets[i].Pilot(db).OneP()
}

// Instead, use Eager Loading!
jets, _ := models.Jets(db, Load("Pilot")).All()
```

Eager loading can be combined with other query mods, and it can also eager load recursively.

```go
// Example of a nested load.
// Each jet will have its pilot loaded, and each pilot will have its languages loaded.
jets, _ := models.Jets(db, Load("Pilot.Languages")).All()
// Note that each level of a nested Load call will be loaded. No need to call Load() multiple times.

// A larger, random example
users, _ := models.Users(db,
  Load("Pets.Vets"),
  Load("Pets.Toys"),
  Load("Property"),
  Where("age > ?", 23),
).All()
```

We provide the following methods for managing relationships on objects:

**To One**
- `SetX()`: Set the foreign key to point to something else: jet.SetPilot(...)
- `RemoveX()`: Null out the foreign key, effectively removing the relationship between these two objects: jet.RemovePilot(...)

**To Many**
- `AddX()`: Add more relationships to the existing set of related Xs: pilot.AddLanguages(...)
- `SetX()`: Remove all existing relationships, and replace them with the provided set: pilot.SetLanguages(...)
- `RemoveX()`: Remove all provided relationships: pilot.RemoveLanguages(...)

**To One** code examples:

```go
  jet, _ := models.FindJet(db, 1)
  pilot, _ := models.FindPilot(db, 1)

  // Set the pilot to an existing pilot
  err := jet.SetPilot(db, false, &pilot)

  pilot = models.Pilot{
    Name: "Erlich",
  }

  // Insert the pilot into the database and assign it to a jet
  err := jet.SetPilot(db, true, &pilot)

  // Remove a relationship. This method only exists for foreign keys that can be NULL.
  err := jet.RemovePilot(db, &pilot)
```

**To Many** code examples:

```go
  pilots, _ := models.Pilots(db).All()
  languages, _ := models.Languages(db).All()

  // Set a group of language relationships
  err := pilots.SetLanguages(db, false, &languages)

  languages := []*models.Language{
    {Language: "Strayan"},
    {Language: "Yupik"},
    {Language: "Pawnee"},
  }

  // Insert new a group of languages and assign them to a pilot
  err := pilots.SetLanguages(db, true, languages...)

  // Add another language relationship to the existing set of relationships
  err := pilots.AddLanguages(db, false, &someOtherLanguage)

  anotherLanguage := models.Language{Language: "Archi"}

  // Insert and then add another language relationship
  err := pilots.AddLanguages(db, true, &anotherLanguage)

  // Remove a group of relationships
  err := pilots.RemoveLanguages(db, languages...)
```

### Hooks

Before and After hooks are available for most operations. If you don't need them you can
shrink the size of the generated code by disabling them with the `--no-hooks` flag.

Every generated package that includes hooks has the following `HookPoints` defined:

```go
const (
  BeforeInsertHook HookPoint = iota + 1
  BeforeUpdateHook
  BeforeDeleteHook
  BeforeUpsertHook
  AfterInsertHook
  AfterSelectHook
  AfterUpdateHook
  AfterDeleteHook
  AfterUpsertHook
)
```

To register a hook for your model you will need to create the hook function, and attach
it with the `AddModelHook` method. Here is an example of a before insert hook:

```go
// Define my hook function
func myHook(exec boil.Executor, p *Pilot) {
  // Do stuff
}

// Register my before insert hook for pilots
models.AddPilotHook(boil.BeforeInsertHook, myHook)
```

Your `ModelHook` will always be defined as `func(boil.Executor, *Model)`

### Transactions

The boil.Executor interface powers all of SQLBoiler. This means anything that conforms
to the three `Exec/Query/QueryRow` methods can be used. `sql.DB`, `sql.Tx` as well as other
libraries (`sqlx`) conform to this interface, and therefore any of these things may be
used as an executor for any query in the system. This makes using transactions very simple:


```go
tx, err := db.Begin()
if err != nil {
  return err
}

users, _ := models.Pilots(tx).All()
users.DeleteAll(tx)

// Rollback or commit
tx.Commit()
tx.Rollback()
```

It's also worth noting that there's a way to take advantage of `boil.SetDB()`
by using the [boil.Begin()](https://godoc.org/github.com/vattle/sqlboiler/boil#Begin) function.
This opens a transaction using the globally stored database.

### Debug Logging

Debug logging will print your generated SQL statement and the arguments it is using.
Debug logging can be toggled on globally by setting the following global variable to `true`:

```go
boil.DebugMode = true

// Optionally set the writer as well. Defaults to os.Stdout
fh, _ := os.Open("debug.txt")
boil.DebugWriter = fh
```

Note: Debug output is messy at the moment. This is something we would like addressed.

### Select

Select is done through [Query Building](#query-building) and [Find](#find). Here's a short example:

```go
// Select one pilot
pilot, err := models.Pilots(db, qm.Where("name=?", "Tim")).One()

// Select specific columns of many jets
jets, err := models.Jets(db, qm.Select("age", "name")).All()
```

### Find

Find is used to find a single row by primary key:

```go
// Retrieve pilot with all columns filled
pilot, err := models.PilotFind(db, 1)

// Retrieve a subset of column values
jet, err := models.JetFind(db, 1, "name", "color")
```

### Insert

The main thing to be aware of with `Insert` is how the `whitelist` operates. If no whitelist
argument is provided, `Insert` will abide by the following rules:

- Insert all columns **without** a database default value.
- Insert all columns with a non-zero value that have a database default value.

On the other hand, if a whitelist is provided, we will only insert the columns specified in the whitelist.

Also note that your object will automatically be updated with any missing default values from the
database after the `Insert` is finished executing. This includes auto-incrementing column values.

```go
var p1 models.Pilot
p1.Name = "Larry"
err := p1.Insert(db) // Insert the first pilot with name "Larry"
// p1 now has an ID field set to 1

var p2 models.Pilot
p2.Name "Boris"
err := p2.Insert(db) // Insert the second pilot with name "Boris"
// p2 now has an ID field set to 2

var p3 models.Pilot
p3.ID = 25
p3.Name = "Rupert"
err := p3.Insert(db) // Insert the third pilot with a specific ID
// The id for this row was inserted as 25 in the database.

var p4 models.Pilot
p4.ID = 0
p4.Name = "Nigel"
err := p4.Insert(db, "id", "name") // Insert the fourth pilot with a zero value ID
// The id for this row was inserted as 0 in the database.
// Note: We had to use the whitelist for this, otherwise
// SQLBoiler would presume you wanted to auto-increment
```

### Update
`Update` can be performed on a single object, a slice of objects or as a [Finisher](#finishers)
for a collection of rows.

`Update` on a single object optionally takes a `whitelist`. The purpose of the
whitelist is to specify which columns in your object should be updated in the database.

If no `whitelist` argument is provided, `Update` will update every column except for
`primary key` columns.

If a `whitelist` argument is provided, `update` will only update the columns specified.

```go
// Find a pilot and update his name
pilot, _ := models.FindPilot(db, 1)
pilot.Name = "Neo"
err := pilot.Update(db)

// Update a slice of pilots to have the name "Smith"
pilots, _ := models.Pilots(db).All()
err := pilots.UpdateAll(db, models.M{"name": "Smith"})

// Update all pilots in the database to to have the name "Smith"
err := models.Pilots(db).UpdateAll(models.M{"name", "Smith"})
```

### Delete

Delete a single object, a slice of objects or specific objects through [Query Building](#query-building).

```go
pilot, _ := models.FindPilot(db, 1)
// Delete the pilot from the database
err := pilot.Delete(db)

// Delete all pilots from the database
err := models.Pilots(db).DeleteAll()

// Delete a slice of pilots from the database
pilots, _ := models.Pilots(db).All()
err := pilots.DeleteAll(db)
```

### Upsert

[Upsert](https://www.postgresql.org/docs/9.5/static/sql-insert.html) allows you to perform an insert
that optionally performs an update when a conflict is found against existing row values.

The `whitelist` operates in the same fashion that it does for [Insert](#insert).

If an insert is performed, your object will be updated with any missing default values from the database,
such as auto-incrementing column values.

```go
var p1 models.Pilot
p1.ID = 5
p1.Name = "Gaben"

// INSERT INTO pilots ("id", "name") VALUES($1, $2)
// ON CONFLICT DO NOTHING
err := p1.Upsert(db, false, nil, nil)

// INSERT INTO pilots ("id", "name") VALUES ($1, $2)
// ON CONFLICT ("id") DO UPDATE SET "name" = EXCLUDED."name"
err := p1.Upsert(db, true, []string{"id"}, []string{"name"})

// Set p1.ID to a zero value. We will have to use the whitelist now.
p1.ID = 0
p1.Name = "Hogan"

// INSERT INTO pilots ("id", "name") VALUES ($1, $2)
// ON CONFLICT ("id") DO UPDATE SET "name" = EXCLUDED."name"
err := p1.Upsert(db, true, []string{"id"}, []string{"name"}, "id", "name")
```

The `updateOnConflict` argument allows you to specify whether you would like Postgres
to perform a `DO NOTHING` on conflict, opposed to a `DO UPDATE`. For MySQL, this param will not be generated.

The `conflictColumns` argument allows you to specify the `ON CONFLICT` columns for Postgres.
For MySQL, this param will not be generated.

Note: Passing a different set of column values to the update component is not currently supported.

### Reload
In the event that your objects get out of sync with the database for whatever reason,
you can use `Reload` and `ReloadAll` to reload the objects using the primary key values
attached to the objects.

```go
pilot, _ := models.FindPilot(db, 1)

// > Object becomes out of sync for some reason, perhaps async processing

// Refresh the object with the latest data from the db
err := pilot.Reload(db)

// Reload all objects in a slice
pilots, _ := models.Pilots(db).All()
err := pilots.ReloadAll(db)
```

Note: `Reload` and `ReloadAll` are not recursive, if you need your relationships reloaded
you will need to call the `Reload` methods on those yourself.

### Exists

```go
jet, err := models.FindJet(db, 1)

// Check if the pilot assigned to this jet exists.
exists, err := jet.Pilot(db).Exists()

// Check if the pilot with ID 5 exists
exists, err := models.Pilots(db, Where("id=?", 5)).Exists()
```

### Enums

If your MySQL or Postgres tables use enums we will generate constants that hold their values
that you can use in your queries. For example:

```
CREATE TYPE workday AS ENUM('monday', 'tuesday', 'wednesday', 'thursday', 'friday');

CREATE TABLE event_one (
  id     serial PRIMARY KEY NOT NULL,
  name   VARCHAR(255),
  day    workday NOT NULL
);
```

An enum type defined like the above, being used by a table, will generate the following enums:

```go
const (
  WorkdayMonday    = "monday"
  WorkdayTuesday   = "tuesday"
  WorkdayWednesday = "wednesday"
  WorkdayThursday  = "thursday"
  WorkdayFriday    = "friday"
)
```

For Postgres we use `enum type name + title cased` value to generate the const variable name.
For MySQL we use `table name + column name + title cased value` to generate the const variable name.

Note: If your enum holds a value we cannot parse correctly due, to non-alphabet characters for example,
it may not be generated. In this event, you will receive errors in your generated tests because
the value randomizer in the test suite does not know how to generate valid enum values. You will
still be able to use your generated library, and it will still work as expected, but the only way
to get the tests to pass in this event is to either use a parsable enum value or use a regular column
instead of an enum.

## FAQ

#### Won't compiling models for a huge database be very slow?

No, because Go's toolchain - unlike traditional toolchains - makes the compiler do most of the work
instead of the linker. This means that when the first `go install` is done it can take
a little bit of time because there is a lot of code that is generated. However, because of this
work balance between the compiler and linker in Go, linking to that code afterwards in the subsequent
compiles is extremely fast.

#### Missing imports for generated package

The generated models might import a couple of packages that are not on your system already, so
`cd` into your generated models directory and type `go get -u -t` to fetch them. You will only need
to run this command once, not per generation.

#### How should I handle multiple schemas?

If your database uses multiple schemas you should generate a new package for each of your schemas.
Note that this only applies to databases that use real, SQL standard schemas (like PostgreSQL), not
fake schemas (like MySQL).

#### How do I use types.BytesArray for Postgres bytea arrays?

Only "escaped format" is supported for types.BytesArray. This means that your byte slice needs to have
a format of "\\x00" (4 bytes per byte) opposed to "\x00" (1 byte per byte). This is to maintain compatibility
with all Postgres drivers. Example:

`x := types.BytesArray{0: []byte("\\x68\\x69")}`

Please note that multi-dimensional Postgres ARRAY types are not supported at this time.

#### Why aren't my time.Time or null.Time fields working in MySQL?

You *must* use a DSN flag in MySQL connections, see: [Requirements](#requirements)

#### Where is the homepage?

The homepage for the [SQLBoiler](https://github.com/vattle/sqlboiler) [Golang ORM](https://github.com/vattle/sqlboiler)
generator is located at: https://github.com/vattle/sqlboiler

## Benchmarks

If you'd like to run the benchmarks yourself check out our [boilbench](https://github.com/vattle/boilbench) repo.

Here are the results (lower is better):

`go test -bench . -benchmem`
```
BenchmarkGORMDelete/gorm-8     	    100000	     15364 ns/op	    5395 B/op	     113 allocs/op
BenchmarkGORPDelete/gorp-8     	   1000000	      1703 ns/op	     304 B/op	      12 allocs/op
BenchmarkXORMDelete/xorm-8     	    100000	     14733 ns/op	    3634 B/op	     107 allocs/op
BenchmarkBoilDelete/boil-8     	   2000000	       986 ns/op	     120 B/op	       7 allocs/op

BenchmarkGORMInsert/gorm-8     	    100000	     19197 ns/op	    8054 B/op	     161 allocs/op
BenchmarkGORPInsert/gorp-8     	    500000	      3413 ns/op	    1008 B/op	      32 allocs/op
BenchmarkXORMInsert/xorm-8     	    100000	     15428 ns/op	    5836 B/op	     131 allocs/op
BenchmarkBoilInsert/boil-8     	    500000	      3041 ns/op	     568 B/op	      21 allocs/op

BenchmarkGORMSelectAll/gorm-8  	     20000	     85422 ns/op	   29912 B/op	     511 allocs/op
BenchmarkGORPSelectAll/gorp-8  	     50000	     35824 ns/op	    8837 B/op	     312 allocs/op
BenchmarkXORMSelectAll/xorm-8  	     30000	     58843 ns/op	   13805 B/op	     298 allocs/op
BenchmarkBoilSelectAll/boil-8  	    100000	     13844 ns/op	    2840 B/op	      61 allocs/op

BenchmarkGORMSelectSubset/gorm-8     10000	    100714 ns/op	   30875 B/op	     517 allocs/op
BenchmarkGORPSelectSubset/gorp-8     30000	     43547 ns/op	    8837 B/op	     312 allocs/op
BenchmarkXORMSelectSubset/xorm-8     30000	     48128 ns/op	   12989 B/op	     282 allocs/op
BenchmarkBoilSelectSubset/boil-8    100000	     12316 ns/op	    2977 B/op	      65 allocs/op

BenchmarkGORMSelectComplex/gorm-8    10000	    133598 ns/op	   49398 B/op	     772 allocs/op
BenchmarkGORPSelectComplex/gorp-8    50000	     40588 ns/op	    9037 B/op	     321 allocs/op
BenchmarkXORMSelectComplex/xorm-8    30000	     56367 ns/op	   14174 B/op	     313 allocs/op
BenchmarkBoilSelectComplex/boil-8   100000	     16941 ns/op	    3821 B/op	      95 allocs/op

BenchmarkGORMUpdate/gorm-8           50000	     25406 ns/op	    9710 B/op	     195 allocs/op
BenchmarkGORPUpdate/gorp-8          300000	      3614 ns/op	    1152 B/op	      34 allocs/op
BenchmarkXORMUpdate/xorm-8          100000	     17510 ns/op	    4458 B/op	     132 allocs/op
BenchmarkBoilUpdate/boil-8          500000	      2958 ns/op	     520 B/op	      16 allocs/op

BenchmarkGORMRawBind/gorm-8    	     10000	    112577 ns/op	   38270 B/op	     595 allocs/op
BenchmarkGORPRawBind/gorp-8    	     30000	     40967 ns/op	    8837 B/op	     312 allocs/op
BenchmarkXORMRawBind/xorm-8    	     30000	     54739 ns/op	   12692 B/op	     273 allocs/op
BenchmarkSQLXRawBind/sqlx-8    	    200000	     13537 ns/op	    4268 B/op	      49 allocs/op
BenchmarkBoilRawBind/boil-8    	    200000	     11144 ns/op	    4334 B/op	      49 allocs/op
```

<img style="margin-right:6px;" src="http://i.imgur.com/TglZGoI.png"/>
<img style="margin-right:6px;" src="http://i.imgur.com/Ktm2ta4.png"/>
<img style="margin-right:6px;" src="http://i.imgur.com/yv8kFPA.png"/>
<img style="margin-right:6px;" src="http://i.imgur.com/890Zswe.png"/>
<img style="margin-right:6px;" src="http://i.imgur.com/qMgoAFJ.png"/>
<img style="margin-right:6px;" src="http://i.imgur.com/sDoNiCN.png"/>
<img style="margin-right:6px;" src="http://i.imgur.com/EvUa4UT.png"/>
