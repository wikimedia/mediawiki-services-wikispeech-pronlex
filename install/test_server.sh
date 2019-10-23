#!/bin/bash 

#############################################################
### SERVER STARTUP SCRIPT
## 1. STARTS A TEST SERVER AND RUNS A SET OF TESTS
## 2. SHUTS DOWN TEST SERVER
## 3. IF NO ERRORS ARE FOUND, THE STANDARD SERVER IS STARTED USING THE DEFAULT PORT


CMD=`basename $0`
PORT="8787"
export GOPATH=`go env GOPATH`
export PATH=$PATH:$GOPATH/bin
PRONLEXPATH=`readlink -f $0 | xargs dirname | xargs dirname` # $GOPATH/src/github.com/stts-se/pronlex

while getopts ":ha:" opt; do
    case $opt in
	h)
	    echo "
[$CMD] SCRIPT TO RUN LEXSERVER INIT TESTS (WITHOUT STARTING THE PROPER SERVER)

Options:
  -h help
  -a appdir (required)

EXAMPLE INVOCATION: $CMD -a lexserver_files
" >&2
	    exit 1
	    ;;
	a)
	    APPDIR=$OPTARG
	    ;;
	\?)
	    echo "Invalid option: -$OPTARG" >&2
	    exit 1
	    ;;
    esac
done

if [ -z "$APPDIR" ] ; then
    echo "[$CMD] APPDIR must be specified using -a!" >&2
    exit 1
fi

if [ -z "$GOPATH" ] ; then
    echo "[$CMD] The GOPATH environment variable is required!" >&2
    exit 1
fi

shift $(expr $OPTIND - 1 )

if [ $# -ne 0 ]; then
    echo "[$CMD] invalid option(s): $*" >&2
    exit 1
fi

APPDIRABS=`readlink -f $APPDIR`


CMDDIR="$PRONLEXPATH/lexserver"
switches="-ss_files $APPDIRABS/symbol_sets/ -db_files $APPDIRABS/db_files/ -static $CMDDIR/static"
cd $PRONLEXPATH/lexserver && go run *.go $switches -test #  && go run *.go $switches $PORT
