#!/usr/bin/env python3
# -*- coding:utf-8 -*- 
#
# @name   : PhoneInfoga - Phone numbers OSINT tool
# @url    : https://github.com/sundowndev
# @author : Raphael Cerveaux (sundowndev)

# dependencies
import sys
import signal
import argparse
import requests
# lib
from lib.banner import banner,__version__
from lib.output import *
from lib.format import *
from lib.vars import *
# scanners
from scanners import numverify
from scanners import localscan
from scanners import ovh
from scanners.osint import osintScan

parser = argparse.ArgumentParser(description="Advanced information gathering tool for phone numbers (https://github.com/sundowndev/PhoneInfoga) version {}".format(__version__),
                                 usage='%(prog)s -n <number> [options]')

parser.add_argument('-n', '--number', metavar='number', type=str,
                    help='The phone number to scan (E164 or international format)')

parser.add_argument('-i', '--input', metavar="input_file", type=argparse.FileType('r'),
                    help='Phone number list to scan (one per line)')

parser.add_argument('-o', '--output', metavar="output_file", type=argparse.FileType('w'),
                    help='Output to save scan results')

parser.add_argument('-s', '--scanner', metavar="scanner", default="all", type=str,
                    help='The scanner to use')

parser.add_argument('--osint', action='store_true',
                    help='Use OSINT reconnaissance')

parser.add_argument('--no-ansi', action='store_true',
                    help='Disable colored output')

parser.add_argument('-v', '--version', action='store_true',
                    help='Show tool version')

args = parser.parse_args()

def scanNumber(InputNumber):
    title("[!] ---- Fetching informations for {} ---- [!]".format(formatNumber(InputNumber)))

    number = localscan.scan(InputNumber)

    print(number)

    if not number:
        throw(("Error: an error occured parsing {}. Skipping.".format(formatNumber(InputNumber))))

    # numverify.scan(number['default'])
    # ovh.scan(number['local'], number['countryIsoCode'])
    osintScan(number)

    info("Scan finished.")

    if not args.no_ansi and not args.output:
        print('\n' + Style.RESET_ALL)
    else:
        print('\n')

def main():
    scanners = ['any', 'all', 'numverify', 'ovh']

    banner()

    if sys.version_info[0] < 3:
        print(
            "\033[1m\033[93m(!) Please run the tool using Python 3" + Style.RESET_ALL)
        sys.exit()

    # If any param is passed, execute help command
    if not len(sys.argv) > 1:
        parser.print_help()
        sys.exit()
    elif args.version:
        print("Version {}".format(__version__))
        sys.exit()

    requests.packages.urllib3.disable_warnings()
    requests.packages.urllib3.util.ssl_.DEFAULT_CIPHERS += 'HIGH:!DH:!aNULL'
    try:
        requests.packages.urllib3.contrib.pyopenssl.DEFAULT_SSL_CIPHER_LIST += 'HIGH:!DH:!aNULL'
    except AttributeError:
        # no pyopenssl support used / needed / available
        pass

    if args.output:
        if args.osint:
            print(
                '[!] OSINT scanner is not available using output option (sorry).')
            sys.exit()

        sys.stdout = args.output
        banner()  # Output banner again in the file

    # Verify scanner option
    if not args.scanner in scanners:
        print(("Error: scanner doesn't exists."))
        sys.exit()

    if args.number:
        scanNumber(args.number)
    elif args.input:
        for line in args.input.readlines():
            scanNumber(line)
    else:
        parser.print_help()
        sys.exit()

    if args.output:
        args.output.close()


def signal_handler(signal, frame):
    print('\n[-] You pressed Ctrl+C! Exiting.')
    sys.exit()


if __name__ == '__main__':
    signal.signal(signal.SIGINT, signal_handler)
    main()
