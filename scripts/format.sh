#!/bin/sh

go fmt $(go list ./... | grep -v /vendor/)