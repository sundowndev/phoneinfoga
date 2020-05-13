#!/bin/sh

echo "Formatting Go files..."

gofmt -s -w -l .

echo "Done."
