#/bin/bash

CMD=`basename $0`

while getopts ":ha:" opt; do
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
" >&2
	exit 1
      ;;
    a)
        APPDIR=$OPTARG
      ;;
    \?)
      echo "Invalid option: -$OPTARG" >&2
      ;;
  esac
done

shift $(expr $OPTIND - 1 )


if [ -z "$APPDIR" ] ; then
    echo "[$CMD] APPDIR must be specified!" >&2
    exit 1
fi

mkdir -p $APPDIR
APPDIRABS=`realpath $APPDIR`
PORT="8787"
USER=`stat -c "%u:%g" $APPDIR`

docker run -u $USER -v $APPDIRABS:/go/appdir -p $PORT:8787 -it sttsse/lexserver $*
