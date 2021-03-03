#!/usr/bin/env bash

echo "Starting Pronlex..."

DIR=`pwd`

export GOROOT=${DIR}/go
export GOPATH=${DIR}/goProjects
export PATH=${GOPATH}/bin:${GOROOT}/bin:${PATH}

cd pronlex

/bin/bash scripts/start_server.sh -a ${DIR}/appdir -e sqlite -p 8787 -r lexserver&

PID=$!
sleep 60
if ! kill -0 ${PID}; then
  echo "ERROR: Service process has prematurely ended!"
  exit 1
fi
wget -O /dev/null -o /dev/null "http://localhost:8787/lexicon/lookup?lexicons=wikispeech_lexserver_demo%3Asv&words=hund%2C%20h%C3%A4st"
EXIT_CODE=$?
kill ${PID}
if [ ${EXIT_CODE} -ne 0 ]; then
  echo "ERROR: Test failed!"
else
  echo "Test successful!"
fi
exit ${EXIT_CODE}
