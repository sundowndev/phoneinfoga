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
        'Cookie': 'Cookie: CGIC=Ij90ZXh0L2h0bWwsYXBwbGljYXRpb24veGh0bWwreG1sLGFwcGxpY2F0aW9uL3htbDtxPTAuOSwqLyo7cT0wLjg; CONSENT=YES+RE.fr+20150809-08-0; 1P_JAR=2018-11-28-14; NID=148=aSdSHJz71rufCokaUC93nH3H7lOb8E7BNezDWV-PyyiHTXqWK5Y5hsvj7IAzhZAK04-QNTXjYoLXVu_eiAJkiE46DlNn6JjjgCtY-7Fr0I4JaH-PZRb7WFgSTjiFqh0fw2cCWyN69DeP92dzMd572tQW2Z1gPwno3xuPrYC1T64wOud1DjZDhVAZkpk6UkBrU0PBcnLWL7YdL6IbEaCQlAI9BwaxoH_eywPVyS9V; SID=uAYeu3gT23GCz-ktdGInQuOSf-5SSzl3Plw11-CwsEYY0mqJLSiv7tFKeRpB_5iz8SH5lg.; HSID=AZmH_ctAfs0XbWOCJ; SSID=A0PcRJSylWIxJYTq_; APISID=HHB2bKfJ-2ZUL5-R/Ac0GK3qtM8EHkloNw; SAPISID=wQoxetHBpyo4pJKE/A2P6DUM9zGnStpIVt; SIDCC=ABtHo-EhFAa2AJrJIUgRGtRooWyVK0bAwiQ4UgDmKamfe88xOYBXM47FoL5oZaTxR3H-eOp7-rE; OTZ=4671861_52_52_123900_48_436380; OGPC=873035776-8:; OGP=-873035776:;'
    }

    try:
        REQ = urlencode({ 'q': req, 'num': stop })
        URL = 'https://www.google.com/search?tbs=li:1&{}&amp;gws_rd=ssl&amp;gl=us'.format(
            REQ)
        r = send('GET', URL + googleAbuseToken, headers=headers)

        while r.status_code != 200:
            warn('You are temporary blacklisted from Google search. Complete the captcha at the following URL and copy/paste the content of GOOGLE_ABUSE_EXEMPTION cookie : {}'.format(URL))
            info('Need help ? Read https://github.com/sundowndev/PhoneInfoga/wiki')
            token = ask('\nGOOGLE_ABUSE_EXEMPTION=')
            googleAbuseToken = '&google_abuse=' + token
            r = send('GET', URL + googleAbuseToken, headers=headers)

        soup = BeautifulSoup(r.text, 'html5lib')

        results = soup.find("div", id="search").find_all("div", class_="g")

        links = []

        for result in results:
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
    options = urlencode({ 'q': req, 'key': google_api_key, 'cx': google_cx_id, 'num': stop })
    r = send('GET', 'https://www.googleapis.com/customsearch/v1?%s' % (options))
    response = r.json()

    if 'error' in response:
        error('Error while fetching Google search API. Maybe usage limit ? Please verify your keys.')
        print(response['error'])
        askForExit()
        return []

    if 'items' not in response:
        return []

    results = response['items']

    links = []

    for result in results:
        if result['link'] is not None:
            links.append(result['link'])

    return links
