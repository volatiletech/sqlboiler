module github.com/volatiletech/sqlboiler/v4

go 1.16

require (
	github.com/DATA-DOG/go-sqlmock v1.4.1
	github.com/Masterminds/sprig/v3 v3.2.2
	github.com/davecgh/go-spew v1.1.1
	github.com/ericlagergren/decimal v0.0.0-20190420051523-6335edbaa640
	github.com/friendsofgo/errors v0.9.2
	github.com/go-sql-driver/mysql v1.6.0
	github.com/gofrs/uuid v4.2.0+incompatible // indirect
	github.com/golang-sql/civil v0.0.0-20220223132316-b832511892a9 // indirect
	github.com/google/go-cmp v0.5.8
	github.com/huandu/xstrings v1.3.2 // indirect
	github.com/imdario/mergo v0.3.13 // indirect
	github.com/inconshreveable/mousetrap v1.0.1 // indirect
	github.com/lib/pq v1.10.6
	github.com/microsoft/go-mssqldb v0.17.0
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/pelletier/go-toml/v2 v2.0.5 // indirect
	github.com/shopspring/decimal v1.3.1 // indirect
	github.com/spf13/afero v1.9.2 // indirect
	github.com/spf13/cast v1.5.0
	github.com/spf13/cobra v1.5.0
	github.com/spf13/viper v1.12.0
	github.com/stretchr/testify v1.8.0
	github.com/subosito/gotenv v1.4.1 // indirect
	github.com/volatiletech/null/v8 v8.1.2
	github.com/volatiletech/randomize v0.0.1
	github.com/volatiletech/strmangle v0.0.5
	golang.org/x/crypto v0.0.0-20220826181053-bd7e27e6170d // indirect
	golang.org/x/sys v0.0.0-20220825204002-c680a09ffe64 // indirect
	golang.org/x/tools v0.1.12 // indirect
	golang.org/x/xerrors v0.0.0-20220609144429-65e65417b02f // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	lukechampine.com/uint128 v1.2.0 // indirect
	modernc.org/cc/v3 v3.36.3 // indirect
	modernc.org/ccgo/v3 v3.16.9 // indirect
	modernc.org/libc v1.17.1 // indirect
	modernc.org/opt v0.1.3 // indirect
	modernc.org/sqlite v1.18.1
	modernc.org/strutil v1.1.3 // indirect
	modernc.org/token v1.0.1 // indirect
)

retract (
	v4.10.0 // Generated models are invalid due to a wrong assignment
	v4.9.0 // Generated code shows v4.8.6, messed up commit tagging and untidy go.mod
)
