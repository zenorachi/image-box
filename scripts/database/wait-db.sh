#!/bin/bash
# Wait for the database to become available

cmd="${!#}"
host="$1"

until pg_isready -h "$host" -U postgres; do
    echo "Waiting for the database to become available..."
    sleep 1
done

echo "DB started"
exec "$cmd"