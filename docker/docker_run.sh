#/bin/bash

CMD=`basename $0`

PORT="8787"

function echo_usage {
	echo "
USAGE:

# SETUP lex server
  $ $CMD -a <APPDIR> -t <DOCKERTAG> setup

# IMPORT lexicon data (optional)
  $ $CMD -a <APPDIR> -t <DOCKERTAG> lex_import

# RUN lex server
  $ $CMD -a <APPDIR> -t <DOCKERTAG>

# BASH inspect lex server
  $ $CMD -a <APPDIR> -t <DOCKERTAG> bash

Options:
  -h help
  -a appdir     (required)
  -p port       (default: $PORT)
  -t docker-tag (required)
" >&2
}

if [ $# -eq 0 ]; then
    echo_usage
    exit 1
fi


while getopts ":ha:p:t:" opt; do
  case $opt in
    h)
	echo_usage
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

CNAME="pronlex"

function shutdown_previos {
    if docker container inspect $CNAME &> /dev/null ; then
	echo -n "[$CMD] STOPPING CONTAINER "
	docker stop $CNAME
	echo -n "[$CMD] DELETING CONTAINER "
	docker rm $CNAME
    fi
}

if [ $# -eq 0 ]; then
    shutdown_previos && docker run --name=pronlex -v $APPDIRABS:/appdir -p $PORT:8787 -it $DOCKERTAG
elif [ $# -eq 1 ] && [ $1 == "setup" ]; then
    shutdown_previos && docker run --name=pronlex -v $APPDIRABS:/appdir -p $PORT:8787 -it $DOCKERTAG setup /appdir
elif [ $# -eq 1 ] && [ $1 == "lex_import" ]; then
    shutdown_previos && docker run --name=pronlex -v $APPDIRABS:/appdir -p $PORT:8787 -it $DOCKERTAG import_all /appdir $APPDIRABS
elif [ $# -eq 1 ] && [ $1 == "bash" ]; then
    shutdown_previos && docker run --name=pronlex -v $APPDIRABS:/appdir -p $PORT:8787 -it $DOCKERTAG bash
else
    echo "[$CMD] Unknown command: $*" >&2
    exit 1
fi

