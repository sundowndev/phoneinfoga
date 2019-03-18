#!/usr/bin/env python
# -*- coding:utf-8 -*- 
#
# @name   : Infoga - Email OSINT
# @url    : http://github.com/m4ll0k
# @author : Momo Outaadi (m4ll0k)

import re

def formatNumber(InputNumber):
    return re.sub("(?:\+)?(?:[^[0-9]*)", "", InputNumber)

def replaceVariables(string, number):
    string = string.replace('$n', number['default'])
    string = string.replace('$i', number['international'])
    string = string.replace('$l', number['local'])

    return string
