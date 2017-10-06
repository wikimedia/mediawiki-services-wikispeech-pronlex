# Download sttsse/wikispeech_base from hub.docker.com | source repository: https://github.com/stts-se/wikispeech_base.git
FROM sttsse/wikispeech_base

WORKDIR "/"

RUN go get github.com/stts-se/pronlex/lexserver 
RUN go install github.com/stts-se/pronlex/lexserver 

RUN go get github.com/stts-se/pronlex/cmd/lexio/createEmptyDB
RUN go install github.com/stts-se/pronlex/cmd/lexio/createEmptyDB

RUN go get github.com/stts-se/pronlex/cmd/lexio/importLex
RUN go install github.com/stts-se/pronlex/cmd/lexio/importLex

ENV APPDIR appdir

RUN export GOPATH=$(go env GOPATH)
RUN export PATH=$PATH:$(go env GOPATH)/bin

RUN ln -s /go/src/github.com/stts-se/pronlex/docker/setup /bin/setup
RUN ln -s /go/src/github.com/stts-se/pronlex/docker/help /bin/help

# import script
RUN ln -s /go/src/github.com/stts-se/pronlex/docker/import_all /bin/import_all

RUN chmod +x /bin/*

# RUNTIME SETTINGS

EXPOSE 8787

# RUN echo "Mount external host dir to $APPDIR"

CMD (setup $APPDIR && lexserver -test -ss_files $APPDIR/symbol_sets -db_files $APPDIR/db_files -static $GOPATH/src/github.com/stts-se/pronlex/lexserver/static && lexserver -ss_files $APPDIR/symbol_sets -db_files $APPDIR/db_files -static $GOPATH/src/github.com/stts-se/pronlex/lexserver/static 8787)

