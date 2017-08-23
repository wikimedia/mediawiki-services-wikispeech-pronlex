#!/bin/bash 

CMD=`basename $0`

#############################################################
### SERVER STARTUP SCRIPT
## 1. STARTS A TEST SERVER AND RUNS A SET OF TESTS
## 2. SHUTS DOWN TEST SERVER
## 3. IF NO ERRORS ARE FOUND, THE STANDARD SERVER IS STARTED USING THE DEFAULT PORT

if [ $# -ne 1 ]; then
    echo "[$CMD] SERVER STARTUP SCRIPT
   1. STARTS A TEST SERVER AND RUNS A SET OF TESTS
   2. SHUTS DOWN TEST SERVER
   3. IF NO ERRORS ARE FOUND, THE STANDARD SERVER IS STARTED USING THE DEFAULT PORT

USAGE: $CMD <APPDIR>
       where <APPDIR> is the folder in which the build script installed the standalone lexserver

EXAMPLE INVOCATION: $CMD lexserver_files
" >&2
    exit 1
fi


APPDIR=$1

if [ -z "$GOPATH" ] ; then
    echo "[$CMD] The GOPATH environment variable is required!" >&2
    exit 1
fi


cd $GOPATH/src/github.com/stts-se/pronlex/lexserver

APPDIR=$1
switches="-ss_files $APPDIR/symbol_sets/ -db_files $APPDIR/db_files/"
go run *.go $switches -test && go run *.go $switches
