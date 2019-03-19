#!/usr/bin/env python3
# -*- coding:utf-8 -*-
#
# @name   : PhoneInfoga - Phone numbers OSINT tool
# @url    : https://github.com/sundowndev
# @author : Raphael Cerveaux (sundowndev)

import re
from lib.args import args
from lib.output import *
from lib.googlesearch import search


def phone_format(phone_number, delimiter):
    clean_phone_number = re.sub('[^0-9]+', '', phone_number)
    formatted_phone_number = re.sub(
        "(\d)(?=(\d{3})+(?!\d))", r"\1" + delimiter, "%d" % int(clean_phone_number[:-1])) + clean_phone_number[-1]
    return formatted_phone_number


def scan(number):
    if not args.recon:
        return -1

    test('Running custom format reconnaissance...')

    if number['countryIsoCode'] == 'US' or number['countryIsoCode'] == 'CA':
        print(1)
    else:
        cc = '33'
        nb = '186481407'
        
        print(phone_format(cc + nb, ' '))

        formats = [
            '+%s01 86 48 14 07' % (cc),
            '+%s0%s' % (cc, nb),
            # '+33018 648 140 7',
            # '(0033)0186481407',
            # '(+33)018 648 140 7',
            '+%s/0%s' % (cc, nb),
            # '(0033)018 648 140 7',
            # '+33018-648-140-7',
            # '(+33)0186481407',
            '(+%s)01 86 48 14 07' % (cc, phone_format(cc + nb, ' ')),
            # '+33/018-648-140-7',
            # '+33/01-86-48-14-07',
            # '+3301-86-48-14-07',
            # '(0033)01 86 48 14 07',
            # '+33/01 86 48 14 07',
            # '(+33)018-648-140-7',
            # '(+33)01-86-48-14-07',
            # '(0033)01-86-48-14-07',
            # '(0033)018-648-140-7',
            # '+33/018 648 140 7'
        ]

    # "6185551212" OR "618 5551212" OR "618-5551212" OR " 1 6185551212" OR " 1 618 5551212"
    for format in formats:
        print('Recon for %s' % (format))
        # for result in search('%s' % (format), stop=5):
        #     plus("URL: " + result)
