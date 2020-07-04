#!/bin/sh

echo "Starting installation of PhoneInfoga for OS: $(uname -s) $(uname -m)"

# TODO: check OS to use systeminfo for Windows
# TODO: check if user have requirements: wget, curl, tar

ARCH="$(uname -s)_$(uname -m)"
LATEST_RELEASE=$(curl -L ...)

echo "Found latest version: $LATEST_RELEASE"

# Download the archive
wget "https://github.com/sundowndev/phoneinfoga/releases/download/$LATEST_RELEASE/phoneinfoga_$ARCH.tar.gz"

# TODO: add a check if arch is actually supported (404 error)
# ----
# echo "Error: your OS and arch does not appear to be supported."
# echo "You may open an issue on GitHub with details about your machine."
# exit 1

# Extract the binary
tar xfv "phoneinfoga_$(uname -s)_$(uname -m).tar.gz"

echo """
Installation finished.

You may now check the binary is working :

$ ./phoneinfoga version
"""

exit 0
