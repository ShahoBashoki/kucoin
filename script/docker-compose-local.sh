#!/bin/bash

set -eux
set -o errexit
set -o pipefail
set -o nounset

make swagger
make build

docker compose down --volumes --remove-orphans
docker compose build --no-cache
docker compose up -d
clear
docker compose logs -f server
