#!/usr/bin/env python3

# Read the file
with open('enhanced_universal_search.py', 'r', encoding='utf-8') as f:
    content = f.read()

# Replace the incorrect NumberFormat.E164 with PhoneNumberFormat.E164
content = content.replace('NumberFormat.E164', 'PhoneNumberFormat.E164')

# Write back to file
with open('enhanced_universal_search.py', 'w', encoding='utf-8') as f:
    f.write(content)

print("Fixed NumberFormat.E164 to PhoneNumberFormat.E164")
