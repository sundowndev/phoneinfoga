#!/usr/bin/env python

import requests
import sys
import hashlib
import json
import argparse
from bs4 import BeautifulSoup
import re

__version__ = '0.3-dev'

print "\n \033[92m"
print "    ___ _                       _____        __                   "
print "   / _ \ |__   ___  _ __   ___  \_   \_ __  / _| ___   __ _  __ _ "
print "  / /_)/ '_ \ / _ \| '_ \ / _ \  / /\/ '_ \| |_ / _ \ / _` |/ _` |"
print " / ___/| | | | (_) | | | |  __/\/ /_ | | | |  _| (_) | (_| | (_| |"
print " \/    |_| |_|\___/|_| |_|\___\____/ |_| |_|_|  \___/ \__, |\__,_|"
print "                                                      |___/       "
print " PhoneInfoga Ver. %s                                              " % __version__
print " Coded by Sundowndev                                              "
print "\033[94m\n"

parser = argparse.ArgumentParser(description=
    "Advanced information gathering tool for phone numbers (https://github.com/sundowndev/PhoneInfoga) version %s" % __version__,
                                 usage='%(prog)s -n <number> [options]')

parser.add_argument('-n', '--number', metavar='number', type=str,
                    help='The phone number to scan (E164 and International format)')

parser.add_argument('-i', '--input', metavar="input_file", type=file,
                    help='Phone number list to scan (one per line)')

parser.add_argument('-o', '--output', metavar="output_file", type=file,
                    help='Output to save scan results')

parser.add_argument('-s', '--scanner', metavar="scanner", default="all", type=str,
                    help='The scanner to use')

parser.add_argument('--osint', action='store_true',
                    help='Use OSINT reconnaissance')

parser.add_argument('-u', '--update', action='store_true',
                    help='Update the tool & databases')

args = parser.parse_args()

# If any param is passed, execute help command
if not len(sys.argv) > 1:
    parser.print_help()

if args.update:
    print 'update'
    sys.exit()

scanners = ['any', 'all', 'numverify', 'ovh', 'whosenumber']

def parseInput(file):
    print 'parse'

def saveToOutput():
    print 'save'

def isNumberValid(PhoneNumber):
    if len(PhoneNumber) < 9 and len(PhoneNumber) > 13:
        return False
    elif not re.match("^[+]*[(]{0,1}[0-9]{1,4}[)]{0,1}[-\s\./0-9]*$", PhoneNumber):
        return False
    else:
        return True

def formatNumber(number):
    PhoneNumber = number.replace("+", "").replace("\n", "").replace(" ", "")
    return PhoneNumber

def searchCountryCode(number):
    #parse code

    #check in json
    print '\033[1;32m[+] Country found : France (FR)'

    #check for area code
    print '\033[1;32m[+] Areas found (approximate) : Bordeaux, Limoges'

    #check for carrier
    #print '\033[1;32m[+] Carrier found:  France Sfr Mobile'
    print '\033[93m[i] This is most likely a landline, or a fixed VoIP.'

def numverifyScan(PhoneNumber):
    if not args.scanner == 'numverify' and not args.scanner == 'all':
        return -1

    print '\033[93m[i] Running Numverify scan...'

    requestSecret = ''
    resp = requests.get('https://numverify.com/')
    soup = BeautifulSoup(resp.text, "html5lib")
    for tag in soup.find_all("input", type="hidden"):
        if tag['name'] == "scl_request_secret":
            requestSecret = tag['value']
            break;

    apiKey = hashlib.md5()
    apiKey.update(PhoneNumber + requestSecret)
    apiKey = apiKey.hexdigest()

    response = requests.get("https://numverify.com/php_helper_scripts/phone_api.php?secret_key=" + apiKey + "&number=" + PhoneNumber)

    if response.content == "Unauthorized" or response.status_code != 200:
        print("[i] An error occured while calling the API (bad request or wrong api key).")
        sys.exit()

    data = json.loads(response.content)

    if data["valid"] == False:
        print("\033[91m[!] Error: Please specify a valid phone number. Example: +6464806649\033[94m")
        sys.exit()

    print "\033[1;32mNumber: (" + data["country_prefix"] + ") " + data["local_format"]
    print("Country: %s (%s)") % (data["country_name"],data["country_code"])
    print("Location: %s") % data["location"]
    print("Carrier: %s") % data["carrier"]
    print("Line type: %s \033[94m") % data["line_type"]

def ovhScan(number):
    if not args.scanner == 'ovh' and not args.scanner == 'all':
        return -1

    print '\033[93m[i] Running OVH scan...'
    print '(!) OVH API credentials missing. Skipping.\033[94m'

def whosenumberScan(number):
    if not args.scanner == 'whosenumber' and not args.scanner == 'all':
        return -1

    print '\033[93m[i] Running Whosenumber scan...\033[94m'

def scanNumber(number):
    PhoneNumber = formatNumber(number)

    print "\033[93m[!] ---- Fetching informations for " + PhoneNumber + " ---- [!]"

    if not isNumberValid(PhoneNumber):
        print("\033[91mError: number " + number + " is not valid. Skipping.")
        sys.exit()

    #check dial code
    searchCountryCode(PhoneNumber)
    #check area code by country
    #if found in area codes -> landline

    numverifyScan(PhoneNumber)
    ovhScan(PhoneNumber)
    whosenumberScan(PhoneNumber)

# Verify scanner
if not args.scanner in scanners:
    print("\033[91mError: scanner doesn't exists.")
    sys.exit()

if args.number:
    scanNumber(args.number)
elif args.input:
    for line in args.input.readlines():
        scanNumber(line)

if args.output:
    args.output.write("Hello World")

    args.output.close()