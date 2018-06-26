# build script for Wikispeech
# mimic travis build tests, always run before pushing!

set -e

SLEEP=60

basedir=`dirname $0`
basedir=`realpath $basedir`
echo $basedir
cd $basedir
mkdir -p .build

go test -v ./... 

mkdir -p .build/appdir

for proc in `ps --sort pid -Af|egrep pronlex| egrep -v  "grep .E"|sed 's/  */\t/g'|cut -f2`; do
    kill $proc || "Couldn't kill $pid"
done

bash install/setup.sh .build/appdir

bash install/start_server.sh -a .build/appdir &
export pid=$!
echo "pronlex server started on pid $pid. wait for $SLEEP seconds before shutting down"
sleep $SLEEP
sh .travis/exit_server_and_fail_if_not_running.sh pronlex $pid
 
docker build . --no-cache -t sttsse/pronlex:buildtest
