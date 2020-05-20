# build script for Wikispeech
# mimic travis build tests, always run before pushing!

set -e

SLEEP=60

if [ $# -ne 0 ]; then
    echo "For developers: If you are developing for Wikispeech, and need to make changes to this repository, make sure you run a test build using build_and_test.sh before you make a pull request. Don't run more than one instance of this script at once, and make sure no pronlex server is already running on the default port."
    exit 0
fi


basedir=`dirname $0`
basedir=`realpath $basedir`
echo $basedir
cd $basedir
mkdir -p .build

#go test -v ./... 

#gosec ./...
#staticcheck ./...

mkdir -p .build/appdir

for proc in `ps --sort pid -Af|egrep pronlex| egrep -v  "grep .E"|sed 's/  */\t/g'|cut -f2`; do
    kill $proc || echo "Couldn't kill $pid"
done

bash scripts/setup.sh -a .build/appdir -e sqlite 

bash scripts/start_server.sh -a .build/appdir -e sqlite &
export pid=$!
echo "pronlex server started on pid $pid. wait for $SLEEP seconds before shutting down"
sleep $SLEEP
sh .travis/exit_server_and_fail_if_not_running.sh pronlex $pid
 
