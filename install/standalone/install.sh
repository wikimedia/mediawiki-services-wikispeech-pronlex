#!/bin/bash 

CMD=`basename $0`

if [ $# -ne 1 ]; then
    echo "USAGE: sh $CMD <APPDIR>

 Pronlex will be installed in the destination folder APPDIR.
">&2
    exit 1
fi

APPDIR=$1

if [ -z "$GOPATH" ] ; then
    echo "[$CMD] The GOPATH environment variable is required!" >&2
    exit 1
fi


## LEXSERVER INSTALL

echo "[$CMD] Installing pronlex/lexserver ... " >&2

go get github.com/stts-se/pronlex/lexserver || exit 1
go install github.com/stts-se/pronlex/lexserver || exit 1

go get github.com/stts-se/pronlex/cmd/lexio/importLex || exit 1
go install github.com/stts-se/pronlex/cmd/lexio/importLex || exit 1

go get github.com/stts-se/pronlex/cmd/lexio/createEmptyDB || exit 1
go install github.com/stts-se/pronlex/cmd/lexio/createEmptyDB || exit 1

### LEXDATA PREPS

echo "[$CMD] Setting up basic files ... " >&2

mkdir -p $APPDIR || exit 1

mkdir -p $APPDIR/.static || exit 1
cp -r $GOPATH/src/github.com/stts-se/pronlex/lexserver/static/* $APPDIR/.static || exit 1

mkdir -p $APPDIR/symbol_sets || exit 1
cp $GOPATH/src/github.com/stts-se/pronlex/lexserver/demo_files/* $APPDIR/symbol_sets || exit 1

echo "[$CMD] Setting up scripts ... " >&2

cp $GOPATH/src/github.com/stts-se/pronlex/install/standalone/import.sh $APPDIR || exit 1
cp $GOPATH/src/github.com/stts-se/pronlex/install/standalone/start_server.sh $APPDIR || exit 1


### COMPLETED

echo "
BUILD COMPLETED! YOU CAN NOW START THE LEXICON SERVER BY INVOKING:
  $ sh $APPDIR/start_server.sh

  USAGE INFO:
  $ sh $APPDIR/start_server.sh -h

OR IMPORT STANDARD LEXICON DATA:
  $ sh $APPDIR/import.sh

  USAGE INFO:
  $ sh $APPDIR/import.sh -h

" >&2

