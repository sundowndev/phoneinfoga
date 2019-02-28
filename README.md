<h1 align="center">PhoneInfoga</h1>

<div align="center">
  <a href="https://travis-ci.org/sundowndev/PhoneInfoga">
    <img src="https://img.shields.io/travis/sundowndev/PhoneInfoga/master.svg?style=flat-square" alt="Build Status" />
  </a>
  <a href="#">
    <img src="https://img.shields.io/badge/python-3.6-blue.svg?style=flat-square" alt="Python version" />
  </a>
  <a href="https://github.com/sundowndev/PhoneInfoga/releases">
    <img src="https://img.shields.io/github/tag/SundownDEV/PhoneInfoga.svg?style=flat-square" alt="Latest version" />
  </a>
  <a href="https://github.com/sundowndev/PhoneInfoga/blob/master/LICENSE">
    <img src="https://img.shields.io/badge/license-MIT-blue.svg?style=flat-square" alt="License" />
  </a>
</div>

<h4 align="center">Information gathering & OSINT reconnaissance tool for phone numbers</h4>

<div align="center">
  <sub>For the love of open source investigations. Built with ❤︎ by
  <a href="https://twitter.com/sundowndev">@sundowndev</a>
</div>

<h3 align="center">
  <a href="https://github.com/sundowndev/PhoneInfoga/wiki">Documentation</a> | 
  <a href="https://medium.com/@SundownDEV/phone-number-scanning-osint-recon-tool-6ad8f0cac27b">OSINT Tutorial</a>
</h3>

## About

PhoneInfoga is one of the most advanced tools to scan phone numbers using only free resources. The goal is to first gather basic information such as country, area, carrier and line type on any international phone numbers with a very good accuracy. Then try to determine the VoIP provider or search for footprints on search engines to try identify the owner.

## Features

- Check if phone number exists and is possible
- Gather standard informations such as country, line type and carrier
- Check several numbers at once
- OSINT reconnaissance using external APIs, Google Hacking, phone books & search engines
- Use custom formatting for more effective OSINT reconnaissance

![](https://i.imgur.com/bWx79dy.png)

## Formats

The tool only accepts E164 and International formats as input.

## Installation

```bash
git clone https://github.com/sundowndev/PhoneInfoga
cd PhoneInfoga/
python3 -m pip install -r requirements.txt
```

## Usage

### [The full usage documentation has been moved to the wiki](https://github.com/sundowndev/PhoneInfoga/wiki)

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

## License

This tool is licensed under the GNU General Public License v3.0.

[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fsundowndev%2FPhoneInfoga.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fsundowndev%2FPhoneInfoga?ref=badge_large)
