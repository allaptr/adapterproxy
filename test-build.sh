#!/usr/bin/env bash
set -x

docker build -t adapterproxy .

docker build -f DockerfileV1 -t testv1 .
docker build -f DockerfileV2 -t testv2 .

set +x