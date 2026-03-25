#!/usr/bin/env bash

# NOTE: this does require sudo.

sql(){
    sudo mysql -v -e "$1"
}

sql_silent(){
    sudo mysql -e "$1"
}

set -e

env=".env"
sql_script_path="sql/db_setup.sql"

if [[ ! -e "$env" ]]; then
    echo ".env file required"
    exit 1
fi

source "$env"

db_pw=${DB_PASSWORD}

sql "source $sql_script_path"

# sql user creation
sql_silent "CREATE USER IF NOT EXISTS 'file_user'@'localhost' IDENTIFIED BY '$db_pw';"
sql "GRANT SELECT, DELETE, INSERT ON MasterLocalCloudStorage.Files TO 'file_user'@'localhost';"