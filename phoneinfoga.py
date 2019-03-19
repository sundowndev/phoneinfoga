#!/usr/bin/env python3
# -*- coding:utf-8 -*-
#
# @name   : PhoneInfoga - Phone numbers OSINT tool
# @url    : https://github.com/sundowndev
# @author : Raphael Cerveaux (sundowndev)

# dependencies
import sys
import signal
# lib
from lib.args import args,parser
from lib.banner import banner, __version__
from lib.output import *
from lib.format import *
from lib.logger import Logger
# scanners
from scanners import numverify
from scanners import localscan
from scanners import ovh
from scanners.footprints import osintScan
from scanners import recon


def scanNumber(InputNumber):
    title("[!] ---- Fetching informations for {} ---- [!]".format(formatNumber(InputNumber)))

    number = localscan.scan(InputNumber)

    if not number:
        throw(("Error: an error occured parsing {}. Skipping.".format(
            formatNumber(InputNumber))))

    numverify.scan(number['default'])
    ovh.scan(number['local'], number['countryIsoCode'])
    recon.scan(number)
    osintScan(number)

    info("Scan finished.\n")


def main():
    scanners = ['any', 'all', 'numverify', 'ovh', 'footprints']

    banner()

    # Ensure the usage of Python3
    if sys.version_info[0] < 3:
        print(
            "(!) Please run the tool using Python 3")
        sys.exit()

    # If any param is passed, execute help command
    if not len(sys.argv) > 1:
        parser.print_help()
        sys.exit()
    elif args.version:
        print("Version {}".format(__version__))
        sys.exit()

    if args.output:
        sys.stdout = Logger()

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
