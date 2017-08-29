#!/usr/bin/env bash

set -xe

DEV_IMAGE_NAME="${DEV_IMAGE_NAME:-integration-kit-dev}"

docker build -t ${DEV_IMAGE_NAME} -f Dockerfile .

docker run --rm --init -it \
  ${DEV_IMAGE_NAME} $*
