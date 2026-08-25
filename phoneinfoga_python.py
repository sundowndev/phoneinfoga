#!/usr/bin/env python3
"""
PhoneInfoga Python Version
Information gathering framework for phone numbers
"""

from flask import Flask, render_template, request, jsonify
from phonenumbers import parse, NumberParseException, is_valid_number, geocoder, carrier
from phonenumbers.phonenumberutil import PhoneNumberFormat
import phonenumbers
from typing import Dict, Any, Optional

class PhoneInfoga:
    """Main PhoneInfoga class for phone number analysis"""
    
    def __init__(self):
        self.scanners = {
            'basic': self._basic_info,
            'google_search': self._google_search,
            'numverify': self._numverify_lookup,
            'ovh': self._ovh_lookup
        }
    
    def validate_number(self, phone_number: str) -> Optional[phonenumbers.PhoneNumber]:
        """Validate and parse phone number"""
        try:
            # Try to parse with default region
            parsed = parse(phone_number, None)
            if is_valid_number(parsed):
                return parsed
        except NumberParseException:
            pass
        
        # Try with common regions if default fails
        for region in ['US', 'GB', 'FR', 'DE', 'RU']:
            try:
                parsed = parse(phone_number, region)
                if is_valid_number(parsed):
                    return parsed
            except NumberParseException:
                continue
        
        return None
    
    def _basic_info(self, phone_number: phonenumbers.PhoneNumber) -> Dict[str, Any]:
        """Get basic phone number information"""
        info = {
            'country_code': phone_number.country_code,
            'national_number': phone_number.national_number,
            'international_format': phonenumbers.format_number(phone_number, PhoneNumberFormat.INTERNATIONAL),
            'national_format': phonenumbers.format_number(phone_number, PhoneNumberFormat.NATIONAL),
            'e164_format': phonenumbers.format_number(phone_number, PhoneNumberFormat.E164),
            'is_valid': is_valid_number(phone_number),
            'country': geocoder.description_for_number(phone_number, 'en'),
            'carrier': carrier.name_for_number(phone_number, 'en')
        }
        return info
    
    def _google_search(self, phone_number: phonenumbers.PhoneNumber) -> Dict[str, Any]:
        """Simulate Google search for phone number (placeholder)"""
        formatted = phonenumbers.format_number(phone_number, PhoneNumberFormat.E164)
        return {
            'google_dorks': [
                f'"{formatted}"',
                f'"{formatted}" site:facebook.com',
                f'"{formatted}" site:linkedin.com',
                f'"{formatted}" site:twitter.com',
                f'"{formatted}" site:instagram.com'
            ],
            'search_query': formatted
        }
    
    def _numverify_lookup(self, phone_number: phonenumbers.PhoneNumber) -> Dict[str, Any]:
        """Numverify API lookup (placeholder - requires API key)"""
        formatted = phonenumbers.format_number(phone_number, PhoneNumberFormat.E164)
        return {
            'service': 'Numverify',
            'requires_api_key': True,
            'endpoint': 'http://apilayer.net/api/validate',
            'phone': formatted,
            'note': 'Requires Numverify API key to work'
        }
    
    def _ovh_lookup(self, phone_number: phonenumbers.PhoneNumber) -> Dict[str, Any]:
        """OVH API lookup (placeholder)"""
        formatted = phonenumbers.format_number(phone_number, PhoneNumberFormat.E164)
        return {
            'service': 'OVH',
            'requires_api_key': True,
            'endpoint': 'https://api.ovh.com/1.0/telephony/number/detailedInfos',
            'phone': formatted,
            'note': 'Requires OVH API key to work'
        }
    
    def scan_number(self, phone_number: str, scanners: list = None) -> Dict[str, Any]:
        """Scan phone number with specified scanners"""
        if scanners is None:
            scanners = ['basic', 'google_search']
        
        # Validate phone number
        parsed = self.validate_number(phone_number)
        if not parsed:
            return {
                'error': 'Invalid phone number',
                'input': phone_number,
                'valid': False
            }
        
        results = {
            'input': phone_number,
            'valid': True,
            'scanners_used': scanners,
            'results': {}
        }
        
        # Run each scanner
        for scanner_name in scanners:
            if scanner_name in self.scanners:
                try:
                    results['results'][scanner_name] = self.scanners[scanner_name](parsed)
                except (ValueError, TypeError, RuntimeError) as e:
                    results['results'][scanner_name] = {
                        'error': str(e),
                        'status': 'failed'
                    }
        
        return results

# Flask Web Application
app = Flask(__name__)
phoneinfoga = PhoneInfoga()

@app.route('/')
def index():
    """Main web interface"""
    return render_template('index.html')

@app.route('/api/scan', methods=['POST'])
def api_scan():
    """API endpoint for scanning phone numbers"""
    data = request.get_json()
    
    if not data or 'phone' not in data:
        return jsonify({'error': 'Phone number is required'}), 400
    
    phone_number = data['phone']
    scanners = data.get('scanners', ['basic', 'google_search'])
    
    result = phoneinfoga.scan_number(phone_number, scanners)
    return jsonify(result)

@app.route('/api/validate', methods=['POST'])
def api_validate():
    """API endpoint for phone number validation"""
    data = request.get_json()
    
    if not data or 'phone' not in data:
        return jsonify({'error': 'Phone number is required'}), 400
    
    phone_number = data['phone']
    parsed = phoneinfoga.validate_number(phone_number)
    
    if parsed:
        return jsonify({
            'valid': True,
            'phone': phone_number,
            'international_format': phonenumbers.format_number(parsed, PhoneNumberFormat.INTERNATIONAL),
            'country': geocoder.description_for_number(parsed, 'en'),
            'carrier': carrier.name_for_number(parsed, 'en')
        })
    
    return jsonify({
        'valid': False,
        'phone': phone_number,
        'error': 'Invalid phone number format'
    })

if __name__ == '__main__':
    print("PhoneInfoga Python Version")
    print("==========================")
    print("Starting web server on http://localhost:5000")
    print("API endpoints:")
    print("  POST /api/scan - Scan phone number")
    print("  POST /api/validate - Validate phone number")
    print()
    
    app.run(host='0.0.0.0', port=5000, debug=True)
