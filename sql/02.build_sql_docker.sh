#!/usr/bin/env bash

# Sets up the tables and the users for the database.
# This must be ran after 01.db_setup.sql, the .env variables is expected to be loaded 
# in docker-compose.yml, and that this is intended to be used only with Docker.

set -e

root_password=${MARIADB_ROOT_PASSWORD}

sql(){
    mariadb -u root -p$root_password -v -e "$1"
}

sql_silent(){
    mariadb -u root -p$root_password -e "$1"
}

db_pw=${FILE_PASSWORD}

db="MasterLocalCloudStorage"
file_table="File"
user_table="UserAccount"

file_user="file_user"

# sql user creation
sql_silent "CREATE USER IF NOT EXISTS '$file_user'@'%' IDENTIFIED BY '$db_pw';"
sql "GRANT SELECT, DELETE, INSERT ON $db.$file_table TO '$file_user'@'%';"