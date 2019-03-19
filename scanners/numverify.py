#!/usr/bin/env python3
# -*- coding:utf-8 -*-
#
# @name   : PhoneInfoga - Phone numbers OSINT tool
# @url    : https://github.com/sundowndev
# @author : Raphael Cerveaux (sundowndev)

from bs4 import BeautifulSoup
import hashlib
import json
from lib.output import *
from lib.request import send


def scan(number):
    if not args.scanner == 'numverify' and not args.scanner == 'all':
        return -1

    test('Running Numverify.com scan...')

    try:
        requestSecret = ''
        res = send('GET', 'https://numverify.com/')
        soup = BeautifulSoup(res.text, "html5lib")
    except Exception as e:
        error('Numverify.com is not available')
        error(e)
        return -1

    for tag in soup.find_all("input", type="hidden"):
        if tag['name'] == "scl_request_secret":
            requestSecret = tag['value']
            break

    apiKey = hashlib.md5((number + requestSecret).encode('utf-8')).hexdigest()

    headers = {
        'Host': 'numverify.com',
        'User-Agent': 'Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:64.0) Gecko/20100101 Firefox/64.0',
        'Accept': 'application/json, text/javascript, */*; q=0.01',
        'Accept-Language': 'fr,fr-FR;q=0.8,en-US;q=0.5,en;q=0.3',
        'Accept-Encoding': 'gzip, deflate, br',
        'Referer': 'https://numverify.com/',
        'X-Requested-With': 'XMLHttpRequest',
        'DNT': '1',
        'Connection': 'keep-alive',
        'Pragma': 'no-cache',
        'Cache-Control': 'no-cache'
    }

    try:
        res = send(
            "GET", "https://numverify.com/php_helper_scripts/phone_api.php?secret_key={}&number={}".format(apiKey, number), headers=headers)

        data = json.loads(res.content.decode('utf-8'))
    except Exception as e:
        error('Numverify.com is not available')
        return -1

    if res.content == "Unauthorized" or res.status_code != 200:
        error(("An error occured while calling the API (bad request or wrong api key)."))
        return -1

    if 'error' in data:
        error('Numverify.com is not available: ' + data['error'])
        return -1

    if data['valid'] == False:
        error(("Error: Please specify a valid phone number. Example: +6464806649"))
        sys.exit()

    InternationalNumber = '({}){}'.format(
        data["country_prefix"], data["local_format"])

    plus(("Number: ({}) {}").format(
        data["country_prefix"], data["local_format"]))
    plus(("Country: {} ({})").format(
        data["country_name"], data["country_code"]))
    plus(("Location: {}").format(data["location"]))
    plus(("Carrier: {}").format(data["carrier"]))
    plus(("Line type: {}").format(data["line_type"]))

    if data["line_type"] == 'landline':
        warn(("This is most likely a landline, but it can still be a fixed VoIP number."))
    elif data["line_type"] == 'mobile':
        warn(("This is most likely a mobile number, but it can still be a VoIP number."))
