#!/bin/bash
# Wait for the database to become available

until pg_isready -h "$DB_HOST" -U "$DB_USERNAME"; do
    echo "Waiting for the database to become available..."
    sleep 1
done

echo "DB started"