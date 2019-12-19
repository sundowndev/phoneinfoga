#!/usr/bin/env python3
# -*- coding:utf-8 -*-
#
# @name   : PhoneInfoga - Phone numbers OSINT tool
# @url    : https://github.com/sundowndev
# @author : Raphael Cerveaux (sundowndev)

import unittest

from lib.banner import banner, __version__


class TestAnswer(unittest.TestCase):
    def test_type(self):
        self.assertIs(type(__version__), type("str"))
