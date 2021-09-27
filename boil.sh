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
    driver=$1

    buildPath=github.com/volatiletech/sqlboiler/v4
    case "${driver}" in
        all)
            set -o xtrace
            go build "${buildPath}"
            { set +o xtrace; } 2>/dev/null
            drivers="psql mysql mssql"
            shift ;;
        psql)
            drivers="psql"
            shift ;;
        mysql)
            drivers="mysql"
            shift ;;
        mssql)
            drivers="mssql"
            shift ;;
        *)
            set -o xtrace
            go build "$@" "${buildPath}"
            { set +o xtrace; } 2>/dev/null
            return ;;
    esac

    for d in $drivers; do
        set -o xtrace
        go build "$@" ${buildPath}/drivers/sqlboiler-"${d}"
        { set +o xtrace; } 2>/dev/null
    done
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
        mssql-docker)
            set -o xtrace
            docker exec --interactive --tty mssql /opt/mssql-tools/bin/sqlcmd -S localhost -U sa -P "${MSSQLPASS}" -Q "create login ${DB_USER} with password = '${MSSQLPASS}';"
            docker exec --interactive --tty mssql /opt/mssql-tools/bin/sqlcmd -S localhost -U sa -P "${MSSQLPASS}" -Q "alter server role sysadmin add member ${DB_USER};"
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
        mssql-docker)
            set -o xtrace
            docker exec --interactive --tty mssql /opt/mssql-tools/bin/sqlcmd -S localhost -U ${DB_USER} -P ${MSSQLPASS} -Q "drop database if exists ${DB_NAME}";
            docker exec --interactive --tty mssql /opt/mssql-tools/bin/sqlcmd -S localhost -U ${DB_USER} -P ${MSSQLPASS} -Q "create database ${DB_NAME}";
            docker exec --interactive --tty mssql /opt/mssql-tools/bin/sqlcmd -S localhost -U ${DB_USER} -P ${MSSQLPASS} -d ${DB_NAME} -i testdata/mssql_test_schema.sql
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
            sqlcmd -S localhost -U sa -P "${MSSQLPASS}" -d "${DRIVER_DB}" -Q "create user ${DRIVER_USER} with password = '${MSSQLPASS}'";
            sqlcmd -S localhost -U sa -P "${MSSQLPASS}" -d "${DRIVER_DB}" -Q "grant alter, control to ${DRIVER_USER}";
            ;;
        mssql-docker)
            set -o xtrace
            docker exec --interactive --tty mssql /opt/mssql-tools/bin/sqlcmd -S localhost -U sa -P "${MSSQLPASS}" -d "${DRIVER_DB}" -Q "create user ${DRIVER_USER} with password = '${MSSQLPASS}'";
            docker exec --interactive --tty mssql /opt/mssql-tools/bin/sqlcmd -S localhost -U sa -P "${MSSQLPASS}" -d "${DRIVER_DB}" -Q "grant alter, control to ${DRIVER_USER}";
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
            sqlcmd -S localhost -U sa -P "${MSSQLPASS}" -Q "create database ${DRIVER_DB};"
            sqlcmd -S localhost -U sa -P "${MSSQLPASS}" -d "${DRIVER_DB}" -Q "exec sp_configure 'contained database authentication', 1;"
            sqlcmd -S localhost -U sa -P "${MSSQLPASS}" -d "${DRIVER_DB}" -Q "reconfigure"
            sqlcmd -S localhost -U sa -P "${MSSQLPASS}" -d "${DRIVER_DB}" -Q "alter database ${DRIVER_DB} set containment = partial;"
            ;;
        mssql-docker)
            set -o xtrace
            docker exec --interactive --tty mssql /opt/mssql-tools/bin/sqlcmd -S localhost -U sa -P "${MSSQLPASS}" -Q "create database ${DRIVER_DB};"
            docker exec --interactive --tty mssql /opt/mssql-tools/bin/sqlcmd -S localhost -U sa -P "${MSSQLPASS}" -d "${DRIVER_DB}" -Q "exec sp_configure 'contained database authentication', 1;"
            docker exec --interactive --tty mssql /opt/mssql-tools/bin/sqlcmd -S localhost -U sa -P "${MSSQLPASS}" -d "${DRIVER_DB}" -Q "reconfigure"
            docker exec --interactive --tty mssql /opt/mssql-tools/bin/sqlcmd -S localhost -U sa -P "${MSSQLPASS}" -d "${DRIVER_DB}" -Q "alter database ${DRIVER_DB} set containment = partial;"
            ;;
        *)
            echo "unknown driver"
            ;;
    esac
}

# ====================================
# MSSQL stuff
# ====================================

mssql_run() {
    if test "attach" = "${1}"; then
        args="--interactive --tty"
    else
        args="--detach"
    fi

    set -o xtrace

    docker run $args --rm \
        --env 'ACCEPT_EULA=Y' --env "SA_PASSWORD=${MSSQLPASS}" \
        --publish 1433:1433 \
        --volume "${PWD}/testdata/mssql_test_schema.sql:/testdata/mssql_test_schema.sql" \
        --name mssql \
        mcr.microsoft.com/mssql/server:2019-latest
}

mssql_stop() {
    set -o xtrace
    docker stop mssql
}

mssql_sqlcmd() {
    set -o xtrace
    docker exec --interactive --tty mssql /opt/mssql-tools/bin/sqlcmd -S localhost -U sa -P "${MSSQLPASS}" -Q "$@"
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

    mssql-run)    mssql_run "$@" ;;
    mssql-stop)   mssql_stop "$@" ;;
    mssql-sqlcmd) mssql_sqlcmd "$@" ;;

    clean) clean ;;
    *)
        echo "./boil.sh command [args]"
        echo "Helpers to build, test and develop on sqlboiler"
        echo
        echo "Commands:"
        echo "  build [all|driver]          - builds the driver (or all drivers if all), or sqlboiler if omitted"
        echo "  gen <driver> [args]         - generates models for driver, passing args along to sqlboiler"
        echo "  test [args]                 - runs model tests, passing args to the test binary"
        echo "  test-user <driver>          - creates a user for the tests (superuser)"
        echo "  test-db <driver>            - recreates test datbase for driver using test-user"
        echo "  driver-test <driver> [args] - runs tests for the driver"
        echo "  driver-test-db <driver>     - create driver db (run before driver-test-user)"
        echo "  driver-test-user <driver>   - creates a user for the driver tests (unprivileged)"
        echo "  mssql-run [attach]          - run mssql docker container, if attach is present will not daemonize"
        echo "  mssql-stop                  - stop mssql docker container"
        echo "  mssql-sqlcmd [args...]      - run sql query using sqlcmd in docker container"
esac
