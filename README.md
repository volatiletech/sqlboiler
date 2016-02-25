<img src="http://i.imgur.com/R5g99sO.png"/>

# SQLBoiler

[![GoDoc](https://godoc.org/github.com/pobri19/sqlboiler?status.svg)](https://godoc.org/github.com/pobri19/sqlboiler)

SQLBoiler is a tool to generate Go boilerplate code for database interactions. So far this includes struct definitions and database statement helper functions.

# Supports?
* Postgres

If anyone wants to create a driver for their own database it's easy very to do. All you need to do is create a driver in the ````/dbdrivers```` package that implements the ````DBDriver```` interface (you can use ````postgres_driver.go```` as an example), and add your driver to the switch statement in the ````initDBDriver()```` function in ````sqlboiler.go````. That's it!

I've included templates for struct definitions and select, delete and insert statement helpers. Editing the output of the existing templates is as easy as modifying the template file, but to add new types of statements you'll need to add a new command and a new template. This is also very easy to do, and you can use any of the existing command files as an example.

# Why?

Because writing boilerplate is annoying, that's why!

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

SQLBoiler connects to your database (defined in your config.toml file) to ascertain the structure of your tables, and builds your Go boilerplate code using the templates defined in the ````/templates```` folder.

Running SQLBoiler without the ````--table```` flag will result in SQLBoiler building boilerplate code for every table in your database marked as ````public````.

For example, on a table with the following schema:

````
  Column  |          Type          | Modifiers | Storage  | Stats target | Description
----------+------------------------+-----------+----------+--------------+-------------
 id       | integer                | not null  | plain    |              |
 guy_name | character varying(255) |           | extended |              |
 guy_age  | bigint                 |           | plain    |              |
````

Running the following command:

````./sqlboiler all --driver="postgres" --table="example_friend"````

Results in the following boilerplate generation:

````
type ExampleFriend struct {
	ID      int64  `db:"example_friend_id",json:"id"`
	GuyName string `db:"example_friend_guy_name",json:"guy_name"`
	GuyAge  int64  `db:"example_friend_guy_age",json:"guy_age"`
}

func DeleteExampleFriend(id int, db *sqlx.DB) error {
	if id == nil {
		return nil, errors.New("No ID provided for ExampleFriend delete")
	}

	err := db.Exec("DELETE FROM example_friend WHERE id=$1", id)
	if err != nil {
		return errors.New("Unable to delete from example_friend: %s", err)
	}

	return nil
}

func InsertExampleFriend(o *ExampleFriend, db *sqlx.DB) (int, error) {
	if o == nil {
		return 0, errors.New("No ExampleFriend provided for insertion")
	}

	var rowID int
	err := db.QueryRow(`
          INSERT INTO example_friend
          (id, guy_name, guy_age)
          VALUES($1, $2, $3)
          RETURNING id
        `)

	if err != nil {
		return 0, fmt.Errorf("Unable to insert example_friend: %s", err)
	}

	return rowID, nil
}

func SelectExampleFriend(id int, db *sqlx.DB) (ExampleFriend, error) {
	if id == 0 {
		return nil, errors.New("No ID provided for ExampleFriend select")
	}

	var exampleFriend ExampleFriend
	err := db.Select(&exampleFriend, `
          SELECT id AS example_friend_id, guy_name AS example_friend_guy_name, guy_age AS example_friend_guy_age
          WHERE id=$1
        `, id)

	if err != nil {
		return nil, fmt.Errorf("Unable to select from example_friend: %s", err)
	}

	return exampleFriend, nil
}
````

# Commands?

You have the ability to generate specific objects, or all objects by using the appropriate command. You can also choose to generate for all tables, a single table, or a group of tables using the ````--table```` flag. By default SQLBoiler outputs to Stdout, unless a file is specified with the ````--out```` flag.

Before you use SQLBoiler make sure you create a ````config.toml```` configuration file with your database details, and specify your database by using the ````--driver```` flag.

For example: ````./sqlboiler all --driver="postgres" --out="filename.go" --table="my_table1, my_table2"````

````
Usage:
  sqlboiler [command]

Available Commands:
  all         Generate all templates from table definitions
  delete      Generate delete statement helpers from table definitions
  insert      Generate insert statement helpers from table definitions
  select      Generate select statement helpers from table definitions
  struct      Generate structs from table definitions

Flags:
  -d, --driver string   The name of the driver in your config.toml (mandatory)
  -h, --help            help for sqlboiler
  -o, --out string      The name of the output file
  -t, --table string    A comma seperated list of table names

Use "sqlboiler [command] --help" for more information about a command.
````

# Templates?

If you wish to modify the boilerplate that SQLBoiler provides you it's as simple as editing the relevant template in the ````/templates```` directory.
