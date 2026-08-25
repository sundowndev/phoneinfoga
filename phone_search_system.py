#!/usr/bin/env python3
"""
Enhanced Phone Number Information Search System
Comprehensive phone number analysis and information gathering
"""

import re
import sqlite3
import requests
import time
from urllib.parse import quote
from flask import Flask, render_template, request, jsonify
from phonenumbers import parse, NumberParseException, is_valid_number, geocoder, carrier
from phonenumbers.phonenumberutil import PhoneNumberFormat
import phonenumbers
from typing import Dict, Any, List, Optional
import logging

# Local breach database (redacted output only)
from data_breaches import DataBreachesParser

# Configure logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

class PhoneSearchSystem:
    """Enhanced phone number search and analysis system"""
    
    def __init__(self):
        self.session = requests.Session()
        self.session.headers.update({
            'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36'
        })
        
        self.search_engines = {
            'google': self._google_search,
            'bing': self._bing_search,
            'duckduckgo': self._duckduckgo_search,
            'yandex': self._yandex_search
        }
        
        self.social_platforms = {
            'facebook': self._facebook_search,
            'instagram': self._instagram_search,
            'twitter': self._twitter_search,
            'linkedin': self._linkedin_search,
            'telegram': self._telegram_search,
            'whatsapp': self._whatsapp_search,
            'vk': self._vk_search,
            'ok': self._ok_search
        }
        
        self.api_services = {
            'numverify': self._numverify_api,
            'abstract_api': self._abstract_api,
            'ipapi': self._ipapi_lookup
        }

        # Local breach database (return only aggregated/redacted output)
        self.data_breaches = DataBreachesParser()

    def _data_breaches_search(self, phone: str) -> Dict[str, Any]:
        """Search phone number in local breach database (redacted output only)."""
        try:
            breach_results = self.data_breaches.search_by_phone(phone)
            if breach_results.get('found'):
                summary = breach_results.get('summary', {}) or {}
                return {
                    'service': 'Data Breaches Database',
                    'description': 'Search through local breach databases (redacted output)',
                    'found': True,
                    'matches': int(breach_results.get('matches', 0) or 0),
                    'data': [],
                    'data_redacted': True,
                    'redaction_note': 'Sensitive personal data has been redacted',
                    'summary': summary,
                    'risk_assessment': {
                        'highest_risk': summary.get('highest_risk', 'LOW'),
                        'total_exposed': int(breach_results.get('matches', 0) or 0),
                        'platforms_affected': summary.get('platforms_affected', [])
                    }
                }
            return {
                'service': 'Data Breaches Database',
                'description': 'Search through local breach databases (redacted output)',
                'found': False,
                'matches': 0,
                'message': 'No records found in breach databases'
            }
        except (ValueError, TypeError, sqlite3.Error, OSError) as e:
            return {
                'service': 'Data Breaches Database',
                'found': False,
                'error': str(e)
            }
    
    def validate_and_parse(self, phone_number: str) -> Optional[phonenumbers.PhoneNumber]:
        """Enhanced phone number validation and parsing"""
        # Clean the input
        phone_number = re.sub(r'[^\d+]', '', phone_number)
        
        # Try to parse with default region
        try:
            parsed = parse(phone_number, None)
            if is_valid_number(parsed):
                return parsed
        except NumberParseException:
            pass
        
        # Try with common regions
        for region in ['US', 'GB', 'FR', 'DE', 'RU', 'IN', 'BR', 'CN', 'JP']:
            try:
                parsed = parse(phone_number, region)
                if is_valid_number(parsed):
                    return parsed
            except NumberParseException:
                continue
        
        return None
    
    def get_basic_info(self, phone_number: phonenumbers.PhoneNumber) -> Dict[str, Any]:
        """Get comprehensive basic phone information"""
        return {
            'country_code': phone_number.country_code,
            'national_number': phone_number.national_number,
            'international_format': phonenumbers.format_number(phone_number, PhoneNumberFormat.INTERNATIONAL),
            'national_format': phonenumbers.format_number(phone_number, PhoneNumberFormat.NATIONAL),
            'e164_format': phonenumbers.format_number(phone_number, PhoneNumberFormat.E164),
            'is_valid': is_valid_number(phone_number),
            'country': geocoder.description_for_number(phone_number, 'en'),
            'carrier': carrier.name_for_number(phone_number, 'en'),
            'type': self._get_phone_type(phone_number),
            'timezone': self._get_timezone(phone_number)
        }
    
    def _get_phone_type(self, phone_number: phonenumbers.PhoneNumber) -> str:
        """Determine phone type (mobile, landline, etc.)"""
        # This is a simplified version - in real implementation you'd use more sophisticated logic
        if phone_number.country_code in [1, 44, 49, 33, 39]:  # Common country codes
            # Simple heuristic for mobile numbers
            if phone_number.national_number >= 7000000000:  # Rough mobile range
                return "Mobile"
            else:
                return "Landline"
        return "Unknown"
    
    def _get_timezone(self, phone_number: phonenumbers.PhoneNumber) -> str:
        """Get timezone for phone number (simplified)"""
        country = geocoder.description_for_number(phone_number, 'en')
        timezones = {
            'United States': 'UTC-5 to UTC-8',
            'United Kingdom': 'UTC+0',
            'France': 'UTC+1',
            'Germany': 'UTC+1',
            'Russia': 'UTC+3 to UTC+12',
            'India': 'UTC+5:30',
            'Brazil': 'UTC-3',
            'China': 'UTC+8',
            'Japan': 'UTC+9'
        }
        return timezones.get(country, 'Unknown')
    
    def _google_search(self, phone_number: str) -> Dict[str, Any]:
        """Google search for phone number"""
        formatted = phone_number
        search_queries = [
            f'"{formatted}"',
            f'"{formatted}" site:facebook.com',
            f'"{formatted}" site:linkedin.com',
            f'"{formatted}" site:twitter.com',
            f'"{formatted}" site:instagram.com',
            f'"{formatted}" site:telegram.me',
            f'"{formatted}" site:wa.me',
            f'"{formatted}" "phone number"',
            f'"{formatted}" "contact"',
            f'"{formatted}" "business"'
        ]
        
        return {
            'engine': 'Google',
            'search_queries': search_queries,
            'google_search_url': f'https://www.google.com/search?q={quote(formatted)}',
            'advanced_dorks': [
                f'intext:"{formatted}"',
                f'inurl:"{formatted}"',
                f'filetype:pdf "{formatted}"',
                f'site:yellowpages.com "{formatted}"',
                f'site:whatsapp.com "{formatted}"'
            ]
        }
    
    def _bing_search(self, phone_number: str) -> Dict[str, Any]:
        """Bing search for phone number"""
        formatted = phone_number
        return {
            'engine': 'Bing',
            'search_url': f'https://www.bing.com/search?q={quote(formatted)}',
            'queries': [
                f'"{formatted}"',
                f'"{formatted}" contact',
                f'"{formatted}" business'
            ]
        }
    
    def _duckduckgo_search(self, phone_number: str) -> Dict[str, Any]:
        """DuckDuckGo search for phone number"""
        formatted = phone_number
        return {
            'engine': 'DuckDuckGo',
            'search_url': f'https://duckduckgo.com/?q={quote(formatted)}',
            'privacy_focused': True,
            'queries': [f'"{formatted}"']
        }
    
    def _yandex_search(self, phone_number: str) -> Dict[str, Any]:
        """Yandex search for phone number"""
        formatted = phone_number
        return {
            'engine': 'Yandex',
            'search_url': f'https://yandex.com/search/?text={quote(formatted)}',
            'region_specific': True,
            'queries': [f'"{formatted}"']
        }
    
    def _facebook_search(self, phone_number: str) -> Dict[str, Any]:
        """Enhanced Facebook search for phone number"""
        formatted = phone_number
        clean_number = formatted.replace('+', '').replace('-', '').replace(' ', '')
        return {
            'platform': 'Facebook',
            'search_url': f'https://www.facebook.com/search/people/?q={quote(formatted)}',
            'direct_url': f'https://www.facebook.com/{quote(clean_number)}',
            'messenger_url': f'https://m.me/{quote(clean_number)}',
            'mobile_app_search': f'fb://search/people?q={quote(formatted)}',
            'directory_search': 'https://www.facebook.com/directory/people/',
            'advanced_search': f'https://www.facebook.com/search/people/?q={quote(formatted)}&filters=eyJyZWFyY2hfb3B0aW9ucyI6eyJoYXNfdmFsdWUiOnsiY2xhc3NfbmFtZSI6IlVzZXJzIiwic3RyaW5nX3ZhbHVlIjoiVXNlcnMifX19fX199fQ==',
            'phone_lookup': f'https://www.facebook.com/login/identify?ctx=rec&search_attempt={quote(formatted)}',
            'note': 'Поиск по номеру телефона может требовать авторизации',
            'privacy_note': 'Facebook ограничивает поиск по номерам из-за политики конфиденциальности',
            'search_tips': [
                'Попробуйте разные форматы номера',
                'Используйте расширенный поиск с фильтрами',
                'Проверьте поиск в группах и на страницах',
                'Попробуйте поиск по первым цифрам номера'
            ]
        }
    
    def _instagram_search(self, phone_number: str) -> Dict[str, Any]:
        """Enhanced Instagram search for phone number"""
        formatted = phone_number
        clean_number = formatted.replace('+', '').replace('-', '').replace(' ', '')
        return {
            'platform': 'Instagram',
            'search_url': f'https://www.instagram.com/explore/tags/{quote(formatted)}/',
            'direct_url': f'https://www.instagram.com/{quote(clean_number)}',
            'people_search': f'https://www.instagram.com/search/people/?q={quote(formatted)}',
            'top_search': f'https://www.instagram.com/search/top/?q={quote(formatted)}',
            'mobile_app_search': f'instagram://search?q={quote(formatted)}',
            'hashtag_search': f'https://www.instagram.com/explore/tags/{quote(formatted)}/',
            'contact_sync': 'https://www.instagram.com/accounts/contacts/',
            'note': 'Поиск по номеру телефона возможен через синхронизацию контактов',
            'privacy_note': 'Instagram не позволяет прямой поиск по номерам телефонов',
            'search_tips': [
                'Используйте синхронизацию контактов в настройках',
                'Попробуйте поиск по хэштегам с номером',
                'Проверьте возможные имена пользователей',
                'Ищите в связанных аккаунтах Facebook'
            ]
        }
    
    def _twitter_search(self, phone_number: str) -> Dict[str, Any]:
        """Twitter search for phone number"""
        formatted = phone_number
        return {
            'platform': 'Twitter',
            'search_url': f'https://twitter.com/search?q={quote(formatted)}',
            'advanced_search': f'https://twitter.com/search-advanced?all={quote(formatted)}',
            'note': 'Search may require login'
        }
    
    def _linkedin_search(self, phone_number: str) -> Dict[str, Any]:
        """LinkedIn search for phone number"""
        formatted = phone_number
        return {
            'platform': 'LinkedIn',
            'search_url': f'https://www.linkedin.com/search/results/all/?keywords={quote(formatted)}',
            'note': 'Professional network search'
        }
    
    def _telegram_search(self, phone_number: str) -> Dict[str, Any]:
        """Telegram search for phone number"""
        formatted = phone_number
        return {
            'platform': 'Telegram',
            'direct_url': f'https://t.me/{quote(formatted)}',
            'search_note': 'Telegram usernames are different from phone numbers',
            'privacy_note': 'Phone numbers are private on Telegram'
        }
    
    def _vk_search(self, phone_number: str) -> Dict[str, Any]:
        """VK.com search for phone number"""
        formatted = phone_number.replace('+', '').replace('-', '').replace(' ', '')
        return {
            'platform': 'VK.com',
            'search_url': f'https://vk.com/search?c[section]=people&c[q]={quote(formatted)}',
            'people_search': f'https://vk.com/search?c[section]=people&c[q]={quote(formatted)}',
            'phone_search': f'https://vk.com/search?c[section]=people&c[q]={quote(formatted)}',
            'advanced_search': f'https://vk.com/search?c[section]=people&c[q]={quote(formatted)}&c[country]=1',
            'mobile_app_search': f'vk://search/people?q={quote(formatted)}',
            'direct_phone_url': f'https://vk.com/restore?email={quote(formatted)}',
            'note': 'Поиск по номеру телефона может требовать авторизации',
            'privacy_note': 'VK скрывает номера телефонов из-за политики конфиденциальности',
            'search_tips': [
                'Попробуйте искать по формату без + и пробелов',
                'Используйте расширенный поиск с фильтрами',
                'Проверьте возможные вариации номера'
            ]
        }
    
    def _whatsapp_search(self, phone_number: str) -> Dict[str, Any]:
        """WhatsApp search for phone number"""
        formatted = phone_number
        return {
            'platform': 'WhatsApp',
            'chat_url': f'https://wa.me/{quote(formatted)}',
            'note': 'Direct WhatsApp chat link',
            'privacy_note': 'Phone number will be visible to recipient'
        }
    
    def _ok_search(self, phone_number: str) -> Dict[str, Any]:
        """OK.ru search for phone number"""
        formatted = phone_number.replace('+', '').replace('-', '').replace(' ', '')
        return {
            'platform': 'OK.ru',
            'search_url': f'https://ok.ru/search?st.query={quote(formatted)}',
            'people_search': f'https://ok.ru/search?st.mode=Users&st.query={quote(formatted)}',
            'groups_search': f'https://ok.ru/search?st.mode=Groups&st.query={quote(formatted)}',
            'advanced_search': f'https://ok.ru/search?st.mode=Users&st.query={quote(formatted)}&st.country=1',
            'mobile_app_search': f'okmobile://search?query={quote(formatted)}',
            'note': 'Поиск по номеру телефона может требовать авторизации',
            'privacy_note': 'OK.ru скрывает номера телефонов из-за политики конфиденциальности',
            'search_tips': [
                'Попробуйте искать по формату без + и пробелов',
                'Используйте фильтры по региону и возрасту',
                'Проверьте поиск в группах и темах'
            ]
        }
    
    def _numverify_api(self, phone_number: str) -> Dict[str, Any]:
        """Numverify API lookup"""
        return {
            'service': 'Numverify',
            'requires_api_key': True,
            'api_url': 'http://apilayer.net/api/validate',
            'api_key_required': 'Get API key from https://numverify.com/',
            'sample_request': f'GET /api/validate?access_key=YOUR_API_KEY&number={phone_number}',
            'response_fields': [
                'valid', 'number', 'local_format', 'international_format',
                'country_prefix', 'country_code', 'country_name', 'location',
                'carrier', 'line_type'
            ]
        }
    
    def _abstract_api(self, phone_number: str) -> Dict[str, Any]:
        """Abstract API phone validation"""
        return {
            'service': 'Abstract API',
            'requires_api_key': True,
            'api_url': 'https://phonevalidation.abstractapi.com/v1/',
            'api_key_required': 'Get API key from https://app.abstractapi.com/api/phone-validation/',
            'sample_request': f'GET /v1/?api_key=YOUR_API_KEY&phone={phone_number}',
            'features': ['Phone validation', 'Carrier info', 'Location data', 'Line type detection']
        }
    
    def _ipapi_lookup(self, phone_number: str) -> Dict[str, Any]:
        """IPAPI phone lookup"""
        return {
            'service': 'IPAPI',
            'requires_api_key': True,
            'api_url': 'https://ipapi.com/phone_api.json',
            'api_key_required': 'Get API key from https://ipapi.com/',
            'sample_request': f'GET /phone_api.json?apikey=YOUR_API_KEY&phone={phone_number}',
            'features': ['Phone validation', 'Carrier detection', 'Timezone info']
        }
    
    def comprehensive_search(self, phone_number: str, search_types: List[str] = None) -> Dict[str, Any]:
        """Comprehensive phone number search"""
        if search_types is None:
            search_types = ['basic', 'google', 'social']
        
        # Validate phone number
        parsed = self.validate_and_parse(phone_number)
        if not parsed:
            return {
                'error': 'Invalid phone number',
                'input': phone_number,
                'valid': False,
                'suggestions': [
                    'Include country code (e.g., +1234567890)',
                    'Check for typos',
                    'Ensure number format is correct'
                ]
            }
        
        formatted = phonenumbers.format_number(parsed, PhoneNumberFormat.E164)
        results = {
            'input': phone_number,
            'formatted': formatted,
            'valid': True,
            'search_types': search_types,
            'timestamp': time.strftime('%Y-%m-%d %H:%M:%S'),
            'results': {}
        }
        
        # Basic information
        if 'basic' in search_types or 'all' in search_types:
            results['results']['basic'] = self.get_basic_info(parsed)
        
        # Search engines
        if 'google' in search_types or 'search' in search_types or 'all' in search_types:
            results['results']['search_engines'] = {}
            for engine in ['google', 'bing', 'duckduckgo', 'yandex']:
                results['results']['search_engines'][engine] = self.search_engines[engine](formatted)
        
        # Social platforms
        if 'social' in search_types or 'all' in search_types:
            results['results']['social_platforms'] = {}
            for platform in ['facebook', 'instagram', 'twitter', 'linkedin', 'telegram', 'whatsapp', 'vk', 'ok']:
                results['results']['social_platforms'][platform] = self.social_platforms[platform](formatted)
        
        # API services
        if 'api' in search_types or 'all' in search_types:
            results['results']['api_services'] = {}
            for service in ['numverify', 'abstract_api', 'ipapi']:
                results['results']['api_services'][service] = self.api_services[service](formatted)

        # Local breach databases (redacted)
        if 'data_breaches' in search_types or 'all' in search_types:
            results['results']['data_breaches'] = self._data_breaches_search(formatted)
        
        return results

# Flask Web Application
app = Flask(__name__)
phone_search = PhoneSearchSystem()

@app.route('/')
def index():
    """Main web interface"""
    return render_template('phone_search.html')

@app.route('/api/search', methods=['POST'])
def api_search():
    """API endpoint for comprehensive phone search"""
    data = request.get_json()
    
    if not data or 'phone' not in data:
        return jsonify({'error': 'Phone number is required'}), 400
    
    phone_number = data['phone']
    search_types = data.get('search_types', ['basic', 'google', 'social'])
    
    result = phone_search.comprehensive_search(phone_number, search_types)
    return jsonify(result)

@app.route('/api/validate', methods=['POST'])
def api_validate():
    """API endpoint for phone number validation"""
    data = request.get_json()
    
    if not data or 'phone' not in data:
        return jsonify({'error': 'Phone number is required'}), 400
    
    phone_number = data['phone']
    parsed = phone_search.validate_and_parse(phone_number)
    
    if parsed:
        basic_info = phone_search.get_basic_info(parsed)
        return jsonify({
            'valid': True,
            'phone': phone_number,
            'info': basic_info
        })
    else:
        return jsonify({
            'valid': False,
            'phone': phone_number,
            'error': 'Invalid phone number format'
        })

@app.route('/api/search_engines', methods=['GET'])
def api_search_engines():
    """Get available search engines"""
    return jsonify({
        'search_engines': list(phone_search.search_engines.keys()),
        'social_platforms': list(phone_search.social_platforms.keys()),
        'api_services': list(phone_search.api_services.keys()),
        'local_databases': ['data_breaches']
    })

if __name__ == '__main__':
    print("Enhanced Phone Number Search System")
    print("==================================")
    print("Starting web server on http://localhost:5000")
    print("Available endpoints:")
    print("  POST /api/search - Comprehensive phone search")
    print("  POST /api/validate - Phone number validation")
    print("  GET /api/search_engines - Available search engines")
    print()
    
    app.run(host='0.0.0.0', port=5000, debug=True)
