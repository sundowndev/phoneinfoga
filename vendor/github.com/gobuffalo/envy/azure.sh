#!/bin/bash

set -xe

cat >> .env << EOF
# This is a comment
# We can use equal or colon notation
DIR: root
FLAVOUR: none
INSIDE_FOLDER=false
EOF
