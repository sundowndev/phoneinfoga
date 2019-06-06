#!/usr/bin/env python3
# -*- coding:utf-8 -*-
#
# @name   : PhoneInfoga - Phone numbers OSINT tool
# @url    : https://github.com/sundowndev
# @author : Raphael Cerveaux (sundowndev)

import sys
import json
from lib.colors import *
from lib.args import args
from lib.logger import Logger


def plus(string):
    if not args.no_ansi and not args.output:
        print("%s[+] %s%s" % (G % 0, string, E))
    else:
        print("[+] %s" % (string))


def warn(string):
    if not args.no_ansi and not args.output:
        print("%s(!) %s%s" % (Y % 0, string, E))
    else:
        print("(!) %s" % (string))


def error(string):
    if not args.no_ansi and not args.output:
        print("%s[!]%s %s%s" % (R % 0, E, string, E))
    else:
        print("[!] %s" % (string))


def test(string):
    if not args.no_ansi and not args.output:
        print("%s[*] %s%s" % (B % 0, string, E))
    else:
        print("[*] %s" % (string))


def info(string):
    if not args.no_ansi and not args.output:
        print("%s[i] %s%s" % (E, string, E))
    else:
        print("[i] %s" % (string))


def more(string):
    if not args.no_ansi and not args.output:
        print(" %s|%s  %s%s" % (W % 0, string, E))
    else:
        print(" | %s" % (string))


def title(string):
    if not args.no_ansi and not args.output:
        print("%s%s%s%s" % (BOLD, Y % 0, string, E))
    else:
        print("%s" % (string))


def throw(string):
    error(string)
    sys.exit()


def askForExit():
    if not args.output:
        user_input = ask('Continue scanning ? (y/N) ')

        if user_input.lower() == 'y' or user_input.lower() == 'yes':
            return -1
        else:
            info("Good bye!")
            sys.exit()

def ask(text):
    if args.output:
        sys.stdout = sys.__stdout__
        res = input(text)
        sys.stdout = Logger()

        return res
    else:
        return input(text)