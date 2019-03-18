#!/usr/bin/env python3
# -*- coding:utf-8 -*- 
#
# @name   : PhoneInfoga - Phone numbers OSINT tool
# @url    : https://github.com/sundowndev
# @author : Raphael Cerveaux (sundowndev)

from lib.request import *
from lib.output import *
from urllib.parse import urlencode
from lib.format import *
from bs4 import BeautifulSoup

numberObj = {}
number = ''
localNumber = ''
internationalNumber = ''
numberCountryCode = ''
customFormatting = ''

googleAbuseToken = ''
customFormatting = ''

def search(req, stop):
    global googleAbuseToken

    chosenUserAgent = random.choice(uagent)

    reqSession = requests.Session()
    headers = {
        'User-Agent': chosenUserAgent,
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
        REQ = urlencode({'q': req})
        URL = 'https://www.google.com/search?tbs=li:1&{}&amp;gws_rd=ssl&amp;gl=us '.format(
            REQ)
        r = reqSession.get(URL + googleAbuseToken, headers=headers)

        while r.status_code != 200:
            warn('You are temporary blacklisted from Google search. Complete the captcha at the following URL and copy/paste the content of GOOGLE_ABUSE_EXEMPTION cookie : {}'.format(URL))
            info('\n' +
                  'Need help ? Read https://github.com/sundowndev/PhoneInfoga/wiki')
            token = input('\nGOOGLE_ABUSE_EXEMPTION=')
            googleAbuseToken = '&google_abuse=' + token
            r = reqSession.get(URL + googleAbuseToken, headers=headers)

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

def osintIndividualScan():
    global numberObj
    global number
    global internationalNumber
    global numberCountryCode
    global customFormatting

    dorks = json.load(open('osint/individuals.json'))

    for dork in dorks:
        if dork['dialCode'] is None or dork['dialCode'] == numberCountryCode:
            if customFormatting:
                dorkRequest = replaceVariables(
                    dork['request'], numberObj) + ' | intext:"{}"'.format(customFormatting)
            else:
                dorkRequest = replaceVariables(dork['request'], numberObj)

            info("Searching for footprints on {}...".format(dork['site']))

            for result in search(dorkRequest, stop=dork['stop']):
                plus("URL: " + result)
        else:
            return -1

def osintReputationScan():
    global numberObj
    global number
    global internationalNumber
    global customFormatting

    dorks = json.load(open('osint/reputation.json'))

    for dork in dorks:
        if customFormatting:
            dorkRequest = replaceVariables(
                dork['request'], numberObj) + ' | intext:"{}"'.format(customFormatting)
        else:
            dorkRequest = replaceVariables(dork['request'], numberObj)

        info("Searching for {}...".format(dork['title']))
        for result in search(dorkRequest, stop=dork['stop']):
            plus("URL: " + result)

def osintSocialMediaScan():
    global numberObj
    global number
    global internationalNumber
    global customFormatting

    dorks = json.load(open('osint/social_medias.json'))

    for dork in dorks:
        if customFormatting:
            dorkRequest = replaceVariables(
                dork['request'], numberObj) + ' | intext:"{}"'.format(customFormatting)
        else:
            dorkRequest = replaceVariables(dork['request'], numberObj)

        info("Searching for footprints on {}...".format(dork['site']))

        for result in search(dorkRequest, stop=dork['stop']):
            plus("URL: " + result)

def osintDisposableNumScan():
    global numberObj
    global number

    dorks = json.load(open('osint/disposable_num_providers.json'))

    for dork in dorks:
        dorkRequest = replaceVariables(dork['request'], numberObj)

        info("Searching for footprints on {}...".format(dork['site']))

        for result in search(dorkRequest, stop=dork['stop']):
            plus("Result found: {}".format(dork['site']))
            plus("URL: " + result)
            askForExit()

def osintScan(numberObject, rerun=False):
    global numberObj
    global number
    global localNumber
    global internationalNumber
    global numberCountryCode
    global customFormatting

    numberObj = numberObject
    number = numberObj['default']
    localNumber = numberObj['local']
    internationalNumber = numberObj['international']
    numberCountryCode = numberObj['countryCode']

    info('Running OSINT footprint reconnaissance...')

    if not rerun:
        # Whitepages
        info("Generating scan URL on 411.com...")
        plus("Scan URL: https://www.411.com/phone/{}".format(
            internationalNumber.replace('+', '').replace(' ', '-')))

        askingCustomPayload = input('Would you like to use an additional format for this number ? (y/N) ')

    if rerun or askingCustomPayload == 'y' or askingCustomPayload == 'yes':
        info('We recommand: {} or {}'.format(internationalNumber,
                                                          internationalNumber.replace(numberCountryCode + ' ', '')))
        customFormatting = input('Custom format: ')

    info('---- Web pages footprints ----')

    info("Searching for footprints on web pages... (limit=10)")
    if customFormatting:
        req = '{} | intext:"{}" | intext:"{}" | intext:"{}"'.format(
            number, number, internationalNumber, customFormatting)
    else:
        req = '{} | intext:"{}" | intext:"{}"'.format(
            number, number, internationalNumber)

    for result in search(req, stop=10):
        plus("Result found: " + result)

    # Documents
    info("Searching for documents... (limit=10)")
    if customFormatting:
        req = '[ext:doc | ext:docx | ext:odt | ext:pdf | ext:rtf | ext:sxw | ext:psw | ext:ppt | ext:pptx | ext:pps | ext:csv | ext:txt | ext:xls] && [intext:"{}"]'.format(
            customFormatting)
    else:
        req = '[ext:doc | ext:docx | ext:odt | ext:pdf | ext:rtf | ext:sxw | ext:psw | ext:ppt | ext:pptx | ext:pps | ext:csv | ext:txt | ext:xls] && [intext:"{}" | intext:"{}"]'.format(
            internationalNumber, localNumber)
    for result in search(req, stop=10):
        plus("Result found: " + result)

    info('---- Reputation footprints ----')

    osintReputationScan()

    info("Generating URL on scamcallfighters.com...")
    plus('http://www.scamcallfighters.com/search-phone-{}.html'.format(number))

    tmpNumAsk = input("Would you like to search for temporary number providers footprints ? (Y/n) ")

    if tmpNumAsk.lower() != 'n' and tmpNumAsk.lower() != 'no':
        info('---- Temporary number providers footprints ----')

        try:
            info("Searching for phone number on tempophone.com...")
            response = requests.request(
                "GET", "https://tempophone.com/api/v1/phones")
            data = json.loads(response.content.decode('utf-8'))
            for voip_number in data['objects']:
                if voip_number['phone'] == formatNumber(number):
                    plus("Found a temporary number provider: tempophone.com")
                    askForExit()
        except Exception as e:
            error("Unable to reach tempophone.com API. Skipping.")

        osintDisposableNumScan()

    info('---- Social media footprints ----')

    osintSocialMediaScan()

    info('---- Phone books footprints ----')

    if numberCountryCode == '+1':
        info("Generating URL on True People... ")
        plus('https://www.truepeoplesearch.com/results?phoneno={}'.format(
            internationalNumber.replace(' ', '')))

    osintIndividualScan()

    retry_input = input("Would you like to rerun OSINT scan ? (e.g to use a different format) (y/N) ")

    if retry_input.lower() == 'y' or retry_input.lower() == 'yes':
        osintScan(numberObj, True)
    else:
        return -1
