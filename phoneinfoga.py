#!/usr/bin/env python

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

import requests
import sys
import hashlib
import json
import argparse
from bs4 import BeautifulSoup
import re
import phonenumbers
from phonenumbers import carrier
from phonenumbers import geocoder
from phonenumbers import timezone

parser = argparse.ArgumentParser(description=
    "Advanced information gathering tool for phone numbers (https://github.com/sundowndev/PhoneInfoga) version %s" % __version__,
                                 usage='%(prog)s -n <number> [options]')

parser.add_argument('-n', '--number', metavar='number', type=str,
                    help='The phone number to scan (E164 or international format)')

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

scanners = ['any', 'all', 'numverify', 'ovh', 'whosenumber', 'freecarrier', '411']

code_info = '\033[97m[*] '
code_warning = '\033[93m(!) '
code_result = '\033[1;32m[+] '
code_error = '\033[91m[!] '

def saveToOutput(output):
    print 'save'

def localScan(number):
    print code_info + 'Running local scan...'

    PhoneNumber = dict();

    FormattedPhoneNumber = number.replace("\n", "").replace("-", "").replace(" ", "")

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

        print code_result + 'Local format: (0)' + PhoneNumber['number']
        print code_result + 'Country code: ' + PhoneNumber['countryCode']
        print code_result + 'Location: %s' % geocoder.description_for_number(PhoneNumberObject, "en")
        print code_result + 'Carrier:  %s' % carrier.name_for_number(PhoneNumberObject, 'en')
        print code_result + 'Area:  %s' % geocoder.description_for_number(PhoneNumberObject, 'en')
        #print '\033[1;32m[+] Timezone:  %s, %s' % (timezone.time_zones_for_number(PhoneNumberObject)[0],timezone.time_zones_for_number(PhoneNumberObject)[1])
        #print code_info + 'This is most likely a landline, or a fixed VoIP.'

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
        sys.exit()

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

def ovhScan(country, number):
    if not args.scanner == 'ovh' and not args.scanner == 'all':
        return -1

    print code_info + 'Running OVH scan...'

    querystring = {"country":country}

    headers = {
        'accept': "application/json",
        'cache-control': "no-cache"
    }

    response = requests.request("GET", "https://api.ovh.com/1.0/telephony/number/detailedZones", data="", headers=headers, params=querystring)

    data = json.loads(response.content)

def whosenumberScan(countryCode, number):
    if not args.scanner == 'whosenumber' and not args.scanner == 'all':
        return -1

    print code_info + 'Running Whosenumber scan...'
    print 'https://whosenumber.info/' + countryCode + number

def repScan(countryCode, number):
    if not args.scanner == '411' and not args.scanner == 'all':
        return -1

    print code_info + 'Running 411.com scan...'
    print 'https://www.411.com/phone/%s-%s' % (countryCode,number)

def freecarrierlookupScan(countryCode, number):
    if not args.scanner == 'freecarrier' and not args.scanner == 'all':
        return -1

    print code_info + 'Running freecarrierlookup.com scan...'

    payload = "phonenum=%s&cc=%s" % (number,countryCode)
    headers = {
        'host': "freecarrierlookup.com",
        'connection': "keep-alive",
        'content-length': "48",
        'accept': "application/json, text/javascript, */*; q=0.01",
        'origin': "https://freecarrierlookup.com",
        'x-requested-with': "XMLHttpRequest",
        'user-agent': "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.99 Safari/537.36",
        'content-type': "application/x-www-form-urlencoded",
        'referer': "https://freecarrierlookup.com/",
        'accept-encoding': "gzip, deflate, br",
        'accept-language': "en-US,en;q=0.9,fr;q=0.8,la;q=0.7,es;q=0.6,zh-CN;q=0.5,zh;q=0.4",
        'cookie': "PHPSESSID=cdifm9u3ch2mqscdnj2pjqjfuq",
        'cache-control': "no-cache",
        'postman-token': "c81a7bb0-f338-c2e5-5f32-e1b94726cce5"
    }

    response = requests.request("POST", "https://freecarrierlookup.com/getcarrier.php", data=payload, headers=headers)

    print response.content

    data = json.loads(response.content)

    if not data["status"] == "success":
        print code_error + '0 result found.'
        return -1

    soup = BeautifulSoup(response.content, "html5lib")
    tags = soup.find_all("p")

    print code_result + 'Phone Number: ' + tags[0].string.replace('<\/p>\\n        <\/div>\\n', '')
    print code_result + 'Carrier: ' + tags[1].string.replace('<\/p>\\n        <\/div>\\n', '')
    print code_result + 'Is Wireless:'
    print code_result + 'SMS Gateway Address: '
    print code_result + 'MMS Gateway Address: '

def scanNumber(number):
    print "\033[1m\033[93m[!] ---- Fetching informations for %s ---- [!]" % number

    print code_info + 'Parsing informations...'

    PhoneNumber = localScan(number)

    if not PhoneNumber:
        print(code_error + "Error: number " + number + " is not valid. Skipping.")
        sys.exit()

    numverifyScan(PhoneNumber['full'])
    ovhScan('fr', PhoneNumber['full']) # TODO: replace 1st parameter to be dynamic
    #freecarrierlookupScan(PhoneNumber['countryCode'], PhoneNumber['number'])
    #whosenumberScan(PhoneNumber['countryCode'], PhoneNumber['number'])
    #repScan(PhoneNumber['countryCode'], PhoneNumber['number'])

    print '\n'

# Verify scanner option
if not args.scanner in scanners:
    print(code_error + "Error: scanner doesn't exists.")
    sys.exit()

if args.number:
    scanNumber(args.number)
elif args.input:
    for line in args.input.readlines():
        scanNumber(line)

if args.output:
    args.output.write("Hello World")

    args.output.close()
