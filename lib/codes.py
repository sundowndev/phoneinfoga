#!/usr/bin/env python3
# -*- coding:utf-8 -*- 
#
# @name   : PhoneInfoga - Phone numbers OSINT tool
# @url    : https://github.com/sundowndev
# @author : Raphael Cerveaux (sundowndev)

if args.no_ansi or args.output:
    code_info = '[-] '
    code_warning = '(!) '
    code_result = '[+] '
    code_error = '[!] '
    code_title = ''
else:
    code_info = Fore.RESET + Style.BRIGHT + '[-] '
    code_warning = Fore.YELLOW + Style.BRIGHT + '(!) '
    code_result = Fore.GREEN + Style.BRIGHT + '[+] '
    code_error = Fore.RED + Style.BRIGHT + '[!] '
    code_title = Fore.YELLOW + Style.BRIGHT