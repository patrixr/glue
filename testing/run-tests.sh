#!/usr/bin/env bash

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
ROOT_DIR="$SCRIPT_DIR/.."

cd $ROOT_DIR

log() {
    echo "$(date '+%Y-%m-%d %H:%M:%S') - $1"
}

function assert {
    local cmd="$1"
    $cmd || (echo "Command failed, exiting" && exit 1)
}

function start_container {
    log "Starting container"
    docker run --rm  --name glue-tests -w /app -v ${ROOT_DIR}:/app golang:latest make test
}

function stop_container {
    log "Stopping container"
    docker compose -f "$SCRIPT_DIR/docker-compose.yml" down -v > /dev/null
}

function  check_container_health {
    local health_status
    health_status=$(docker inspect -f '{{.State.Health.Status}}' "glue-test-container" 2>/dev/null)
    if [ "$health_status" = "healthy" ]; then
        return 0
    else
        return 1
    fi
}

start_container

# if [ $RESULT -ne 0 ]; then
#     log "Tests failed"
# else
#     log "Tests completed successfully"
# fi

# exit $RESULT
