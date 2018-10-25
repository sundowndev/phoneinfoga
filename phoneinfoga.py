#!/usr/bin/env python

try:
	import requests
	import sys
	import hashlib
	import json
	from bs4 import BeautifulSoup
except:
    print "Request library not found, please install it before proceeding\n"
    sys.exit()

print "\n \033[92m"
print "    ___ _                       _____        __                   "
print "   / _ \ |__   ___  _ __   ___  \_   \_ __  / _| ___   __ _  __ _ "
print "  / /_)/ '_ \ / _ \| '_ \ / _ \  / /\/ '_ \| |_ / _ \ / _` |/ _` |"
print " / ___/| | | | (_) | | | |  __/\/ /_ | | | |  _| (_) | (_| | (_| |"
print " \/    |_| |_|\___/|_| |_|\___\____/ |_| |_|_|  \___/ \__, |\__,_|"
print "                                                      |___/       "
print " PhoneInfoga Ver. 0.2.1                                          "
print " Coded by Raphael Cerveaux <raphael@crvx.fr>                      "
print "\033[94m\n\n"

def help():
	print "Usage: PhoneInfoga options \n"
	print "       -n|--number: Phone number to search"
	print "       -h|--help: Help command"

def getRequestSecret():
	requestSecret = ''
	resp = requests.get('https://numverify.com/')
	soup = BeautifulSoup(resp.text, "html5lib")
	for tag in soup.find_all("input", type="hidden"):
		if tag['name'] == "scl_request_secret":
			requestSecret = tag['value']
			break;
	
	return requestSecret

def getInformations(PhoneNumber):
	# verify input type
	if str.isdigit(PhoneNumber) != True:
		print("\033[31mError: please enter a valid integer.")
		sys.exit()

	print("Fetching information for number +" + PhoneNumber + "...")

	apiKey = hashlib.md5(PhoneNumber + getRequestSecret()).hexdigest()

	response = requests.get("https://numverify.com/php_helper_scripts/phone_api.php?secret_key=" + apiKey + "&number=" + PhoneNumber)
	if response.content == "Unauthorized" or response.status_code != 200:
		print("An error occured while calling the API (bad request or wrong api key).")
		sys.exit()

	data = json.loads(response.content)
	
	try:
		data["valid"] == True
	except:
		print("\033[31mError: Please specify a phone number. " + PhoneNumber + " is not valid.")
		print("Example: 14158586273\033[94m")
		sys.exit()
	else:
		print "\n"
		print "\033[1;32m1 result found for (" + data["country_prefix"] + ") " + data["local_format"]
		print "\n"
		print("[Country] " + data["country_name"] + "(" + data["country_code"] + ")")
		print("[Carrier] " + data["carrier"])
		print("[Line type] " + data["line_type"])

try:
	sys.argv[1:][0] == "-n" or sys.argv[1:][0] == "--number"
except:
	help()
	sys.exit()
else:
	PhoneNumber = sys.argv[1:][1]
	getInformations(PhoneNumber)
