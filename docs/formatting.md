# Formatting phone numbers

## Basics

The tool only accepts E164 and International formats as input.

- E164: +3396360XXXX
- International: +33 9 63 60 XX XX
- National: 09 63 60 XX XX
- RFC3966: tel:+33-9-63-60-XX-XX
- Out-of-country format from US: 011 33 9 63 60 XX XX

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


## Custom formatting

Sometimes the phone number has footprints but is used with a different formatting. This is a problem because for example if we search for "+15417543010", we'll not find web pages that write it that way : "(541) 754–3010". So the tool use a (optional) custom formatting given by the user to find further and more accurate results. To use this feature properly and make the results more valuable, try to use a format that someone of the number' country would usually use to share the phone number online.

For example, French people usually write numbers that way online : `06.20.30.40.50` or `06 20 30 40 50`.

For US-based numbers, the most common format is : `543-456-1234`.

### Examples

Here are some examples of custom formatting for US-based phone numbers : 

- `+1 618-555-xxxx`
- `(+1)618-555-xxxx`
- `+1/618-555-xxxx`
- `(618) 555xxxx`
- `(618) 555-xxxx`
- `(618) 555.xxxx`
- `(618)555xxxx`
- `(618)555-xxxx`
- `(618)555.xxxx`

For European countries (France as example) : 

- `+3301 86 48 xx xx`
- `+33018648xxxx`
- `+33018 648 xxx x`
- `(0033)018648xxxx`
- `(+33)018 648 xxx x`
- `+33/018648xxxx`
- `(0033)018 648 xxx x`
- `+33018-648-xxx-x`
- `(+33)018648xxxx`
- `(+33)01 86 48 xx xx`
- `+33/018-648-xxx-x`
- `+33/01-86-48-xx-xx`
- `+3301-86-48-xx-xx`
- `(0033)01 86 48 xx xx`
- `+33/01 86 48 xx xx`
- `(+33)018-648-xxx-x`
- `(+33)01-86-48-xx-xx`
- `(0033)01-86-48-xx-xx`
- `(0033)018-648-xxx-x`
- `+33/018 648 xxx x`

## Local scan formatting (developers)

The local scan create several variables to be usable in OSINT footprinting and other scanners.

Examples : 

```
{
    'input': '+1 512-618-xxxx',
    'default': '1512618xxxx',
    'local': '512618xxxx',
    'international': '+1 512-618-xxxx',
    'country': 'United States',
    'countryCode': '+1',
    'countryIsoCode': 'US',
    'location': 'Texas',
    'carrier': ''
}
```

```
{
    'input': '+33 066651xxxx',
    'default': '3366651xxxx',
    'local': '66651xxxx',
    'international': '+33 66651xxxx',
    'country': 'France',
    'countryCode': '+33',
    'countryIsoCode': 'FR',
    'location': 'France',
    'carrier': 'Bouygues'
}
```