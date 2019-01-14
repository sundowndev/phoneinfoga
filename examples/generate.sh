#!/bin/bash

scriptDir=$(dirname -- "$(readlink -f -- "$BASH_SOURCE")")

python3 $scriptDir/../phoneinfoga.py -i $scriptDir/input.txt -o $scriptDir/output_from_input.txt -s any

python3 $scriptDir/../phoneinfoga.py -n "+86 591 2284 8571" -s all -o $scriptDir/output_single.txt
