if [ $# -ne 1 ]; then
    echo "USAGE: sh $0 <APPDIR>"
    exit 1
fi

APPDIR=$1


## LEXSERVER INSTALL

echo "[$0] Installing pronlex/lexserver ... "

go get github.com/stts-se/pronlex/lexserver
go install github.com/stts-se/pronlex/lexserver

go get github.com/stts-se/pronlex/cmd/lexio/importLex
go install github.com/stts-se/pronlex/cmd/lexio/importLex

go get github.com/stts-se/pronlex/cmd/lexio/createEmptyDB
go install github.com/stts-se/pronlex/cmd/lexio/createEmptyDB

### LEXDATA PREPS

echo "[$0] Setting up basic files ... "

mkdir -p $APPDIR

mkdir -p $APPDIR/static
cp -r $GOPATH/src/github.com/stts-se/pronlex/lexserver/static/* $APPDIR/static/

mkdir -p $APPDIR/symbol_sets
cp $GOPATH/src/github.com/stts-se/pronlex/lexserver/demo_files/* $APPDIR/symbol_sets

echo "[$0] Setting up scripts ... "

cp $GOPATH/src/github.com/stts-se/pronlex/install/standalone/import.sh $APPDIR
cp $GOPATH/src/github.com/stts-se/pronlex/install/standalone/start_server.sh $APPDIR


### COMPLETED

echo "
BUILD COMPLETED. YOU CAN NOW START THE LEXICON SERVER BY INVOKING
  $ sh $APPDIR/start_server.sh $APPDIR
"
