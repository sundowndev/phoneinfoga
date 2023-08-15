#! /bin/bash

# Bash installation script for UNIX systems only
# Use it : curl -sSL https://raw.githubusercontent.com/sundowndev/phoneinfoga/master/support/scripts/install | bash

OS="$(uname -s)_$(uname -m)"
PHONEINFOGA_VERSION=
AUTOMATIC=
SKIP_CHECKSUM=1
CURL_INSTALLED=false
WGET_INSTALLED=false


choose_wget_or_curl() {

    which curl > /dev/null
    [ $? -eq 1 ] || CURL_INSTALLED=true

    if $CURL_INSTALLED; then
        PHONEINFOGA_VERSION=$(curl -s https://api.github.com/repos/sundowndev/phoneinfoga/releases/latest | grep tag_name | cut -d '"' -f 4)
        return
    fi

    which wget > /dev/null
    [ $? -eq 1 ] || WGET_INSTALLED=true

    if $WGET_INSTALLED; then
        PHONEINFOGA_VERSION=$(wget -q --output-document - https://api.github.com/repos/sundowndev/phoneinfoga/releases/latest | grep tag_name | cut -d '"' -f 4)
        return
    fi

    echo "Error: You need to have either curl or wget installed to be able to use this script"
    exit 1

}


usage() {

    echo "PhoneInfoga Installer for version $PHONEINFOGA_VERSION"
    echo
    echo "DESCRIPTION: An installer script for downloading the latest release of PhoneInfoga. Without any arguments, $0 will detect your operating system and attempt to download the corresponding version of PhoneInfoga. Please submit an issue on GitHub if you encounter issues with this script." | fold -s
    echo
    echo "USAGE: $0 [flag...] (-h|-m|-s)"
    echo
    echo "      -h | --help             Print this message and exit"
    echo
    echo "      -m | --manual           Manually select version to download."
    echo "                              Useful when $0 is unable to detect"
    echo "                              your OS automatically"
    echo
    echo "      -s | --skip-checksum    [Not Reccomended] Skip checksum validation"
    echo "                              Only use this option if $0"
    echo "                              is failing due to missing software."
    echo
}


validate_OS_type() {

    if [[ -z "${OS}" ]]; then
        echo "ERROR: Unable to determine your system type."
        exit 1
    fi
}


validate_supported_OS() {

    SUPPORTED_OS_TYPES=(
        "Darwin_arm64"
        "Darwin_x86_64"
        "Linux_arm64"
        "Linux_armv6"
        "Linux_armv7"
        "Linux_i386"
        "Linux_x86_64"
        "Windows_arm64"
        "Windows_armv6"
        "Windows_armv7"
        "Windows_i386"
        "Windows_x86_64"
    )

    for OS_TYPE in "${SUPPORTED_OS_TYPES[@]}"; do
        if [[ "$OS_TYPE" == "$OS" ]]; then
            AUTOMATIC=0
            return
        fi
    done

    echo "Error: $OS is not supported, installation will attempt to proceed manually."
    echo
    echo "Please check the releases page for a list of supported systems."
    echo "https://github.com/sundowndev/phoneinfoga/releases"
    echo
    echo "Read more at https://sundowndev.github.io/phoneinfoga/install/"
    echo
    AUTOMATIC=1

}


validate_checksum() {

    echo "Validating checksum ..."
    if $CURL_INSTALLED; then
        curl --progress-bar -LOC - "https://github.com/sundowndev/phoneinfoga/releases/download/$PHONEINFOGA_VERSION/phoneinfoga_checksums.txt"
    elif $WGET_INSTALLED; then
        wget --quiet --show-progress "https://github.com/sundowndev/phoneinfoga/releases/download/$PHONEINFOGA_VERSION/phoneinfoga_checksums.txt"
    fi

    local SHA256SUM_INSTALLED=false
    which sha256sum > /dev/null
    [ $? -eq 1 ] || SHA256SUM_INSTALLED=true

    if $SHA256SUM_INSTALLED; then
        sha256sum --ignore-missing -c phoneinfoga_checksums.txt
        [ $? -eq 0 ] || exit 1
    elif ! $SHA256SUM_INSTALLED; then
        shasum --ignore-missing -c phoneinfoga_checksums.txt
        [ $? -eq 0 ] || exit 1
    fi

    rm phoneinfoga_checksums.txt
}


download_latest_release() {

    echo "Downloading PhoneInfoga Version $PHONEINFOGA_VERSION for $OS"

    if $CURL_INSTALLED; then
        curl --progress-bar -LOC - "https://github.com/sundowndev/phoneinfoga/releases/download/$PHONEINFOGA_VERSION/phoneinfoga_$OS.tar.gz"

    elif $WGET_INSTALLED; then
        wget --quiet --show-progress "https://github.com/sundowndev/phoneinfoga/releases/download/$PHONEINFOGA_VERSION/phoneinfoga_$OS.tar.gz"

    fi
}


manual_download_latest_release() {

    if $CURL_INSTALLED; then
        ASSETS=$(curl -s https://api.github.com/repos/sundowndev/phoneinfoga/releases/latest | jq --raw-output '.assets[] | .name' | grep --color=never .tar.gz)
    elif $WGET_INSTALLED; then
        ASSETS=$(wget -q --output-document - https://api.github.com/repos/sundowndev/phoneinfoga/releases/latest | jq --raw-output '.assets[] | .name' | grep --color=never .tar.gz)
    fi

    echo "Please select a version of PhoneInfoga to download."
    PS3="Enter a number: "
    select OPT in $ASSETS; do
        case $OPT in

            '')
                echo "Error: Invalid option, please enter a number from the list."
                ;;

            *)
                echo "Downloading $OPT version $PHONEINFOGA_VERSION"
                echo "https://github.com/sundowndev/phoneinfoga/releases/download/$PHONEINFOGA_VERSION/$OPT"

                if $CURL_INSTALLED; then
                    curl --progress-bar -LOC - "https://github.com/sundowndev/phoneinfoga/releases/download/$PHONEINFOGA_VERSION/$OPT"

                elif $WGET_INSTALLED; then
                    wget --quiet --show-progress "https://github.com/sundowndev/phoneinfoga/releases/download/$PHONEINFOGA_VERSION/$OPT"

                fi

                break
                ;;
        esac
    done
}


unpack_and_cleanup() {

    local FILE
    FILE=$(find . -name "phoneinfoga_*.tar.gz")

    tar -xzf $FILE
    [ $? -eq 0 ] || exit 1

    echo "Cleaning up ..."
    rm $FILE
}


main() {

    validate_OS_type
    validate_supported_OS
    choose_wget_or_curl

    if [[ -z "$1" ]]; then
        echo
    else
        while [[ "$1" == -* ]]; do
            case "$1" in

                -h | --help)
                    usage
                    exit 0
                    ;;

                -m | --manual)
                    AUTOMATIC=1
                    ;;

                -s | --skip-checksum)
                    SKIP_CHECKSUM=0
                    ;;

                *)
                    echo
                    echo "ERROR: Unrecognized command '$1', run '$0 --help' for help."
                    exit 1
                    ;;

                esac
            shift
        done
    fi

    if [[ $AUTOMATIC == 0 ]]; then
        download_latest_release
    elif [[ $AUTOMATIC == 1 ]]; then
        manual_download_latest_release
    else
        echo "something bad happened"
        exit 1
    fi

    if [[ $SKIP_CHECKSUM == 1 ]]; then
        validate_checksum
    elif [[ $SKIP_CHECKSUM == 0 ]]; then
        echo
        echo "WARNING: Skipping checksum validation. Please be sure to verify your download manually. You can find the checksum for your software version here:" | fold -s
        echo
        echo "https://github.com/sundowndev/phoneinfoga/releases/download/$PHONEINFOGA_VERSION/phoneinfoga_checksums.txt"
    fi

    unpack_and_cleanup

    echo
    echo "Installation completed successfully."
    echo "To check the version installed: ./phoneinfoga version"
    echo
    echo "Add it to your path by running"
    echo "'sudo mv ./phoneinfoga /usr/local/bin/phoneinfoga'"

    exit 0
}


main "${@}"
