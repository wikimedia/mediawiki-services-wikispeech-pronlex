#!/usr/bin/env bash

echo "Starting Pronlex..."

DIR=`pwd`

cd pronlex

if [[ -z "${PRONLEX_MARIADB_URI}" ]]; then
  /bin/bash scripts/start_server.sh -a ${DIR}/appdir -e sqlite -p 8787 -r lexserver
else
  /bin/bash scripts/start_server.sh -a ${DIR}/appdir -e mariadb -l "${PRONLEX_MARIADB_URI}" -p 8787 -r lexserver
fi
