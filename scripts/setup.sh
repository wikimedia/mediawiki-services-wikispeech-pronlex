#!/bin/bash 

set -e

CMD=`basename $0`
SCRIPTDIR=`dirname $0`
PRONLEXPATH=`readlink -f $0 | xargs dirname | xargs dirname`
GOBINARIES=0


print_help() {
    echo "
USAGE: bash $CMD [options]

Options:
  -h print help/options
  -e db engine (required)
  -a application folder (required)
  -l db location (required for mariadb; for sqlite default is appdir)
  -b use go binaries (optional, as opposed to 'go run' with source code)

EXAMPLE INVOCATIONS:
 $CMD -e sqlite -a ~/wikispeech/sqlite
 $CMD -e mariadb -l 'speechoid:@tcp(127.0.0.1:3306)' -a ~/wikispeech/mariadb

Setup files will be added to the application folder.
" >&2
}


while getopts "hbe:l:a:" opt; do
    case $opt in
	h)
	    print_help
	    exit 1
	    ;;
	e)
	    DBENGINE=$OPTARG
	    if [ "$DBENGINE" != "sqlite" ] && [ "$DBENGINE" != "mariadb" ] ; then
		echo "[$CMD] Invalid db engine: $DBENGINE" >&2
		print_help
		exit 1
	    fi
	    ;;
	a)
	    APPDIR=$OPTARG
	    ;;
	l)
	    DBLOCATION=$OPTARG
	    ;;
	b)
	    GOBINARIES=1
	    ;;
	\?)
	    echo "Invalid option: -$OPTARG" >&2
	    exit 1
	    ;;
    esac
done

if [ -z "$DBENGINE" ] ; then
    echo "[$CMD] DBENGINE must be specified using -e" >&2
    print_help
    exit 1
fi
if [ -z "$APPDIR" ] ; then
    echo "[$CMD] APPDIR must be specified using -a" >&2
    print_help
    exit 1
fi
if [ -z "$DBLOCATION" ] ; then
    if [ $DBENGINE == "sqlite" ]; then
	DBLOCATION=$APPDIR
    else
	echo "[$CMD] DBLOCATION must be specified using -l" >&2
	print_help
	exit 1
    fi
fi

shift $(expr $OPTIND - 1 )


if [ $# -ne 0 ]; then
    echo "[$CMD] Invalid option(s): $*" >&2
    exit 1
fi


echo "[$CMD] OPTIONS:" >&2
echo "[$CMD] appdir: $APPDIR" >&2
echo "[$CMD] db engine: $DBENGINE" >&2
echo "[$CMD] db location: $DBLOCATION" >&2
echo "[$CMD] go binaries: $GOBINARIES" >&2

LEXDB=speechoid_demo
SS_FILES="$APPDIR/symbol_sets"
DEMOFILES=$PRONLEXPATH/lexserver/demo_files
CMDDIR="$PRONLEXPATH/cmd/lexio"

initial_setup() {
    
    if [ -f $SS_FILES ]; then
	echo "[$CMD] Application folder $APPDIR is already configured. No setup needed." >&2
	exit 0
    fi
    if [ -f $APPDIR/lexserver_testdb.db ] ; then
	echo "[$CMD] Application folder $APPDIR is already configured. No setup needed." >&2
	exit 0
    fi
    if [ -f $APPDIR/${LEXDB}.db ] ; then
	echo "[$CMD] Application folder $APPDIR is already configured. No setup needed." >&2
	exit 0
    fi
    # TODO: same checks for mariadb $LEXDB?


    ### LEXDATA PREPS

    echo "[$CMD] Setting up basic files in folder $APPDIR... " >&2

    mkdir -p $APPDIR || exit 1
    mkdir -p $SS_FILES || exit 1
    #mkdir -p $APPDIR/db_files || exit 1

    cp $DEMOFILES/*.sym $SS_FILES/ || exit 1
    cp $DEMOFILES/*.cnv $SS_FILES/ || exit 1
    cp $DEMOFILES/mappers.txt $SS_FILES/ || exit 1

    cp $PRONLEXPATH/scripts/import.sh $APPDIR || exit 1
    cp $PRONLEXPATH/scripts/start_server.sh $APPDIR || exit 1

    if [ "$DBENGINE" == "sqlite" ]; then
	DBLOCATION=`readlink -f $DBLOCATION`
    fi
    
}

function run_go_command {
    cmd=$1
    args=${@:2}
    echo $args
    if [ $GOBINARIES -eq 1 ]; then
	$cmd $args
    else
	#echo "go run $CMDDIR/$cmd/$cmd.go $args"
	go run $CMDDIR/$cmd/$cmd.go $args
    fi
}

### INITIAL SETUP
initial_setup


### LEXDATA PATHS

# scripts/.. is the root folder
cd $SCRIPTDIR/..


if [ $DBENGINE == "mariadb" ]; then
    mysql -u root -e "create database $LEXDB ; GRANT ALL PRIVILEGES ON $LEXDB.* TO 'speechoid'@'localhost' "
fi
if run_go_command createEmptyDB -db_engine $DBENGINE -db_location $DBLOCATION -db_name $LEXDB ; then
    echo "[$CMD] Created empty db $LEXDB @ $DBLOCATION for $DBENGINE" >&2
else
    echo "[$CMD] couldn't create empty db $LEXDB @ $DBLOCATION for $DBENGINE" >&2
    exit 1
fi



### LEXDATA IMPORT

echo "" >&2
echo "IMPORT: $LEXDB:sv" >&2

if run_go_command importLex -db_engine $DBENGINE -db_location $DBLOCATION -db_name $LEXDB -lex_name sv -locale sv_SE -lex_file $DEMOFILES/swe030224NST.pron-ws.utf8.for_testing -symbolset $SS_FILES/sv-se_ws-sampa.sym ; then
    echo -n ""
else
    echo "$LEXDB:sv FAILED" >&2
    exit 1
fi

echo "" >&2
if run_go_command importLex -db_engine $DBENGINE -db_location $DBLOCATION -db_name $LEXDB -lex_name nb -locale nb_NO -lex_file $DEMOFILES/nor030224NST.pron-ws.utf8.for_testing -symbolset $SS_FILES/nb-no_ws-sampa.sym ; then
    echo -n ""
else
    echo "$LEXDB:nb FAILED" >&2
    exit 1
fi

echo "" >&2
echo "IMPORT: $LEXDB:enu" >&2
if run_go_command importLex -db_engine $DBENGINE -db_location $DBLOCATION -db_name $LEXDB -lex_name enu -locale en_US -lex_file $DEMOFILES/cmudict-0.7b-ws.utf8.for_testing -symbolset $SS_FILES/en-us_ws-sampa.sym ; then
    echo -n ""
else
    echo "$LEXDB:enu FAILED" >&2
    exit 1
    fi

echo "" >&2
echo "IMPORT: $LEXDB:ar" >&2
if  run_go_command importLex -db_engine $DBENGINE -db_location $DBLOCATION -db_name $LEXDB -lex_name ar -locale ar_AR -lex_file $DEMOFILES/ar_TEST.pron-ws.utf8.for_testing -symbolset $SS_FILES/ar_ws-sampa.sym ; then
    echo -n ""
else
    echo "$LEXDB:ar FAILED" >&2
    exit 1
fi

echo "[$CMD] Clearing lexdata cache" >&2
rm -fr $APPDIR/lexdata.git


echo "
BUILD COMPLETED! YOU CAN NOW START THE LEXICON SERVER BY INVOKING:
  $ bash $SCRIPTDIR/start_server.sh -e $DBENGINE -l $DBLOCATION

  USAGE INFO:
  $ bash $SCRIPTDIR/start_server.sh -h


OR IMPORT STANDARD LEXICON DATA FROM MASTER BRANCH:
  $ bash $SCRIPTDIR/import.sh -e $DBENGINE -l $DBLOCATION -L <LEXDATA-GIT> master

  USAGE INFO:
  $ bash $SCRIPTDIR/import.sh -h

" >&2
