if [ $# -ne 1 ]; then
    echo "USAGE: sh $CMD <FULL DOCKERTAG>" >&2
    exit 1
fi

DOCKERTAG=$1

docker login && docker push $DOCKERTAG
