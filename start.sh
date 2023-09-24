# !/bin/sh

set -e

# echo "run db migration"
# /app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up

echo "starting the app..."
exec "$@" # take all parameters and run in this script