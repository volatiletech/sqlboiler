module github.com/volatiletech/sqlboiler/v4

go 1.16

require (
	github.com/DATA-DOG/go-sqlmock v1.4.1
	github.com/Masterminds/sprig/v3 v3.2.2
	github.com/davecgh/go-spew v1.1.1
	github.com/ericlagergren/decimal v0.0.0-20181231230500-73749d4874d5
	github.com/friendsofgo/errors v0.9.2
	github.com/go-sql-driver/mysql v1.5.0
	github.com/google/go-cmp v0.5.6
	github.com/lib/pq v1.2.1-0.20191011153232-f91d3411e481
	github.com/microsoft/go-mssqldb v0.15.0
	github.com/spf13/cast v1.4.1
	github.com/spf13/cobra v1.2.1
	github.com/spf13/viper v1.9.0
	github.com/stretchr/testify v1.7.0
	github.com/volatiletech/null/v8 v8.1.2
	github.com/volatiletech/randomize v0.0.1
	github.com/volatiletech/strmangle v0.0.4
	modernc.org/sqlite v1.14.5
)

retract (
	v4.10.0 // Generated models are invalid due to a wrong assignment
	v4.9.0 // Generated code shows v4.8.6, messed up commit tagging and untidy go.mod
)
