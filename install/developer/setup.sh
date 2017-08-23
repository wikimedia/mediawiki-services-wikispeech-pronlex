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

### LEXDATA PREPS

echo "[$CMD] Setting up basic files ... " >&2

mkdir -p $APPDIR || exit 1

mkdir -p $APPDIR/symbol_sets || exit 1
cp $GOPATH/src/github.com/stts-se/pronlex/lexserver/demo_files/* $APPDIR/symbol_sets || exit 1


### COMPLETED

echo "
BUILD COMPLETED. YOU CAN NOW START THE LEXICON SERVER BY INVOKING
  $ sh start_server.sh $APPDIR
" >&2
