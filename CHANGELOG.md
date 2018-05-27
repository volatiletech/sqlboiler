# Changelog

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

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

- Rename postgresql driver to psql
- Drivers are now responsible for their own config validation
- Drivers are now responsible for their own default config values
- Drivers can now take any configuration via environment or config files. This
  allows drivers to be much more configurable and there's no constraints
  imposed by the main/testing of sqlboiler itself.
- Schema is no longer a sqlboiler level configuration setting. Instead it is
  passed back to sqlboiler when the driver is asked to produce the information
  for the database we're generating for.
- Imports is now it's own package and is cleaned up significantly
- Replace templates (hidden flag) now use the ; separator not : for windows
  compatibility.
- Rows affected is now on by default, this is a breaking change but easily
  fixed by using the flag to turn it off.

### Removed

- The concept of TestMain is now gone from both templates and imports. It's
  been superceded by the new driver abilities to supply templates and imports.
  The drivers add their mains to the TestSingleton templates.
