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
gobinaries=0

print_help(){
	    echo "
[$CMD] SERVER STARTUP SCRIPT
   1. STARTS A TEST SERVER AND RUNS A SET OF TESTS
   2. SHUTS DOWN TEST SERVER
   3. IF NO ERRORS ARE FOUND, THE STANDARD SERVER IS STARTED USING THE DEFAULT PORT

Options:
  -h help
  -a appdir (required)
  -p port   (default: $PORT)
  -b use go binaries (optional, as opposed to 'go run' with source code)

EXAMPLE INVOCATION: $CMD -a lexserver_files
" >&2
}

while getopts ":hp:a:b" opt; do
    case $opt in
	h)
	    print_help
	    exit 1
	    ;;
	a)
	    APPDIR=$OPTARG
	    ;;
	p)
	    PORT=$OPTARG
	    ;;
	b)
	    gobinaries=1
	    ;;
	\?)
	    echo "Invalid option: -$OPTARG" >&2
	    exit 1
	    ;;
    esac
done

if [ -z "$APPDIR" ] ; then
    echo "[$CMD] APPDIR must be specified using -a!" >&2
    print_help
    exit 1
fi

if [ -z "$GOPATH" ] && [ $gobinaries -eq 0 ] ; then
    echo "[$CMD] The GOPATH environment variable is required!" >&2
    exit 1
fi

shift $(expr $OPTIND - 1 )

if [ $# -ne 0 ]; then
    echo "[$CMD] invalid option(s): $*" >&2
    exit 1
fi

APPDIRABS=`readlink -f $APPDIR`


CMDDIR="$GOPATH/src/github.com/stts-se/pronlex/lexserver"


switches="-ss_files $APPDIRABS/symbol_sets/ -db_files $APPDIRABS/db_files/ -static $CMDDIR/static"
#cd $GOPATH/src/github.com/stts-se/pronlex/lexserver && go run *.go $switches -test && go run *.go $switches $PORT

echo "[$CMD] Go binaries: $gobinaries" >&2
if [ $gobinaries -eq 1 ]; then
    lexserver $switches -test && lexserver $switches $PORT
else
    cd $GOPATH/src/github.com/stts-se/pronlex/lexserver && go run *.go $switches -test && go run *.go $switches $PORT
fi
