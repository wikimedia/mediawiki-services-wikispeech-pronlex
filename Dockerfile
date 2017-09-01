FROM golang

RUN apt-get update -y && apt-get upgrade -y && apt-get install apt-utils -y

RUN apt-get update -y && apt-get upgrade -y && apt-get install sqlite3 -y && apt-get install git -y && apt-get install gcc -y && apt-get install build-essential -y

RUN go get github.com/stts-se/pronlex/lexserver 
RUN go install github.com/stts-se/pronlex/lexserver 

RUN go get github.com/stts-se/pronlex/cmd/lexio/createEmptyDB
RUN go install github.com/stts-se/pronlex/cmd/lexio/createEmptyDB

RUN go get github.com/stts-se/pronlex/cmd/lexio/importLex
RUN go install github.com/stts-se/pronlex/cmd/lexio/importLex

ENV APPDIR appdir

RUN export GOPATH=$(go env GOPATH)
RUN export PATH=$PATH:$(go env GOPATH)/bin

RUN echo $GOPATH

## CHANGE USER FROM ROOT TO SYSTEM USER
#ARG USER
#RUN useradd $USER
#RUN groupadd docker
#USER $user

# setup script
RUN echo "#!/bin/bash" > bin/setup
RUN echo "cp -r /go/src/github.com/stts-se/pronlex/lexserver/demo_files/* $APPDIR/symbol_sets/" >> bin/setup

# import script
RUN ln -s /go/src/github.com/stts-se/pronlex/install/standalone/import.sh bin/import_all0
RUN echo "#!/bin/bash" > bin/import_all
RUN echo "sh bin/import_all0 -a $APPDIR" >> bin/import_all

RUN chmod +x bin/*

EXPOSE 8787

RUN echo "Mount external host dir to $APPDIR"

CMD (lexserver -test -ss_files $APPDIR/symbol_sets -db_files $APPDIR/db_files -static $GOPATH/src/github.com/stts-se/pronlex/lexserver/static && lexserver -ss_files $APPDIR/symbol_sets -db_files $APPDIR/db_files -static $GOPATH/src/github.com/stts-se/pronlex/lexserver/static 8787)

