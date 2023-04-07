#!/bin/bash

set -eux
set -o errexit
set -o pipefail
set -o nounset

docker compose down --remove-orphans
docker compose build --no-cache
docker compose up -d
clear
docker compose logs -f server
