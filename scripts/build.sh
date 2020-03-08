#!/bin/sh

PROJECT="phoneinfoga"

echo ">> Building binaries for darwin and linux"

# Linux
echo ">>> Building for arch linux/amd64"
GOOS="linux"
GOARCH="amd64"
$GOPATH/bin/packr2 build -o "${PROJECT}_${GOOS}_${GOARCH}"

echo ">>> Building for arch linux/arm64"
GOOS="linux"
GOARCH="arm64"
$GOPATH/bin/packr2 build -o "${PROJECT}_${GOOS}_${GOARCH}"

echo ">>> Building for arch linux/armv6"
GOOS="linux"
GOARCH="armv6"
$GOPATH/bin/packr2 build -o "${PROJECT}_${GOOS}_${GOARCH}"

echo ">>> Building for arch linux/i386"
GOOS="linux"
GOARCH="i386"
$GOPATH/bin/packr2 build -o "${PROJECT}_${GOOS}_${GOARCH}"

echo ">>> Building for arch linux/x86_64"
GOOS="linux"
GOARCH="x86_64"
$GOPATH/bin/packr2 build -o "${PROJECT}_${GOOS}_${GOARCH}"

# Darwin
echo ">>> Building for arch darwin/x86_64"
GOOS="darwin"
GOARCH="x86_64"
$GOPATH/bin/packr2 build -o "${PROJECT}_${GOOS}_${GOARCH}"

echo ">>> Building for arch darwin/i386"
GOOS="darwin"
GOARCH="i386"
$GOPATH/bin/packr2 build -o "${PROJECT}_${GOOS}_${GOARCH}"

echo ">>> Cleaning packr files"
$GOPATH/bin/packr2 clean

echo ">> Done."