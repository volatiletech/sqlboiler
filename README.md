# SQLBoiler

[![License](https://img.shields.io/badge/license-BSD-blue.svg)](https://github.com/vattle/sqlboiler/blob/master/LICENSE)
[![GoDoc](https://godoc.org/github.com/pobri19/sqlboiler?status.svg)](https://godoc.org/github.com/pobri19/sqlboiler)
[![CircleCI](https://circleci.com/gh/vattle/sqlboiler.svg?style=shield)](https://circleci.com/gh/vattle/sqlboiler)
[![Go Report Card](https://goreportcard.com/badge/kubernetes/helm)](http://goreportcard.com/report/vattle/sqlboiler)

SQLBoiler is a tool to generate a Go ORM tailored to your database schema.

#### Config

Before you use SQLBoiler make sure you create a `sqlboiler.toml` configuration file containing your database details.

The configuration file loader checks the working directory, `$HOME/sqlboiler/.config` directory and `$XDG_CONFIG_HOME/sqlboiler` directory to locate your configuration file. The working directory takes precedence.

`sqlboiler.toml` example file:

```
[postgres]
  host="localhost"
  port=5432
  user="dbusername"
  pass="dbpassword"
  dbname="dbname"
  sslmode="require"
```

The following config fields are optional, and default to:

```
postgres.port=5432
postgres.sslmode="require"
```

To disable SSLMode, use `sslmode="disable"`.


#### How

SQLBoiler connects to your database (defined in your sqlboiler.toml file) to ascertain the structure of your tables, and builds your Go ORM code using the templates defined in the `/templates` and `/templates_test` folders.

Running SQLBoiler without the `--table` flag will result in SQLBoiler building boilerplate code for every table in your database  the `public` schema. Tables created in Postgres will default to the `public` schema.



````
Usage:
  sqlboiler [flags] <driver>

Examples:
sqlboiler postgres

Flags:
  -o, --output string    The name of the folder to output to (default "models")
  -p, --pkgname string   The name you wish to assign to your generated package (default "models")
````
