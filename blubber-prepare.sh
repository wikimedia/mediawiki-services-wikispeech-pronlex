#!/usr/bin/env bash

DIR=`pwd`

m_error() {
  echo $1
  exit 2
}

install_go() {
  cd ${DIR}/blubber
  if [ ! -f /tmp/go1.13.linux-amd64.tar.gz ]; then
   if ! wget https://dl.google.com/go/go1.13.linux-amd64.tar.gz -O /tmp/go1.13.linux-amd64.tar.gz; then
     m_error "Unable to download Go lang 1.13 from Google!"
   fi
  fi
  if [ ! -d ${DIR}/blubber/go ]; then
    tar xvfz /tmp/go1.13.linux-amd64.tar.gz
  fi
  echo "Go installed"
}

install_pronlex() {
  cd ${DIR}/blubber/pronlex

  if ! go build ./...; then
    m_error "Failed to build Pronlex!"
  fi

  # todo consider download testdata and test pronlex

  if [ -d ${DIR}/blubber/appdir ]; then
    rm -rf ${DIR}/blubber/appdir
  fi

  echo "Setting up Pronlex"
  cd ${DIR}/blubber/pronlex
  /bin/bash scripts/setup.sh -a ${DIR}/blubber/appdir -e sqlite


  echo "Setting up Lexdata"
  if [ -d ${DIR}/blubber/lexdata ]; then
    cd ${DIR}/blubber/lexdata
    if ! git pull; then
      m_error "Unable to update Lexdata from Git repo"
    fi
  else
    cd ${DIR}/blubber
    if ! git clone https://github.com/stts-se/lexdata.git; then
      m_error "Unable to close Lexdata from Git repo"
    fi
  fi

  echo "Importing Lexdata"
  cd ${DIR}/blubber/pronlex

  /bin/bash scripts/import.sh -a ${DIR}/blubber/appdir -e sqlite -f ${DIR}/blubber/lexdata

  echo "Starting Pronlex server. Will wait a minute for it to start up and download any dependencies, and then kill it."
  /bin/bash scripts/start_server.sh -a ${DIR}/blubber/appdir -p 8080 -e sqlite &
  PRONLEX_PID=$!
  for i in $(seq 1 6); do
    if ! kill -0 ${PRONLEX_PID}; then
      echo "Warning! Pronlex process has prematurely ended!"
      break
    fi
    sleep 10
    echo "${i}0/60 seconds slept before killing server..."
  done
  kill ${PRONLEX_PID}
}

if [ ! -d blubber ]; then
  cp -r . /tmp/pronlex_tmp
  mkdir blubber
  cd blubber
  mv /tmp/pronlex_tmp pronlex
fi

if [ ! -d ${DIR}/blubber/go ]; then
  install_go
fi

export GOROOT=${DIR}/blubber/go
export GOPATH=${DIR}/blubber/goProjects
export PATH=${GOPATH}/bin:${GOROOT}/bin:${PATH}

install_pronlex

echo "Successfully prepared Pronlex! Now run ./blubber-build.sh"
