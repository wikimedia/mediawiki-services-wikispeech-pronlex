FROM golang

RUN apt-get update -y && apt-get upgrade -y && apt-get install apt-utils -y

RUN apt-get update -y && apt-get upgrade -y && apt-get install sqlite3 -y && apt-get install git -y && apt-get install gcc -y && apt-get install build-essential -y

RUN go get github.com/stts-se/pronlex/lexserver 
RUN go install github.com/stts-se/pronlex/lexserver 

RUN go get github.com/stts-se/pronlex/cmd/lexio/createEmptyDB
RUN go install github.com/stts-se/pronlex/cmd/lexio/createEmptyDB

RUN go get github.com/stts-se/pronlex/cmd/lexio/importLex
RUN go install github.com/stts-se/pronlex/cmd/lexio/importLex

RUN ln -s $GOPATH/src/github.com/stts-se/pronlex/install/standalone/import.sh import_lex0
RUN echo "sh import_lex0 -a lexserver_files" > import_lex
#ENV PATH /usr/local/nginx/bin:$PATH:.
#RUN "echo $PATH" > .profile
#RUN /bin/bash 'source .profile'
#RUN /bin/bash 'chmod +x import_lex'

EXPOSE 8787

RUN echo "Mount external host dir to /go/lexserver_files"

#CMD lexserver -ss_files lexserver_files/symbol_sets -db_files lexserver_files/db_files -static lexserver_files/static

CMD (lexserver -test -ss_files lexserver_files/symbol_sets -db_files lexserver_files/db_files -static lexserver_files/static && lexserver -ss_files lexserver_files/symbol_sets -db_files lexserver_files/db_files -static lexserver_files/static)

## TODO: root should not own created files! fix this!
## TODO: set path so that we can run "import_lex" without the "sh" prefix
## TODO: specify APPDIR variable or similar, instead of hardwired 'lexserver_files'
