# Everything in this make file assumes you have the ability
# to run the commands as they are written (meaning .my.cnf and .pgpass
# files set up with admin users) as well as the mssql db set up with
# the sa / Sqlboiler@1234 credentials. See testdata/env.sh.

-include testdata/env.mk

# Automatically import test related environments
.PHONY: testdata/env.mk
testdata/env.mk: testdata/env.sh
	sed 's/"//g ; s/=/:=/' < $< > $@

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
	psql --host localhost --username postgres --command "create user $(SQLBOILER_TEST_USER) with superuser password '$(SQLBOILER_TEST_PASS)';"

.PHONY: test-user-mysql
test-user-mysql:
	mysql --host localhost --execute "create user $(SQLBOILER_TEST_USER) identified by '$(SQLBOILER_TEST_PASS)';" 
	mysql --host localhost --execute "grant all privileges on *.* to $(SQLBOILER_TEST_USER);" 

# Must be run after a database is created
.PHONY: test-user-mssql
test-user-mssql:
	sqlcmd -S localhost -U sa -P $(SQLBOILER_TEST_PASS) -Q "create login $(SQLBOILER_TEST_USER) with password = '$(SQLBOILER_TEST_PASS)';"
	sqlcmd -S localhost -U sa -P $(SQLBOILER_TEST_PASS) -Q "alter server role sysadmin add member $(SQLBOILER_TEST_USER);"

.PHONY: test-db-psql
test-db-psql:
	env PGPASSWORD=$(SQLBOILER_TEST_PASS) dropdb --host localhost --username $(SQLBOILER_TEST_USER) --if-exists $(SQLBOILER_TEST_DB)
	env PGPASSWORD=$(SQLBOILER_TEST_PASS) createdb --host localhost --owner $(SQLBOILER_TEST_USER) --username $(SQLBOILER_TEST_USER) $(SQLBOILER_TEST_DB)
	env PGPASSWORD=$(SQLBOILER_TEST_PASS) psql --host localhost --username $(SQLBOILER_TEST_USER) --file testdata/psql_test_schema.sql $(SQLBOILER_TEST_DB)

.PHONY: test-db-mysql
test-db-mysql:
	mysql --host localhost --user $(SQLBOILER_TEST_USER) --password=$(SQLBOILER_TEST_PASS) --execute "drop database if exists $(SQLBOILER_TEST_DB);"
	mysql --host localhost --user $(SQLBOILER_TEST_USER) --password=$(SQLBOILER_TEST_PASS) --execute "create database $(SQLBOILER_TEST_DB);"
	mysql --host localhost --user $(SQLBOILER_TEST_USER) --password=$(SQLBOILER_TEST_PASS) $(SQLBOILER_TEST_DB) < testdata/mysql_test_schema.sql

.PHONY: test-db-mssql
test-db-mssql:
	sqlcmd -S localhost -U $(SQLBOILER_TEST_USER) -P $(SQLBOILER_TEST_PASS) -Q "drop database if exists $(SQLBOILER_TEST_DB)";
	sqlcmd -S localhost -U $(SQLBOILER_TEST_USER) -P $(SQLBOILER_TEST_PASS) -Q "create database $(SQLBOILER_TEST_DB)";
	sqlcmd -S localhost -U $(SQLBOILER_TEST_USER) -P $(SQLBOILER_TEST_PASS) -d $(SQLBOILER_TEST_DB) -i testdata/mssql_test_schema.sql

.PHONY: test-generate-psql
test-generate-psql:
	printf "[psql]\nhost=\"localhost\"\nport=5432\nuser=\"%s\"\npass=\"%s\"\ndbname=\"%s\"\nsslmode=\"disable\"\n" $(SQLBOILER_TEST_USER) $(SQLBOILER_TEST_PASS) $(SQLBOILER_TEST_DB) > sqlboiler.toml
	./sqlboiler --wipe psql

.PHONY: test-generate-mysql
test-generate-mysql:
	printf "[mysql]\nhost=\"localhost\"\nport=3306\nuser=\"%s\"\npass=\"%s\"\ndbname=\"%s\"\nsslmode=\"false\"\n" $(SQLBOILER_TEST_USER) $(SQLBOILER_TEST_PASS) $(SQLBOILER_TEST_DB) > sqlboiler.toml
	./sqlboiler --wipe mysql

.PHONY: test-generate-mssql
test-generate-mssql:
	printf "[mssql]\nhost=\"localhost\"\nport=1433\nuser=\"%s\"\npass=\"%s\"\ndbname=\"%s\"\nsslmode=\"disable\"\n" $(SQLBOILER_TEST_USER) $(SQLBOILER_TEST_PASS) $(SQLBOILER_TEST_DB) > sqlboiler.toml
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
	env PGPASSWORD=$(SQLBOILER_TEST_PASS) createdb --host localhost --username $(SQLBOILER_TEST_USER) --owner $(SQLBOILER_TEST_USER) $(SQLBOILER_TEST_DRIVER_DB)

.PHONY: driver-db-mysql
driver-db-mysql:
	mysql --host localhost --user $(SQLBOILER_TEST_USER) --password=$(SQLBOILER_TEST_PASS) --execute "create database $(SQLBOILER_TEST_DRIVER_DB);"

.PHONY: driver-db-mssql
driver-db-mssql:
	sqlcmd -S localhost -U sa -P $(SQLBOILER_TEST_PASS) -Q "create database $(SQLBOILER_TEST_DRIVER_DB);"
	sqlcmd -S localhost -U sa -P $(SQLBOILER_TEST_PASS) -d $(SQLBOILER_TEST_DRIVER_DB) -Q "exec sp_configure 'contained database authentication', 1;"
	sqlcmd -S localhost -U sa -P $(SQLBOILER_TEST_PASS) -d $(SQLBOILER_TEST_DRIVER_DB) -Q "reconfigure"
	sqlcmd -S localhost -U sa -P $(SQLBOILER_TEST_PASS) -d $(SQLBOILER_TEST_DRIVER_DB) -Q "alter database $(SQLBOILER_TEST_DRIVER_DB) set containment = partial;"

.PHONY: driver-user-psql
driver-user-psql:
	env PGPASSWORD=$(SQLBOILER_TEST_PASS) psql --host localhost --username $(SQLBOILER_TEST_USER) --command "create role $(SQLBOILER_TEST_DRIVER_USER) login nocreatedb nocreaterole password '$(SQLBOILER_TEST_PASS)';" $(SQLBOILER_TEST_DRIVER_DB)
	env PGPASSWORD=$(SQLBOILER_TEST_PASS) psql --host localhost --username $(SQLBOILER_TEST_USER) --command "alter database $(SQLBOILER_TEST_DRIVER_DB) owner to $(SQLBOILER_TEST_DRIVER_USER);" $(SQLBOILER_TEST_DRIVER_DB)

.PHONY: driver-user-mysql
driver-user-mysql:
	mysql --host localhost --execute "create user $(SQLBOILER_TEST_DRIVER_USER) identified by '$(SQLBOILER_TEST_PASS)';" 
	mysql --host localhost --execute "grant all privileges on $(SQLBOILER_TEST_DRIVER_DB).* to $(SQLBOILER_TEST_DRIVER_USER);" 

.PHONY: driver-user-mssql
driver-user-mssql:
	sqlcmd -S localhost -U sa -P $(SQLBOILER_TEST_PASS) -d $(SQLBOILER_TEST_DRIVER_DB) -Q "create user $(SQLBOILER_TEST_DRIVER_USER) with password = '$(SQLBOILER_TEST_PASS)'";
	sqlcmd -S localhost -U sa -P $(SQLBOILER_TEST_PASS) -d $(SQLBOILER_TEST_DRIVER_DB) -Q "grant alter, control to $(SQLBOILER_TEST_DRIVER_USER)";

.PHONY: driver-test-psql
driver-test-psql:
	DRIVER_DB=$(SQLBOILER_TEST_DRIVER_DB) DRIVER_USER=$(SQLBOILER_TEST_DRIVER_USER) DRIVER_PASS=$(SQLBOILER_TEST_DRIVER_PASS) go test -v -race github.com/volatiletech/sqlboiler/drivers/sqlboiler-psql/driver

.PHONY: driver-test-mysql
driver-test-mysql:
	DRIVER_DB=$(SQLBOILER_TEST_DRIVER_DB) DRIVER_USER=$(SQLBOILER_TEST_DRIVER_USER) DRIVER_PASS=$(SQLBOILER_TEST_DRIVER_PASS) go test -v -race github.com/volatiletech/sqlboiler/drivers/sqlboiler-mysql/driver

.PHONY: driver-test-mssql
driver-test-mssql:
	DRIVER_DB=$(SQLBOILER_TEST_DRIVER_DB) DRIVER_USER=$(SQLBOILER_TEST_DRIVER_USER) DRIVER_PASS=$(SQLBOILER_TEST_DRIVER_PASS) go test -v -race github.com/volatiletech/sqlboiler/drivers/sqlboiler-mssql/driver
	
# ====================================
# Clean operations
# ====================================

.PHONY: clean
clean:
	rm -f ./sqlboiler
	rm -f ./sqlboiler-*
