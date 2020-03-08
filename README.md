<p align="center">
  <img src="https://i.imgur.com/LtUGnF3.png" width=500 />
</p>

<div align="center">
  <a href="https://github.com/sundowndev/PhoneInfoga/actions">
    <img src="https://img.shields.io/endpoint.svg?url=https://actions-badge.atrox.dev/sundowndev/PhoneInfoga/badge?ref=master&style=flat-square" alt="build status" />
  </a>
  <a href="https://goreportcard.com/report/github.com/sundowndev/PhoneInfoga">
    <img src="https://goreportcard.com/badge/github.com/sundowndev/PhoneInfoga" alt="go report" />
  </a>
  <a href="https://codeclimate.com/github/sundowndev/PhoneInfoga/maintainability">
    <img src="https://api.codeclimate.com/v1/badges/3259feb1c68df1cd4f71/maintainability" />
  </a>
  <a href="https://github.com/sundowndev/PhoneInfoga/releases">
    <img src="https://img.shields.io/github/release/SundownDEV/PhoneInfoga.svg?style=flat-square" alt="Latest version" />
  </a>
  <a href="https://github.com/sundowndev/PhoneInfoga/blob/master/LICENSE">
    <img src="https://img.shields.io/github/license/sundowndev/PhoneInfoga.svg?style=flat-square" alt="License" />
  </a>
</div>

<h4 align="center">Information gathering & OSINT reconnaissance tool for phone numbers</h4>

<p align="center">
  <a href="https://sundowndev.github.io/PhoneInfoga/">Documentation</a> •
  <a href="https://sundowndev.github.io/PhoneInfoga/usage/">Basic usage</a> •
  <a href="https://sundowndev.github.io/PhoneInfoga/resources/">OSINT resources</a> •
  <a href="https://medium.com/@SundownDEV/phone-number-scanning-osint-recon-tool-6ad8f0cac27b">Related blog post</a>
</p>

![](./docs/images/screenshot.png)

## About

PhoneInfoga is one of the most advanced tools to scan phone numbers using only free resources. The goal is to first gather standard information such as country, area, carrier and line type on any international phone numbers with a very good accuracy. Then search for footprints on search engines to try to find the VoIP provider or identify the owner.

## Features

- Check if phone number exists and is possible
- Gather standard informations such as country, line type and carrier
- OSINT footprinting using external APIs, Google Hacking, phone books & search engines
- Check for reputation reports, social media, disposable numbers and more
- Scan several numbers at once
- Use custom formatting for more effective OSINT reconnaissance
- **NEW**: Serve a web client along with a REST API to run scans from the browser
- **NEW**: Run your own web instance as a service
- **NEW**: Programmatic usage with Go modules

![Footprinting process](https://i.imgur.com/qCkgzz8.png)

## License

This tool is licensed under the GNU General Public License v3.0.

[Icon](https://www.flaticon.com/free-icon/fingerprint-search-symbol-of-secret-service-investigation_48838) made by <a href="https://www.freepik.com/" title="Freepik">Freepik</a> from <a href="https://www.flaticon.com/" title="Flaticon">flaticon.com</a> is licensed by <a href="http://creativecommons.org/licenses/by/3.0/" title="Creative Commons BY 3.0" target="_blank">CC 3.0 BY</a>.
