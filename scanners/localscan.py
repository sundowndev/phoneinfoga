#!/usr/bin/env python3
# -*- coding:utf-8 -*-
#
# @name   : PhoneInfoga - Phone numbers OSINT tool
# @url    : https://github.com/sundowndev
# @author : Raphael Cerveaux (sundowndev)

import phonenumbers
from phonenumbers import carrier
from phonenumbers import geocoder
from phonenumbers import timezone
from lib.output import *
from lib.format import *


def scan(InputNumber, print_results=True):
    test('Running local scan...')

    FormattedPhoneNumber = "+" + formatNumber(InputNumber)

    try:
        PhoneNumberObject = phonenumbers.parse(FormattedPhoneNumber, None)
    except Exception as e:
        throw(e)
    else:
        if not phonenumbers.is_valid_number(PhoneNumberObject):
            return False

        number = phonenumbers.format_number(
            PhoneNumberObject, phonenumbers.PhoneNumberFormat.E164).replace('+', '')
        numberCountryCode = phonenumbers.format_number(
            PhoneNumberObject, phonenumbers.PhoneNumberFormat.INTERNATIONAL).split(' ')[0]
        numberCountry = phonenumbers.region_code_for_country_code(
            int(numberCountryCode))

        localNumber = phonenumbers.format_number(
            PhoneNumberObject, phonenumbers.PhoneNumberFormat.E164).replace(numberCountryCode, '')
        internationalNumber = phonenumbers.format_number(
            PhoneNumberObject, phonenumbers.PhoneNumberFormat.INTERNATIONAL)

        country = geocoder.country_name_for_number(PhoneNumberObject, "en")
        location = geocoder.description_for_number(PhoneNumberObject, "en")
        carrierName = carrier.name_for_number(PhoneNumberObject, 'en')

        if print_results:
            plus('International format: {}'.format(internationalNumber))
            plus('Local format: {}'.format(localNumber))
            plus('Country found: {} ({})'.format(country, numberCountryCode))
            plus('City/Area: {}'.format(location))
            plus('Carrier: {}'.format(carrierName))
            for timezoneResult in timezone.time_zones_for_number(PhoneNumberObject):
                plus('Timezone: {}'.format(timezoneResult))

            if phonenumbers.is_possible_number(PhoneNumberObject):
                info('The number is valid and possible.')
            else:
                warn('The number is valid but might not be possible.')

    numberObj = {}
    numberObj['input'] = InputNumber
    numberObj['default'] = number
    numberObj['local'] = localNumber
    numberObj['international'] = internationalNumber
    numberObj['country'] = country
    numberObj['countryCode'] = numberCountryCode
    numberObj['countryIsoCode'] = numberCountry
    numberObj['location'] = location
    numberObj['carrier'] = carrierName

    return numberObj
