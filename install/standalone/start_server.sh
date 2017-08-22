#############################################################
### SERVER STARTUP SCRIPT
## 1. STARTS A TEST SERVER AND RUNS A SET OF TESTS
## 2. SHUTS DOWN TEST SERVER
## 3. IF NO ERRORS ARE FOUND, THE STANDARD SERVER IS STARTED USING THE DEFAULT PORT

if [ $# -ne 1 ]; then
    echo "[$0] SERVER STARTUP SCRIPT
   1. STARTS A TEST SERVER AND RUNS A SET OF TESTS
   2. SHUTS DOWN TEST SERVER
   3. IF NO ERRORS ARE FOUND, THE STANDARD SERVER IS STARTED USING THE DEFAULT PORT

USAGE: $0 <APPDIR>
       where <APPDIR> is the folder in which the build script installed the standalone lexserver

EXAMPLE INVOCATION: $0 lexserver_files
"
    exit 1
fi


APPDIR=$1
switches="-ss_files $APPDIR/symbol_sets/ -db_files $APPDIR/db_files/ -static $APPDIR/static/"
lexserver $switches -test && lexserver $switches
