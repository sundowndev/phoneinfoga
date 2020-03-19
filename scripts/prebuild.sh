#!/bin/sh

PROJECT="phoneinfoga"

echo ">> Building web client"

(cd client && yarn && yarn build)

echo ">> Building assets"
$GOPATH/bin/packr2

echo "> Done."