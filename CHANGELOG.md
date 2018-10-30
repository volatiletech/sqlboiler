# Changelog

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [v3.1.0] - 2018-10-29

### Added

- Add extra text to clarify the conditions required for a transparent join table
  in the readme.
- Add extra development instructions to the CONTRIBUTING.md (thanks @gemscng)
- Add `UpdateAllG` function variant (thanks @gemscng)
- Add context to Find examples in README.md (thanks @jones77)

### Fixed

- Fix boil.sh did not go-generate sqlboiler when using `all`, it now does
  `sqlboiler` and all drivers as well.
- Fix a bug with MSSQL exists finisher, it now uses `*` instead of the schema
  name and a star (like the count query).
- Fix a panic in aliases code (thanks @nadilas)
- Fix bug when eager loading with null.Uint64 ids
- Fix dead links to drivers in README (thanks @DenLilleMand)
- Fix a problem with eager loading where null foreign keys would create bad IDs
  and cause general problems.
- Fix mysql bigint signed type to not use an unsigned int in Go (thanks @nazo)
- Fix an issue where generated tests failed when output directory was too
  far away from the config file (thanks @gedorinku).

## [v3.0.1] - 2018-08-24

### Fixed

- Fix a DSN formatting issue when connecting to sql server instances
- Fix missing BindG function (should have been there)
- Fix blacklist using the whitelist instead of blacklist in the mssql driver

## [v3.0.0] - 2018-08-13

### Added

- Add constant for relationship names, as well as a helper to use them in
  query mods: See `ModelNameRels` for the constants and `qm.Rels` for the helper
  (thanks @glerchundi).
- Add support for PSQL Identity columns (thanks @gencer)
- Add a new syntax for import maps in the config file. This allows us to
  sidestep viper's constant downcasing of config keys. This is the
  exact same fix as happened with aliases previously.

### Changed

- Change querymods to an interface much like the http.Handler interface. This
  is a breaking change but facilitates being able to actually test querymods
  as well as a more flexible method of being able to create them. For
  compatibility with older querymods you've created, use the `qm.QueryModFunc`
  to convert them to a function type that implements the interface, just like
  `http.HandlerFunc` (thanks @glerchundi)
- Change the `seed` parameter to a `func () int64` function that when called
  safely produces a new thread-safe pseudorandom sequential integer. This was
  done to remove the dependency on `sqlboiler` code.
- Stop using gopkg.in for versioning for the null package. Please ensure you
  have null package v8.0.0 or higher after this change.
- MySQL driver was erroneously using time.Time for the `time` type but the
  most prolific driver does not support this, use string instead. This change
  was PR'd to v2 but never to v3. (thanks @ceshihao)
- Remove satori.uuid in favor of a properly maintained fork, this should not
  be a breaking change.
- Hstore now uses the null package in order to have nicer JSON serialization
- Changed bindata fork to https://github.com/kevinburke/go-bindata

### Fixed

- MySQL now correctly imports the null package for null json fields
  (thanks @mickeyreiss)
- Pass is now optional as in the README, except mssql (thanks @izumin5210)
- menus now singularizes correctly (thanks @jonas747)
- Randomize the time as a string for mysql
- Remove generation disclaimer for non-go files which prevents proper parsing
  of languages like html

## [v3.0.0-rc9]

### Fixed

- Fix a bug where rows may not be closed when bind failed in some way
- Fix aliasing of primary key column names in exist template
- Fix a bug in the psql driver where `double precision` was converted into
  a decimal type when it indeed should be float64.
- Fix a bug in the psql driver where `double precision` and `real` arrays were
  converted into decimal arrays, now they are converted into float64 arrays.
  For `real` this may not be ideal but it's a better fix for now since we don't
  have a float32 array.

## [v3.0.0-rc8]

### Added

- Add alternative syntax that's case-sensitive for defining aliases. This is
  particularly useful if toml key syntax is not good enough or viper keeps
  lowercasing all your keys (and you have uppercase names in say mssql).
- Added `aliasCols`, a template helper that can be used in conjunction with
  stringMap to transform a slice of column names into their aliased names.

### Fixed

- Fix several places that referenced primary keys in the templates that were
  not alias aware.
- Fix strmangle's CamelCase such that it forces the first char to be lowercase.

## [v3.0.0-rc7]

### Fixed

- Fix a bug where relationship to_many eager had an old argument in the query
  call.
- Fix a bug where relationship to_many eager call for afterSelectHooks was not
  using aliases.

### Removed

- Remove dependency on spew package in tests for boilingcore

## [v3.0.0-rc6]

### Fixed

- Fix a badly templated value in the mssql upsert template when using context
- Fix a bug where the database schema name was not properly being
  given to the templates. This broke explicit schemas where implicit schemas
  like in mysql, psql public etc. would still work.

## [v3.0.0-rc5]

### Fixed

- Fix generation failure on windows due to path manipulation issues (#314)

## [v3.0.0-rc4]

### Fixed

- Fix an issue in mysql driver where the null json field was getting the wrong type (#311)
- Fix compilation failures when using a nullable created_at column (#313)

## [v3.0.0-rc3]

### Added

- Add CockroachDB out of band driver link to readme (thanks @glerchundi)

### Fixed

- Fix an issue in psql driver where if you had a primary key with the same name
  in two different tables, it could possibly get the wrong one (#308)
- Allow errcheck to run successfully (exclude fmt.Fprint(|ln|f) function) on models
- Fix an issue where debug output wouldn't output on crash (exactly when you want it)
- Fix issue with boolean values 'true' and 'false' not being accepted (thanks @glerchundi)

## [v3.0.0-rc2]

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
