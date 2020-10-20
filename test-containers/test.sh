#!/bin/bash

function run_tests () {

    echo "Testing if env vars are still set..."
    [[ "$(docker-compose exec timer sh -c 'echo -n $TESTVAR')" == "stuff" ]] || return 2

    echo "Testing if mounts are still valid..."
    [[ "$(docker-compose exec timer cat /test-volume/test-file-2)"  == "stufferson" ]] || return 3

    echo "Testing if the network still works..."
    docker-compose exec test curl timer:8998 > /dev/null || return 1

    echo "Testing if exposed ports are still exposed..."  
    curl localhost:9812 > /dev/null 2>&1 || return 1

    return 0
}

if run_tests; then
    echo "Success!"
else 
    echo "Failed!"
fi


