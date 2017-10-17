#!/bin/bash 

CMD=`basename $0`
export GOPATH=`go env GOPATH`
export PATH=$PATH:$GOPATH/bin

echo "[$CMD] temporarily disabled, please use docker for now" >&2
exit 1

if [ $# -ne 2 ]; then
    echo "USAGE: sh $CMD <LEXDATA-GIT> <APPDIR>

Imports lexicon data for Swedish, Norwegian, US English and a small test file for Arabic from the lexdata repository.
If the <LEXDATA-GIT> folder doesn't exist, it will be downloaded from github.
" >&2
    exit 1
fi

LEXDATA=$1
APPDIR=`realpath $2`

### LEXDATA SETUP


if [ ! -d $APPDIR/symbol_sets ] ; then
	echo "[$CMD] $APPDIR is not configured for the lexserver. Run setup.sh first!" >&2
    exit 1
fi

if [ -d $LEXDATA ]; then
    cd $LEXDATA && git pull && cd - || exit 1
else
    git clone https://github.com/stts-se/lexdata.git $LEXDATA || exit 1
fi

mkdir -p $APPDIR/db_files || exit 1
mkdir -p $APPDIR/symbol_sets || exit 1

cp $LEXDATA/*/*/*.sym $APPDIR/symbol_sets/ || exit 1
echo "" >> $APPDIR/symbol_sets/mappers.txt || exit 1
cat $LEXDATA/mappers.txt >> $APPDIR/symbol_sets/mappers.txt || exit 1
cp $LEXDATA/converters/*.cnv $APPDIR/symbol_sets/ || exit 1


### LEXDATA IMPORT

SVLEX=sv_se_nst_lex.db
NOBLEX=no_nob_nst_lex.db
AMELEX=en_am_cmu_lex.db
ARLEX=ar_ar_tst_lex.db

if [ -e $APPDIR/db_files/$SVLEX ]; then
    echo "[$CMD] cannot create db if it already exists: $APPDIR/db_files/$SVLEX" >&2
    exit 1
fi
if [ -e $APPDIR/db_files/$NOBLEX ]; then
    echo "[$CMD] cannot create db if it already exists: $APPDIR/db_files/$NOBLEX" >&2
    exit 1
fi
if [ -e $APPDIR/db_files/$AMELEX ]; then
    echo "[$CMD] cannot create db if it already exists: $APPDIR/db_files/$AMELEX" >&2
    exit 1
fi
if [ -e $APPDIR/db_files/$ARLEX ]; then
    echo "[$CMD] cannot create db if it already exists: $APPDIR/db_files/$ARLEX" >&2
    exit 1
fi


CMDDIR="$GOPATH/src/github.com/stts-se/pronlex/cmd/lexio"

echo "" >&2
echo "IMPORT: $SVLEX" >&2
if  go run $CMDDIR/createEmptyDB/createEmptyDB.go $APPDIR/db_files/$SVLEX ; then
    go run $CMDDIR/importLex/importLex.go $APPDIR/db_files/$SVLEX sv-se.nst $LEXDATA/sv-se/nst/swe030224NST.pron-ws.utf8.gz sv-se_ws-sampa $APPDIR/symbol_sets
else
    echo "$SVLEX FAILED" >&2
    exit 1
fi

echo "" >&2
echo "IMPORT: $NOBLEX" >&2
if go run $CMDDIR/createEmptyDB/createEmptyDB.go $APPDIR/db_files/$NOBLEX ; then
    go run $CMDDIR/importLex/importLex.go $APPDIR/db_files/$NOBLEX nb-no.nst $LEXDATA/nb-no/nst/nor030224NST.pron-ws.utf8.gz nb-no_ws-sampa $APPDIR/symbol_sets
else
    echo "$NOBLEX FAILED" >&2
    exit 1
fi

echo "" >&2
echo "IMPORT: $AMELEX" >&2
if go run $CMDDIR/createEmptyDB/createEmptyDB.go $APPDIR/db_files/$AMELEX ; then 
    go run $CMDDIR/importLex/importLex.go $APPDIR/db_files/$AMELEX en-us.cmu $LEXDATA/en-us/cmudict/cmudict-0.7b-ws.utf8 en-us_ws-sampa $APPDIR/symbol_sets
else
    echo "$AMELEX FAILED" >&2
    exit 1
fi

echo "" >&2
echo "IMPORT: $ARLEX" >&2
if go run $CMDDIR/createEmptyDB/createEmptyDB.go $APPDIR/db_files/$ARLEX ; then
    go run $CMDDIR/importLex/importLex.go $APPDIR/db_files/$ARLEX ar-test $LEXDATA/ar/TEST/ar_TEST.pron-ws.utf8 ar_ws-sampa $APPDIR/symbol_sets
else
    echo "$ARLEX FAILED" >&2
    exit 1
fi
