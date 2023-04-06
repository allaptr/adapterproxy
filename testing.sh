#!/usr/bin/env bash
set -x
file=${1}
while read -r line; do
# echo -e "$line\n"
res=$(curl -v "localhost:9000/company?country_iso=${2}&id=$line")
echo $res
sleep 1
done <$file
set +x

