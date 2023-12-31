#!/bin/bash

echo "Helps generate files in this folder like 20220325123456_create_users.sql"
echo "Enter the name of the name for the param 20220325123456_[create_users].sql"
read first

# Get the current timestamp in the format YYYYMMDDHHMMSS
timestamp=$(date "+%Y%m%d%H%M%S")

filename="$timestamp"_"$first".sql

touch "$filename"

echo "PRAGMA USER_VERSION = X;" >> "$filename"

echo "Migration file created: $filename"
