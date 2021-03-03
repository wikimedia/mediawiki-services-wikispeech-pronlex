#!/bin/bash 

#############################################################
### SERVER STARTUP SCRIPT

set -e

CMD=`basename $0`
GOCMD="lexserver"
PORT="8787"
PREFIX=""
PRONLEXPATH=`readlink -f $0 | xargs dirname | xargs dirname`
GOBINARIES=0
SERVERHELP=0
LOGGER="stderr"

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
  -r explicit server prefix (default: empty)
  -a application folder (required)
  -l db location (required for mariadb; for sqlite default is application folder)
  -p lexserver port (default: $PORT)
  -b use go binaries (optional, as opposed to 'go run' with source code)
  -o system logger (stderr, syslog or filename)
  -t test mode (default: $TESTMODE)
     $TESTOFF: no tests
     $TESTON: test before starting lexserver, exit on error
     $TESTONLY: exit after tests

EXAMPLE INVOCATIONS:
 bash $0 -a ~/wikispeech/sqlite -e sqlite
 bash $0 -a ~/wikispeech/mariadb -e mariadb -l 'speechoid:@tcp(127.0.0.1:3306)'
" >&2
}

while getopts "hHbt:p:l:e:a:o:r:" opt; do
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
	r)
	    PREFIX=$OPTARG
	    ;;
	o)
	    LOGGER=$OPTARG
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

STATIC=`realpath $APPDIR/static`

echo "[$CMD] OPTIONS:" >&2
echo "[$CMD] application folder: $APPDIR" >&2
echo "[$CMD] db engine: $DBENGINE" >&2
echo "[$CMD] db location: $DBLOCATION" >&2
echo "[$CMD] lexserver port: $PORT" >&2
if [ "<$PREFIX>" != "<>" ]; then
   echo "[$CMD] lexserver prefix: $PREFIX" >&2
fi
echo "[$CMD] static: $STATIC" >&2
echo "[$CMD] logger: $LOGGER" >&2
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

switches="-logger $LOGGER -db_engine $DBENGINE -db_location $DBLOCATION -static $STATIC"
if [ "<$PREFIX>" != "<>" ]; then
    switches="$switches -prefix $PREFIX"
fi
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

