#!/bin/sh
set -e

echo "migrations started"
/app/bin/migrate -command up
echo "migrations completed"

echo "api server started"
exec /app/bin/api
