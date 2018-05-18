# Changelog

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Add support for geotyepes for psql (thanks @saulortega)
- Add a boil.sh that includes many commands that help build and test sqlboiler
- Add a flag to set the config file (thanks @l4u)
- Add ability for drivers to provide their own templates/replacement templates
- Add support for citext to psql (thanks @boxofrad)
- Ignore virtual columns in mysql (thanks @Nykakin)

### Changed

- Rename postgresql driver to psql
- Driver split

    Drivers are now separate binaries. A lot of the reason for this is because
    having to keep all the drivers together inside sqlboiler was going to later
    cause more churn that it was worth. sqlite showed us this immediately with
    it's cgo dependencies as an example.

- Drivers are now responsible for their own config validation
- Drivers are now responsible for their own default config values
- Drivers can now take any configuration via environment or config files. This
  allows drivers to be much more configurable and there's no constraints
  imposed by the main/testing of sqlboiler itself.
- Schema is no longer a sqlboiler level configuration setting. Instead it is
  passed back to sqlboiler when the driver is asked to produce the information
  for the database we're generating for.
- Imports is now it's own package and is cleaned up significantly
