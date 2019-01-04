# PhoneInfoga

![](https://img.shields.io/badge/python-3-blue.svg)
![](https://img.shields.io/github/tag/SundownDEV/PhoneInfoga.svg)
![](https://img.shields.io/badge/license-MIT-blue.svg)

Information gathering & OSINT reconnaissance tool for phone numbers.

One of the most advanced tools to scan phone numbers using only free resources. The goal is to first gather basic information such as country, area, carrier and line type on any international phone numbers with a very good accuracy. Then try to determine the VoIP provider or search for footprints on search engines to try identify the owner.

### [OSINT Tutorial: Building an OSINT Reconnaissance Tool from Scratch](https://medium.com/@SundownDEV/phone-number-scanning-osint-recon-tool-6ad8f0cac27b)

## Features

- Check if phone number exists and is possible
- Gather standard informations such as country, line type and carrier
- Check several numbers at once
- OSINT reconnaissance using external APIs, Google Hacking, phone books & search engines
- Use custom formatting for more effective OSINT reconnaissance

![](https://i.imgur.com/bWx79dy.png)

## Formats

The tool only accepts E164 and International formats as input.

- E164: +3396360XXXX
- International: +33 9 63 60 XX XX
- National: 09 63 60 XX XX
- RFC3966: tel:+33-9-63-60-XX-XX
- Out-of-country format from US: 011 33 9 63 60 XX XX

## Available scanners

Use `any` to disable this feature. Default value: `all`

- numverify
- ovh

## Installation

```bash
git clone https://github.com/sundowndev/PhoneInfoga
cd PhoneInfoga/
python3 -m pip install -r requirements.txt
```

## Usage

```
usage: phoneinfoga.py -n <number> [options]

Advanced information gathering tool for phone numbers
(https://github.com/sundowndev/PhoneInfoga)

optional arguments:
  -h, --help            show this help message and exit
  -n number, --number number
                        The phone number to scan (E164 or International
                        format)
  -i input_file, --input input_file
                        Phone number list to scan (one per line)
  -o output_file, --output output_file
                        Output to save scan results
  -s scanner, --scanner scanner (any to skip, default: all)
                        The scanner to use
  --osint               Use OSINT reconnaissance
  -u, --update          Update the tool & databases
```

Example (quotes are optional, use it when typing special formats) :

```
python3 phoneinfoga.py -n "(+42)837544833"
```

Check for a number range on OVH :

```
python3 phoneinfoga.py -n +42837544833 -s ovh
```

Check several numbers at once :

```
python3 phoneinfoga.py -i numbers.txt -o results.txt
```

**Note: `--osint` is not compatible with `--output` option.**

Use all scanners and run OSINT reconnaissance :

```
python3 phoneinfoga.py -n +42837544833 -s all --osint
```

## Formatting

E.164 formatting for phone numbers entails the following:

- A + (plus) sign
- International Country Calling code
- Local Area code
- Local Phone number

For example, here’s a US-based number in standard local formatting: (415) 555-2671

![](https://i.imgur.com/0e2SMdL.png)

Here’s the same phone number in E.164 formatting: +14155552671

![](https://i.imgur.com/KfrvacR.png)

In the UK, and many other countries internationally, local dialing may require the addition of a '0' in front of the subscriber number. With E.164 formatting, this '0' must usually be removed.

For example, here’s a UK-based number in standard local formatting: 020 7183 8750

Here’s the same phone number in E.164 formatting: +442071838750

## Dealing with Google captcha

PhoneInfo use a workaround to handle Google bot detection. When running OSINT scan, you will usually be blacklisted very easily by Google, which will ask the tool to complete a captcha.

>When you search on Google using custom requests (Google Dorks), you get very easily blacklisted. So Google shows up a page where you have to complete a captcha to continue. As soon as the captcha is completed, Google create a cookie named "GOOGLE_ABUSE_EXEMPTION" which is used to whitelist your browser and IP address for some minutes. This temporary whitelist is enough to let you gather a lot of information from many sources. So I decided to add a simple user manipulation to bypass this bot detection. [...] So I'll just try make requests and wait until I get a 503 error, which means I got blacklisted. Then I ask the user to follow an URL to manually complete the captcha and copy the whitelist token to paste it in the CLI. The tool is now able to continue to scan!

![](https://i.imgur.com/qbFZa1m.png)

### How to handle captcha
- Follow the URL
- Complete the captcha if needed
- Open the dev tool (F12 on most browsers)
- Go to **Storage**, then **Cookies**
- Copy the value of the *GOOGLE_ABUSE_EXEMPTION* cookie and simply paste it in the CLI

![](https://i.imgur.com/KkE1EM5.png)

### Troubleshooting

The cookie should be created after you complete the captcha. If there's no captcha and *GOOGLE_ABUSE_EXEMPTION* cookie, try pressing F5 to refresh the page. The cookie should've been created. If refreshing the page does not help, change the query to something different (change the number or add text). Google will not necessarily ask you to complete a captcha if your request is the exact same as the previous one, because it'll usually be cached.

## Custom formatting

Sometimes the phone number has footprints but is used with a different formatting. This is a problem because for example if we search for "+15417543010", we'll not find web pages that write it that way : "(541) 754–3010". So the tool use a (optional) custom formatting given by the user to find further and more accurate results. To use this feature properly and make the results more valuable, try to use a format that someone of the number' country would usually use to share the phone number online. For example, French people usually write numbers that way online : *06.20.30.40.50*, *06 20 30 40 50*.

## License

This tool is licensed under the GNU General Public License v3.0.

[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fsundowndev%2FPhoneInfoga.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fsundowndev%2FPhoneInfoga?ref=badge_large)
