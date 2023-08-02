#!/bin/bash
# Wait for the database to become available

host="$1"
#cmd="$2"

until pg_isready -h "$host" -U postgres; do
    echo "Waiting for the database to become available..."
    sleep 1
done

echo "DB started"
exec "/dlv exec image-box-app"