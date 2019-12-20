#!/usr/bin/env python3
# -*- coding:utf-8 -*-
#
# @name   : PhoneInfoga - Phone numbers OSINT tool
# @url    : https://github.com/sundowndev
# @author : Raphael Cerveaux (sundowndev)

import unittest
from unittest.mock import patch

from lib.format import formatNumber, replaceVariables


class TestFormat(unittest.TestCase):
    """
    We want to be sure re.sub is imported and called
    """

    @patch("re.sub")
    def test_formatNumberCallingSub(self, subMock):
        subMock.return_value = ""

        result = formatNumber("000")

        subMock.assert_called_with("(?:\+)?(?:[^[0-9]*)", "", "000")
        self.assertEqual(result, "")

    """
    We verify the function formats the number the right way
    """

    def test_formatNumber(self):
        self.assertEqual(formatNumber("+33 81495357"), "3381495357")
        self.assertEqual(formatNumber("0 81 49 53 57"), "081495357")
        self.assertEqual(formatNumber("+1 555-444-888"), "1555444888")

    def test_replaceVariables(self):
        number = {
            "input": "+33651580074",
            "default": "33651580074",
            "local": "651580074",
            "international": "+33 6 51 58 00 74",
            "country": "France",
            "countryCode": "+33",
            "countryIsoCode": "FR",
            "location": "France",
            "carrier": "",
        }

        self.assertEqual(replaceVariables("test $n", number), "test 33651580074")
        self.assertEqual(
            replaceVariables("test $i", number), "test +33 6 51 58 00 74"
        )
        self.assertEqual(replaceVariables("test $l", number), "test 6 51 58 00 74")
