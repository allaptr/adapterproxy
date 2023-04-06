#!/usr/bin/env bash
proxyimage=$1
v1image=$2
v2image=$3
set -x
docker container run --network="host" ${proxyimage} timeout=150 loglevel=DebugLevel us=http://localhost:9002 ru=http://localhost:9003 &
docker container run --network="host" ${v1image} &
docker container run --network="host" ${v2image} delay=250 &
set +x
