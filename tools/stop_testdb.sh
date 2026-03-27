#!/usr/bin/env bash

set -e

source "./tools/support/testdb/sql_vars.sh"
source "./tools/support/testdb/sql_utils.sh"

container_status=$(check_container_status $c_name)

if [[ "$container_status" == "true" ]]; then
    echo "Stopping container $c_name..."

    stop_docker_container $c_name
else
    echo "error: Container $c_name is not running"
fi