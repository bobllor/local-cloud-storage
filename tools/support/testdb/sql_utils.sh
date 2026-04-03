#!/usr/bin/env bash

# sql executes SQL commands passed in in the first argument.
sql(){
    mysql -h 127.0.0.1 -P $host_port -u root -e "$1"
}

# check_container_status checks for the container status and returns true or false.
check_container_status(){
    docker container ls | grep "$1" > /dev/null && echo true || echo false
}

# stop_docker_container stops a given container and removes it.
stop_docker_container(){
    docker container stop "$1" > /dev/null
    docker container rm "$1" > /dev/null

    echo "Container $c_name stopped"
}