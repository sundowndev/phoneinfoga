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


def phone_us_format(phone_number, delimiter):
    clean_phone_number = re.sub('[^0-9]+', '', phone_number)
    formatted_phone_number = re.sub(
        "(\d)(?=(\d{3})+(?!\d))", r"\1" + delimiter, "%d" % int(clean_phone_number[:-1])) + clean_phone_number[-1]
    return formatted_phone_number


def phone_format(phone_number, delimiter):
    clean_phone_number = re.sub('[^0-9]+', '', phone_number)
    formatted_phone_number = re.sub(
        "(\d)(?=(\d{2})+(?!\d))", r"\1" + delimiter, "%d" % int(clean_phone_number[:-1])) + clean_phone_number[-1]
    return formatted_phone_number


def scan(number):
    if not args.recon:
        return -1

    test('Running custom format reconnaissance...')

    cc = number['countryCode'].replace('+', '')
    nb = number['local']

    if number['countryIsoCode'] == 'US' or number['countryIsoCode'] == 'CA':
        segments = phone_us_format(cc + nb, ' ').split(' ')

        seg1 = segments[-3]
        seg2 = segments[-2]
        seg3 = segments[-1]

        formats = [
            '%s%s%s' % (seg1, seg2, seg3),
            '%s %s%s%s' % (cc, seg1, seg2, seg3),
            '%s %s %s%s' % (cc, seg1, seg2, seg3),
            '%s %s%s' % (seg1, seg2, seg3),
            '%s-%s%s' % (seg1, seg2, seg3),
            '%s-%s-%s' % (seg1, seg2, seg3),
            '+%s %s-%s-%s' % (cc, seg1, seg2, seg3),
            '(+%s)%s-%s-%s' % (cc, seg1, seg2, seg3),
            '+%s/%s-%s-%s' % (cc, seg1, seg2, seg3),
            '(%s) %s%s' % (seg1, seg2, seg3),
            '(%s) %s-%s' % (seg1, seg2, seg3),
            '(%s) %s.%s' % (seg1, seg2, seg3),
            '(%s)%s%s' % (seg1, seg2, seg3),
            '(%s)%s-%s' % (seg1, seg2, seg3),
            '(%s)%s.%s' % (seg1, seg2, seg3)
        ]
    else:
        formated_number = number['international'].replace(number['countryCode'] + ' ', '').split(' ')

        segments = []
        
        for seg in formated_number:
          segments.append(seg)
        
        formats = [
            '+%s0%s' % (cc, nb),
            '(00%s)0%s' % (cc, number['local']),
            '+%s/0%s' % (cc, nb),
            '+%s0%s' % (cc, '-'.join(segments)),
            '(+%s)0%s' % (cc, '-'.join(segments)),
            '(00%s)0%s' % (cc, '-'.join(segments)),
            '(00%s)0%s' % (cc, '-'.join(segments)),
            '+%s/0%s' % (cc, ' '.join(segments)),
            '+%s0%s' % (cc, ' '.join(segments)),
            '(00%s)0%s' % (cc, ' '.join(segments)),
            '(+%s)0%s' % (cc, ''.join(segments)),
            '(+%s)0%s' % (cc, ' '.join(segments)),
            '+%s/0%s' % (cc, '-'.join(segments)),
            '+%s/0%s' % (cc, ' '.join(segments)),
        ]
        
    for format in formats:
        print('Footprint reconnaissance for %s' % (format))
        for result in search('"%s"' % (format), stop=5):
            plus("URL: " + result)
