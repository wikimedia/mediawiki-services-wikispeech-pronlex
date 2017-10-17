#!/bin/bash 

CMD=`basename $0`
export GOPATH=`go env GOPATH`
export PATH=$PATH:$GOPATH/bin


if [ $# -ne 1 ]; then
    echo "USAGE: sh $CMD <APPDIR>

 Pronlex will be installed in the destination folder APPDIR.
">&2
    exit 1
fi

APPDIR=$1
LEXDB=wikispeech_testdb.db

if [ -z "$GOPATH" ] ; then
    echo "[$CMD] The GOPATH environment variable is required!" >&2
    exit 1
fi


## LEXSERVER INSTALL

echo "[$CMD] Installing pronlex/lexserver ... " >&2

go get -u github.com/stts-se/pronlex/lexserver || exit 1
go install github.com/stts-se/pronlex/lexserver || exit 1

go get -u github.com/stts-se/pronlex/cmd/lexio/importLex || exit 1
go install github.com/stts-se/pronlex/cmd/lexio/importLex || exit 1

go get -u github.com/stts-se/pronlex/cmd/lexio/createEmptyDB || exit 1
go install github.com/stts-se/pronlex/cmd/lexio/createEmptyDB || exit 1

function initial_setup {
    if [ -f $APPDIR/db_files/$LEXDB ] ; then
	echo "[$CMD] Docker folder $APPDIR is already configured. No setup needed." >&2
	exit 0
    fi
    if [ -f $APPDIR/db_files/lexserver_testdb.db ] ; then
	echo "[$CMD] Docker folder $APPDIR is already configured. No setup needed." >&2
	exit 0
    fi


    if [ -z "$GOPATH" ] ; then
	echo "[$CMD] The GOPATH environment variable is required!" >&2
	exit 1
    fi

    ### LEXDATA PREPS

    echo "[$CMD] Setting up basic files in docker folder $APPDIR... " >&2

    mkdir -p $APPDIR || exit 1
    mkdir -p $APPDIR/symbol_sets || exit 1
    mkdir -p $APPDIR/db_files || exit 1

    mkdir -p $APPDIR/.static || exit 1
    cp -r $GOPATH/src/github.com/stts-se/pronlex/lexserver/static/* $APPDIR/.static || exit 1

    cp $GOPATH/src/github.com/stts-se/pronlex/lexserver/demo_files/*.sym $APPDIR/symbol_sets/ || exit 1
    cp $GOPATH/src/github.com/stts-se/pronlex/lexserver/demo_files/*.cnv $APPDIR/symbol_sets/ || exit 1
    echo "" >> $APPDIR/symbol_sets/mappers.txt || exit 1
    cat $GOPATH/src/github.com/stts-se/pronlex/lexserver/demo_files/mappers.txt >> $APPDIR/symbol_sets/mappers.txt || exit 1

    cp $GOPATH/src/github.com/stts-se/pronlex/install/standalone/import.sh $APPDIR || exit 1
    cp $GOPATH/src/github.com/stts-se/pronlex/install/standalone/start_server.sh $APPDIR || exit 1
}



### INITIAL SETUP
initial_setup


### LEXDATA PATHS

if createEmptyDB $APPDIR/db_files/$LEXDB ; then
    echo "[$CMD] Created empty db in docker: $APPDIR/db_files/$LEXDB" >&2
else
    echo "[$CMD] couldn't create empty db in docker: $APPDIR/db_files/$LEXDB" >&2
    exit 1
fi


nConverters=`ls $APPDIR/lexdata/converters/*.cnv 2> /dev/null |wc -l`
nSymSets=`ls $APPDIR/lexdata/*/*/*.sym 2> /dev/null |wc -l`
if [ -d "$APPDIR/lexdata" ] && [ $nSymSets != 0 ] || [ $nConverters != 0 ]; then
    cd $APPDIR/lexdata && git pull && cd -
    KEEP=1
else
    if git clone https://github.com/stts-se/lexdata.git $APPDIR/lexdata.git; then
	echo -n "" # OK
    else
	echo "[$CMD] git clone failed" >&2
	exit 1
    fi
fi

mkdir -p $APPDIR/db_files || exit 1
mkdir -p $APPDIR/symbol_sets || exit 1

cp $APPDIR/lexdata.git/*/*/*.sym $APPDIR/symbol_sets/ || exit 1
echo "" >> $APPDIR/symbol_sets/mappers.txt || exit 1
cat $APPDIR/lexdata.git/mappers.txt >> $APPDIR/symbol_sets/mappers.txt || exit 1
cp $APPDIR/lexdata.git/converters/*.cnv $APPDIR/symbol_sets/ || exit 1


### LEXDATA IMPORT

echo "" >&2
echo "IMPORT: $LEXDB:sv" >&2

if zcat $APPDIR/lexdata.git/sv-se/nst/swe030224NST.pron-ws.utf8.gz | egrep -wi "apa|hund|färöarna|det|här|är|ett|test" > $APPDIR/lexdata.git/sv-se/nst/swe030224NST.pron-ws.utf8.for_testing && importLex $APPDIR/db_files/$LEXDB sv $APPDIR/lexdata.git/sv-se/nst/swe030224NST.pron-ws.utf8.for_testing sv-se_ws-sampa $APPDIR/symbol_sets ; then
    echo -n ""
else
    echo "$LEXDB:sv FAILED" >&2
    exit 1
fi

echo "" >&2
echo "IMPORT: $LEXDB:nb" >&2
if zcat $APPDIR/lexdata.git/nb-no/nst/nor030224NST.pron-ws.utf8.gz | egrep -wi "apa|hund|det|ær|test|banebrytende" > $APPDIR/lexdata.git/nb-no/nst/nor030224NST.pron-ws.utf8.for_testing && importLex $APPDIR/db_files/$LEXDB nb $APPDIR/lexdata.git/nb-no/nst/nor030224NST.pron-ws.utf8.for_testing nb-no_ws-sampa $APPDIR/symbol_sets ; then
    echo -n ""
else
    echo "$LEXDB:nb FAILED" >&2
    exit 1
fi

echo "" >&2
echo "IMPORT: $LEXDB:enu" >&2
if cat $APPDIR/lexdata.git/en-us/cmudict/cmudict-0.7b-ws.utf8 | egrep -w "test|a|children" > $APPDIR/lexdata.git/en-us/cmudict/cmudict-0.7b-ws.utf8.for_testing && importLex $APPDIR/db_files/$LEXDB enu $APPDIR/lexdata.git/en-us/cmudict/cmudict-0.7b-ws.utf8.for_testing en-us_ws-sampa $APPDIR/symbol_sets ; then
    echo -n ""
else
    echo "$LEXDB:enu FAILED" >&2
    exit 1
    fi

echo "" >&2
echo "IMPORT: $LEXDB:ar" >&2
if importLex $APPDIR/db_files/$LEXDB ar $APPDIR/lexdata.git/ar/TEST/ar_TEST.pron-ws.utf8 ar_ws-sampa $APPDIR/symbol_sets ; then
    echo -n ""
else
    echo "$LEXDB:ar FAILED" >&2
    exit 1
fi

echo "[$CMD] Clearing lexdata cache" >&2
rm -fr $APPDIR/lexdata.git



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

