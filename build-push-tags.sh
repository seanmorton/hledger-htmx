#!/bin/bash
if [ "$#" -ne 1 ]; then
  echo "verison tag is required"
  exit 1
fi

docker build --tag seanmorton/hledger-htmx:latest --platform linux/amd64 .
docker tag seanmorton/hledger-htmx:latest seanmorton/hledger-htmx:"$1" &&
docker push seanmorton/hledger-htmx:latest &&
docker push seanmorton/hledger-htmx:"$1" &&
git tag "$1" &&
git push --tags
