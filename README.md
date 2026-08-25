<p align="center">
  <img src="./docs/images/banner.png" width=500  alt="project logo"/>
</p>

<div align="center">
  <a href="https://github.com/sundowndev/phoneinfoga/actions">
    <img src="https://github.com/sundowndev/phoneinfoga/workflows/Build/badge.svg" alt="build status" />
  </a>
  <a href="https://goreportcard.com/report/github.com/sundowndev/phoneinfoga/v2">
    <img src="https://goreportcard.com/badge/github.com/sundowndev/phoneinfoga/v2" alt="go report" />
  </a>
  <a href="https://codeclimate.com/github/sundowndev/phoneinfoga/maintainability">
    <img src="https://api.codeclimate.com/v1/badges/3259feb1c68df1cd4f71/maintainability"  alt="code climate badge"/>
  </a>
  <a href='https://coveralls.io/github/sundowndev/phoneinfoga'>
    <img src='https://coveralls.io/repos/github/sundowndev/phoneinfoga/badge.svg' alt='Coverage Status' />
  </a>
  <a href="https://github.com/sundowndev/phoneinfoga/releases">
    <img src="https://img.shields.io/github/release/SundownDEV/phoneinfoga.svg" alt="Latest version" />
  </a>
  <a href="https://hub.docker.com/r/sundowndev/phoneinfoga">
    <img src="https://img.shields.io/docker/pulls/sundowndev/phoneinfoga.svg" alt="Docker pulls" />
  </a>
</div>

<h4 align="center">Information gathering framework for phone numbers</h4>

<p align="center">
  <a href="https://sundowndev.github.io/phoneinfoga/">Documentation</a> •
  <a href="https://petstore.swagger.io/?url=https://raw.githubusercontent.com/sundowndev/phoneinfoga/master/web/docs/swagger.yaml">API documentation</a> •
  <a href="https://medium.com/@SundownDEV/phone-number-scanning-osint-recon-tool-6ad8f0cac27b">Related blog post</a>
</p>

## About

PhoneInfoga is one of the most advanced tools to scan international phone numbers. It allows you to first gather basic information such as country, area, carrier and line type, then use various techniques to try to find the VoIP provider or identify the owner. It works with a collection of scanners that must be configured in order for the tool to be effective. PhoneInfoga doesn't automate everything, it's just there to help investigating on phone numbers.

## Current status

This project is stable but unmaintained. Upcoming bugs won't be fixed and repository could be archived at any time.

## Features

- Check if phone number exists
- Gather basic information such as country, line type and carrier
- OSINT footprinting using external APIs, phone books & search engines
- Check for reputation reports, social media, disposable numbers and more
- Use the graphical user interface to run scans from the browser
- Programmatic usage with the [REST API](https://petstore.swagger.io/?url=https://raw.githubusercontent.com/sundowndev/phoneinfoga/master/web/docs/swagger.yaml) and [Go modules](https://pkg.go.dev/github.com/sundowndev/phoneinfoga/v2)

## Telegram bot (workspace addition)

This workspace includes a lightweight Telegram bot in `telegram_bot.py` for quick phone checks.

### What it does

- Supports `/start`, `/help`, `/search <phone>`, `/ip <address>`, `/email <address>`
- Accepts plain text phone numbers in chat
- Uses local search engine from `universal_search_system.py`
- Returns basic phone info, local owner candidates, and redacted breach summary
- Supports optional chat allowlist via `TELEGRAM_ALLOWED_CHAT_IDS`

### Quick start

1. Install dependencies from `requirements.txt`
2. Set `TELEGRAM_BOT_TOKEN` in `.env`
   - Optional: set `TELEGRAM_ALLOWED_CHAT_IDS` as comma-separated Telegram chat IDs
3. Run the bot:
  - direct: `python telegram_bot.py`
  - integrated CLI: `python osint_cli.py telegram-bot`
  - make target: `make telegram-bot`

## Anti-features

- Does not claim to provide relevant or verified data, it's just a tool !
- Does not allow to "track" a phone or its owner in real time
- Does not allow to get the precise phone location
- Does not allow to hack a phone

## License

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fsundowndev%2FPhoneInfoga.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fsundowndev%2FPhoneInfoga?ref=badge_shield)

This tool is licensed under the GNU General Public License v3.0.

[Icon](https://www.flaticon.com/free-icon/fingerprint-search-symbol-of-secret-service-investigation_48838) made by <a href="https://www.freepik.com/" title="Freepik">Freepik</a> from <a href="https://www.flaticon.com/" title="Flaticon">flaticon.com</a> is licensed by <a href="http://creativecommons.org/licenses/by/3.0/" title="Creative Commons BY 3.0" target="_blank">CC 3.0 BY</a>.

## Support

Support me by signing up to DigitalOcean using my link ($200 free credits)

[![DigitalOcean Referral Badge](https://web-platforms.sfo2.cdn.digitaloceanspaces.com/WWW/Badge%203.svg)](https://www.digitalocean.com/?refcode=31f5ef768eb3&utm_campaign=Referral_Invite&utm_medium=Referral_Program&utm_source=badge)

<div align="center">
  <img src="https://github.com/sundowndev/static/raw/main/sponsors.svg?v=c68eba9" width="100%" heigh="auto" />
</div>
