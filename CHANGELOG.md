# Changelog

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic
Versioning](http://semver.org/spec/v2.0.0.html).

## [v4.8.3] - 2021-11-16

### Fixed

- Fix bad use of titlecase in mysql enum name generation

## [v4.8.2] - 2021-11-16

### Fixed

- Fix regression in enum name generation

## [v4.8.1] - 2021-11-14

### Fixed

- Fix a regression in the soft delete test template generation introduced in
  4.8.1

## [v4.8.0] - 2021-11-14

### Added

- Add `--add-enum-types` to create distinct enum types instead of strings
  (thanks @stephenamo)

### Fixed

- Fix a regression in soft delete generation introduced in 4.7.1
  (thanks @stephenamo)

## [v4.7.1] - 2021-09-30

### Changed

- Change template locations to templates/{main,test}. This unfortunate move
  is necessary to preserve old behavior.

### Fixed

- Revert change to boilingcore.New() both in behavior and function signature

## [v4.7.0] - 2021-09-26

### Added

- Add configuration for overriding custom timestamp column names
  (thanks @stephanafamo)
- Add support for arguments to order by (thanks @emwalker and @alexdor)
- Add support for comments to mysql (thanks @Wuvist)

### Fixed

- Fix CVEs in transitive dependencies by bumping cobra & viper
- Fix inconsistent generation of IsNull/IsNotNull where helpers for types that
  appear both as null and not null in the database.
- JSON unmarshalling null into types.NullDecimal will no longer panic. String
  and format have been overridden to return "nil" when the underlying decimal
  is nil instead of crashing.

### Removed

- Removed bindata in favor of go:embed. This is not a breaking change as there
  are no longer supported versions of Go that do not support go:embed.

## [v4.6.0] - 2021-06-06

### Added

- Add `models.TableColumns.MODELNAME` which has the table.column name, useful
  for custom printf style queries (thanks @sadayuki-matsuno)

### Fixed

- Fix limit 0 queries (no longer omits limit clause) (thanks @longngn)
- Fix ordering issue when doing where clause on `deleted_at` and also trying to
  query for deleted_at
- Fix filename generation for tables that begin with `_`
- Add MarshalJSON implementation to NullDecimal to fix marshalling this type
  when nil.
- Fix issue with Go 1.16 compatibility for mssql driver by bumping mssql version
  (thanks @stefkampen)
- Fix Remove set operations for to-many relationships error when passing in nil
  or empty arrays of related models, it's now a no-op.

## [v4.5.0] - 2021-03-14

### Added

- Add new query mod WithDeleted to sidestep soft deletes in queries that
  support query mods (note there still is no way to do this for exists/find
  operations, see #854 for details)
- Add select hooks to the Find() methods, this was an accidental omission
  in previous versions (thanks @jakecoffman)

### Changed

- Change go-bindata to v3.22.0
- Change datetimeoffset and uniqueidentifier types in mssql this is a breaking
  change if you are using these types, but at least in the case of
  uniquedidentifier it was not possible to use without this change
  (thanks @severedsea)

### Fixed

- Fix unnecessary copies in JSON type which improves performance (thanks @bouk)
- Fix inclusion of foreign key constraints that target generated pg columns
  (thanks @chochihim)
- Fix generation failure bug in delete template when using --no-context
  --add-global-variants and --add-soft-deletes
- Fix cross-schema psql enum generation bug (thanks @csueiras)
- Fix bug where column alias was not respected in Load names (thanks @jalan)
- Fix bug with large uint64 values in eager loading (thanks @maku693)

## [v4.4.0] - 2020-12-16

### Added

- Add support for a qm.Comment query mod to add comments to queries that will
  be given to the server for tracing purposes (thanks @Pilatuz)

### Fixed

- Fix compatibility with ANSI_QUOTES in mysql (thanks @alexsander-souza)

## [v4.3.1] - 2020-11-16

### Fixed

- Fix case sensitive table name lookup in psql driver (thanks @severedsea)

## [v4.3.0] - 2020-11-03

### Added

- Add comments to generated code from db for psql driver (thanks @vladvelici)
- Add boil.None() to help with `DO NOTHING` upsert in mssql/mysql
  (thanks @emmanual099)

### Fixed

- Fix qm.WhereNotIn/qm.AndNotIn/qm.OrNotIn generating the wrong types of
  clauses (thanks @peterIdowns)
- Fix an issue where order of columns can change during eager loading which
  could cause errors (thanks @inoc603)
- Fix longstanding naming conflict when not using suffixes for foreign keys
  (thanks @yuzuy)
- Fix auto-generated timestamp columns not respecting aliases
  (thanks @while-loop)
- Fix upsert bug using schema names in mysql/mssql (thanks @emmanuel099)
- Fix blacklist/whitelist as environment variables being clobbered by incorrect
  values (thanks @Amandeepsinghghai)

## [v4.2.0] - 2020-07-03

### Added

- Add types.DecimalContext to control the context with which new decimals
  are created. This is important if your application uses a special context
  with more precision or requires the Go operating mode for example.
- Add WhereNotIn/AndNotIn/OrNotIn query mods to help solve a bug
- Add alias struct case type (uses the columns alias) (thanks @Darkclainer)
- Add ability to type replace on tables (thanks @stephenafamo)

### Changed

- Change the way column cache keys are created and looked up, improving memory
  performance of sqlboiler's caching layer (thanks @zikaeroh)

### Fixed

- Fix the psql driver to correctly ignore generated columns (thanks @chochihim)
- Fix an issue with mariadb ssl-mode when running generated tests
  (thanks @tooolbox)
- Fix $1 placeholder in mysql DeleteAll() when using soft delete
  (thanks @mfzy602)
- Fix boilingcore tests to use current module via replace instead of the
  published v4 module

## [v4.1.2] - 2020-05-18

### Fixed

- Fix $1 placeholder in mysql Delete() when using soft delete

## [v4.1.1] - 2020-05-05

### Fixed

- Fix mysql generation error made by previous commit

## [v4.1.0] - 2020-05-04

### Added

- Add support for postgresql `oid` type (thanks @ImVexed)
- Add a new `--relation-tag` option to control the tag name of the relationship
  struct in generated structs - this can expose loaded relationships to APIs
  (thanks @speatzle)

### Fixed

- Fix issue that caused horrible mysql generation performance (thanks @oderwat)

## [v4.0.1] - 2020-05-02

### Fixed

- Fix missing soft-delete pieces for the P/G variants (thanks @psucodervn)

## [v4.0.0] - 2020-04-26

**NOTE:** Your database drivers must be rebuilt upon moving to this release

### Added

- Add a `--add-soft-deletes` that changes the templates to use soft deletion
  (thanks @namco1992)
- Add a `--no-back-referencing` flag to disable setting backreferences in
  relationship helpers and eager loading. (thanks @namco1992)
- Add not in helpers (thanks @RichardLindhout)

### Changed

- Changed dependency scheme to go modules
- Changed randomize/strmangle to be external dependencies to avoid module
  cycles with the null package.
- Changed the way comparisons work for keying caches which dramatically speeds
  up cache lookups (thanks @zikaeroh)

### Fixed

- Fix postgres tests failing on partioned tables (thanks @troyanov)
- Fix enums with spaces being disallowed (thanks @razor-1)
- Fix postgresql looking at other types of tables (eg. views) and reporting that
  they do not have primary keys (thanks @chochihim)

## [v3.7.1] - 2020-04-26

### Fixed

- Fix bug in --version command reporting old version

## [v3.7.0] - 2020-03-31

### Added

- Add a whitelist of table.column names to set "-" struct tags for to ignore
  them during serialization (thanks @bogdanpradnj)
- Add 'where in' helpers for all primitive Go types (thanks @nwidger)
- Add a usage example of accessing the .R field (thanks @tooolbox)
- Add a check that safely sidesteps empty WhereIn queries like in other ORMs.
  (thanks @rekki)
- Add LeftOuter/RightOuter/FullOuter join query mods. Keep in mind this has no
  direct Bind() support as of yet so as with inner joins you -must- use custom
  data structs. (thanks @tzachshabtay)
- Add `distinct` query mod (thanks @tzachshabtay)

### Fixed

- Fix an idempotency issue with primary key column ordering
- Fix the plural/singular helpers for the word schema (thanks @Mushus)
- Fix some panics in Remove relationship set operations
- Fix panic when using WhereNullEQ with NullDecimal: implemented
  qmhelper.Nullable for NullDecimal

## [v3.6.1] - 2019-11-08

### Fixed

- Fix bug where mysql instead of mssql got the fix for offset rows which broke
  mysql users. (thanks @kkoudev)

## [v3.6.0] - 2019-10-21

### Added

- Log lines now can be controlled by context (both their presence and output)
  see the new boil.WithDebug/boil.WithDebugWriter (thanks @zikaeroh)
- OrderBy can now optionally take arguments (thanks @emwalker)
- Driver templates can now be disabled to allow further re-use of the existing
  drivers without having to recompile custom ones. (thanks @kurt-stolle)
- Add "whereIn" helpers for int64 (thanks @sadayuki-matsuno)

### Fixed

- Fix bug in mssql offset clause, it now properly includes the "ROWS" suffix
  as required by the official T-SQL spec.
- Fix unmarshalling json into a 0-value Decimal or NullDecimal type when
  the `*big.Decimal` itself is nil (thanks @ericlagergren)
- Fix a bug where mysql/sqlite3 would make spurious select id from table
  calls on insert when the id was already known.
- Fix various quality issues such as Replace(..., -1) -> ReplaceAll, use of
  EqualFold, etc. (thanks @sosiska)
- Fix an issue where using environment variables to configure sqlboiler's
  drivers the port would fail to parse (thanks @letientai299)

### Changed

- Changed the github.com/pkg/errors library for github.com/friendsofgo/errors
  This change is backwards compatible, it simply needs the new dependency
  to be downloaded. It also provides compatible with the new Go 1.13 error
  handling idioms with which we'd like for users to be able to use.
  (thanks @jwilner)

## [v3.5.0] - 2019-08-26

### Added

- Add wherein helpers for string and int column types (thanks @bugzpodder)
- Add sqlboiler version number in output

### Changed

- Rewrite alias documentation for relationships to be a bit more clear and
  concise.

### Fixed

- Fix a bug where relationship helpers were not quoting properly
- Fix issue where eager loads could produce ambiguous wherein clauses
- Fix import for sqlmock (thanks @zikaeroh)
- Fix issue with quoting identifiers that contain -
- Fix unsigned mediumint overflows for mysql tests
- Fix bad reference in mssql templates (thanks @cliedeman)

## [v3.4.0] - 2019-05-27

### Added

- Add domain name to psql column data from driver (thanks @natsukagami)
- Add full_db_type to psql column data from driver (thanks @vincentserpoul)

### Fixed

- Fix problems with idempotency patch messing up the lists of columns, now to
  achieve idempotency in column lists we sort based on information_schema's
  ordinal index for columns.
- Fix code compilation failure due to whitespace when certain struct tag
  options were used.

## [v3.3.0] - 2019-05-20

### Added

- Add title as an option for struct tag casing configuration
- CTE Support for query-mod built queries: see With query mod (thanks @lucaslsl)
- Extra documentation to explain how local/foreign work in aliases
  (thanks @NickyMateev).
- When re-running the sqlboiler command to dump a schema, all tables, columns,
  and foreign keys are now selected in a predictable sorted order. This means
  that if you run the command against the same schema twice you should get
  exactly the same output each time. This is useful if you want to check in
  your generated code, as it avoids pointless churn. It is also helpful if you
  want to test that the checked-in generated code is up to date. You can now
  regenerate the code and simply check that nothing has changed. Note that
  with MS SQL this only works if you provide explicit names for all of your
  foreign keys, as MS SQL generates names with a random component
  otherwise. (thanks @autarch)

### Changed

- Change error return on ModelSlice.DeleteAll when ModelSlice is nil to simply
  return 0, nil. This is not really an error state and is burdensome to users.

### Fixed

- Fix identifiers in Where helpers not being quoted nor having table
  specifications.
- Fix a bug where types declared within a non-global scope could cause cache
  key collisions and create undesirable behavior within Bind().
- Fix a bug where decimal types in the database without decimal points
  would be provided by the driver as an int64 which failed to Scan() into
  a types.Decimal
- A Postgres domain that was created from an array ("CREATE DOMAIN foo AS
  INT[]") would cause a runtime panic when trying to generate model
  code. These are now handled correctly when the array's type is a
  built-in. If the array type is itself a UDT then this will be treated as a
  string, which will not be correct in some cases. (thanks @autarch)
- Fix doc typo around qm.Load/qm.Rels (thanks @KopiasCsaba)
- Fix a bug where hstore did not Value() properly
- Fix an issue related to interface{} keys from yaml type replacements

## [v3.2.0] - 2019-01-22

### Added

- Type-safe where clauses can now be created, see README for details. It is
  highly recommended that this be used at all times.
- Type-safe where clauses can now be combined with Or2 for setting or.
- The new Expr query mod will allow you to group where statements manually
  (this turns off all automatic paretheses in the where clause for the query).
- Driver specific commands (eg. pg_dump) that are run for test scaffolding
  will now output their error messages to stderr where they were previously
  silently failing (thanks @LukasAuerbeck)
- Add skipsqlcmd to generated test code for each driver. This allows skipping
  the whole drop/create database cycle while testing so you may point sqlboiler
  at a pre-setup test database. (thanks @gemscng)
- Add a way to skip hook execution for a given query (boil.SkipHooks)
- Add a way to skip timestamp updating for a given query (boil.SkipTimestamps)
- Add note about mysql minimum version requirement to README (thanks @jlarusso)

### Fixed

- Fix panic on eager load with nullable foreign keys
- Fix bug that prevented 'where' and 'in' from being mixed naturally as 'in'
  query mods would always be rendered at the end of the query resulting in
  an unintentional problem.
- Fix an issue where 'in' query mods were not being automatically grouped in
  parentheses like 'where' statements.
- Fix bug where mysql columns can sometimes be selected out of order in
  certain internal queries. (thanks @cpickett-ml)
- Fix bug where an incorrect query could be built while eager loading nullable
  relations (thanks @parnic)
- Fix bug where aliases weren't used in many-to-many eager loading
  (thanks @nwidger for suggested fix)
- Fix bug where mysql driver would look outside the current database for
  indexes that applied to tables and columns named the same and apply those
  constraints to the generated schema.
- Fix MSSQL Link in Readme (thanks @philips)
- Fix bug where psql upsert would error when not doing an update
- Fix bug where mysql upsert did not have quotes around the table name
- Fix bug where yaml config files would panic due to type assertions (thanks
  @ch3rub1m)
- Fix bug where a table name that was a Go keyword could cause test failures
- Fix missing boil.Columns type in README (thanks @tatane616)

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
