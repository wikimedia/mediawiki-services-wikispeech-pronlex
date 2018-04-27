export PATH=$PATH:/usr/local/go/bin/go
export GOPATH=`go env GOPATH`
export PATH=$PATH:$GOPATH/bin

bash start_server.sh $*
