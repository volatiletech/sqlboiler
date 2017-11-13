# Everything in this make file assumes you have the ability
# to run the commands as they are written (meaning .my.cnf and .pgpass
# files set up with admin users) as well as the mssql db set up with
# the sa / Sqlboiler@1234 credentials (see docker run below for example)

.PHONY: \
	build        build-psql       build-mysql       build-mssql \
	tests        test-psql        test-mysql        test-mssql \
	test-dbs     test-db-psql     test-db-mysql     test-db-mssql \
	test-users   test-user-psql   test-user-mysql   test-user-mssql \
	driver-tests driver-test-psql driver-test-mysql driver-test-mssql \
	driver-users driver-user-psql driver-user-mysql driver-user-mssql \
	driver-dbs   driver-db-psql   driver-db-mysql   driver-db-mssql \
	run-mssql

USER=sqlboiler_root_user
DB=sqlboiler_model_test
PASS=sqlboiler
MSSQLPASS=Sqlboiler@1234

DRIVER_USER=sqlboiler_driver_user
DRIVER_DB=sqlboiler_driver_test

# Builds all software and runs model tests
tests: test-psql test-mysql test-mssql
# Creates super users
test-users: test-user-psql test-user-mysql test-user-mssql
# Creates databases
test-dbs: test-db-psql test-db-mysql test-db-mssql

# Runs tests on the drivers
driver-tests: driver-test-psql driver-test-mysql driver-test-mssql
# Creates regular databases
driver-dbs: driver-db-psql driver-db-mysql driver-db-mssql
# Creates regular users with access to only one DB, may modify the DB
# to ensure write access to the user
driver-users: driver-user-psql driver-user-mysql driver-user-mssql

run-mssql:
	docker run --detach --rm --env 'ACCEPT_EULA=Y' --env 'SA_PASSWORD=Sqlboiler@1234' --publish 1433:1433 --name mssql microsoft/mssql-server-linux:2017-latest
kill-mssql:
	docker rm --force mssql

# ====================================
# Building operations
# ====================================

build:
	go build github.com/volatiletech/sqlboiler

build-psql:
	go build github.com/volatiletech/sqlboiler/drivers/sqlboiler-psql

build-mysql:
	go build github.com/volatiletech/sqlboiler/drivers/sqlboiler-mysql

build-mssql:
	go build github.com/volatiletech/sqlboiler/drivers/sqlboiler-mssql

# ====================================
# Testing operations
# ====================================

test-psql: build build-psql
	./sqlboiler --wipe psql
	go test ./models

test-mysql: build build-mysql
	./sqlboiler --wipe mysql
	go test ./models

test-mssql: build build-mssql
	./sqlboiler --wipe mssql
	go test ./models

test-user-psql:
	# Can't use createuser because it interactively promtps for a password
	psql --command "create user $(USER) with superuser password '$(PASS)';" postgres

test-user-mysql:
	mysql --execute "create user $(USER) identified by '$(PASS)';" 
	mysql --execute "grant all privileges on *.* to $(USER);" 

# Must be run after a database is created
test-user-mssql:
	sqlcmd -S localhost -U sa -P $(MSSQLPASS) -Q "create login $(USER) with password = '$(MSSQLPASS)';"
	sqlcmd -S localhost -U sa -P $(MSSQLPASS) -Q "alter server role sysadmin add member $(USER);"

test-db-psql:
	env PGPASSWORD=$(PASS) dropdb --username $(USER) --if-exists $(DB)
	env PGPASSWORD=$(PASS) createdb --owner $(USER) --username $(USER) $(DB)
	env PGPASSWORD=$(PASS) psql --username $(USER) --file testdata/psql_test_schema.sql $(DB)

test-db-mysql:
	mysql --user $(USER) --password=$(PASS) --execute "drop database if exists $(DB);"
	mysql --user $(USER) --password=$(PASS) --execute "create database $(DB);"
	mysql --user $(USER) --password=$(PASS) $(DB) < testdata/mysql_test_schema.sql

test-db-mssql:
	sqlcmd -S localhost -U $(USER) -P $(MSSQLPASS) -Q "drop database if exists $(DB)";
	sqlcmd -S localhost -U $(USER) -P $(MSSQLPASS) -Q "create database $(DB)";
	sqlcmd -S localhost -U $(USER) -P $(MSSQLPASS) -d $(DB) -i testdata/mssql_test_schema.sql

# ====================================
# Driver operations
# ====================================

driver-test-psql:
	cd drivers/sqlboiler-psql/driver && go test
driver-test-mysql:
	go test github.com/volatiletech/sqlboiler/drivers/sqlboiler-mysql/driver
driver-test-mssql:
	go test github.com/volatiletech/sqlboiler/drivers/sqlboiler-mssql/driver

driver-db-psql:
	createdb $(DRIVER_DB)

driver-db-mysql:
	mysql --execute "create database $(DRIVER_DB);"

driver-db-mssql:
	sqlcmd -S localhost -U sa -P $(MSSQLPASS) -Q "create database $(DRIVER_DB);"
	sqlcmd -S localhost -U sa -P $(MSSQLPASS) -d $(DRIVER_DB) -Q "exec sp_configure 'contained database authentication', 1;"
	sqlcmd -S localhost -U sa -P $(MSSQLPASS) -d $(DRIVER_DB) -Q "reconfigure"
	sqlcmd -S localhost -U sa -P $(MSSQLPASS) -d $(DRIVER_DB) -Q "alter database $(DRIVER_DB) set containment = partial;"

driver-user-psql:
	psql --command "create role $(DRIVER_USER) login nocreatedb nocreaterole nocreateuser password '$(PASS)';" $(DRIVER_DB)
	psql --command "alter database $(DRIVER_DB) owner to $(DRIVER_USER);" $(DRIVER_DB)

driver-user-mysql:
	mysql --execute "create user $(DRIVER_USER) identified by '$(PASS)';" 
	mysql --execute "grant all privileges on $(DRIVER_DB).* to $(DRIVER_USER);" 

driver-user-mssql:
	sqlcmd -S localhost -U sa -P $(MSSQLPASS) -d $(DRIVER_DB) -Q "create user $(DRIVER_USER) with password = '$(MSSQLPASS)'";
	sqlcmd -S localhost -U sa -P $(MSSQLPASS) -d $(DRIVER_DB) -Q "grant alter, control to $(DRIVER_USER)";
	
# ====================================
# Clean operations
# ====================================

clean:
	rm ./sqlboiler
	rm ./sqlboiler-psql
	rm ./sqlboiler-mysql
	rm ./sqlboiler-mssql
