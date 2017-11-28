#!/bin/bash 

CMD=`basename $0`
SCRIPTDIR=`dirname $0`

if [ $# -ne 2 ]; then
    echo "USAGE: bash $CMD <LEXDATA-DIR> <DEST-DIR>" >&2
    exit 1
fi

LEXDATA=$1
DESTDIR=$2

echo "Converting:sv" >&2

if zcat $LEXDATA/sv-se/nst/swe030224NST.pron-ws.utf8.gz | egrep -wi "apa|hund|färöarna|det|här|är|ett|test" > $DESTDIR/swe030224NST.pron-ws.utf8.for_testing; then
    echo -n ""
else
    echo "sv FAILED" >&2
    exit 1
fi

echo "Converting:nb" >&2
if zcat $LEXDATA/nb-no/nst/nor030224NST.pron-ws.utf8.gz | egrep -wi "apa|hund|det|ær|test|banebrytende" > $DESTDIR/nor030224NST.pron-ws.utf8.for_testing; then
    echo -n ""
else
    echo "$LEXDB:nb FAILED" >&2
    exit 1
fi

echo "Converting:enu" >&2
if cat $LEXDATA/en-us/cmudict/cmudict-0.7b-ws.utf8 | egrep -w "test|a|children" > $DESTDIR/cmudict-0.7b-ws.utf8.for_testing; then
    echo -n ""
else
    echo "enu FAILED" >&2
    exit 1
    fi

echo "Converting:ar" >&2
if  cp $LEXDATA/ar/TEST/ar_TEST.pron-ws.utf8 $DESTDIR/ar_TEST.pron-ws.utf8.for_testing ; then
    echo -n ""
else
    echo "ar FAILED" >&2
    exit 1
fi

echo "Copying symbol sets" >&2
cp $LEXDATA/*/*/*.sym $DESTDIR

echo "Copying converters" >&2
cp $LEXDATA/converters/*.cnv $DESTDIR

#echo "Copying mappers" >&2
# cp $LEXDATA/mappers.txt $DESTDIR
