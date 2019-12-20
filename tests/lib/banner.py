#!/usr/bin/env python3
# -*- coding:utf-8 -*-
#
# @name   : PhoneInfoga - Phone numbers OSINT tool
# @url    : https://github.com/sundowndev
# @author : Raphael Cerveaux (sundowndev)

import unittest
from unittest.mock import patch

from lib.banner import banner, __version__


class TestBanner(unittest.TestCase):
    def test_version(self):
        self.assertIs(type(__version__), type("str"))

    @patch("builtins.print")
    def test_banner(self, printMock):
        banner()

        self.assertEqual(printMock.call_count, 9)
