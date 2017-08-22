if [ $# -ne 1 ]; then
    echo "USAGE: sh $0 <APPDIR>"
    exit 1
fi

APPDIR=$1


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

rm -fr $APPDIR/lexdata


### LEXDATA IMPORT

SVLEX=sv_se_nst_lex.db
NOBLEX=no_nob_nst_lex.db
AMELEX=en_am_cmu_lex.db
ARLEX=ar_ar_tst_lex.db

echo ""
echo "IMPORT: $SVLEX"
if createEmptyDB $APPDIR/db_files/$SVLEX ; then
    importLex $APPDIR/db_files/$SVLEX sv-se.nst $APPDIR/lexdata/sv-se/nst/swe030224NST.pron-ws.utf8.gz sv-se_ws-sampa $APPDIR/symbol_sets
else
    echo "$SVLEX FAILED"
    exit 1
fi

echo ""
echo "IMPORT: $NOBLEX"
if createEmptyDB $APPDIR/db_files/$NOBLEX ; then
    importLex $APPDIR/db_files/$NOBLEX nb-no.nst $APPDIR/lexdata/nb-no/nst/nor030224NST.pron-ws.utf8.gz nb-no_ws-sampa $APPDIR/symbol_sets
else
    echo "$NOBLEX FAILED"
    exit 1
fi

echo ""
echo "IMPORT: $AMELEX"
if createEmptyDB $APPDIR/db_files/$AMELEX ; then 
    importLex $APPDIR/db_files/$AMELEX en-us.cmu $APPDIR/lexdata/en-us/cmudict/cmudict-0.7b-ws.utf8 en-us_ws-sampa $APPDIR/symbol_sets
else
    echo "$AMELEX FAILED"
    exit 1
fi

echo ""
echo "IMPORT: $ARLEX"
if createEmptyDB $APPDIR/db_files/$ARLEX ; then
    importLex $APPDIR/db_files/$ARLEX ar-test $APPDIR/lexdata/ar/TEST/ar_TEST.pron-ws.utf8 ar_ws-sampa $APPDIR/symbol_sets
else
    echo "$ARLEX FAILED"
    exit 1
fi

rm -fr $APPDIR/lexdata


