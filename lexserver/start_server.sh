#!/bin/bash

extragofiles=`ls *go|egrep -v lexserver.go`

# start server for testing
>&2 echo "[$0] STARTING SERVER FOR TESTING"
go run lexserver.go $extragofiles 8799 TEST &
    
pid=`echo $!`
>&2 echo "[$0] SERVER FOR TESTING RUNNING WITH PID $pid"

sleep 1

kill_tests() {
    # shutdown server
    >&2 echo "[$0] SHUTTING DOWN TEST SERVER WITH PID $pid"
    kill -9 $pid
}


# run tests
>&2 echo "[$0] RUNNING TESTS ..."

if go run demotest/demotest.go 8799; then
    echo -n "" # OK
else 
    >&2 echo "[$0] TESTS FAILED!"
    kill_tests
    exit 1
fi

>&2 echo "[$0] TESTS COMPLETED ..."

kill_tests

# start server again
>&2 echo "[$0] STARTING STANDARD SERVER ..."
go run lexserver.go $extragofiles $*
