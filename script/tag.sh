#!/bin/bash

set -eux
set -o errexit
set -o pipefail
set -o nounset

PACKAGE=$(echo ${1} | sed -En 's/(development|production)\/(kucoin)-v(([0-9]{4}.[0-9]{2}.[0-9]{2})(-rc|-rc(\.[1-9]+))?)/\2/gp')
FULL_VERSION=$(echo ${1} | sed -En 's/(development|production)\/(kucoin)-v(([0-9]{4}.[0-9]{2}.[0-9]{2})(-rc|-rc(\.[1-9]+))?)/\3/gp')
VERSION=$(echo ${1} | sed -En 's/(development|production)\/(kucoin)-v(([0-9]{4}.[0-9]{2}.[0-9]{2})(-rc|-rc(\.[1-9]+))?)/\4/gp')

COMMIT_SHA=${2:-""}

GIT_COMMITTER_DATE=$(git --no-pager log --format="%aD" --max-count=1 "${COMMIT_SHA}") git tag -a -m "${PACKAGE} version ${FULL_VERSION}" "${1}" "${COMMIT_SHA}"

git push --tags
