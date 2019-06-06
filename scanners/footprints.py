#!/usr/bin/env python3
# -*- coding:utf-8 -*-
#
# @name   : PhoneInfoga - Phone numbers OSINT tool
# @url    : https://github.com/sundowndev
# @author : Raphael Cerveaux (sundowndev)

from urllib.parse import urlencode
from lib.output import *
from lib.format import *
from lib.request import send
from lib.googlesearch import search

numberObj = {}
number = ''
localNumber = ''
internationalNumber = ''
numberCountryCode = ''
customFormatting = ''


def osintIndividualScan():
    global numberObj
    global number
    global internationalNumber
    global numberCountryCode
    global customFormatting

    info('---- Phone books footprints ----')

    if numberCountryCode == '+1':
        info("Generating URL on True People... ")
        plus('https://www.truepeoplesearch.com/results?phoneno={}'.format(
            internationalNumber.replace(' ', '')))

    dorks = json.load(open('osint/individuals.json'))

    for dork in dorks:
        if dork['dialCode'] is None or dork['dialCode'] == numberCountryCode:
            if customFormatting:
                dorkRequest = replaceVariables(
                    dork['request'], numberObj) + ' OR "{}"'.format(customFormatting)
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

    info('---- Reputation footprints ----')

    dorks = json.load(open('osint/reputation.json'))

    for dork in dorks:
        if customFormatting:
            dorkRequest = replaceVariables(
                dork['request'], numberObj) + ' OR "{}"'.format(customFormatting)
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

    info('---- Social media footprints ----')

    dorks = json.load(open('osint/social_medias.json'))

    for dork in dorks:
        if customFormatting:
            dorkRequest = replaceVariables(
                dork['request'], numberObj) + ' OR "{}"'.format(customFormatting)
        else:
            dorkRequest = replaceVariables(dork['request'], numberObj)

        info("Searching for footprints on {}...".format(dork['site']))

        for result in search(dorkRequest, stop=dork['stop']):
            plus("URL: " + result)


def osintDisposableNumScan():
    global numberObj
    global number

    info('---- Temporary number providers footprints ----')

    try:
        info("Searching for phone number on tempophone.com...")
        response = send(
            "GET", "https://tempophone.com/api/v1/phones")
        data = json.loads(response.content.decode('utf-8'))
        for voip_number in data['objects']:
            if voip_number['phone'] == formatNumber(number):
                plus("Found a temporary number provider: tempophone.com")
                askForExit()
    except Exception as e:
        error("Unable to reach tempophone.com API. Skipping.")

    dorks = json.load(open('osint/disposable_num_providers.json'))

    for dork in dorks:
        dorkRequest = replaceVariables(dork['request'], numberObj)

        info("Searching for footprints on {}...".format(dork['site']))

        for result in search(dorkRequest, stop=dork['stop']):
            plus("Result found: {}".format(dork['site']))
            plus("URL: " + result)
            askForExit()


def osintScan(numberObject, rerun=False):
    if not args.scanner == 'footprints' and not args.scanner == 'all':
        return -1

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

    test('Running OSINT footprint reconnaissance...')

    if not rerun:
        # Whitepages
        info("Generating scan URL on 411.com...")
        plus("Scan URL: https://www.411.com/phone/{}".format(
            internationalNumber.replace('+', '').replace(' ', '-')))

        askingCustomPayload = ask(
            'Would you like to use an additional format for this number ? (y/N) ')

    if rerun or askingCustomPayload == 'y' or askingCustomPayload == 'yes':
        info('We recommand: {} or {}'.format(internationalNumber,
                                             internationalNumber.replace(numberCountryCode + ' ', '')))
        customFormatting = ask('Custom format: ')

    info('---- Web pages footprints ----')

    info("Searching for footprints on web pages... (limit=10)")
    if customFormatting:
        req = '{} OR "{}" OR "{}" OR "{}"'.format(
            number, number, internationalNumber, customFormatting)
    else:
        req = '{} OR "{}" OR "{}"'.format(
            number, number, internationalNumber)

    for result in search(req, stop=10):
        plus("Result found: " + result)

    # Documents
    info("Searching for documents... (limit=10)")
    if customFormatting:
        req = '[ext:doc OR ext:docx OR ext:odt OR ext:pdf OR ext:rtf OR ext:sxw OR ext:psw OR ext:ppt OR ext:pptx OR ext:pps OR ext:csv OR ext:txt OR ext:xls] AND ["{}"]'.format(
            customFormatting)
    else:
        req = '[ext:doc OR ext:docx OR ext:odt OR ext:pdf OR ext:rtf OR ext:sxw OR ext:psw OR ext:ppt OR ext:pptx OR ext:pps OR ext:csv OR ext:txt OR ext:xls] AND ["{}" OR "{}"]'.format(
            internationalNumber, localNumber)
    for result in search(req, stop=10):
        plus("Result found: " + result)

    osintReputationScan()

    info("Generating URL on scamcallfighters.com...")
    plus('http://www.scamcallfighters.com/search-phone-{}.html'.format(number))

    tmpNumAsk = ask(
        "Would you like to search for temporary number providers footprints ? (Y/n) ")

    if tmpNumAsk.lower() != 'n' and tmpNumAsk.lower() != 'no':
        osintDisposableNumScan()

    osintSocialMediaScan()

    osintIndividualScan()

    retry_input = ask(
        "Would you like to rerun OSINT scan ? (e.g to use a different format) (y/N) ")

    if retry_input.lower() == 'y' or retry_input.lower() == 'yes':
        osintScan(numberObj, True)
    else:
        return -1
