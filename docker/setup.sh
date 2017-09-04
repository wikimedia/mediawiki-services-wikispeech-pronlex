#!/bin/bash 

CMD=`basename $0`
SCRIPTDIR=`dirname $0`
export GOPATH=`go env GOPATH`
export PATH=$PATH:$GOPATH/bin

if [ $# -ne 1 ]; then
    echo "USAGE: sh $CMD <APPDIR>

Setup files will be added to the destination folder APPDIR.
" >&2
    exit 1
fi

APPDIR=$1

if [ -d $APPDIR/symbol_sets ] ; then
    echo "[$CMD] $APPDIR is already configured. No setup needed." >&2
    exit 0
fi


if [ -z "$GOPATH" ] ; then
    echo "[$CMD] The GOPATH environment variable is required!" >&2
    exit 1
fi

### LEXDATA PREPS

echo "[$CMD] Setting up basic files ... " >&2

mkdir -p $APPDIR || exit 1
mkdir -p $APPDIR/symbol_sets || exit 1

cp $GOPATH/src/github.com/stts-se/pronlex/lexserver/demo_files/*.sym $APPDIR/symbol_sets/ || exit 1
cp $GOPATH/src/github.com/stts-se/pronlex/lexserver/demo_files/*.cnv $APPDIR/symbol_sets/ || exit 1
echo "" >> $APPDIR/symbol_sets/mappers.txt || exit 1
cat $GOPATH/src/github.com/stts-se/pronlex/lexserver/demo_files/mappers.txt >> $APPDIR/symbol_sets/mappers.txt || exit 1


### COMPLETED

echo "
BUILD COMPLETED! YOU CAN NOW START THE LEXICON SERVER BY INVOKING:
  $ sh $SCRIPTDIR/start_server.sh -a $APPDIR

  USAGE INFO:
  $ sh $SCRIPTDIR/start_server.sh -h


OR IMPORT STANDARD LEXICON DATA:
  $ sh $SCRIPTDIR/import.sh <LEXDATA-GIT> $APPDIR

  USAGE INFO:
  $ sh $SCRIPTDIR/import.sh

" >&2
