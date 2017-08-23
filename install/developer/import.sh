if [ $# -ne 1 ]; then
    echo "USAGE: sh $0 <APPDIR>

Imports lexicon data for Swedish, Norwegian, US English and a small test file for Arabic.
"
    exit 1
fi

APPDIR=`realpath $1`


### LEXDATA SETUP

mkdir -p $APPDIR

if [ -d "$APPDIR/lexdata" ]; then
    cd $APPDIR/lexdata && git pull && cd -
    else
	git clone https://github.com/stts-se/lexdata.git $APPDIR/lexdata
fi

mkdir -p $APPDIR/db_files
mkdir -p $APPDIR/symbol_sets

cp $APPDIR/lexdata/*/*/*.sym $APPDIR/symbol_sets/
cp $APPDIR/lexdata/mappers.txt $APPDIR/symbol_sets/
cp $APPDIR/lexdata/converters/*.cnv $APPDIR/symbol_sets/


### LEXDATA IMPORT

SVLEX=sv_se_nst_lex.db
NOBLEX=no_nob_nst_lex.db
AMELEX=en_am_cmu_lex.db
ARLEX=ar_ar_tst_lex.db

CMDDIR="$GOPATH/src/github.com/stts-se/pronlex/cmd/lexio"

echo ""
echo "IMPORT: $SVLEX"
if  go run $CMDDIR/createEmptyDB/createEmptyDB.go $APPDIR/db_files/$SVLEX ; then
    go run $CMDDIR/importLex/importLex.go $APPDIR/db_files/$SVLEX sv-se.nst $APPDIR/lexdata/sv-se/nst/swe030224NST.pron-ws.utf8.gz sv-se_ws-sampa $APPDIR/symbol_sets
else
    echo "$SVLEX FAILED"
    exit 1
fi

echo ""
echo "IMPORT: $NOBLEX"
if go run $CMDDIR/createEmptyDB/createEmptyDB.go $APPDIR/db_files/$NOBLEX ; then
    go run $CMDDIR/importLex/importLex.go $APPDIR/db_files/$NOBLEX nb-no.nst $APPDIR/lexdata/nb-no/nst/nor030224NST.pron-ws.utf8.gz nb-no_ws-sampa $APPDIR/symbol_sets
else
    echo "$NOBLEX FAILED"
    exit 1
fi

echo ""
echo "IMPORT: $AMELEX"
if go run $CMDDIR/createEmptyDB/createEmptyDB.go $APPDIR/db_files/$AMELEX ; then 
    go run $CMDDIR/importLex/importLex.go $APPDIR/db_files/$AMELEX en-us.cmu $APPDIR/lexdata/en-us/cmudict/cmudict-0.7b-ws.utf8 en-us_ws-sampa $APPDIR/symbol_sets
else
    echo "$AMELEX FAILED"
    exit 1
fi

echo ""
echo "IMPORT: $ARLEX"
if go run $CMDDIR/createEmptyDB/createEmptyDB.go $APPDIR/db_files/$ARLEX ; then
    go run $CMDDIR/importLex/importLex.go $APPDIR/db_files/$ARLEX ar-test $APPDIR/lexdata/ar/TEST/ar_TEST.pron-ws.utf8 ar_ws-sampa $APPDIR/symbol_sets
else
    echo "$ARLEX FAILED"
    exit 1
fi

rm -fr $APPDIR/lexdata


