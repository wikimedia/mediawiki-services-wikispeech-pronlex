#!/bin/bash 

CMD=`basename $0`

if [ $# -ne 1 ]; then
    echo "USAGE: sh $CMD <APPDIR>" >&2
    exit 1
fi

APPDIR=$1

if [ -z "$GOPATH" ] ; then
    echo "[$CMD] The GOPATH environment variable is required!" >&2
    exit 1
fi


## LEXSERVER INSTALL

echo "[$CMD] Installing pronlex/lexserver ... " >&2

go get github.com/stts-se/pronlex/lexserver
go install github.com/stts-se/pronlex/lexserver

go get github.com/stts-se/pronlex/cmd/lexio/importLex
go install github.com/stts-se/pronlex/cmd/lexio/importLex

go get github.com/stts-se/pronlex/cmd/lexio/createEmptyDB
go install github.com/stts-se/pronlex/cmd/lexio/createEmptyDB

### LEXDATA PREPS

echo "[$CMD] Setting up basic files ... " >&2

mkdir -p $APPDIR

mkdir -p $APPDIR/static
cp -r $GOPATH/src/github.com/stts-se/pronlex/lexserver/static/* $APPDIR/static/

mkdir -p $APPDIR/symbol_sets
cp $GOPATH/src/github.com/stts-se/pronlex/lexserver/demo_files/* $APPDIR/symbol_sets

echo "[$CMD] Setting up scripts ... " >&2

cp $GOPATH/src/github.com/stts-se/pronlex/install/standalone/import.sh $APPDIR
cp $GOPATH/src/github.com/stts-se/pronlex/install/standalone/start_server.sh $APPDIR


### COMPLETED

echo "
BUILD COMPLETED. YOU CAN NOW START THE LEXICON SERVER BY INVOKING
  $ sh $APPDIR/start_server.sh
" >&2

