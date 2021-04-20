#!/usr/bin/env bash

#
# This script is executed from within the docker image during Blubber build.
#

mkdir pronlex
mv * pronlex

m_error() {
  echo $1
  exit 2
}

install_go() {
  cd /srv
  if [ ! -f /tmp/go1.13.linux-amd64.tar.gz ]; then
   if ! wget https://dl.google.com/go/go1.13.linux-amd64.tar.gz -O /tmp/go1.13.linux-amd64.tar.gz; then
     m_error "Unable to download Go lang 1.13 from Google!"
   fi
  fi
  tar xvfz /tmp/go1.13.linux-amd64.tar.gz
  echo "Go installed"
}

install_pronlex() {
  cd /srv/pronlex

  if ! go build ./...; then
    m_error "Failed to build Pronlex!"
  fi

  echo "Setting up Pronlex"
  cd /srv/pronlex
  /bin/bash scripts/setup.sh -a /srv/appdir -e sqlite


  echo "Cloning Lexdata"
  cd /srv
  if ! git clone https://github.com/stts-se/wikispeech-lexdata.git lexdata; then
    m_error "Unable to clone Lexdata from Git repo"
  fi

  echo "Importing Lexdata"
  cd /srv/pronlex
  /bin/bash scripts/import.sh -a /srv/appdir -e sqlite -f /srv/lexdata

  echo "Starting Pronlex server. Will wait a minute for it to start up and download any dependencies, and then kill it."
  /bin/bash scripts/start_server.sh -a /srv/appdir -p 8080 -e sqlite &
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

install_go

export GOROOT=/srv/go
export GOPATH=/srv/goProjects
export PATH=${GOPATH}/bin:${GOROOT}/bin:${PATH}

install_pronlex

echo "Successfully prepared Pronlex!"
