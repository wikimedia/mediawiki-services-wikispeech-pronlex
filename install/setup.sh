#!/bin/bash 

CMD=`basename $0`
SCRIPTDIR=`dirname $0`
export GOPATH=`go env GOPATH`
export PATH=$PATH:$GOPATH/bin
gobinaries=0

function print_help {
    echo "USAGE: bash $CMD [options] <APPDIR>

OPTIONS: -b use go binaries (optional, as opposed to 'go run' with source code)

Setup files will be added to the destination folder APPDIR.
" >&2
}

if [ $# -gt 0 ] && [ $1 == "-h" ]; then
    print_help
    exit 1
fi

if [ $# -eq 2 ] && [ $1 == "-b" ]; then
    gobinaries=1
    shift
fi

if [ $# -ne 1 ]; then
    print_help
    exit 1
fi

APPDIRREL=$1

APPDIR=`readlink -f $APPDIRREL`
if [ "<>" == "<$APPDIR>" ]; then
APPDIR=`realpath $APPDIRREL`
fi
LEXDB=wikispeech_testdb.db

DEMOFILES=$GOPATH/src/github.com/stts-se/pronlex/lexserver/demo_files
CMDDIR="$GOPATH/src/github.com/stts-se/pronlex/cmd/lexio"


if [ -z "$GOPATH" ] && [ $gobinaries -eq 0 ] ; then
    echo "[$CMD] The GOPATH environment variable is required!" >&2
    exit 1
fi

function initial_setup {
    echo "[$CMD] Go binaries: $gobinaries" >&2

    if [ -f $APPDIR/db_files/$LEXDB ] ; then
	echo "[$CMD] Folder $APPDIR is already configured. No setup needed." >&2
	exit 0
    fi
    if [ -f $APPDIR/db_files/lexserver_testdb.db ] ; then
	echo "[$CMD] Folder $APPDIR is already configured. No setup needed." >&2
	exit 0
    fi


    if [ -z "$GOPATH" ] ; then
	echo "[$CMD] The GOPATH environment variable is required!" >&2
	exit 1
    fi

    ### LEXDATA PREPS

    echo "[$CMD] Setting up basic files in folder $APPDIR... " >&2

    mkdir -p $APPDIR || exit 1
    mkdir -p $APPDIR/symbol_sets || exit 1
    mkdir -p $APPDIR/db_files || exit 1

    cp $DEMOFILES/*.sym $APPDIR/symbol_sets/ || exit 1
    cp $DEMOFILES/*.cnv $APPDIR/symbol_sets/ || exit 1
    cp $DEMOFILES/mappers.txt $APPDIR/symbol_sets/ || exit 1

    cp $GOPATH/src/github.com/stts-se/pronlex/install/import.sh $APPDIR || exit 1
    cp $GOPATH/src/github.com/stts-se/pronlex/install/start_server.sh $APPDIR || exit 1
}

function run_go_command {
    cmd=$1
    args=${@:2}
    if [ $gobinaries -eq 1 ]; then
	$cmd $args
    else
	#echo "go run $CMDDIR/$cmd/$cmd.go $args"
	go run $CMDDIR/$cmd/$cmd.go $args
    fi
}

### INITIAL SETUP
initial_setup


### LEXDATA PATHS

# install/.. is the root folder
cd $SCRIPTDIR/..


#if go run cmd/lexio/createEmptyDB/createEmptyDB.go $APPDIR/db_files/$LEXDB ; then
if run_go_command createEmptyDB $APPDIR/db_files/$LEXDB ; then
    echo "[$CMD] Created empty db: $APPDIR/db_files/$LEXDB" >&2
else
    echo "[$CMD] couldn't create empty db: $APPDIR/db_files/$LEXDB" >&2
    exit 1
fi



### LEXDATA IMPORT

echo "" >&2
echo "IMPORT: $LEXDB:sv" >&2

#if go run $CMDDIR/importLex/importLex.go $APPDIR/db_files/$LEXDB sv sv_SE $DEMOFILES/swe030224NST.pron-ws.utf8.for_testing $APPDIR/symbol_sets/sv-se_ws-sampa.sym ; then
if run_go_command importLex $APPDIR/db_files/$LEXDB sv sv_SE $DEMOFILES/swe030224NST.pron-ws.utf8.for_testing $APPDIR/symbol_sets/sv-se_ws-sampa.sym ; then
    echo -n ""
else
    echo "$LEXDB:sv FAILED" >&2
    exit 1
fi

echo "" >&2
echo "IMPORT: $LEXDB:nb" >&2
#if go run $CMDDIR/importLex/importLex.go $APPDIR/db_files/$LEXDB nb nb_NO $DEMOFILES/nor030224NST.pron-ws.utf8.for_testing $APPDIR/symbol_sets/nb-no_ws-sampa.sym ; then
if run_go_command importLex $APPDIR/db_files/$LEXDB nb nb_NO $DEMOFILES/nor030224NST.pron-ws.utf8.for_testing $APPDIR/symbol_sets/nb-no_ws-sampa.sym ; then
    echo -n ""
else
    echo "$LEXDB:nb FAILED" >&2
    exit 1
fi

echo "" >&2
echo "IMPORT: $LEXDB:enu" >&2
#if run_go_command importLex $APPDIR/db_files/$LEXDB enu en_US $DEMOFILES/cmudict-0.7b-ws.utf8.for_testing $APPDIR/symbol_sets/en-us_ws-sampa.sym ; then
if run_go_command importLex $APPDIR/db_files/$LEXDB enu en_US $DEMOFILES/cmudict-0.7b-ws.utf8.for_testing $APPDIR/symbol_sets/en-us_ws-sampa.sym ; then
    echo -n ""
else
    echo "$LEXDB:enu FAILED" >&2
    exit 1
    fi

echo "" >&2
echo "IMPORT: $LEXDB:ar" >&2
#if  run_go_command importLex $APPDIR/db_files/$LEXDB ar ar_AR $DEMOFILES/ar_TEST.pron-ws.utf8.for_testing $APPDIR/symbol_sets/ar_ws-sampa.sym ; then
if  run_go_command importLex $APPDIR/db_files/$LEXDB ar ar_AR $DEMOFILES/ar_TEST.pron-ws.utf8.for_testing $APPDIR/symbol_sets/ar_ws-sampa.sym ; then
    echo -n ""
else
    echo "$LEXDB:ar FAILED" >&2
    exit 1
fi

echo "[$CMD] Clearing lexdata cache" >&2
rm -fr $APPDIR/lexdata.git


echo "
BUILD COMPLETED! YOU CAN NOW START THE LEXICON SERVER BY INVOKING:
  $ bash $SCRIPTDIR/start_server.sh -a $APPDIR

  USAGE INFO:
  $ bash $SCRIPTDIR/start_server.sh -h


OR IMPORT STANDARD LEXICON DATA FROM MASTER BRANCH:
  $ bash $SCRIPTDIR/import.sh <LEXDATA-GIT> $APPDIR master

  USAGE INFO:
  $ bash $SCRIPTDIR/import.sh

" >&2
