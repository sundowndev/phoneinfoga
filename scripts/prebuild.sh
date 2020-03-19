#!/bin/sh

if !(hash curl 2>/dev/null); then
        echo "Curl is needed to run prebuild."
        exit 1;
fi;

if !(hash node 2>/dev/null); then
        VERSION=v12.16.1
        DISTRO=linux-x64
        sudo mkdir -p /usr/local/lib/nodejs
        sudo tar -xJvf node-$VERSION-$DISTRO.tar.xz -C /usr/local/lib/nodejs
fi;

if !(hash yarn 2>/dev/null); then
        curl -o- -L https://yarnpkg.com/install.sh | bash
fi;

echo ">> Building web client"
(cd client && $HOME/.yarn/bin/yarn && $HOME/.yarn/bin/yarn build)

echo ">> Building static assets"
$GOPATH/bin/packr2

echo "> Done."
