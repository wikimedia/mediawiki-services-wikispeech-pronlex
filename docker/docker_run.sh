#/bin/bash

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
  -t docker-tag (required)
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
    echo "[$CMD] FAILED: APPDIR must be specified using -a!" >&2
    exit 1
fi

if [ -z "$DOCKERTAG" ] ; then
    echo "[$CMD] FAILED: DOCKERTAG must be specified using -t!" >&2
    exit 1
fi

echo "[$CMD] APPDIR    : $APPDIR" >&2
echo "[$CMD] PORT      : $PORT" >&2
echo "[$CMD] DOCKERTAG : $DOCKERTAG" >&2

#mkdir -p $APPDIR

APPDIRABS=`realpath $APPDIR`

NETWORKARGS="--network=wikispeech"

CNAME="pronlex"
if docker container inspect $CNAME &> /dev/null ; then
    echo -n "STOPPING CONTAINER "
    docker stop $CNAME
    echo -n "DELETING CONTAINER "
    docker rm $CNAME
fi


docker run --name=pronlex $NETWORKARGS -v $APPDIRABS:/go/appdir -p $PORT:8787 -it $DOCKERTAG $*
