#!/usr/bin/env bash
proxyimage=$1
v1image=$2
v2image=$3
set -x
docker container run --network="host" ${proxyimage} --timeout=150ms --loglevel=DebugLevel us=http://localhost:9002 ru=http://localhost:9003 &
if  ${v1image}!="", then 
docker container run --network="host" ${v1image} &
elif ${v2image}!="", then
docker container run --network="host" ${v2image} delay=250 &
fi
set +x
