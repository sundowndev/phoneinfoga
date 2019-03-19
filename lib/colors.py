#!/usr/bin/env python3
# -*- coding:utf-8 -*-
#
# @name   : PhoneInfoga - Phone numbers OSINT tool
# @url    : https://github.com/sundowndev
# @author : Raphael Cerveaux (sundowndev)

import sys
import colorama

if sys.platform == 'win32':
    colorama.init()

R = "\033[%s;31m"
B = "\033[%s;34m"
G = "\033[%s;32m"
W = "\033[%s;38m"
Y = "\033[%s;33m"
E = "\033[0m"
BOLD = "\033[1m"
