#/bin/bash

DOCKERNAME="stts-lexserver-local"

CMD=`basename $0`

PORT="8787"

while getopts ":ha:p:" opt; do
  case $opt in
    h)
	echo "
USAGES:

# SETUP lex server
  $ $CMD -a <APPDIR> setup

# IMPORT lexicon data (optional)
  $ $CMD -a <APPDIR> lex_import

# RUN lex server
  $ $CMD -a <APPDIR> 

# BASH inspect lex server
  $ $CMD -a <APPDIR> bash

Options:
  -h help
  -a appdir (required)
  -p port   (default: $PORT)
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


if [ -z "$APPDIR" ] ; then
    echo "[$CMD] FAILED: APPDIR must be specified!" >&2
    exit 1
fi

mkdir -p $APPDIR
chgrp docker $APPDIR
APPDIRABS=`realpath $APPDIR`

docker run -u $USER -v $APPDIRABS:/go/appdir -p $PORT:8787 -it $DOCKERNAME $*
