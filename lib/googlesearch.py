#!/usr/bin/env python3
# -*- coding:utf-8 -*-
#
# @name   : PhoneInfoga - Phone numbers OSINT tool
# @url    : https://github.com/sundowndev
# @author : Raphael Cerveaux (sundowndev)

import re
import json
from urllib.parse import urlencode
from bs4 import BeautifulSoup
from lib.output import *
from lib.request import send
from config import *

googleAbuseToken = ''
customFormatting = ''


def search(req, stop):
    global googleAbuseToken

    if google_api_key and google_cx_id:
        return searchApi(req, stop)

    headers = {
        'Accept': 'text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8',
        'Accept-Language': 'en-us,en;q=0.5',
        'Accept-Encoding': 'gzip,deflate',
        'Accept-Charset': 'ISO-8859-1,utf-8;q=0.7,*;q=0.7',
        'Keep-Alive': '115',
        'Connection': 'keep-alive',
        'Cache-Control': 'no-cache',
    }

    try:
        REQ = urlencode({'q': req})
        URL = 'https://www.google.com/search?tbs=li:1&{}&amp;gws_rd=ssl&amp;gl=us '.format(
            REQ)
        r = send('GET', URL + googleAbuseToken, headers=headers)

        while r.status_code != 200:
            warn('You are temporary blacklisted from Google search. Complete the captcha at the following URL and copy/paste the content of GOOGLE_ABUSE_EXEMPTION cookie : {}'.format(URL))
            info('Need help ? Read https://github.com/sundowndev/PhoneInfoga/wiki')
            token = input('\nGOOGLE_ABUSE_EXEMPTION=')
            googleAbuseToken = '&google_abuse=' + token
            r = send('GET', URL + googleAbuseToken, headers=headers)

        soup = BeautifulSoup(r.text, 'html5lib')

        results = soup.find("div", id="search").find_all("div", class_="g")

        links = []
        counter = 0

        for result in results:
            counter += 1

            if int(counter) > int(stop):
                break

            url = result.find("a").get('href')
            url = re.sub(r'(?:\/url\?q\=)', '', url)
            url = re.sub(r'(?:\/url\?url\=)', '', url)
            url = re.sub(r'(?:\&sa\=)(?:.*)', '', url)
            url = re.sub(r'(?:\&rct\=)(?:.*)', '', url)

            if re.match(r"^(?:\/search\?q\=)", url) is not None:
                url = 'https://google.com' + url

            if url is not None:
                links.append(url)

        return links
    except Exception as e:
        error('Request failed. Please retry or open an issue on https://github.com/sundowndev/PhoneInfoga.')
        print(e)
        return []


def searchApi(req, stop):
    options = urlencode({'q': req, 'key': google_api_key, 'cx': google_cx_id})
    r = send('GET', 'https://www.googleapis.com/customsearch/v1?%s' % (options))
    response = r.json()

    if 'error' in response:
        error('Error while fetching Google search API. Please verify your keys.')

    if 'items' not in response:
        return []

    results = response['items']

    links = []
    counter = 0

    for result in results:
        counter += 1

        if int(counter) > int(stop):
            break

        if result['link'] is not None:
            links.append(result['link'])

    return links
