### INSTALL LEXSERVER

go get github.com/stts-se/pronlex/lexserver
go install github.com/stts-se/pronlex/lexserver

go get github.com/stts-se/pronlex/cmd/lexio/importLex
go install github.com/stts-se/pronlex/cmd/lexio/importLex

go get github.com/stts-se/pronlex/cmd/lexio/createEmptyDB
go install github.com/stts-se/pronlex/cmd/lexio/createEmptyDB


### LEXDATA

APP_DIR=lexserver_files
SVLEX=sv_se_nst_lex.db
NOBLEX=no_nob_nst_lex.db
AMELEX=en_am_cmu_lex.db
ARLEX=ar_ar_tst_lex.db

mkdir -p $APP_DIR

if [ -d "$APP_DIR/lexdata" ]; then
    cd $APP_DIR/lexdata && git pull && cd -
    else
	git clone https://github.com/stts-se/lexdata.git $APP_DIR/lexdata
fi

mkdir -p $APP_DIR/db_files
mkdir -p $APP_DIR/symbol_sets

cp $APP_DIR/lexdata/*/*/*.sym $APP_DIR/symbol_sets/
cp $APP_DIR/lexdata/mappers.txt $APP_DIR/symbol_sets/
cp $APP_DIR/lexdata/converters/*.cnv $APP_DIR/symbol_sets/


echo ""
echo "IMPORT: $SVLEX"
if createEmptyDB $APP_DIR/db_files/$SVLEX ; then
    importLex $APP_DIR/db_files/$SVLEX sv-se.nst $APP_DIR/lexdata/sv-se/nst/swe030224NST.pron-ws.utf8.gz sv-se_ws-sampa $APP_DIR/symbol_sets
else
    echo "FAILED"
    exit 1
fi

echo ""
echo "IMPORT: $NOBLEX"
if createEmptyDB $APP_DIR/db_files/$NOBLEX ; then
    importLex $APP_DIR/db_files/$NOBLEX nb-no.nst $APP_DIR/lexdata/nb-no/nst/nor030224NST.pron-ws.utf8.gz nb-no_ws-sampa $APP_DIR/symbol_sets
else
    echo "FAILED"
    exit 1
fi

echo ""
echo "IMPORT: $AMELEX"
if createEmptyDB $APP_DIR/db_files/$AMELEX ; then 
    importLex $APP_DIR/db_files/$AMELEX en-us.cmu $APP_DIR/lexdata/en-us/cmudict/cmudict-0.7b-ws.utf8 en-us_ws-sampa $APP_DIR/symbol_sets
else
    echo "FAILED"
    exit 1
fi

echo ""
echo "IMPORT: $ARLEX"
if createEmptyDB $APP_DIR/db_files/$ARLEX ; then
    importLex $APP_DIR/db_files/$ARLEX ar-test $APP_DIR/lexdata/ar/TEST/ar_TEST.pron-ws.utf8 ar_ws-sampa $APP_DIR/symbol_sets
else
    echo "FAILED"
    exit 1
fi

rm -fr $APP_DIR/lexdata



