Here is the documentation for CLI usage.

```shell
$ phoneinfoga

PhoneInfoga is one of the most advanced tools to scan phone numbers using only free resources.

Usage:
  phoneinfoga [command]

Examples:
phoneinfoga scan -n <number>

Available Commands:
  help        Help about any command
  scan        Scan a phone number
  serve       Serve web client
  version     Print current version of the tool

Flags:
  -h, --help   help for phoneinfoga

Use "phoneinfoga [command] --help" for more information about a command.
```

### Running a scan

Use the `scan` command with the `-n` (or `--number`) option.

```
phoneinfoga scan -n "+1 (555) 444-1212"
phoneinfoga scan -n "+33 06 79368229"
phoneinfoga scan -n "33679368229"
```

Special chars such as `( ) - +` will be escaped so typing US-based numbers stay easy : 

```
phoneinfoga scan -n "+1 555-444-3333"
```

!!! note "Note that the country code is essential. You don't know which country code to use ? [Find it here](https://www.countrycode.org/)"

<!--
#### Input & output file

Check several numbers at once and send results to a file.

```
phoneinfoga scan -i numbers.txt -o results.txt
```

Input file must contain one phone number per line. Invalid numbers will be skipped.

#### Footprinting

```
phoneinfoga scan -n +42837544833 -s footprints
```

#### Custom format reconnaissance

You don't know where to search and what custom format to use ? Let the tool try several custom formats based on the country code for you.

```
phoneinfoga recon -n +42837544833 
```
-->

## Available scanners

- Numverify
- Google search
- OVH

**Numverify** provide standard but useful informations such as number's country code, location, line type and carrier.

**OVH** is, besides being a web and cloud hosting company, a telecom provider with several VoIP numbers in Europe. Thanks to their API-key free [REST API](https://api.ovh.com/), we are able to tell if a number is owned by OVH Telecom or not.

**Google search** uses Google search engine and [Google Dorks](https://en.wikipedia.org/wiki/Google_hacking) to search phone number's footprints everywhere on the web. It allows you to search for scam reports, social media profiles, documents and more. **This scanner does only one thing:** generating several Google search links from a given phone number. You then have to manually open them in your browser to see results. You may therefore have links that do not return any results.

## Launch web client & REST API

Run the tool through a REST API with a web client. The API has been written in Go and web client in Vue.js.

```shell
phoneinfoga serve
phoneinfoga serve -p 8080 # default port is 5000
```

You should then be able to see the web client at `http://localhost:<port>`.

![](./images/screenshot.png)

### Run the REST API only

You can choose to only run the REST API without the web client :

```
phoneinfoga serve --no-client
```
