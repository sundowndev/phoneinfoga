#!/usr/bin/env python
# -*- coding:utf-8 -*- 
#
# @name   : Infoga - Email OSINT
# @url    : http://github.com/m4ll0k
# @author : Momo Outaadi (m4ll0k)

import sys
import colorama
import atexit

if sys.platform == 'win32':
    colorama.init()

# TODO: if args.no_ansi or args.output:
if True:
    R = "\033[%s;31m"
    B = "\033[%s;34m"
    G = "\033[%s;32m"
    W = "\033[%s;38m"
    Y = "\033[%s;33m"
    E = "\033[0m"
    BOLD = ""
else:
    R = ""
    B = ""
    G = ""
    W = ""
    Y = ""
    E = ""
    BOLD = ""
