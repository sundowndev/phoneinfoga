```
$ python3 phoneinfoga.py -h
usage: phoneinfoga.py -n <number> [options]

Advanced information gathering tool for phone numbers
(https://github.com/sundowndev/PhoneInfoga) version v1.6.8

optional arguments:
  -h, --help            show this help message and exit
  -n number, --number number
                        The phone number to scan (E164 or international
                        format)
  -i input_file, --input input_file
                        Phone number list to scan (one per line)
  -o output_file, --output output_file
                        Output to save scan results
  -s scanner, --scanner scanner
                        The scanner to use
  --recon               Launch custom format reconnaissance
  --no-ansi             Disable colored output
  -v, --version         Show tool version
```

#### Basic scan

```
python3 phoneinfoga.py -n "(+42) 837544833"
```

Country code and special chars such as `( ) - +` will be escaped so typing US-based numbers stay easy : 

```
python3 phoneinfoga.py -n "+1 555-444-3333"
```

!!! note "Note that the country code is essential. You don't know which country code to use ? [Find it here](https://www.countrycode.org/)"

#### Output file

Check several numbers at once and send results to a file. Optionally, ensure no color code is used with `--no-ansi`

```
python3 phoneinfoga.py -i numbers.txt -o results.txt --no-ansi
```

Input file must contain one phone number per line. Invalid numbers will be skipped.

#### Footprinting

```
python3 phoneinfoga.py -n +42837544833 -s footprints
```

#### Custom format reconnaissance

You don't know where to search and what custom format to use ? Let the tool try several custom formats based on the country code for you.

```
python3 phoneinfoga.py -n +42837544833 -s any --recon
```

## Available scanners

Use `any` to disable this feature. Default value: `all`

- numverify
- ovh
- footprints

**Numverify** provide standard but useful informations such as number's country code, location, line type and carrier.

**OVH** is, besides being a web and cloud hosting company, a telecom provider with several VoIP numbers in the Europe. Thanks to their API-key free [REST API](https://api.ovh.com/), we are able to tell if a number is owned by OVH Telecom or not.

**Footprints** scanner uses Google search engine and [Google Dorks](https://en.wikipedia.org/wiki/Google_hacking) to search phone number's footprints everywhere on the web. It allows you to search for scam reports, social media profiles, documents and more.

## Examples

Check for a number range on OVH :

```
python3 phoneinfoga.py -n "+33 01 88 33 40 32" -s ovh
```

Output : 

```
[!] ---- Fetching informations for 330188334032 ---- [!]
[*] Running local scan...
[+] International format: +33 1 88 33 40 32
[+] Local format: 188334032
[+] Country found: France (+33)
[+] City/Area: France
[+] Carrier: 
[+] Timezone: Europe/Paris
[i] The number is valid and possible.
[*] Running OVH scan...
[+] 1 result found in OVH database
[+] Number range: 018833xxxx
[+] City: Paris
[+] Zip code: 
Continue scanning ? (y/N) 
[i] Good bye!
```