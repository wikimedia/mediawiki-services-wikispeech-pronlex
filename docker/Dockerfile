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

# setup script
RUN echo "#!/bin/bash" > bin/setup
RUN echo "cp -r $GOPATH/src/github.com/stts-se/pronlex/lexserver/static $APPDIR/static/ && cp -r $GOPATH/src/github.com/stts-se/pronlex/lexserver/demo_files $APPDIR/symbol_sets" >> bin/setup

# import_lex script
RUN ln -s $GOPATH/src/github.com/stts-se/pronlex/install/standalone/import.sh bin/import_lex0
RUN echo "#!/bin/bash" > bin/import_lex
RUN echo "sh bin/import_lex0 -a $APPDIR" >> bin/import_lex

RUN chmod +x bin/*

EXPOSE 8787

RUN echo "Mount external host dir to $APPDIR"

#CMD lexserver -ss_files $APPDIR/symbol_sets -db_files $APPDIR/db_files -static $APPDIR/static

CMD (lexserver -test -ss_files $APPDIR/symbol_sets -db_files $APPDIR/db_files -static $APPDIR/static && lexserver -ss_files $APPDIR/symbol_sets -db_files $APPDIR/db_files -static $APPDIR/static 8787)

