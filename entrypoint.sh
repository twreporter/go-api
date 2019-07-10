#!/bin/sh

# Wait for mysql server to be ready
sh -c "dockerize -wait tcp://$GOAPI_DB_MYSQL_ADDRESS:$GOAPI_DB_MYSQL_PORT -timeout 30s"

# Upgrade MEMBERSHIP Database to latest version
migrate -database "mysql://$GOAPI_DB_MYSQL_MEMBERSHIP_MIGRATE_USER:$GOAPI_DB_MYSQL_MEMBERSHIP_MIGRATE_PASSWORD@tcp($GOAPI_DB_MYSQL_ADDRESS:$GOAPI_DB_MYSQL_PORT)/$GOAPI_DB_MYSQL_NAME" -path $MIGRATION_DIR up

# Execute pangolin
exec "$@"
