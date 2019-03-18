#!/usr/bin/env python
# -*- coding:utf-8 -*- 
#
# @name   : Infoga - Email OSINT
# @url    : http://github.com/m4ll0k
# @author : Momo Outaadi (m4ll0k)

import sys
import json
from lib.colors import *

def plus(string): print("%s[+] %s%s" % (G%0, string, E))
def warn(string): print("%s(!) %s%s" % (Y%0, string, E))
def error(string): print("%s[!]%s %s%s" % (R%0, E, string, E))
def test(string): print("%s[*] %s%s" % (B%0, string, E))
def info(string): print("%s[i] %s%s" % (E, string, E))
def more(string): print(" %s|%s  %s%s" % (W%0, string, E))
def title(string): print("%s%s%s" % (Y%0, string, E))
def throw(string):
    error(string)
    sys.exit()

def askForExit():
    # TODO: parse args
    # if not args.output
    if not False:
        user_input = input("Continue scanning ? (y/N) ")

        if user_input.lower() == 'y' or user_input.lower() == 'yes':
            return -1
        else:
            info("Good bye!")
            sys.exit()