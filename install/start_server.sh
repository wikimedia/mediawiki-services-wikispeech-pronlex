#!/bin/bash 

#############################################################
### SERVER STARTUP SCRIPT
## 1. STARTS A TEST SERVER AND RUNS A SET OF TESTS
## 2. SHUTS DOWN TEST SERVER
## 3. IF NO ERRORS ARE FOUND, THE STANDARD SERVER IS STARTED USING THE DEFAULT PORT


CMD=`basename $0`
PORT="8787"
export GOPATH=`go env GOPATH`

while getopts ":hp:a:g:" opt; do
    case $opt in
	h)
	    echo "
[$CMD] SERVER STARTUP SCRIPT
   1. STARTS A TEST SERVER AND RUNS A SET OF TESTS
   2. SHUTS DOWN TEST SERVER
   3. IF NO ERRORS ARE FOUND, THE STANDARD SERVER IS STARTED USING THE DEFAULT PORT

Options:
  -h help
  -a appdir (required)
  -p port   (default: $PORT)
  -g gopath (default: $GOPATH)

EXAMPLE INVOCATION: $CMD -a lexserver_files
" >&2
	    exit 1
	    ;;
	a)
	    APPDIR=$OPTARG
	    ;;
	p)
	    PORT=$OPTARG
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

export PATH=$PATH:$GOPATH/bin

shift $(expr $OPTIND - 1 )

if [ $# -ne 0 ]; then
    echo "[$CMD] invalid option(s): $*" >&2
    exit 1
fi

APPDIRABS=`readlink -f $APPDIR`


CMDDIR="$GOPATH/src/github.com/stts-se/pronlex/lexserver"
switches="-ss_files $APPDIRABS/symbol_sets/ -db_files $APPDIRABS/db_files/ -static $CMDDIR/static"
cd $GOPATH/src/github.com/stts-se/pronlex/lexserver && go run *.go $switches -test && go run *.go $switches $PORT
