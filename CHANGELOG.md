# Changelog

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Add ability to rename tables, columns, and relationships.
- Add support for rows affected to Update/Delete calls.
- Add support for geotypes for psql (thanks @saulortega)
- Add a flag to set the config file (thanks @l4u)
- Add support for citext to psql (thanks @boxofrad)
- Add virtual columns ignoring for mysql (thanks @Nykakin)
- Add ability for drivers to provide their own templates/replacement templates
- Add ability for drivers to specify imports
- Add a boil.sh that includes many commands that help build and test sqlboiler
- Add many more driver dialect flags to be able to remove all the DriverName
  comparisons inside sqlboiler templates. This allows us to more cleanly support
  more drivers.
- Add ability to override imports via config file
- Add ability to replace types using the config file
- Add ability to re-use queries with SetArgs
- Add ability to specify blacklisted or whitelisted columns in a table by
  using the syntax tablename.columnname in the driver's whitelist/blacklist
  config setting.
- Add way to create a relationship struct: `modelName.R.NewStruct()`
- MySQL numeric, fixed, dec types now generate using the new types.(Null)Decimal
- Use bindata as the default method of accessing templates, this prevents many
  bug reports we've had in the past. During development or to otherwise opt out
  the --templates flag will not load from bindata (drivers are the exception).
- Add the ability to generate non-go files. The new --templates flag allows you to override
  the default bindata templates/ and templates_test/ directory that ship with sqlboiler.
- Export the queries.BuildQuery method for public use. This allows building
  queries without executing.

### Changed

#### Driver split

Drivers are now separate binaries. A lot of the reason for this is because
having to keep all the drivers together inside sqlboiler was going to later
cause more churn that it was worth. sqlite showed us this immediately with
it's cgo dependencies as an example. This was the source of a huge number
of changes, many of them breaking but other than the new workflow of using
a new binary instead of just a string, users shouldn't actually see too much
difference.

#### Smaller changes

- Insert, Update and Upsert now take a `boil.Columns` instead of a
  `whitelist []string`. This allows more flexibility in that we can use
  different column list kinds to control what gets inserted and updated.
  See `boil.Infer`, `boil.Whitelist`, `boil.Blacklist`, `boil.Greylist`.
  This is a breaking change to the syntax of Insert and Update as this argument
  is no longer optional. To migrate an app and keep the same behavior,
  add `boil.Infer()` where there's a missing argument, and add
  `boil.Whitelist("columns", "here")` where columns were previously specified.
- Eager loading can now accept QueryMods. This is a breaking change. To migrate
  an application simply take calls that were of the form: `Load("a.b", "a.c")`
  and break them into two separate `Load()` query mods. The second argument
  is now a variadic slice of query mods, hence the breaking change. See the docs
  for `Load()` for more details.
- Eager loading now attaches both sides of the relationship in the `R` structs.
  This is consistent with the way the set relationship helpers work.
- Queries no longer keep a handle to the database and therefore the db
  parameter is no longer passed in where you create a query, instead it's passed
  in where you execute the query.
- context.Context is now piped through everything, this is a breaking change
  syntactically. This can be controlled by the `--no-context` flag.
- Rename postgresql driver to psql
- MySQL driver no longer accepts `tinyint_as_bool` via command line. Instead
  the flag has been inverted (`tinyint_as_int`) to stop the automatic bool
  conversion and it's now passed in as a driver configuration setting (via the
  env or config).
- Schema is no longer an sqlboiler level configuration setting. Instead it is
  passed back to sqlboiler when the driver is asked to produce the information
  for the database we're generating for. This is a breaking change.
- Rows affected is now returned by update and deletes by default. This is
  syntactically a breaking change. This can be controlled by the
  `--no-rows-affected` flag.
- Panic and Global variants are now hidden behind flags. This reduces the size
  of the generated code by 15%. This is a breaking change if you were using
  them.
- Randomize now counts on an interface to do proper randomization for
  custom types.
- Imports is now it's own package and is cleaned up significantly
- Changed all `$dot` references to `$.` as is proper
- Drivers are now responsible for their own config validation
- Drivers are now responsible for their own default config values
- Upsert is no longer a core part of sqlboiler, and is currently being provided
  by all three supported drivers, but other drivers may choose not to implement
  it.
- Drivers can now take any configuration via environment or config files. This
  allows drivers to be much more configurable and there's no constraints
  imposed by the main/testing of sqlboiler itself.
- Replace templates (hidden flag) now use the `;` separator not `:` for windows
  compatibility.
- Postgres numeric, fixed, decimal types now generate using the new
  types.(Null)Decimal instead of float64, this is a breaking change,
  but a very correct and important one.

### Removed

- The concept of TestMain is now gone from both templates and imports. It's
  been superceded by the new driver abilities to supply templates and imports.
  The drivers add their mains to the TestSingleton templates.
- Remove the reliance on the null package in the templates. Instead favor
  using the sql.Scanner and driver.Valuer interface that most types will
  implement anyway (and the null package does too).

### Fixed

- Fixed a bug in Bind() where all the given `*sql.Rows` would be consumed
  even in the event where we were binding to a single object. This allows
  for slower scanning of a `*sql.Rows` object.
- Fixed a problem in eager loading where the same object would be queried
  for multiple times. This was simply inefficient and so we now de-duplicate
  the ids before creating the query.
- Fix a bug with insert where if the columns to insert were nil, values
  would never be loaded into the query (RETURNING clause missing).