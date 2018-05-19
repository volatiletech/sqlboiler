#!/bin/sh

DB_USER="sqlboiler_root_user"
DB_NAME="sqlboiler_model_test"
DB_PASS="sqlboiler"
MSSQLPASS="Sqlboiler@1234"

DRIVER_USER="sqlboiler_driver_user"
DRIVER_DB="sqlboiler_driver_test"

# ====================================
# Building
# ====================================

build() {
    subcommand=$1

    case "${subcommand}" in
        psql)  driver=1; shift ;;
        mysql) driver=1; shift ;;
        mssql) driver=1; shift ;;
    esac

    path=github.com/volatiletech/sqlboiler
    if test "${driver}"; then
        path="${path}/drivers/sqlboiler-${subcommand}"
    fi

    set -o xtrace
    go build "$@" ${path}
}

# ====================================
# Generation
# ====================================

gen() {
    db=$1
    shift
    set -o xtrace
    ./sqlboiler "$@" --wipe "${db}"
}

# ====================================
# Testing
# ====================================

runtest() {
    set -o xtrace
    go test -v -race "$@" ./models
}

# ====================================
# Test users
# ====================================

test_user() {
    driver=$1

    case "${driver}" in
        psql)
            # Can't use createuser because it interactively promtps for a password
            set -o xtrace
            psql --host localhost --username postgres --command "create user ${DB_USER} with superuser password '${DB_PASS}';"
            ;;
        mysql)
            set -o xtrace
            mysql --host localhost --execute "create user ${DB_USER} identified by '${DB_PASS}';" 
            mysql --host localhost --execute "grant all privileges on *.* to ${DB_USER};" 
            ;;
        mssql)
            set -o xtrace
            sqlcmd -S localhost -U sa -P "${MSSQLPASS}" -Q "create login ${DB_USER} with password = '${MSSQLPASS}';"
            sqlcmd -S localhost -U sa -P "${MSSQLPASS}" -Q "alter server role sysadmin add member ${DB_USER};"
            ;;
        *)
            echo "unknown driver"
            ;;
    esac
}

# ====================================
# Test DBs
# ====================================

test_db() {
    driver=$1

    case "${driver}" in
        psql)
            set -o xtrace
            env PGPASSWORD=${DB_PASS} dropdb --host localhost --username ${DB_USER} --if-exists ${DB_NAME}
            env PGPASSWORD=${DB_PASS} createdb --host localhost --owner ${DB_USER} --username ${DB_USER} ${DB_NAME}
            env PGPASSWORD=${DB_PASS} psql --host localhost --username ${DB_USER} --file testdata/psql_test_schema.sql ${DB_NAME}
            ;;
        mysql)
            set -o xtrace
            mysql --host localhost --user ${DB_USER} --password=${DB_PASS} --execute "drop database if exists ${DB_NAME};"
            mysql --host localhost --user ${DB_USER} --password=${DB_PASS} --execute "create database ${DB_NAME};"
            mysql --host localhost --user ${DB_USER} --password=${DB_PASS} ${DB_NAME} < testdata/mysql_test_schema.sql
            ;;
        mssql)
            set -o xtrace
            sqlcmd -S localhost -U ${DB_USER} -P ${MSSQLPASS} -Q "drop database if exists ${DB_NAME}";
            sqlcmd -S localhost -U ${DB_USER} -P ${MSSQLPASS} -Q "create database ${DB_NAME}";
            sqlcmd -S localhost -U ${DB_USER} -P ${MSSQLPASS} -d ${DB_NAME} -i testdata/mssql_test_schema.sql
            ;;
        *)
            echo "unknown driver"
            ;;
    esac
}

# ====================================
# Driver test
# ====================================

driver_test() {
    driver=$1
    shift

    cd "drivers/sqlboiler-${driver}/driver"
    set -o xtrace
    go test "$@"
}

# ====================================
# Driver test users
# ====================================

driver_test_user() {
    driver=$1

    case "${driver}" in
        psql)
            set -o xtrace
            env PGPASSWORD=${DB_PASS} psql --host localhost --username ${DB_USER} --command "create role ${DRIVER_USER} login nocreatedb nocreaterole password '${DB_PASS}';" ${DRIVER_DB}
            env PGPASSWORD=${DB_PASS} psql --host localhost --username ${DB_USER} --command "alter database ${DRIVER_DB} owner to ${DRIVER_USER};" ${DRIVER_DB}
            ;;
        mysql)
            set -o xtrace
            mysql --host localhost --execute "create user ${DRIVER_USER} identified by '${DB_PASS}';" 
            mysql --host localhost --execute "grant all privileges on ${DRIVER_DB}.* to ${DRIVER_USER};" 
            ;;
        mssql)
            set -o xtrace
            sqlcmd -S localhost -U sa -P ${MSSQLPASS} -d ${DRIVER_DB} -Q "create user ${DRIVER_USER} with password = '${MSSQLPASS}'";
            sqlcmd -S localhost -U sa -P ${MSSQLPASS} -d ${DRIVER_DB} -Q "grant alter, control to ${DRIVER_USER}";
            ;;
        *)
            echo "unknown driver"
            ;;
    esac
}

# ====================================
# Driver test databases
# ====================================

driver_test_db() {
    driver=$1

    case "${driver}" in
        psql)
            set -o xtrace
            env PGPASSWORD=${DB_PASS} createdb --host localhost --username ${DB_USER} --owner ${DB_USER} ${DRIVER_DB}
            ;;
        mysql)
            set -o xtrace
            mysql --host localhost --user ${DB_USER} --password=${DB_PASS} --execute "create database ${DRIVER_DB};"
            ;;
        mssql)
            set -o xtrace
            sqlcmd -S localhost -U sa -P ${MSSQLPASS} -Q "create database ${DRIVER_DB};"
            sqlcmd -S localhost -U sa -P ${MSSQLPASS} -d ${DRIVER_DB} -Q "exec sp_configure 'contained database authentication', 1;"
            sqlcmd -S localhost -U sa -P ${MSSQLPASS} -d ${DRIVER_DB} -Q "reconfigure"
            sqlcmd -S localhost -U sa -P ${MSSQLPASS} -d ${DRIVER_DB} -Q "alter database ${DRIVER_DB} set containment = partial;"
            ;;
        *)
            echo "unknown driver"
            ;;
    esac
}

# ====================================
# Clean
# ====================================

clean() {
    set -o xtrace
    rm -f ./sqlboiler
    rm -f ./sqlboiler-*
}

# ====================================
# CLI
# ====================================

command=$1
shift

case "${command}" in
    build)     build "$@" ;;
    gen)       gen "$@" ;;

    test)      runtest "$@" ;;
    test-user) test_user "$1" ;;
    test-db)   test_db "$1" ;;

    driver-test)      driver_test "$@" ;;
    driver-test-db)   driver_test_db "$1" ;;
    driver-test-user) driver_test_user "$1" ;;

    clean) clean ;;
    *)
        echo "./boil.sh command [args]"
        echo "Helpers to build, test and develop on sqlboiler"
        echo
        echo "Commands:"
        echo "  build [driver]              - builds the driver, or sqlboiler if driver omitted"
        echo "  gen <driver> [args]         - generates models for driver, passing args along to sqlboiler"
        echo "  test [args]                 - runs model tests, passing args to the test binary"
        echo "  test-user <driver>          - creates a user for the tests (superuser)"
        echo "  test-db <driver>            - recreates test datbase for driver using test-user"
        echo "  driver-test <driver> [args] - runs tests for the driver"
        echo "  driver-test-db <driver>     - create driver db (run before driver-test-user)"
        echo "  driver-test-user <driver>   - creates a user for the driver tests (unprivileged)"
esac
