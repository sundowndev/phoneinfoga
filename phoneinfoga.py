#!/usr/bin/env python

__version__ = '0.6-dev'

def banner():
    print "    ___ _                       _____        __                   "
    print "   / _ \ |__   ___  _ __   ___  \_   \_ __  / _| ___   __ _  __ _ "
    print "  / /_)/ '_ \ / _ \| '_ \ / _ \  / /\/ '_ \| |_ / _ \ / _` |/ _` |"
    print " / ___/| | | | (_) | | | |  __/\/ /_ | | | |  _| (_) | (_| | (_| |"
    print " \/    |_| |_|\___/|_| |_|\___\____/ |_| |_|_|  \___/ \__, |\__,_|"
    print "                                                      |___/       "
    print " PhoneInfoga Ver. %s                                              " % __version__
    print " Coded by Sundowndev                                              "

print "\n \033[92m"
banner()
print "\033[94m\n"

import sys
import argparse

parser = argparse.ArgumentParser(description=
    "Advanced information gathering tool for phone numbers (https://github.com/sundowndev/PhoneInfoga) version %s" % __version__,
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

parser.add_argument('-u', '--update', action='store_true',
                    help='Update the tool & databases')

args = parser.parse_args()

# If any param is passed, execute help command
if not len(sys.argv) > 1:
    parser.print_help()
    sys.exit();

if args.update:
    print 'update'
    sys.exit()

import requests
import hashlib
import json
from bs4 import BeautifulSoup
import re
import phonenumbers
from phonenumbers import carrier
from phonenumbers import geocoder
from phonenumbers import timezone

scanners = ['any', 'all', 'numverify', 'ovh']

def formatNumber(number):
    return re.sub("(?:\+)?(?:[^[0-9]*)", "", number)

def localScan(number):
    print code_info + 'Running local scan...'

    PhoneNumber = dict();

    FormattedPhoneNumber = "+" + formatNumber(number)

    try:
        PhoneNumberObject = phonenumbers.parse(FormattedPhoneNumber, None)
    except:
        return False
    else:
        if not phonenumbers.is_valid_number(PhoneNumberObject):
            return False

        PhoneNumber['full'] = phonenumbers.format_number(PhoneNumberObject, phonenumbers.PhoneNumberFormat.E164).replace('+', '')
        PhoneNumber['countryCode'] = phonenumbers.format_number(PhoneNumberObject, phonenumbers.PhoneNumberFormat.INTERNATIONAL).split(' ')[0]
        PhoneNumber['number'] = phonenumbers.format_number(PhoneNumberObject, phonenumbers.PhoneNumberFormat.E164).replace(PhoneNumber['countryCode'], '')

        print code_result + 'Local format: (0)%s' % PhoneNumber['number']
        print code_result + 'Country code: %s' % PhoneNumber['countryCode']
        print code_result + 'Location: %s' % geocoder.description_for_number(PhoneNumberObject, "en")
        print code_result + 'Carrier: %s' % carrier.name_for_number(PhoneNumberObject, 'en')
        print code_result + 'Area: %s' % geocoder.description_for_number(PhoneNumberObject, 'en')
        for timezoneResult in timezone.time_zones_for_number(PhoneNumberObject):
            print code_result + 'Timezone: %s' % (timezoneResult)

        if phonenumbers.is_possible_number(PhoneNumberObject):
            print code_info + 'The number is valid and possible.'
        else:
            print code_warning + 'The number is valid but might not be possible.'

        return PhoneNumber

def numverifyScan(PhoneNumber):
    if not args.scanner == 'numverify' and not args.scanner == 'all':
        return -1

    print code_info + 'Running Numverify.com scan...'

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

    headers = {
        'host': "numverify.com",
        'connection': "keep-alive",
        'content-length': "49",
        'accept': "application/json, text/javascript, */*; q=0.01",
        'origin': "https://numverify.com",
        'x-requested-with': "XMLHttpRequest",
        'user-agent': "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.99 Safari/537.36",
        'content-type': "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW",
        'referer': "https://numverify.com/",
        'accept-encoding': "gzip, deflate, br",
        'accept-language': "en-US,en;q=0.9,fr;q=0.8,la;q=0.7,es;q=0.6,zh-CN;q=0.5,zh;q=0.4",
        'cache-control': "no-cache"
    }

    response = requests.request("GET", "https://numverify.com/php_helper_scripts/phone_api.php?secret_key=" + apiKey + "&number=" + PhoneNumber, data="", headers=headers)

    if response.content == "Unauthorized" or response.status_code != 200:
        print(code_error + "An error occured while calling the API (bad request or wrong api key).")
        return -1

    data = json.loads(response.content)

    if data["valid"] == False:
        print(code_error + "Error: Please specify a valid phone number. Example: +6464806649")
        sys.exit()

    InternationalNumber = '('+data["country_prefix"]+')' + data["local_format"]

    print(code_result + "Number: (%s) %s") % (data["country_prefix"],data["local_format"])
    print(code_result + "Country: %s (%s)") % (data["country_name"],data["country_code"])
    print(code_result + "Location: %s") % data["location"]
    print(code_result + "Carrier: %s") % data["carrier"]
    print(code_result + "Line type: %s") % data["line_type"]

    if data["line_type"] == 'landline':
        print(code_warning + "This is most likely a land line, but it can still be a fixed VoIP.")
    elif data["line_type"] == 'mobile':
        print(code_warning + "This is most likely a mobile, but it can still be a VoIP.")

def ovhScan(countryCode, number):
    if not args.scanner == 'ovh' and not args.scanner == 'all':
        return -1

    print code_info + 'Running OVH scan...'

    #cc

    querystring = {"country":"fr"}

    headers = {
        'accept': "application/json",
        'cache-control': "no-cache"
    }

    response = requests.request("GET", "https://api.ovh.com/1.0/telephony/number/detailedZones", data="", headers=headers, params=querystring)

    data = json.loads(response.content)

    if isinstance(data, list):
        askedNumber = "0" + number.replace(number[-4:], 'xxxx')

        for voip_number in data:
            if voip_number['number'] == askedNumber:
                print(code_info + "1 result found in OVH database")
                print(code_result + "Number range: " + voip_number['number'])
                print(code_result + "City: " + voip_number['city'])
                print(code_result + "Zip code: " + voip_number['zipCode'] if voip_number['zipCode'] is not None else '')

def osintScan(countryCode, number):
    if not args.osint:
        return -1

    from googlesearch import search

    print code_info + 'Running OSINT reconnaissance...'
    # OSINT recon

    # social profiles: facebook, twitter, linkedin, instagram
    # websites
    # emails

    # Whitepages
    print(code_info + "Searching for owner on 411.com...")
    #https://www.411.com/phone/33-6-79-36-82-33

    # Reputation
    print(code_info + "Searching for reputation page on whosenumber.info...")
    for result in search('site:whosenumber.info intext:"%s" intitle:"who called"' % number, stop=1):
        if result:
            print(code_result + "Found 1 result on whosenumber.info.")
            print(code_info + "This usually means you are not the first to search about this number. Check the URL for eventual comments.")
            print(code_result + "URL: " + result)

    # VoIP providers
    print(code_info + "Searching for results on hs3x.com...")
    for result in search('site:"hs3x.com" intext:"+%s"' % number, stop=1):
        if result:
            print(code_result + "Found 1 result on hs3x.com.")
            print(code_info + "This number seems to be a VoIP number from hs3x.")
            print(code_result + "URL: " + result)

def scanNumber(number):
    print code_title + "[!] ---- Fetching informations for %s ---- [!]" % formatNumber(number)

    PhoneNumber = localScan(number)

    if not PhoneNumber:
        print(code_error + "Error: number " + formatNumber(number) + " is not valid. Skipping.")
        sys.exit()

    numverifyScan(PhoneNumber['full'])
    ovhScan(PhoneNumber['countryCode'], PhoneNumber['number']) # TODO: replace 1st parameter to be dynamic
    osintScan(PhoneNumber['countryCode'], PhoneNumber['full'])

    print '\n'

# Verify scanner option
if not args.scanner in scanners:
    print(code_error + "Error: scanner doesn't exists.")
    sys.exit()

if args.output:
    code_info = '[*] '
    code_warning = '(!) '
    code_result = '[+] '
    code_error = '[!] '
    code_title = ''

    sys.stdout = args.output
else:
    code_info = '\033[97m[*] '
    code_warning = '\033[93m(!) '
    code_result = '\033[1;32m[+] '
    code_error = '\033[91m[!] '
    code_title = '\033[1m\033[93m'

if args.number:
    scanNumber(args.number)
elif args.input:
    for line in args.input.readlines():
        scanNumber(line)

if args.output:
    args.output.close()
