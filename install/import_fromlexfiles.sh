#!/bin/bash 

CMD=`basename $0`
export GOPATH=`go env GOPATH`
export PATH=$PATH:$GOPATH/bin
PRONLEXPATH=`readlink -f $0 | xargs dirname | xargs dirname` # $GOPATH/src/github.com/stts-se/pronlex

if [ $# -ne 3 ]; then
    echo "USAGE: bash $CMD <LEXDATA-GIT> <APPDIR> <LEXDATA-RELEASE-TAG>

Imports lexicon data for Swedish, Norwegian, US English and a small test file for Arabic from the lexdata repository.
Imports from plain text lexicon files.
If the <LEXDATA-GIT> folder doesn't exist, it will be downloaded from github: https://github.com/stts-se/lexdata

If you don't know what release tag you should use, you should probably use master.
" >&2
    exit 1
fi

LEXDATA=$1
APPDIR=`readlink -f $2`
RELEASETAG=$3

### LEXDATA SETUP


if [ ! -d $APPDIR/symbol_sets ] ; then
	echo "[$CMD] $APPDIR is not configured for the lexserver. Run setup.sh first!" >&2
    exit 1
fi

if [ -d $LEXDATA ]; then
    cd $LEXDATA && git pull && git checkout $RELEASETAG && cd - || exit 1
else
    git clone https://github.com/stts-se/lexdata.git --branch $RELEASETAG --single-branch $LEXDATA || exit 1
fi

mkdir -p $APPDIR/db_files || exit 1
mkdir -p $APPDIR/symbol_sets || exit 1

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


### COPY REQUIRED FILES
cp $LEXDATA/*/*/*.sym $APPDIR/symbol_sets/ || exit 1
echo "" >> $APPDIR/symbol_sets/mappers.txt || exit 1
cat $LEXDATA/mappers.txt >> $APPDIR/symbol_sets/mappers.txt || exit 1
cp $LEXDATA/converters/*.cnv $APPDIR/symbol_sets/ || exit 1



CMDDIR="$PRONLEXPATH/cmd/lexio"

echo "" >&2
echo "IMPORT: $SVLEX" >&2
if  go run $CMDDIR/createEmptyDB/createEmptyDB.go $APPDIR/db_files/$SVLEX &&
	go run $CMDDIR/importLex/importLex.go $APPDIR/db_files/$SVLEX sv-se.nst sv_SE $LEXDATA/sv-se/nst/swe030224NST.pron-ws.utf8.gz $APPDIR/symbol_sets/sv-se_ws-sampa.sym ; then
    echo -n ""
else
    echo "$SVLEX FAILED" >&2
    exit 1
fi

echo "" >&2
echo "IMPORT: $NOBLEX" >&2
if go run $CMDDIR/createEmptyDB/createEmptyDB.go $APPDIR/db_files/$NOBLEX &&
	go run $CMDDIR/importLex/importLex.go $APPDIR/db_files/$NOBLEX nb-no.nst nb_NO $LEXDATA/nb-no/nst/nor030224NST.pron-ws.utf8.gz $APPDIR/symbol_sets/nb-no_ws-sampa.sym ; then
    echo -n ""
else
    echo "$NOBLEX FAILED" >&2
    exit 1
fi

echo "" >&2
echo "IMPORT: $AMELEX" >&2
if go run $CMDDIR/createEmptyDB/createEmptyDB.go $APPDIR/db_files/$AMELEX &&
    go run $CMDDIR/importLex/importLex.go $APPDIR/db_files/$AMELEX en-us.cmu en_US $LEXDATA/en-us/cmudict/cmudict-0.7b-ws.utf8 $APPDIR/symbol_sets/en-us_ws-sampa.sym ; then
    echo -n ""
else
    echo "$AMELEX FAILED" >&2
    exit 1
fi

echo "" >&2
echo "IMPORT: $ARLEX" >&2
if go run $CMDDIR/createEmptyDB/createEmptyDB.go $APPDIR/db_files/$ARLEX &&
	go run $CMDDIR/importLex/importLex.go $APPDIR/db_files/$ARLEX ar-test ar_AR $LEXDATA/ar/TEST/ar_TEST.pron-ws.utf8 $APPDIR/symbol_sets/ar_ws-sampa.sym ; then
    echo -n ""
else
    echo "$ARLEX FAILED" >&2
    exit 1
fi
