#!/bin/bash



{
    # command which may fail and give an error 
    scriptDir=$(dirname -- "$(readlink -f -- "$BASH_SOURCE")")

    python3 $scriptDir/../phoneinfoga.py -n "+86 591 2284 8571" -h

    python3 $scriptDir/../phoneinfoga.py -n "+86 591 2284 8571" -s any --no-ansi

    python3 $scriptDir/../phoneinfoga.py -i $scriptDir/input.txt -o $scriptDir/output_from_input.txt -s any

    python3 $scriptDir/../phoneinfoga.py -n "+86 591 2284 8571" -s numverify -o $scriptDir/output_single.txt
    
    echo "Test script executed."
    
    python4 
} || {
       # command which should be run instead of the above failing command
       echo "Command exited with an error."
       exit 1
}