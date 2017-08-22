FROM golang

RUN apt-get update -y && apt-get upgrade -y && apt-get install apt-utils -y

RUN apt-get update -y && apt-get upgrade -y && apt-get install sqlite3 -y && apt-get install git -y && apt-get install gcc -y && apt-get install build-essential -y

#ENV HOST_DIR lexserver_files

RUN go get github.com/stts-se/pronlex/lexserver 
RUN go install github.com/stts-se/pronlex/lexserver 

RUN go get github.com/stts-se/pronlex/cmd/lexio/createEmptyDB
RUN go install github.com/stts-se/pronlex/cmd/lexio/createEmptyDB

RUN go get github.com/stts-se/pronlex/cmd/lexio/createEmptyDB
RUN go install github.com/stts-se/pronlex/cmd/lexio/import

#ADD init_lexica.sh .

#CMD bash init_lexica.sh

EXPOSE 8787

RUN echo "Mount external host dir to /go/lexserver_files"

#CMD lexserver -ss_files lexserver_files/symbol_sets -db_files lexserver_files/db_files -static lexserver_files/static

CMD (lexserver -test -ss_files lexserver_files/symbol_sets -db_files lexserver_files/db_files -static lexserver_files/static && lexserver -ss_files lexserver_files/symbol_sets -db_files lexserver_files/db_files -static lexserver_files/static)
