#!/usr/bin/env bash

set -e

c_name="testdb"
host_port="3307"
c_port="3306"
sql_path="./sql/scripts/00.testdb_setup.sql"

# executes SQL commands
sql(){
    mariadb -h 127.0.0.1 -P $host_port -u root -e "$1"
}

stop_container(){
    docker container stop "$c_name" > /dev/null
    docker container rm "$c_name" > /dev/null
}

container_exists=$(docker container ls | grep "$c_name" && echo true || echo false)

if [[ "$container_exists" == "false" ]]; then
    echo "Starting container $c_name"

    docker run --detach \
        --name "$c_name" \
        -p "$host_port:$c_port" \
        --env MARIADB_ALLOW_EMPTY_ROOT_PASSWORD=1 \
        mariadb:lts > /dev/null

    init_status=false
    echo "Waiting for server connection..."
    for i in {1..7}; do
        if ! sql "source $sql_path" 2> /dev/null; then
            sleep 1.5
            continue
        else
            echo "Server connection establised"
            echo "Database initialized"
            init_status=true
            break
        fi
    done

    if [[ $init_status == false ]]; then
        echo "error: Failed to establish connection to server, stopping $c_name"
        stop_container
        exit 1
    fi

    echo "Container $c_name started"
else
    echo "Stopping container $c_name..."

    stop_container

    echo "Container $c_name stopped"
fi