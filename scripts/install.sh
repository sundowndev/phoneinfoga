#!/bin/bash

os="$(uname -s)_$(uname -m)"

if [ $os == "Linux_x86_64" ] || [ $os == "Linux_armv6" ] || [ $os == "Darwin_x86_64" ] || [ $os == "Linux_arm64" ] || [ $os == "Linux_i386" ]]
then
	echo "Installing Phoneinfoga........."
	phoneinfoga_version=$(curl -s https://api.github.com/repos/sundowndev/PhoneInfoga/releases/latest | grep tag_name | cut -d '"' -f 4)
	wget "https://github.com/sundowndev/phoneinfoga/releases/download/$phoneinfoga_version/phoneinfoga_$os.tar.gz"
	tar xfv "phoneinfoga_$os.tar.gz"
	mv ./phoneinfoga /usr/bin/phoneinfoga
	rm -rf phoneinfoga_$os.tar.gz
	echo "Installation completed Successfully."

else
	echo "Your OS/arch is not Supported."
fi

exit 0
