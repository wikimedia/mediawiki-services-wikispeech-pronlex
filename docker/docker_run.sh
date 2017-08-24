#/bin/bash

DOCKERTAG="stts-lexserver-local"

CMD=`basename $0`

PORT="8787"

while getopts ":ha:p:t:" opt; do
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
  -a appdir     (required)
  -p port       (default: $PORT)
  -t docker-tag (default: $DOCKERTAG)
" >&2
	exit 1
      ;;
    a)
        APPDIR=$OPTARG
      ;;
    p)
        PORT=$OPTARG
      ;;
    t)
        DOCKERTAG=$OPTARG
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

echo "[$CMD] APPDIR    : $APPDIR" >%2
echo "[$CMD] PORT      : $PORT" >%2
echo "[$CMD] DOCKERTAG : $DOCKERTAg" >%2

#mkdir -p $APPDIR
#chgrp docker $APPDIR

APPDIRABS=`realpath $APPDIR`

## => use system user inside container
# docker run -u $USER -v $APPDIRABS:/go/appdir -p $PORT:8787 -it $DOCKERTAG $*

## => root user
docker run -v $APPDIRABS:/go/appdir -p $PORT:8787 -it $DOCKERTAG $*
