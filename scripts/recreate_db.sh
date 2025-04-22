#!/bin/bash

# Load environment variables
if [ -f .env ]; then
  export $(cat .env | grep -v '^#' | xargs)
else
  echo "No .env file found. Using default values."
  export DB_NAME="validra"
  export DB_USER="postgres"
  export DB_PASSWORD="postgres"
  export DB_HOST="localhost"
  export DB_PORT="5432"
fi

# Confirm action
echo "This will drop and recreate the database: $DB_NAME"
echo "All data will be lost. Are you sure? (y/n)"
read confirmation

if [[ $confirmation != "y" && $confirmation != "Y" ]]; then
  echo "Operation canceled."
  exit 0
fi

# Drop and recreate database
psql -h $DB_HOST -p $DB_PORT -U $DB_USER -c "DROP DATABASE IF EXISTS $DB_NAME;" postgres
psql -h $DB_HOST -p $DB_PORT -U $DB_USER -c "CREATE DATABASE $DB_NAME;" postgres

echo "Database $DB_NAME has been recreated."

# Apply schema directly
echo "Applying schema via Atlas..."
atlas schema apply --env dev --auto-approve

echo "Database schema has been applied."
echo "You can now run migrations with: make atlas-migrate-apply" 