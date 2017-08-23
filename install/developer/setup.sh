## WORK IN PROGRESS

if [ $# -ne 1 ]; then
    echo "USAGE: sh $0 <APPDIR>"
    exit 1
fi

APPDIR=$1

if [ -z "$GOPATH" ] ; then
    echo "[$0] The GOPATH environment variable is required!"
    exit 1
fi

### LEXDATA PREPS

echo "[$0] Setting up basic files ... "

mkdir -p $APPDIR

mkdir -p $APPDIR/symbol_sets
cp $GOPATH/src/github.com/stts-se/pronlex/lexserver/demo_files/* $APPDIR/symbol_sets


### COMPLETED

echo "
BUILD COMPLETED. YOU CAN NOW START THE LEXICON SERVER BY INVOKING
  $ sh start_server.sh $APPDIR
"
