#!/bin/bash
docker tag smorton517/hledger-htmx:latest smorton517/hledger-htmx:"$1" &&
docker push smorton517/hledger-htmx:"$1" &&
git tag "$1" &&
git push --tags
