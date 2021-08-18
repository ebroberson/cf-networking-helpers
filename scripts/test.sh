#!/bin/bash

set -eu
set -o pipefail

cd $(dirname $0)/..

BIN_DIR="${PWD}/bin"
mkdir -p "${BIN_DIR}"
export PATH="${PATH}:${BIN_DIR}"

go build -o "$BIN_DIR/ginkgo" github.com/onsi/ginkgo/ginkgo

if [ "${1:-""}" = "" ]; then
  extraArgs=""
else
  extraArgs="${@}"
fi

if [ ${DB:-"none"} = "mysql" ] || [ ${DB:-"none"} = "mysql-5.6" ]; then
  # bootMysql
  ginkgo -r --race -randomizeAllSpecs ${extraArgs} db/timeouts
elif [ ${DB:-"none"} = "postgres" ]; then
  # bootPostgres
  ginkgo -r --race -randomizeAllSpecs ${extraArgs} db/timeouts
else
  echo "skipping database"
  extraArgs="-skipPackage=db ${extraArgs}"
fi

ginkgo -r -p --race -randomizeAllSpecs -randomizeSuites -skipPackage=timeouts ${extraArgs}
