#/bin/bash

APPDIR="appdir"
mkdir -p $APPDIR
APPDIRABS=`realpath $APPDIR`
PORT="8787"
USER=`stat -c "%u:%g" $APPDIR`

docker run -u $USER -v $APPDIRABS:/go/appdir -p $PORT:8787 -it sttsse/lexserver $*
