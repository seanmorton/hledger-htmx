#!/bin/bash
docker tag seanmorton/hledger-htmx:latest seanmorton/hledger-htmx:"$1" &&
docker push seanmorton/hledger-htmx:latest &&
docker push seanmorton/hledger-htmx:"$1" &&
git tag "$1" &&
git push --tags
