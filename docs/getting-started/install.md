#!/usr/bin/env python3
"""
phone_lookup.py

Public, non-invasive phone number information lookup.

Requirements:
  pip install phonenumbers requests

Usage:
  - Fill in API credentials for Twilio or Numverify if you want to use them.
  - Run: python phone_lookup.py +8801793526778
"""

import os
import sys
import json
import phonenumbers
from phonenumbers import geocoder, carrier, timezone, NumberParseException
import requests
from typing import Optional

def basic_libphonenumber_lookup(number: str, default_region: Optional[str]=None) -> dict:
    """Parse a phone number with libphonenumber and return public metadata."""
    out = {"input": number}
    try:
        if default_region:
            pn = phonenumbers.parse(number, default_region)
        else:
            pn = phonenumbers.parse(number, None)
    except NumberParseException as e:
        out["error"] = f"parse error: {e}"
        return out

    out.update({
        "e164": phonenumbers.format_number(pn, phonenumbers.PhoneNumberFormat.E164),
        "international": phonenumbers.format_number(pn, phonenumbers.PhoneNumberFormat.INTERNATIONAL),
        "national": phonenumbers.format_number(pn, phonenumbers.PhoneNumberFormat.NATIONAL),
        "country_code": pn.country_code,
        "national_number": pn.national_number,
        "valid": phonenumbers.is_valid_number(pn),
        "possible": phonenumbers.is_possible_number(pn),
        "number_type": phonenumbers.number_type(pn).name,  # e.g., MOBILE, FIXED_LINE
        "region": geocoder.description_for_number(pn, "en"),  # rough region/country
        "carrier": carrier.name_for_number(pn, "en") or None,
        "time_zones": timezone.time_zones_for_number(pn) or [],
    })
    return out

def twilio_lookup(e164_number: str, account_sid: str, auth_token: str, types=("carrier", "caller-name")) -> dict:
    """
    Query Twilio Lookup API.
    Requires a Twilio account with Lookup product enabled. This is a paid endpoint for some data.
    Returns JSON or raises for network/errors.
    """
    base = f"https://lookups.twilio.com/v1/PhoneNumbers/{e164_number}"
    params = []
    for t in types:
        params.append(("Type", t))
    # requests will encode repeated params properly if given list of tuples
    resp = requests.get(base, params=params, auth=(account_sid, auth_token), timeout=10)
    resp.raise_for_status()
    return resp.json()

def numverify_lookup(number: str, access_key: str) -> dict:
    """
    Example for Numverify (apilayer). Replace endpoint and param names if service changes.
    Note: This API is paid/limited. Returns country_code, carrier (sometimes), line_type, etc.
    """
    url = "http://apilayer.net/api/validate"
    params = {
        "access_key": access_key,
        "number": number,
        "format": 1,
    }
    resp = requests.get(url, params=params, timeout=10)
    resp.raise_for_status()
    return resp.json()

def pretty_print(d: dict):
    print(json.dumps(d, indent=2, sort_keys=True, ensure_ascii=False))

def main():
    if len(sys.argv) < 2:
        print("Usage: python phone_lookup.py <phone-number> [default_region]")
        print("Example: python phone_lookup.py +8801793526778")
        sys.exit(1)

    number = sys.argv[1]
    default_region = sys.argv[2] if len(sys.argv) > 2 else None

    print("=== libphonenumber (local) lookup ===")
    local_info = basic_libphonenumber_lookup(number, default_region)
    pretty_print(local_info)

    # Optional: Twilio Lookup example (commented out unless you provide credentials)
    TWILIO_SID = os.getenv("TWILIO_ACCOUNT_SID")
    TWILIO_TOKEN = os.getenv("TWILIO_AUTH_TOKEN")
    if TWILIO_SID and TWILIO_TOKEN and local_info.get("e164"):
        try:
            print("\n=== Twilio Lookup (carrier & caller-name) ===")
            tw = twilio_lookup(local_info["e164"], TWILIO_SID, TWILIO_TOKEN, types=("carrier","caller-name"))
            pretty_print(tw)
        except Exception as e:
            print("Twilio lookup error:", e)
    else:
        print("\nSkipping Twilio lookup (set TWILIO_ACCOUNT_SID and TWILIO_AUTH_TOKEN env vars to enable)")

    # Optional: Numverify (apilayer) example
    NUMVERIFY_KEY = os.getenv("NUMVERIFY_ACCESS_KEY")
    if NUMVERIFY_KEY:
        try:
            print("\n=== Numverify Lookup ===")
            nv = numverify_lookup(number, NUMVERIFY_KEY)
            pretty_print(nv)
        except Exception as e:
            print("Numverify lookup error:", e)
    else:
        print("\nSkipping Numverify lookup (set NUMVERIFY_ACCESS_KEY env var to enable)")

if __name__ == "__main__":
    main()
