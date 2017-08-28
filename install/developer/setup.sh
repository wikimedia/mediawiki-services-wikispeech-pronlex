#!/bin/bash 

CMD=`basename $0`
SCRIPTDIR=`dirname $0`
GOPATH=`go env GOPATH`

if [ $# -ne 1 ]; then
    echo "USAGE: sh $CMD <APPDIR>

Setup files will be added to the destination folder APPDIR.
" >&2
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
BUILD COMPLETED! YOU CAN NOW START THE LEXICON SERVER BY INVOKING:
  $ sh $SCRIPTDIR/start_server.sh $APPDIR

  USAGE INFO:
  $ sh $SCRIPTDIR/start_server.sh -h


OR IMPORT STANDARD LEXICON DATA:
  $ sh $SCRIPTDIR/import.sh <LEXDATA-GIT> $APPDIR

  USAGE INFO:
  $ sh $SCRIPTDIR/import.sh

" >&2
