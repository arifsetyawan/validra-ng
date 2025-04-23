#!/bin/bash

# Fix checksum issues in Atlas migrations
echo "Fixing checksums for migrations..."
atlas migrate hash --env dev

echo "Checksums have been recalculated."
echo "Now you can run migrations with: make atlas-migrate-apply" 