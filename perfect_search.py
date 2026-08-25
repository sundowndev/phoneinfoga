#!/usr/bin/env python3
"""
Perfect Universal Search System
Clean, efficient, and reliable OSINT tool
"""

import os
import json
import requests
import re
import logging
import time
import hashlib
from datetime import datetime, timedelta
from typing import Dict, Any, List, Optional, Tuple
from urllib.parse import quote
from functools import wraps

from flask import Flask, render_template, request, jsonify
from phonenumbers import parse, NumberParseException, is_valid_number, geocoder, carrier
from phonenumbers.phonenumberutil import PhoneNumberFormat
import phonenumbers

# Import data breaches parser
from data_breaches import DataBreachesParser
from sherlock_report import SherlockReportGenerator

# Configuration
class Config:
    CACHE_TTL = 3600  # 1 hour
    RATE_LIMIT = 100  # requests per hour
    MAX_CONTENT_LENGTH = 16 * 1024 * 1024  # 16MB

# Enhanced logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
    handlers=[
        logging.FileHandler('perfect_search.log'),
        logging.StreamHandler()
    ]
)
logger = logging.getLogger(__name__)

# Simple cache
class SimpleCache:
    def __init__(self):
        self.cache = {}
    
    def get(self, key: str) -> Optional[Dict[str, Any]]:
        if key in self.cache:
            data, timestamp = self.cache[key]
            if datetime.now().timestamp() - timestamp < Config.CACHE_TTL:
                return data
            else:
                del self.cache[key]
        return None
    
    def set(self, key: str, value: Dict[str, Any]) -> None:
        self.cache[key] = (value, datetime.now().timestamp())
    
    def clear(self) -> None:
        self.cache.clear()

# Rate limiting
class RateLimiter:
    def __init__(self):
        self.requests = {}
    
    def is_allowed(self, ip: str) -> Tuple[bool, Dict[str, Any]]:
        now = datetime.now().timestamp()
        hour_ago = now - 3600
        
        if ip not in self.requests:
            self.requests[ip] = []
        
        # Clean old requests
        self.requests[ip] = [req_time for req_time in self.requests[ip] if req_time > hour_ago]
        
        if len(self.requests[ip]) >= Config.RATE_LIMIT:
            return False, {
                'allowed': False,
                'limit': Config.RATE_LIMIT,
                'remaining': 0,
                'reset_time': int(min(self.requests[ip]) + 3600)
            }
        
        self.requests[ip].append(now)
        return True, {
            'allowed': True,
            'limit': Config.RATE_LIMIT,
            'remaining': Config.RATE_LIMIT - len(self.requests[ip]),
            'reset_time': int(now + 3600)
        }

# Input validator
class InputValidator:
    @staticmethod
    def validate_phone_number(phone: str) -> Tuple[bool, str, str]:
        """Validate and format phone number"""
        if not phone:
            return False, "Phone number is required", ""
        
        # Clean the phone number
        phone = re.sub(r'[^\d+]', '', phone)
        
        if not phone.startswith('+'):
            if phone.startswith('7') and len(phone) == 11:
                phone = '+' + phone
            elif phone.startswith('380') and len(phone) == 12:
                phone = '+' + phone
            else:
                return False, "Invalid phone number format. Use +7, +380, etc.", ""
        
        try:
            parsed = parse(phone)
            if not is_valid_number(parsed):
                return False, "Invalid phone number", phone
            
            return True, "Valid phone number", phone
            
        except NumberParseException:
            return False, "Invalid phone number format", ""
    
    @staticmethod
    def validate_fio(fio: str) -> Tuple[bool, str]:
        """Validate FIO (Full Name)"""
        if not fio:
            return False, "FIO is required"
        
        # Clean and validate
        fio = re.sub(r'[^\w\s]', '', fio).strip()
        words = fio.split()
        
        if len(words) < 2:
            return False, "Please provide at least first and last name"
        
        if len(fio) < 3:
            return False, "Name too short"
        
        return True, ' '.join(word.capitalize() for word in words)
    
    @staticmethod
    def validate_search_types(search_types: List[str]) -> Tuple[bool, str]:
        """Validate search types"""
        if not search_types:
            return False, "At least one search type must be specified"
        
        valid_types = {
            'basic', 'search_engines', 'social', 'owlsint', 'data_breaches', 'sherlock', 'all'
        }
        
        invalid_types = [t for t in search_types if t not in valid_types]
        if invalid_types:
            return False, f"Invalid search types: {', '.join(invalid_types)}"
        
        return True, "Valid search types"

# Perfect Search Engine
class PerfectSearch:
    def __init__(self):
        self.cache = SimpleCache()
        self.rate_limiter = RateLimiter()
        self.validator = InputValidator()
        self.data_breaches = DataBreachesParser()  # Initialize data breaches parser
        self.sherlock = SherlockReportGenerator()  # Initialize Sherlock report generator
        
        # Working search engines
        self.search_engines = {
            'google': lambda query: f"https://www.google.com/search?q={quote(query)}",
            'yandex': lambda query: f"https://yandex.com/search/?text={quote(query)}",
            'bing': lambda query: f"https://www.bing.com/search?q={quote(query)}",
            'duckduckgo': lambda query: f"https://duckduckgo.com/?q={quote(query)}"
        }
        
        # Working social platforms
        self.social_platforms = {
            'vk': lambda phone: f"https://vk.com/search?c[section]=people&c[q]={phone}",
            'telegram': lambda phone: f"https://t.me/{phone.replace('+', '')}",
            'whatsapp': lambda phone: f"https://wa.me/{phone.replace('+', '')}",
            'instagram': lambda phone: f"https://www.instagram.com/explore/tags/{phone.replace('+', '')}/",
            'facebook': lambda phone: f"https://www.facebook.com/search/people/?q={phone}",
            'linkedin': lambda phone: f"https://www.linkedin.com/search/results/all/?keywords={phone}"
        }
    
    def _get_client_ip(self) -> str:
        """Get client IP address"""
        if request.environ.get('HTTP_X_FORWARDED_FOR'):
            return request.environ['HTTP_X_FORWARDED_FOR'].split(',')[0]
        return request.environ.get('REMOTE_ADDR', 'unknown')
    
    def _generate_cache_key(self, search_type: str, query: str, params: List[str] = None) -> str:
        """Generate cache key"""
        key_data = f"{search_type}:{query}:{':'.join(params) if params else ''}"
        return hashlib.md5(key_data.encode()).hexdigest()
    
    def _basic_phone_info(self, phone: str) -> Dict[str, Any]:
        """Get basic phone information"""
        try:
            parsed = parse(phone)
            
            return {
                'valid': True,
                'e164_format': phonenumbers.format_number(parsed, PhoneNumberFormat.E164),
                'international_format': phonenumbers.format_number(parsed, PhoneNumberFormat.INTERNATIONAL),
                'national_format': phonenumbers.format_number(parsed, PhoneNumberFormat.NATIONAL),
                'country_code': phonenumbers.region_code_for_number(parsed),
                'country': geocoder.description_for_number(parsed, "en"),
                'carrier': carrier.name_for_number(parsed, "en"),
                'is_mobile': phonenumbers.number_type(parsed) == 1,
                'is_valid': phonenumbers.is_valid_number(parsed)
            }
        except Exception as e:
            return {
                'valid': False,
                'error': str(e)
            }
    
    def _search_engines_info(self, phone: str) -> Dict[str, Any]:
        """Generate search engine URLs"""
        queries = [
            phone,
            f'"{phone}"',
            f'"{phone}" contact',
            f'"{phone}" owner',
            f'"{phone}" address'
        ]
        
        results = {}
        for engine_name, engine_func in self.search_engines.items():
            results[engine_name] = {
                'engine': engine_name.title(),
                'search_urls': [engine_func(query) for query in queries],
                'note': f'Search {engine_name} for phone number'
            }
        
        return results
    
    def _social_platforms_info(self, phone: str) -> Dict[str, Any]:
        """Generate social platform URLs"""
        results = {}
        for platform_name, platform_func in self.social_platforms.items():
            results[platform_name] = {
                'platform': platform_name.title(),
                'search_url': platform_func(phone),
                'note': f'Search {platform_name} for phone number'
            }
        
        return results
    
    def _owlsint_advanced_info(self, phone: str) -> Dict[str, Any]:
        """Owl-sint advanced phone tracking"""
        try:
            parsed = parse(phone)
            
            return {
                'service': 'Owl-sint Advanced Tracking',
                'version': '1.2',
                'phone_number': phone,
                'basic_info': {
                    'international_format': phonenumbers.normalize_digits_only(parsed),
                    'national_format': phonenumbers.national_significant_number(parsed),
                    'valid_number': phonenumbers.is_valid_number(parsed),
                    'country_code': phonenumbers.region_code_for_number(parsed),
                    'location': geocoder.description_for_number(parsed, "en"),
                    'carrier': carrier.name_for_number(parsed, "en"),
                    'is_mobile': phonenumbers.number_type(parsed) == 1
                },
                'social_links': {
                    'whatsapp': f"https://wa.me/{phonenumbers.normalize_digits_only(parsed)}",
                    'telegram': f"https://t.me/{phonenumbers.normalize_digits_only(parsed)}",
                    'viber': f"viber://chat?number=%2B{phonenumbers.normalize_digits_only(parsed)}",
                    'signal': f"signal://send?phone=%2B{phonenumbers.normalize_digits_only(parsed)}"
                },
                'validation_formats': {
                    'e164': phonenumbers.format_number(parsed, PhoneNumberFormat.E164),
                    'international': phonenumbers.format_number(parsed, PhoneNumberFormat.INTERNATIONAL),
                    'national': phonenumbers.format_number(parsed, PhoneNumberFormat.NATIONAL)
                },
                'success': True
            }
        except Exception as e:
            return {
                'service': 'Owl-sint Advanced Tracking',
                'error': str(e),
                'success': False
            }
    
    def _data_breaches_search(self, phone: str) -> Dict[str, Any]:
        """Search for phone number in data breaches database"""
        try:
            # Search in breach database
            breach_results = self.data_breaches.search_by_phone(phone)
            
            if breach_results['found']:
                return {
                    'service': 'Data Breaches Database',
                    'description': 'Search through leaked databases',
                    'found': True,
                    'matches': breach_results['matches'],
                    # Do not return raw breach records (PII). Keep only aggregated summary.
                    'data': [],
                    'data_redacted': True,
                    'redaction_note': 'Sensitive personal data has been redacted',
                    'summary': breach_results['summary'],
                    'risk_assessment': {
                        'highest_risk': breach_results['summary']['highest_risk'],
                        'total_exposed': breach_results['matches'],
                        'platforms_affected': breach_results['summary']['platforms_affected']
                    }
                }
            else:
                return {
                    'service': 'Data Breaches Database',
                    'description': 'Search through leaked databases',
                    'found': False,
                    'phone': phone,
                    'message': 'No records found in breach databases'
                }
                
        except Exception as e:
            return {
                'service': 'Data Breaches Database',
                'error': str(e),
                'found': False
            }
    
    def _data_breaches_fio_search(self, fio: str) -> Dict[str, Any]:
        """Search for FIO in data breaches database"""
        try:
            # Search in breach database
            breach_results = self.data_breaches.search_by_name(fio)
            
            if breach_results['found']:
                return {
                    'service': 'Data Breaches Database',
                    'description': 'Search through leaked databases',
                    'found': True,
                    'matches': breach_results['matches'],
                    # Do not return raw breach records (PII). Keep only aggregated summary.
                    'data': [],
                    'data_redacted': True,
                    'redaction_note': 'Sensitive personal data has been redacted',
                    'summary': breach_results['summary'],
                    'risk_assessment': {
                        'highest_risk': breach_results['summary']['highest_risk'],
                        'total_exposed': breach_results['matches'],
                        'platforms_affected': breach_results['summary']['platforms_affected']
                    }
                }
            else:
                return {
                    'service': 'Data Breaches Database',
                    'description': 'Search through leaked databases',
                    'found': False,
                    'fio': fio,
                    'message': 'No records found in breach databases'
                }
                
        except Exception as e:
            return {
                'service': 'Data Breaches Database',
                'error': str(e),
                'found': False
            }
    
    def _sherlock_report_search(self, phone: str) -> Dict[str, Any]:
        """Generate Sherlock-style detailed report"""
        try:
            # Generate complete Sherlock report
            sherlock_report = self.sherlock.generate_sherlock_report(phone, redact=True)
            
            # Check if we actually found any data
            has_profiles = sherlock_report['total_profiles'] > 0
            has_sections = len(sherlock_report['sections']) > 0
            
            return {
                'service': 'Sherlock Report Generator',
                'description': 'Complete detailed report in Sherlock format',
                'found': has_profiles or has_sections,
                'report': sherlock_report,
                'total_profiles': sherlock_report['total_profiles'],
                'total_sources': sherlock_report['total_sources'],
                'sections_count': len(sherlock_report['sections'])
            }
                
        except Exception as e:
            return {
                'service': 'Sherlock Report Generator',
                'error': str(e),
                'found': False
            }
    
    def _has_meaningful_results(self, results: Dict[str, Any]) -> bool:
        """Check if results contain meaningful information"""
        if not isinstance(results, dict):
            return False
        
        for key, value in results.items():
            if key == 'basic':
                if isinstance(value, dict) and value.get('valid', False):
                    return True
            elif key == 'sherlock_report':
                if isinstance(value, dict) and value.get('found', False):
                    return True
            elif key == 'data_breaches':
                if isinstance(value, dict) and value.get('found', False):
                    return True
            elif key in ['search_engines', 'social_platforms']:
                return True  # Always meaningful (URLs)
            elif key == 'owlsint':
                if isinstance(value, dict) and value.get('success', False):
                    return True
        
        return False
    
    def universal_phone_search(self, phone_number: str, search_types: List[str] = None) -> Dict[str, Any]:
        """Perfect universal phone search"""
        try:
            # Rate limiting
            ip = self._get_client_ip()
            allowed, rate_info = self.rate_limiter.is_allowed(ip)
            if not allowed:
                return {
                    'error': 'Rate limit exceeded',
                    'error_type': 'rate_limit',
                    'rate_limit': rate_info
                }
            
            # Input validation
            is_valid, message, formatted = self.validator.validate_phone_number(phone_number)
            if not is_valid:
                return {
                    'error': message,
                    'error_type': 'validation',
                    'input': phone_number,
                    'valid': False,
                    'rate_limit': rate_info
                }
            
            # Default search types
            if search_types is None:
                search_types = ['basic', 'search_engines', 'social']
            
            # Validate search types
            is_valid, message = self.validator.validate_search_types(search_types)
            if not is_valid:
                return {
                    'error': message,
                    'error_type': 'validation',
                    'input': phone_number,
                    'valid': True,
                    'rate_limit': rate_info
                }
            
            # Check cache
            cache_key = self._generate_cache_key('phone', formatted, search_types)
            cached_result = self.cache.get(cache_key)
            if cached_result:
                cached_result['cached'] = True
                cached_result['rate_limit'] = rate_info
                return cached_result
            
            # Perform search
            result = {
                'input': phone_number,
                'formatted': formatted,
                'valid': True,
                'search_types': search_types,
                'timestamp': datetime.now().isoformat(),
                'cached': False,
                'rate_limit': rate_info,
                'results': {}
            }
            
            # Basic info
            if 'basic' in search_types or 'all' in search_types:
                result['results']['basic'] = self._basic_phone_info(formatted)
            
            # Search engines
            if 'search_engines' in search_types or 'all' in search_types:
                result['results']['search_engines'] = self._search_engines_info(formatted)
            
            # Social platforms
            if 'social' in search_types or 'all' in search_types:
                result['results']['social_platforms'] = self._social_platforms_info(formatted)
            
            # Owl-sint advanced tracking
            if 'owlsint' in search_types or 'all' in search_types:
                result['results']['advanced_tracking'] = self._owlsint_advanced_info(formatted)
            
            # Data breaches search
            if 'data_breaches' in search_types or 'all' in search_types:
                result['results']['data_breaches'] = self._data_breaches_search(formatted)
            
            # Sherlock report
            if 'sherlock' in search_types or 'all' in search_types:
                result['results']['sherlock_report'] = self._sherlock_report_search(formatted)
            
            # Cache result
            self.cache.set(cache_key, result)
            
            # Return only if meaningful results exist
            if self._has_meaningful_results(result['results']):
                return result
            else:
                return {
                    'input': phone_number,
                    'formatted': formatted,
                    'valid': True,
                    'search_types': search_types,
                    'timestamp': datetime.now().isoformat(),
                    'cached': False,
                    'rate_limit': rate_info,
                    'results': {},
                    'message': 'No meaningful information found for this phone number',
                    'note': 'Search completed but no actual data was found in any source'
                }
        
        except Exception as e:
            logger.error(f"Unexpected error in phone search: {e}")
            return {
                'error': 'Internal server error',
                'error_type': 'internal',
                'input': phone_number
            }
    
    def universal_fio_search(self, fio: str, search_types: List[str] = None) -> Dict[str, Any]:
        """Perfect FIO search"""
        try:
            # Rate limiting
            ip = self._get_client_ip()
            allowed, rate_info = self.rate_limiter.is_allowed(ip)
            if not allowed:
                return {
                    'error': 'Rate limit exceeded',
                    'error_type': 'rate_limit',
                    'rate_limit': rate_info
                }
            
            # Input validation
            is_valid, message = self.validator.validate_fio(fio)
            if not is_valid:
                return {
                    'error': message,
                    'error_type': 'validation',
                    'input': fio,
                    'valid': False,
                    'rate_limit': rate_info
                }
            
            # Default search types
            if search_types is None:
                search_types = ['search_engines', 'social']
            
            # Validate search types
            is_valid, message = self.validator.validate_search_types(search_types)
            if not is_valid:
                return {
                    'error': message,
                    'error_type': 'validation',
                    'input': fio,
                    'valid': True,
                    'rate_limit': rate_info
                }
            
            # Check cache
            cache_key = self._generate_cache_key('fio', message, search_types)
            cached_result = self.cache.get(cache_key)
            if cached_result:
                cached_result['cached'] = True
                cached_result['rate_limit'] = rate_info
                return cached_result
            
            # Perform search
            result = {
                'input': fio,
                'cleaned_fio': message,
                'valid': True,
                'search_types': search_types,
                'timestamp': datetime.now().isoformat(),
                'cached': False,
                'rate_limit': rate_info,
                'results': {}
            }
            
            # Search engines
            if 'search_engines' in search_types or 'all' in search_types:
                queries = [
                    message,
                    f'"{message}"',
                    f'"{message}" контакт',
                    f'"{message}" телефон',
                    f'"{message}" адрес'
                ]
                
                search_results = {}
                for engine_name, engine_func in self.search_engines.items():
                    search_results[engine_name] = {
                        'engine': engine_name.title(),
                        'search_urls': [engine_func(query) for query in queries],
                        'note': f'Search {engine_name} for person'
                    }
                
                result['results']['search_engines'] = search_results
            
            # Social platforms
            if 'social' in search_types or 'all' in search_types:
                social_results = {}
                for platform_name, platform_func in self.social_platforms.items():
                    social_results[platform_name] = {
                        'platform': platform_name.title(),
                        'search_url': platform_func(message),
                        'note': f'Search {platform_name} for person'
                    }
                
                result['results']['social_platforms'] = social_results
            
            # Data breaches search
            if 'data_breaches' in search_types or 'all' in search_types:
                result['results']['data_breaches'] = self._data_breaches_fio_search(message)
            
            # Cache result
            self.cache.set(cache_key, result)
            
            # Return only if meaningful results exist
            if self._has_meaningful_results(result['results']):
                return result
            else:
                return {
                    'input': fio,
                    'cleaned_fio': message,
                    'valid': True,
                    'search_types': search_types,
                    'timestamp': datetime.now().isoformat(),
                    'cached': False,
                    'rate_limit': rate_info,
                    'results': {},
                    'message': 'No meaningful information found for this FIO',
                    'note': 'Search completed but no actual data was found in any source'
                }
        
        except Exception as e:
            logger.error(f"Unexpected error in FIO search: {e}")
            return {
                'error': 'Internal server error',
                'error_type': 'internal',
                'input': fio
            }

# Initialize Flask app
app = Flask(__name__)
app.config['MAX_CONTENT_LENGTH'] = Config.MAX_CONTENT_LENGTH

# Initialize search engine
perfect_search = PerfectSearch()

# Routes
@app.route('/')
def index():
    """Main page"""
    return render_template('perfect_search.html')

@app.route('/api/phone_search', methods=['POST'])
def api_phone_search():
    """Phone search API"""
    try:
        data = request.get_json()
        if not data:
            return jsonify({
                'error': 'No data provided',
                'error_type': 'validation'
            }), 400
        
        phone = data.get('phone', '').strip()
        search_types = data.get('search_types', None)
        
        if not phone:
            return jsonify({
                'error': 'Phone number is required',
                'error_type': 'validation'
            }), 400
        
        result = perfect_search.universal_phone_search(phone, search_types)
        
        if 'error' in result:
            status_code = 400 if result.get('error_type') == 'validation' else 500
            return jsonify(result), status_code
        
        return jsonify(result)
    
    except Exception as e:
        logger.error(f"Phone search API error: {e}")
        return jsonify({
            'error': 'Internal server error',
            'error_type': 'internal'
        }), 500

@app.route('/api/fio_search', methods=['POST'])
def api_fio_search():
    """FIO search API"""
    try:
        data = request.get_json()
        if not data:
            return jsonify({
                'error': 'No data provided',
                'error_type': 'validation'
            }), 400
        
        fio = data.get('fio', '').strip()
        search_types = data.get('search_types', None)
        
        if not fio:
            return jsonify({
                'error': 'FIO is required',
                'error_type': 'validation'
            }), 400
        
        result = perfect_search.universal_fio_search(fio, search_types)
        
        if 'error' in result:
            status_code = 400 if result.get('error_type') == 'validation' else 500
            return jsonify(result), status_code
        
        return jsonify(result)
    
    except Exception as e:
        logger.error(f"FIO search API error: {e}")
        return jsonify({
            'error': 'Internal server error',
            'error_type': 'internal'
        }), 500

@app.route('/api/sherlock/txt_report')
def api_sherlock_txt_report():
    """Generate Sherlock TXT report"""
    try:
        phone = request.args.get('phone', '').strip()
        if not phone:
            return jsonify({
                'error': 'Phone number is required',
                'error_type': 'validation'
            }), 400
        
        # Validate phone
        is_valid, message, formatted = perfect_search.validator.validate_phone_number(phone)
        if not is_valid:
            return jsonify({
                'error': message,
                'error_type': 'validation'
            }), 400
        
        # Generate TXT report
        txt_report = perfect_search.sherlock.generate_txt_report(formatted, redact=True)
        
        from flask import Response
        return Response(
            txt_report,
            mimetype='text/plain',
            headers={'Content-Disposition': f'attachment; filename=sherlock_report_{formatted}.txt'}
        )
        
    except Exception as e:
        logger.error(f"Sherlock TXT report API error: {e}")
        return jsonify({
            'error': 'Internal server error',
            'error_type': 'internal'
        }), 500

@app.route('/api/breaches/statistics')
def api_breaches_statistics():
    """Data breaches statistics"""
    try:
        stats = perfect_search.data_breaches.get_breach_statistics()
        return jsonify(stats)
    except Exception as e:
        logger.error(f"Breaches statistics API error: {e}")
        return jsonify({
            'error': 'Internal server error',
            'error_type': 'internal'
        }), 500

@app.route('/api/status')
def api_status():
    """System status"""
    return jsonify({
        'status': 'running',
        'version': '1.0.0',
        'timestamp': datetime.now().isoformat(),
        'cache_size': len(perfect_search.cache.cache),
        'rate_limit': Config.RATE_LIMIT,
        'cache_ttl': Config.CACHE_TTL
    })

@app.route('/api/cache/clear', methods=['POST'])
def api_cache_clear():
    """Clear cache"""
    perfect_search.cache.clear()
    return jsonify({
        'message': 'Cache cleared successfully',
        'timestamp': datetime.now().isoformat()
    })

if __name__ == '__main__':
    print("""
    ╔══════════════════════════════════════════════════════════════╗
    ║                Perfect Universal Search System v1.0.0              ║
    ║                                                              ║
    ║  Clean • Efficient • Reliable • Fast                              ║
    ║                                                              ║
    ║  Features:                                                    ║
    ║  • Phone number validation and analysis                           ║
    ║  • Search engine integration                                   ║
    ║  • Social platform search                                       ║
    ║  • Owl-sint advanced tracking                                  ║
    ║  • Smart caching and rate limiting                              ║
    ║  • Meaningful results filtering                                 ║
    ║                                                              ║
    ║  Server: http://localhost:5000                                 ║
    ╚══════════════════════════════════════════════════════════════╝
    """)
    
    app.run(host='0.0.0.0', port=5000, debug=False)
