# Everything in this make file assumes you have the ability
# to run the commands as they are written (meaning .my.cnf and .pgpass
# files set up with admin users) as well as the mssql db set up with
# the sa / Sqlboiler@1234 credentials (see docker run below for example)

USER=sqlboiler_root_user
DB=sqlboiler_model_test
PASS=sqlboiler
MSSQLPASS=Sqlboiler@1234

DRIVER_USER=sqlboiler_driver_user
DRIVER_DB=sqlboiler_driver_test

# Builds all software and runs model tests
.PHONY: tests
tests: test-psql test-mysql test-mssql

# Creates super users
.PHONY: test-users
test-users: test-user-psql test-user-mysql test-user-mssql

# Creates databases
.PHONY: test-dbs
test-dbs: test-db-psql test-db-mysql test-db-mssql

# Runs tests on the drivers
.PHONY: driver-tests
driver-tests: driver-test-psql driver-test-mysql driver-test-mssql

# Creates regular databases
.PHONY: driver-dbs
driver-dbs: driver-db-psql driver-db-mysql driver-db-mssql

# Creates regular users with access to only one DB, may modify the DB
# to ensure write access to the user
.PHONY: driver-users
driver-users: driver-user-psql driver-user-mysql driver-user-mssql

# ====================================
# Building operations
# ====================================

.PHONY: build
build:
	go build github.com/volatiletech/sqlboiler

.PHONY: build-psql
build-psql:
	go build github.com/volatiletech/sqlboiler/drivers/sqlboiler-psql

.PHONY: build-mysql
build-mysql:
	go build github.com/volatiletech/sqlboiler/drivers/sqlboiler-mysql

.PHONY: build-mssql
build-mssql:
	go build github.com/volatiletech/sqlboiler/drivers/sqlboiler-mssql

# ====================================
# Testing operations
# ====================================

.PHONY: test
test:
	go test -v -race $(go list ./... | grep -v /drivers/ | grep -v /vendor/)

.PHONY: test-user-psql
test-user-psql:
	# Can't use createuser because it interactively promtps for a password
	psql --host localhost --username postgres --command "create user $(USER) with superuser password '$(PASS)';"

.PHONY: test-user-mysql
test-user-mysql:
	mysql --host localhost --execute "create user $(USER) identified by '$(PASS)';" 
	mysql --host localhost --execute "grant all privileges on *.* to $(USER);" 

# Must be run after a database is created
.PHONY: test-user-mssql
test-user-mssql:
	sqlcmd -S localhost -U sa -P $(MSSQLPASS) -Q "create login $(USER) with password = '$(MSSQLPASS)';"
	sqlcmd -S localhost -U sa -P $(MSSQLPASS) -Q "alter server role sysadmin add member $(USER);"

.PHONY: test-db-psql
test-db-psql:
	env PGPASSWORD=$(PASS) dropdb --host localhost --username $(USER) --if-exists $(DB)
	env PGPASSWORD=$(PASS) createdb --host localhost --owner $(USER) --username $(USER) $(DB)
	env PGPASSWORD=$(PASS) psql --host localhost --username $(USER) --file testdata/psql_test_schema.sql $(DB)

.PHONY: test-db-mysql
test-db-mysql:
	mysql --host localhost --user $(USER) --password=$(PASS) --execute "drop database if exists $(DB);"
	mysql --host localhost --user $(USER) --password=$(PASS) --execute "create database $(DB);"
	mysql --host localhost --user $(USER) --password=$(PASS) $(DB) < testdata/mysql_test_schema.sql

.PHONY: test-db-mssql
test-db-mssql:
	sqlcmd -S localhost -U $(USER) -P $(MSSQLPASS) -Q "drop database if exists $(DB)";
	sqlcmd -S localhost -U $(USER) -P $(MSSQLPASS) -Q "create database $(DB)";
	sqlcmd -S localhost -U $(USER) -P $(MSSQLPASS) -d $(DB) -i testdata/mssql_test_schema.sql

.PHONY: test-generate-psql
test-generate-psql:
	printf "[psql]\nhost=\"localhost\"\nport=5432\nuser=\"%s\"\npass=\"%s\"\ndbname=\"%s\"\nsslmode=\"disable\"\n" $(USER) $(PASS) $(DB) > sqlboiler.toml
	./sqlboiler --wipe psql

.PHONY: test-generate-mysql
test-generate-mysql:
	printf "[mysql]\nhost=\"localhost\"\nport=3306\nuser=\"%s\"\npass=\"%s\"\ndbname=\"%s\"\nsslmode=\"false\"\n" $(USER) $(PASS) $(DB) > sqlboiler.toml
	./sqlboiler --wipe mysql

.PHONY: test-generate-mssql
test-generate-mssql:
	printf "[mssql]\nhost=\"localhost\"\nport=1433\nuser=\"%s\"\npass=\"%s\"\ndbname=\"%s\"\nsslmode=\"disable\"\n" $(USER) $(MSSQLPASS) $(DB) > sqlboiler.toml
	./sqlboiler --wipe mssql 

.PHONY: test-psql
test-psql:
	go test -v -race ./models

.PHONY: test-mysql
test-mysql:
	go test -v -race ./models

.PHONY: test-mssql
test-mssql:
	go test -v -race ./models

# ====================================
# Driver operations
# ====================================

.PHONY: driver-db-psql
driver-db-psql:
	env PGPASSWORD=$(PASS) createdb --host localhost --username $(USER) --owner $(USER) $(DRIVER_DB)

.PHONY: driver-db-mysql
driver-db-mysql:
	mysql --host localhost --user $(USER) --password=$(PASS) --execute "create database $(DRIVER_DB);"

.PHONY: driver-db-mssql
driver-db-mssql:
	sqlcmd -S localhost -U sa -P $(MSSQLPASS) -Q "create database $(DRIVER_DB);"
	sqlcmd -S localhost -U sa -P $(MSSQLPASS) -d $(DRIVER_DB) -Q "exec sp_configure 'contained database authentication', 1;"
	sqlcmd -S localhost -U sa -P $(MSSQLPASS) -d $(DRIVER_DB) -Q "reconfigure"
	sqlcmd -S localhost -U sa -P $(MSSQLPASS) -d $(DRIVER_DB) -Q "alter database $(DRIVER_DB) set containment = partial;"

.PHONY: driver-user-psql
driver-user-psql:
	env PGPASSWORD=$(PASS) psql --host localhost --username $(USER) --command "create role $(DRIVER_USER) login nocreatedb nocreaterole password '$(PASS)';" $(DRIVER_DB)
	env PGPASSWORD=$(PASS) psql --host localhost --username $(USER) --command "alter database $(DRIVER_DB) owner to $(DRIVER_USER);" $(DRIVER_DB)

.PHONY: driver-user-mysql
driver-user-mysql:
	mysql --host localhost --execute "create user $(DRIVER_USER) identified by '$(PASS)';" 
	mysql --host localhost --execute "grant all privileges on $(DRIVER_DB).* to $(DRIVER_USER);" 

.PHONY: driver-user-mssql
driver-user-mssql:
	sqlcmd -S localhost -U sa -P $(MSSQLPASS) -d $(DRIVER_DB) -Q "create user $(DRIVER_USER) with password = '$(MSSQLPASS)'";
	sqlcmd -S localhost -U sa -P $(MSSQLPASS) -d $(DRIVER_DB) -Q "grant alter, control to $(DRIVER_USER)";

.PHONY: driver-test-psql
driver-test-psql:
	go test -v -race github.com/volatiletech/sqlboiler/drivers/sqlboiler-psql/driver -hostname localhost -username $(DRIVER_USER) -password $(PASS) -database $(DRIVER_DB)

.PHONY: driver-test-mysql
driver-test-mysql:
	go test -v -race github.com/volatiletech/sqlboiler/drivers/sqlboiler-mysql/driver -hostname localhost -username $(DRIVER_USER) -password $(PASS) -database $(DRIVER_DB)

.PHONY: driver-test-mssql
driver-test-mssql:
	go test -v -race github.com/volatiletech/sqlboiler/drivers/sqlboiler-mssql/driver -hostname localhost -username $(DRIVER_USER) -password $(MSSQLPASS) -database $(DRIVER_DB)
	
# ====================================
# Clean operations
# ====================================

.PHONY: clean
clean:
	rm -f ./sqlboiler
	rm -f ./sqlboiler-*
