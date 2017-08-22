if [ $# -ne 1 ]; then
    echo "USAGE: sh $0 <APPDIR>"
    exit 1
fi

APPDIR=$1

### LEXSERVER INSTALL

install.sh $APPDIR

### LEXDATA SETUP

import.sh $APPDIR

### COMPLETED

echo "
BUILD COMPLETED. YOU CAN NOW START THE LEXICON SERVER BY INVOKING
  $ sh run_standalone.sh $APPDIR
"
