# PhoneInfoga

Advanced information gathering tool for phone numbers.

**This tool requires python 2.x**

## Features

- Check if phone number exists
- Gather standard informations such as country, line type and carrier
- Check several numbers at once
- Set an output for result(s)
- Check if number is from a VoIP provider
- Get informations about special numbers
- Phone book search

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
- any
- all

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
Usage: PhoneInfoga options 

       -n|--number: Phone number to search
       -i|--input: Phone number to search
       -o|--output: Phone number to search
       -s|--scanner: Only use a specific scanner
       -h|--help: Help command
       --update: Update the tool & databases
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