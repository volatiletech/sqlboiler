# SQLBoiler

[![GoDoc](https://godoc.org/github.com/pobri19/sqlboiler?status.svg)](https://godoc.org/github.com/pobri19/sqlboiler)
![CircleCI](https://circleci.com/gh/nullbio/sqlboiler.svg?style=shield&circle-token=:circle-token)

SQLBoiler is a tool to generate a Go ORM tailored to your database schema.

#### Config

config.toml

````
[postgres]
  host="localhost"
  port=5432
  user="dbusername"
  pass="dbpassword"
  dbname="dbname"
````

#### How

SQLBoiler connects to your database (defined in your config.toml file) to ascertain the structure of your tables, and builds your Go ORM code using the templates defined in the ````/templates```` folder.

Running SQLBoiler without the ````--table```` flag will result in SQLBoiler building boilerplate code for every table in your database marked as ````public````.

Before you use SQLBoiler make sure you create a ````config.toml```` configuration file with your database details, and specify your database by using the ````--driver```` flag.


````
Usage:
  sqlboiler [flags]

Flags:
  -d, --driver string    The name of the driver in your config.toml (mandatory)
  -f, --folder string    The name of the output folder (default "output")
  -p, --pkgname string   The name you wish to assign to your generated package (default "model")
  -t, --table string     A comma seperated list of table names (generates all tables if not provided)
````
