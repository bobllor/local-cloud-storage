#!/usr/bin/env bash

# Builds the database locally, it creates the tables and users
# needed for the database.
# If Docker is preferred, read the Docker documentation.

sql(){
    sudo mysql -v -e "$1"
}

sql_silent(){
    sudo mysql -e "$1"
}

set -e

env=".env"
sql_script_path="sql/01.db_setup.sql"

if [[ ! -e "$env" ]]; then
    echo "error: $env not found"
    exit 1
fi

source "$env"

db_pw=${DB_PASSWORD}

sql "source $sql_script_path"

# sql user creation
sql_silent "CREATE USER IF NOT EXISTS 'file_user'@'localhost' IDENTIFIED BY '$db_pw';"
sql "GRANT SELECT, DELETE, INSERT ON MasterLocalCloudStorage.Files TO 'file_user'@'localhost';"