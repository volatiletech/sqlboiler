<img src="http://i.imgur.com/R5g99sO.png"/>

# STILL IN DEVELOPMENT. ETA RELEASE: 1 MONTH

# SQLBoiler

[![GoDoc](https://godoc.org/github.com/pobri19/sqlboiler?status.svg)](https://godoc.org/github.com/pobri19/sqlboiler)

SQLBoiler is a tool to generate a Go ORM tailored to your database schema.

# Config?

To use SQLBoiler you need to create a ````config.toml```` in SQLBoiler's root directory. The file format looks like the following:

````
[postgres]
  host="localhost"
  port=5432
  user="dbusername"
  pass="dbpassword"
  dbname="dbname"
````

# How?

SQLBoiler connects to your database (defined in your config.toml file) to ascertain the structure of your tables, and builds your Go ORM code using the templates defined in the ````/templates```` folder.

Running SQLBoiler without the ````--table```` flag will result in SQLBoiler building boilerplate code for every table in your database marked as ````public````.

Before you use SQLBoiler make sure you create a ````config.toml```` configuration file with your database details, and specify your database by using the ````--driver```` flag.


````
Complete documentation is available at http://github.com/pobri19/sqlboiler

Usage:
  sqlboiler [flags]

Flags:
  -d, --driver string    The name of the driver in your config.toml (mandatory)
  -f, --folder string    The name of the output folder (default "output")
  -p, --pkgname string   The name you wish to assign to your generated package (default "model")
  -t, --table string     A comma seperated list of table names (generates all tables if not provided)
````
