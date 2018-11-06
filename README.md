# PhoneInfoga

Advanced information gathering tool & OSINT reconnaissance for phone numbers.

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

## Number format by countries

#### Europe

- Belgium : 9 digits for land lines and 10 for mobile
- Denmark : 8 digits
- Germany : 10 digits
- Greece : 10 digits
- Hungary : 10 digits
- Iceland : 10 digits
- Ireland : 10 digits
- Italy : 10 digits
- Netherlands : 10 digits
- Norway : 10 digits
- Hungary : 10 digits

## Available scanners

- ovh
- annu
- numverify

## Installation

```bash
git clone https://github.com/sundowndev/PhoneInfoga
cd ./PhoneInfoga
pip install -r requirements.txt
python ./phoneinfoga.py -h
```

Then set APIs credentials in `secrets.py`.

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
python phoneinfoga.py -n 0428375448
```

Check several numbers at once :

```
python ./phoneinfoga.py -i numbers.txt -o results.txt
```

Check for a number range on OVH (just put some zeros) :

```
python phoneinfoga.py -n 0428370000 -s ovh
```

## License

This tool is MIT licensed.

## Resources

Regular expression : `^[+]*[(]{0,1}[0-9]{1,4}[)]{0,1}[-\s\./0-9]*$`

- https://api.ovh.com/console/#/telephony
- https://countrycode.org/
- http://www.countryareacode.net/en/
- http://whitepages.fr/phonesystem/
- http://directory.didww.com/area-prefixes

### Scanners
- https://www.phonevalidator.com/
- https://freecarrierlookup.com/
- https://www.411.com/

### OSINT
- https://osintframework.com/
