#!/bin/sh

if !(hash curl 2>/dev/null); then
        echo "Curl is needed to run prebuild."
        exit 1;
fi;

if !(hash yarn 2>/dev/null); then
        curl -o- -L https://yarnpkg.com/install.sh | bash
fi;

echo ">> Building web client"
(cd client && yarn && yarn build)

echo ">> Building static assets"
$GOPATH/bin/packr2

echo "> Done."