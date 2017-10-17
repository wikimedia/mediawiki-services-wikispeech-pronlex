#!/bin/bash 

CMD=`basename $0`
APPDIR=`dirname $0`
KEEP=0
export GOPATH=`go env GOPATH`
export PATH=$PATH:$GOPATH/bin

#echo "[$CMD] temporarily disabled, please use docker for now" >&2
#exit 1

while getopts ":hka:" opt; do
  case $opt in
    h)
	echo "
USAGE: sh $CMD <OPTIONS>

Imports lexicon data for Swedish, Norwegian, US English and a small test file for Arabic.

Options:
  -h help
  -a appdir (default: this script's folder)
  -k keep lexdata files after import
" >&2
	exit 1
      ;;
    a)
        APPDIR=$OPTARG
      ;;
    k)
        KEEP=1
      ;;
    \?)
      echo "Invalid option: -$OPTARG" >&2
      ;;
  esac
done

shift $(expr $OPTIND - 1 )

if [ $# -ne 0 ]; then
    echo "[$CMD] invalid option(s): $*" >&2
    exit 1
fi

### LEXDATA SETUP

if [ ! -d $APPDIR/symbol_sets ] ; then
    echo "FAILED: $APPDIR is not configured for the lexserver!" >&2
    exit 1
fi

if [ -d "$APPDIR/lexdata.git" ]; then
    cd $APPDIR/lexdata.git && git pull && cd -
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

echo "" >&2
echo "IMPORT: $SVLEX" >&2
if createEmptyDB $APPDIR/db_files/$SVLEX ; then
    importLex $APPDIR/db_files/$SVLEX sv-se.nst $APPDIR/lexdata.git/sv-se/nst/swe030224NST.pron-ws.utf8.gz sv-se_ws-sampa $APPDIR/symbol_sets
else
    echo "$SVLEX FAILED" >&2
    exit 1
fi

echo "" >&2
echo "IMPORT: $NOBLEX" >&2
if createEmptyDB $APPDIR/db_files/$NOBLEX ; then
    importLex $APPDIR/db_files/$NOBLEX nb-no.nst $APPDIR/lexdata.git/nb-no/nst/nor030224NST.pron-ws.utf8.gz nb-no_ws-sampa $APPDIR/symbol_sets
else
    echo "$NOBLEX FAILED" >&2
    exit 1
fi

echo "" >&2
echo "IMPORT: $AMELEX" >&2
if createEmptyDB $APPDIR/db_files/$AMELEX ; then 
    importLex $APPDIR/db_files/$AMELEX en-us.cmu $APPDIR/lexdata.git/en-us/cmudict/cmudict-0.7b-ws.utf8 en-us_ws-sampa $APPDIR/symbol_sets
else
    echo "$AMELEX FAILED" >&2
    exit 1
fi

echo "" >&2
echo "IMPORT: $ARLEX" >&2
if createEmptyDB $APPDIR/db_files/$ARLEX ; then
    importLex $APPDIR/db_files/$ARLEX ar-test $APPDIR/lexdata.git/ar/TEST/ar_TEST.pron-ws.utf8 ar_ws-sampa $APPDIR/symbol_sets
else
    echo "$ARLEX FAILED" >&2
    exit 1
fi

if [ $KEEP -eq 0 ]; then
    echo "[$CMD] Clearing lexdata cache" >&2
    rm -fr $APPDIR/lexdata.git
else
    echo "[$CMD] Keeping lexdata cache" >&2  
fi

echo "[$CMD] Done. BYE!" >&2


