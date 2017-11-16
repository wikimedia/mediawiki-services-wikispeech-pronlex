# Download sttsse/wikispeech_base from hub.docker.com | source repository: https://github.com/stts-se/wikispeech_base.git
FROM sttsse/wikispeech_base

RUN mkdir -p /wikispeech/bin
WORKDIR "/wikispeech"

RUN go get github.com/stts-se/pronlex/lexserver 
RUN go install github.com/stts-se/pronlex/lexserver 

RUN go get github.com/stts-se/pronlex/cmd/lexio/createEmptyDB
RUN go install github.com/stts-se/pronlex/cmd/lexio/createEmptyDB

RUN go get github.com/stts-se/pronlex/cmd/lexio/importLex
RUN go install github.com/stts-se/pronlex/cmd/lexio/importLex

ENV APPDIR /appdir

RUN export GOPATH=$(go env GOPATH)
RUN export PATH=$PATH:$(go env GOPATH)/bin

RUN ln -s /go/src/github.com/stts-se/pronlex/docker/setup /wikispeech/bin/setup
RUN ln -s /go/src/github.com/stts-se/pronlex/docker/import_all /wikispeech/bin/import_all
RUN chmod +x /wikispeech/bin/*

# BUILD INFO
RUN echo -n "Build timestamp: " > /wikispeech/.pronlex_build_info.txt
RUN date --utc "+%Y-%m-%d %H:%M:%S %Z" >> /wikispeech/.pronlex_build_info.txt
RUN echo "Built by: docker" >> /wikispeech/.pronlex_build_info.txt
RUN echo "Application name: pronlex"  >> /wikispeech/.pronlex_build_info.txt


# RUNTIME SETTINGS

EXPOSE 8787

# RUN echo "Mount external host dir to $APPDIR"

CMD (/wikispeech/bin/setup $APPDIR && lexserver -test -ss_files $APPDIR/symbol_sets -db_files $APPDIR/db_files -static $GOPATH/src/github.com/stts-se/pronlex/lexserver/static && lexserver -ss_files $APPDIR/symbol_sets -db_files $APPDIR/db_files -static $GOPATH/src/github.com/stts-se/pronlex/lexserver/static 8787)

