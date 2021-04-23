#!/usr/bin/env bash

echo "Starting Pronlex..."

DIR=`pwd`

export GOROOT=${DIR}/go
export GOPATH=${DIR}/goProjects
export PATH=${GOPATH}/bin:${GOROOT}/bin:${PATH}

cd pronlex

if [[ -z "${PRONLEX_MARIADB_URI}" ]]; then
  /bin/bash scripts/start_server.sh -a ${DIR}/appdir -e sqlite -p 8787 -r lexserver
else
  /bin/bash scripts/start_server.sh -a ${DIR}/appdir -e mariadb -l "${PRONLEX_MARIADB_URI}" -p 8787 -r lexserver
fi
