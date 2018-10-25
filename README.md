# PhoneInfoga

Information gathering tool for phone numbers.

**This tool requires python 2.x**

## Features

- Check if phone number exists
- Gather standard informations such as country, line type and carrier

##### Up coming
- Check several numbers at once
- Set an output for result(s)
- Check if number is from a VoIP provider
- Get informations for special numbers (emergency)
- Phone book search

## Installation

```bash
git clone https://github.com/sundowndev/PhoneInfoga
cd ./PhoneInfoga
pip install -r requirements.txt
python phoneinfoga.py -h
```

## Usage

```
Usage: PhoneInfoga options 

       -n|--number: Phone number to search
       -h|--help: Help command
```

Example :

```
python phoneinfoga.py -n 447700900409
```

## License

This tool is MIT licensed.
