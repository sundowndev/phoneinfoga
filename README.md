# PhoneInfoga

Advanced information gathering tool & OSINT reconnaissance for phone numbers.

## The project

Building the most advanced tool to scan phone numbers using only free resources. The goal is to first identify basic informations such as country, area, carrier and line type on any international phone numbers with a very good accuracy, and then detect the VoIP provider or search for footprints on search engines to try identify the owner.

**This tool requires python 2.x**

## Features

- Check if phone number exists
- Gather standard informations such as country, line type and carrier
- Check several numbers at once
- Set an output for result(s)
- Check if number is from a VoIP provider
- OSINT reconnaissance using external APIs, Google Hacking, phone books & search engines

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
cd ./PhoneInfoga
pip install -r requirements.txt
python ./phoneinfoga.py -h
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

Example :

```
python phoneinfoga.py -n +42837544833
```

Check several numbers at once :

```
python ./phoneinfoga.py -i numbers.txt -o results.txt
```

Check for a number range on OVH (just put some zeros) :

```
python phoneinfoga.py -n +42837544833 -s ovh
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

![](https://i.imgur.com/WdXKSZY.png)

Here’s the same phone number in E.164 formatting: +442071838750

![](https://i.imgur.com/Ovso0w2.png)

## License

This tool is MIT licensed.

## Resources

Regular expression : `^[+]*[(]{0,1}[0-9]{1,4}[)]{0,1}[-\s\./0-9]*$`

### Docs

- http://whitepages.fr/phonesystem/
- https://support.twilio.com/hc/en-us/articles/223183008-Formatting-International-Phone-Numbers
- https://en.wikipedia.org/wiki/National_conventions_for_writing_telephone_numbers

### open data

- https://api.ovh.com/console/#/telephony
  - `/telephony/number/ranges`
  - `/telephony/number/detailedZones`
  - `/telephony/number/specificNumbers`
- https://countrycode.org/
- http://www.countryareacode.net/en/
- http://directory.didww.com/area-prefixes
- http://hs3x.com/

### Scanners

- https://www.phonevalidator.com/
- https://freecarrierlookup.com/
- https://www.411.com/
- https://www.washington.edu/home/peopledir/

### OSINT

- https://osintframework.com/

#### Google dork requests

- `insubject:"+XXXXXXXXX" | insubject:"+XXXXX" | insubject:"XXXXX XXX XXX`
- `insubject:"{number}" | intitle:"{number}"`
- `intext:"{number}" ext:doc | ext:docx | ext:odt | ext:pdf | ext:rtf | ext:sxw | ext:psw | ext:ppt | ext:pptx | ext:pps | ext:csv | ext:txt | ext:html`
