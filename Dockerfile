FROM golang

RUN apt-get update -y && apt-get upgrade -y && apt-get install apt-utils -y

RUN apt-get install -y sqlite3 git gcc build-essential

# FOR DEBUGGING
RUN apt-get install -y libnet-ifconfig-wrapper-perl/stable

RUN go get github.com/stts-se/pronlex/lexserver 
RUN go install github.com/stts-se/pronlex/lexserver 

RUN go get github.com/stts-se/pronlex/cmd/lexio/createEmptyDB
RUN go install github.com/stts-se/pronlex/cmd/lexio/createEmptyDB

RUN go get github.com/stts-se/pronlex/cmd/lexio/importLex
RUN go install github.com/stts-se/pronlex/cmd/lexio/importLex

ENV APPDIR appdir

RUN export GOPATH=$(go env GOPATH)
RUN export PATH=$PATH:$(go env GOPATH)/bin

RUN ln -s /go/src/github.com/stts-se/pronlex/docker/setup bin/setup
RUN ln -s /go/src/github.com/stts-se/pronlex/docker/help bin/help

# import script
RUN ln -s /go/src/github.com/stts-se/pronlex/docker/import bin/import_all0
RUN echo "#!/bin/bash" > bin/import_all
RUN echo "setup $APPDIR && import_all0 -a $APPDIR" >> bin/import_all

RUN chmod --silent +x bin/*

EXPOSE 8787

RUN echo "Mount external host dir to $APPDIR"

CMD (setup $APPDIR && lexserver -test -ss_files $APPDIR/symbol_sets -db_files $APPDIR/db_files -static $GOPATH/src/github.com/stts-se/pronlex/lexserver/static && lexserver -ss_files $APPDIR/symbol_sets -db_files $APPDIR/db_files -static $GOPATH/src/github.com/stts-se/pronlex/lexserver/static 8787)

