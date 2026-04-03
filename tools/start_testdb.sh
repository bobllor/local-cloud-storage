#!/usr/bin/env bash

set -e

source "./tools/support/testdb/sql_vars.sh"
source "./tools/support/testdb/sql_utils.sh"

container_status=$(check_container_status $c_name)
declare -i max_attempts=7

if [[ "$container_status" == "false" ]]; then
    echo "Starting container $c_name"

    # volume will also share the container name
    docker volume create "$c_name" | xargs -I x echo "Created test volume x"

    docker run --detach \
        --name "$c_name" \
        -p "$host_port:$c_port" \
        --env MYSQL_ALLOW_EMPTY_PASSWORD=yes \
        --mount type=volume,src=$c_name,dst=/var/lib/mysql \
        --mount type=bind,src=/etc/timezone,dst=/etc/timezone,readonly \
        --mount type=bind,src=/etc/localtime,dst=/etc/localtime,readonly \
        mysql:lts-oracle 2>&1

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
        stop_docker_container $c_name
        exit 1
    fi

    echo "Container $c_name started"
else
    echo "error: Container $c_name is already running"
fi