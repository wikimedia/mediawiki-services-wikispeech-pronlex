#!/bin/bash 

#############################################################
### SERVER STARTUP SCRIPT

set -e

CMD=`basename $0`
GOCMD="lexserver"
PORT="8787"
PRONLEXPATH=`readlink -f $0 | xargs dirname | xargs dirname`
GOBINARIES=0
SERVERHELP=0

# Test mode:
TESTON="on"
TESTOFF="off"
TESTONLY="test"
TESTMODE=$TESTON

print_help(){
	    echo "
USAGE bash $CMD [options]

Options:
  -h print help/options and exit
  -H call $GOCMD help and exit
  -e db engine (required)
  -a application folder (required)
  -l db location (required for mariadb; for sqlite default is application folder)
  -p lexserver port (default: $PORT)
  -b use go binaries (optional, as opposed to 'go run' with source code)
  -t test mode (default: $TESTMODE)
     $TESTOFF: no tests
     $TESTON: test before starting lexserver, exit on error
     $TESTONLY: exit after tests

EXAMPLE INVOCATIONS:
 bash $0 -a ~/wikispeech/sqlite -e sqlite
 bash $0 -a ~/wikispeech/mariadb -e mariadb -l 'speechoid:@tcp(127.0.0.1:3306)'
" >&2
}

while getopts "hHbt:p:l:e:a:" opt; do
    case $opt in
	h)
	    print_help
	    exit 1
	    ;;
	H)
	    SERVERHELP=1
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
	p)
	    PORT=$OPTARG
	    ;;
	b)
	    GOBINARIES=1
	    ;;
	t)
	    TESTMODE=$OPTARG
	    if [ "$TESTMODE" != "$TESTON" ] && [ "$TESTMODE" != "$TESTOFF" ] && [ "$TESTMODE" != "$TESTONLY" ] ; then
		echo "[$CMD] Invalid test mode: $TESTMODE" >&2
		print_help
		exit 1
	    fi
	    ;;
	\?)
	    echo "Invalid option: -$OPTARG" >&2
	    exit 1
	    ;;
    esac
done

if [ -z "$DBENGINE" ] ; then
    echo "[$CMD] db engine must be specified using -e" >&2
    print_help
    exit 1
fi
if [ -z "$APPDIR" ] ; then
    echo "[$CMD] application folder must be specified using -a" >&2
    print_help
    exit 1
fi
if [ -z "$DBLOCATION" ] ; then
    if [ $DBENGINE == "sqlite" ]; then
	DBLOCATION=$APPDIR
    else
	echo "[$CMD] db location must be specified using -l" >&2
	print_help
	exit 1
    fi
fi


shift $(expr $OPTIND - 1 )

if [ $# -ne 0 ]; then
    echo "[$CMD] Invalid option(s): $*" >&2
    exit 1
fi

if [ "$DBENGINE" == "sqlite" ]; then
    DBLOCATION=`readlink -f $DBLOCATION`
fi


echo "[$CMD] OPTIONS:" >&2
echo "[$CMD] applciation folder: $APPDIR" >&2
echo "[$CMD] db engine: $DBENGINE" >&2
echo "[$CMD] db location: $DBLOCATION" >&2
echo "[$CMD] lexserver port: $PORT" >&2
echo "[$CMD] go binaries: $GOBINARIES" >&2
echo "[$CMD] test mode: $TESTMODE" >&2

function run_go_cmd {
    args=${@:1}
    if [ $GOBINARIES -eq 1 ]; then
	$GOCMD $args
    else
	cd $PRONLEXPATH/lexserver
	go run *.go $args
    fi
}

switches="-db_engine $DBENGINE -db_location $DBLOCATION -static $APPDIR/static "
if [ $SERVERHELP -eq 1 ]; then
    switches="-help"
    echo "[$CMD] Calling lexserver help and exit" >&2
    echo "" >&2
fi

if [ $TESTMODE == $TESTON ] || [ $TESTMODE == $TESTONLY ] ; then
    run_go_cmd $switches -test
fi
if [ $TESTMODE != $TESTONLY ]; then
    run_go_cmd $switches $PORT
fi

