#!/bin/bash 

#############################################################
### SERVER STARTUP SCRIPT
## 1. STARTS A TEST SERVER AND RUNS A SET OF TESTS
## 2. SHUTS DOWN TEST SERVER
## 3. IF NO ERRORS ARE FOUND, THE STANDARD SERVER IS STARTED USING THE DEFAULT PORT


CMD=`basename $0`
APPDIR=`dirname $0`
PORT="8787"

while getopts ":hp:a:" opt; do
  case $opt in
    h)
    echo "
[$CMD] SERVER STARTUP SCRIPT
   1. STARTS A TEST SERVER AND RUNS A SET OF TESTS
   2. SHUTS DOWN TEST SERVER
   3. IF NO ERRORS ARE FOUND, THE STANDARD SERVER IS STARTED USING THE DEFAULT PORT

Options:
  -h help
  -a appdir (default: this script's folder)
  -p port   (default: 8787)

EXAMPLE INVOCATION: $CMD -a lexserver_files
" >&2
	echo "USAGE: sh $CMD <OPTIONS>

Imports lexicon data for Swedish, Norwegian, US English and a small test file for Arabic.

Options:
  -h help
  -a appdir (default: this script's folder)
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
      ;;
  esac
done

shift $(expr $OPTIND - 1 )

if [ $# -ne 0 ]; then
    echo "[$CMD] invalid option(s): $*" >&2
    exit 1
fi


APPDIR=$1
switches="-ss_files $APPDIR/symbol_sets/ -db_files $APPDIR/db_files/ -static $APPDIR/static/"
echo "lexserver $switches -test && lexserver $switches -port $PORT"
