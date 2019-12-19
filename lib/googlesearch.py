#!/usr/bin/env python3
# -*- coding:utf-8 -*-
#
# @name   : PhoneInfoga - Phone numbers OSINT tool
# @url    : https://github.com/sundowndev
# @author : Raphael Cerveaux (sundowndev)

import os
import re
import json
import time
from urllib.parse import urlencode
from bs4 import BeautifulSoup
from lib.output import *
from lib.request import send
import os
from config import *

from selenium import webdriver
from selenium.webdriver.firefox.options import DesiredCapabilities
from selenium.webdriver.common.proxy import Proxy,ProxyType
from selenium.webdriver.firefox.firefox_binary import FirefoxBinary

browser = None
count = 0
fo = webdriver.FirefoxOptions()
fo.add_argument("log-level=3")
PROXIES = []
fo.add_argument("--headless")

def closeBrowser():
    if browser is not None:
        browser.quit()
        
def get_proxies(fo=fo):
    driver = webdriver.Firefox(firefox_options=fo)
    driver.get("https://free-proxy-list.net/")
    proxies = driver.find_elements_by_css_selector("tr[role='row']")
    for p in proxies:
        result = p.text.split(" ")
        if result[7] == "yes":
            PROXIES.append(result[0]+":"+result[1])
    driver.close()
    return PROXIES

def proxy_driver(PROXIES,fo=fo, binary = None):
    prox = Proxy()
    if PROXIES:
        pxy = PROXIES[-1]
    else:
        info("Proxies used up (%s)" % len(PROXIES))
        PROXIES = get_proxies()
        print(PROXIES)
        pxy = PROXIES[-1]
    info('Trying proxy '+pxy)
    proxy = Proxy({
        'proxyType': ProxyType.MANUAL,
        'httpProxy': pxy,
        'ftpProxy': pxy,
        'sslProxy': pxy,
        'noProxy': '' 
    })
    driver = None
    if binary==None:
        driver  = webdriver.Firefox(firefox_options=fo,proxy=proxy)
    else:
        driver = webdriver.Firefox(firefox_options=fo, proxy=proxy,firefox_binary=binary)
    return driver

def search(req, stop,count=0):
    global browser
    global PROXIES
    if count == 1:
        if browser:
            browser.close()
        browser = None
        if len(PROXIES) > 0:
            PROXIES.pop()
    if len(PROXIES) == 0:
        PROXIES = get_proxies()
    if google_api_key and google_cx_id:
        return searchApi(req, stop)

    if browser is None:
        if os.environ.get('webdriverRemote'):
            browser = webdriver.Remote(os.environ.get('webdriverRemote'), webdriver.DesiredCapabilities.FIREFOX.copy())
        else:
            if firefox_exe_path.lstrip() == '':
                browser = proxy_driver(PROXIES)
                #browser = webdriver.Firefox()
            else:
                binary = FirefoxBinary(firefox_exe_path)
                browser = proxy_driver(PROXIES,binary=binary)
                #browser = webdriver.Firefox(firefox_binary=binary)

    try:
        REQ = urlencode({ 'q': req, 'num': stop, 'hl': 'en' })
        URL = 'https://www.google.com/search?tbs=li:1&{}&amp;gws_rd=ssl&amp;gl=us'.format(
            REQ)
        browser.get(URL)
        htmlBody = browser.find_element_by_css_selector("body").get_attribute('innerHTML')

        soup = BeautifulSoup(htmlBody, 'html5lib')

        if soup.find("div", id="recaptcha") is not None:
            warn('Temporary blacklisted from Google search.')
            return search(req,stop,count=1)
            token = ask('>')
            htmlBody = browser.find_element_by_css_selector("body").get_attribute('innerHTML')
            soup = BeautifulSoup(htmlBody, 'html5lib')

        results = soup.find("div", id="search").find_all("div", class_="g")

        links = []

        for result in results:
            url = result.find("a").get('href')
            url = re.sub(r'(?:\/url\?q\=)', '', url)
            url = re.sub(r'(?:\/url\?url\=)', '', url)
            url = re.sub(r'(?:\&sa\=)(?:.*)', '', url)
            url = re.sub(r'(?:\&rct\=)(?:.*)', '', url)

            if re.match(r"^(?:\/search\?q\=)", url) is not None:
                url = 'https://google.com' + url

            if url is not None:
                links.append(url)

        return links
    except Exception as e:
        error('Request failed. Please retry or open an issue on https://github.com/sundowndev/PhoneInfoga.')
        print(e)
        return []


def searchApi(req, stop):
    options = urlencode({ 'q': req, 'key': google_api_key, 'cx': google_cx_id, 'num': stop })
    r = send('GET', 'https://www.googleapis.com/customsearch/v1?%s' % (options))
    response = r.json()

    if 'error' in response:
        error('Error while fetching Google search API. Maybe usage limit ? Please verify your keys.')
        print(response['error'])
        askForExit()
        return []

    if 'items' not in response:
        return []

    results = response['items']

    links = []

    for result in results:
        if result['link'] is not None:
            links.append(result['link'])

    return links
