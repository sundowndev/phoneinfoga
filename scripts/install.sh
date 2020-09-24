#!/bin/bash
os=`echo "$(uname -s)_$(uname -m)"`
if [ $os == "Linux_x86_64" ] || [ $os == "Linux_armv6" ] || [ $os == "Windows_x86_64" ] || [ $os == "Darwin_x86_64" ] || [ $os == "Linux_arm64" ] || [ $os == "Windows_i386" ] || [ $os == "Linux_i386" ]]
then
        echo Installing Phoneinfoga.........
        wget "https://github.com/sundowndev/phoneinfoga/releases/download/v2.3.2/phoneinfoga_$(uname -s)_$(uname -m).tar.gz"
        tar xfv "phoneinfoga_$(uname -s)_$(uname -m).tar.gz"
        mv ./phoneinfoga /usr/bin/phoneinfoga
        echo Installation completed  Successfully.
else
        echo Your OS/arch is not Supported.
fi
exit 
