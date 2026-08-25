#!/usr/bin/env python3
"""
Universal Search System
Comprehensive OSINT tool for phone numbers and photos with all possible sources
"""

import os
import sqlite3
import requests
from flask import Flask, render_template, request, jsonify
from werkzeug.utils import secure_filename
from PIL import Image, ExifTags
from phonenumbers import parse, NumberParseException, is_valid_number, geocoder, carrier
from phonenumbers.phonenumberutil import PhoneNumberFormat
import phonenumbers
from typing import Dict, Any, List, Optional
import logging
import time
import uuid
from urllib.parse import quote

# Local data-breaches database (redacted output only)
from data_breaches import DataBreachesParser
from directory_db import search_records, stats_by_city_and_category

# Optional .env support (do not hard-require python-dotenv)
try:
    from dotenv import load_dotenv  # type: ignore
except ImportError:
    load_dotenv = None  # type: ignore

if load_dotenv:
    # If python-dotenv is installed, automatically load local .env.
    load_dotenv()

# Safe subset of X-osint-like utilities
from xosint_toolkit import XOsintToolkit

# Configure logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

class UniversalSearchSystem:
    """Universal OSINT search system for phones and photos"""
    
    def __init__(self):
        self.session = requests.Session()
        self.session.headers.update({
            'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36'
        })
        
        # File upload settings
        self.upload_folder = 'uploads'
        self.allowed_extensions = {'png', 'jpg', 'jpeg', 'gif', 'bmp', 'webp'}
        self.max_file_size = 16 * 1024 * 1024  # 16MB
        
        if not os.path.exists(self.upload_folder):
            os.makedirs(self.upload_folder)
        
        # Phone search engines
        self.phone_search_engines = {
            'google': self._google_phone_search,
            'bing': self._bing_phone_search,
            'duckduckgo': self._duckduckgo_phone_search,
            'yandex': self._yandex_phone_search,
            'baidu': self._baidu_phone_search,
            'yahoo': self._yahoo_phone_search,
            'ask': self._ask_phone_search,
            'startpage': self._startpage_phone_search
        }
        
        # Social platforms (extended)
        self.social_platforms = {
            'facebook': self._facebook_search,
            'instagram': self._instagram_search,
            'twitter': self._twitter_search,
            'linkedin': self._linkedin_search,
            'telegram': self._telegram_search,
            'whatsapp': self._whatsapp_search,
            'vk': self._vk_search,
            'ok': self._ok_search,
            'tiktok': self._tiktok_search,
            'youtube': self._youtube_search,
            'reddit': self._reddit_search,
            'pinterest': self._pinterest_search,
            'snapchat': self._snapchat_search,
            'discord': self._discord_search
        }
        
        # Photo search engines
        self.photo_search_engines = {
            'google': self._google_photo_search,
            'yandex': self._yandex_photo_search,
            'bing': self._bing_photo_search,
            'tineye': self._tineye_search,
            'baidu': self._baidu_photo_search,
            'sogou': self._sogou_photo_search,
            'yandex_reverse': self._yandex_reverse_search,
            'iqdb': self._iqdb_search,
            'saucenao': self._saucenao_search
        }
        
        # API services
        self.api_services = {
            'numverify': self._numverify_api,
            'abstract_api': self._abstract_api,
            'ipapi': self._ipapi_lookup,
            'twilio': self._twilio_lookup,
            'infobel': self._infobel_lookup,
            'globalphone': self._globalphone_lookup
        }
        
        # Facial recognition services
        self.facial_services = {
            'face_recognition': self._face_recognition_analysis,
            'facepp': self._facepp_analysis,
            'kairos': self._kairos_analysis,
            'amazon_rekognition': self._amazon_rekognition_analysis,
            'azure_face': self._azure_face_analysis,
            'google_vision': self._google_vision_analysis
        }

        # Local breach database (NOTE: we return redacted/aggregated info only)
        self.data_breaches = DataBreachesParser()

        # Optional X-osint-style toolkit (safe subset)
        self.xosint = XOsintToolkit(self.session)

    def data_breaches_search(self, phone: str) -> Dict[str, Any]:
        """Search phone number in local breach database.

        Security/privacy: this method intentionally does NOT return raw leaked records (PII).
        It returns only aggregated summary and counts.
        """
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
        except (ValueError, TypeError, sqlite3.Error) as e:
            return {
                'service': 'Data Breaches Database',
                'found': False,
                'error': str(e)
            }

    def _data_breaches_search(self, phone: str) -> Dict[str, Any]:
        """Backward-compatible alias for internal/external callers."""
        return self.data_breaches_search(phone)
    
    # Phone search methods
    def _google_phone_search(self, phone_number: str) -> Dict[str, Any]:
        """Google phone search"""
        return {
            'engine': 'Google',
            'search_url': f'https://www.google.com/search?q={quote(phone_number)}',
            'phonebook_url': f'https://www.google.com/search?q={quote(phone_number)}&tbm=ppl',
            'maps_url': f'https://www.google.com/maps/search/{quote(phone_number)}',
            'advanced_dorks': [
                f'intext:"{phone_number}"',
                f'inurl:"{phone_number}"',
                f'filetype:pdf "{phone_number}"',
                f'site:linkedin.com "{phone_number}"',
                f'site:facebook.com "{phone_number}"',
                f'site:instagram.com "{phone_number}"',
                f'site:twitter.com "{phone_number}"',
                f'site:vk.com "{phone_number}"',
                f'site:ok.ru "{phone_number}"'
            ]
        }
    
    def _bing_phone_search(self, phone_number: str) -> Dict[str, Any]:
        """Bing phone search"""
        return {
            'engine': 'Bing',
            'search_url': f'https://www.bing.com/search?q={quote(phone_number)}',
            'people_url': f'https://www.bing.com/search?q={quote(phone_number)}&qs=n&form=QBRE',
            'advanced_query': f'"{phone_number}" contact business'
        }
    
    def _yandex_phone_search(self, phone_number: str) -> Dict[str, Any]:
        """Yandex phone search"""
        return {
            'engine': 'Yandex',
            'search_url': f'https://yandex.com/search/?text={quote(phone_number)}',
            'people_url': f'https://yandex.com/search/?text={quote(phone_number)}&lr=213',  # Moscow region
            'phonebook_url': f'https://yandex.com/search/?text={quote(phone_number)}&local=1'
        }
    
    def _baidu_phone_search(self, phone_number: str) -> Dict[str, Any]:
        """Baidu phone search"""
        return {
            'engine': 'Baidu',
            'search_url': f'https://www.baidu.com/s?wd={quote(phone_number)}',
            'note': 'Chinese search engine, good for Asian numbers'
        }
    
    def _yahoo_phone_search(self, phone_number: str) -> Dict[str, Any]:
        """Yahoo phone search"""
        return {
            'engine': 'Yahoo',
            'search_url': f'https://search.yahoo.com/search?p={quote(phone_number)}',
            'people_url': f'https://search.yahoo.com/search?p={quote(phone_number)}&fr=yfp-t'
        }
    
    def _duckduckgo_phone_search(self, phone_number: str) -> Dict[str, Any]:
        """DuckDuckGo phone search"""
        return {
            'engine': 'DuckDuckGo',
            'search_url': f'https://duckduckgo.com/?q={quote(phone_number)}',
            'privacy_focused': True
        }
    
    def _ask_phone_search(self, phone_number: str) -> Dict[str, Any]:
        """Ask.com phone search"""
        return {
            'engine': 'Ask.com',
            'search_url': f'https://www.ask.com/web?q={quote(phone_number)}',
            'note': 'Question-answer based search engine'
        }
    
    def _startpage_phone_search(self, phone_number: str) -> Dict[str, Any]:
        """Startpage phone search"""
        return {
            'engine': 'Startpage',
            'search_url': f'https://www.startpage.com/do/search?query={quote(phone_number)}',
            'privacy_focused': True
        }
    
    # Extended social media search
    def _tiktok_search(self, phone_number: str) -> Dict[str, Any]:
        """TikTok search"""
        return {
            'platform': 'TikTok',
            'search_url': f'https://www.tiktok.com/search?q={quote(phone_number)}',
            'hashtag_url': f'https://www.tiktok.com/tag/{quote(phone_number)}',
            'note': 'TikTok rarely shows phone numbers in search'
        }
    
    def _youtube_search(self, phone_number: str) -> Dict[str, Any]:
        """YouTube search"""
        return {
            'platform': 'YouTube',
            'search_url': f'https://www.youtube.com/results?search_query={quote(phone_number)}',
            'channel_url': f'https://www.youtube.com/search?q={quote(phone_number)}&sp=EgIQAgAUAFBABCAAYAFgAkgBEggBMAU%3D',
            'note': 'Search for videos and channels mentioning the number'
        }
    
    def _reddit_search(self, phone_number: str) -> Dict[str, Any]:
        """Reddit search"""
        return {
            'platform': 'Reddit',
            'search_url': f'https://www.reddit.com/search?q={quote(phone_number)}',
            'user_search': f'https://www.reddit.com/search?q={quote(phone_number)}&type=user',
            'post_search': f'https://www.reddit.com/search?q={quote(phone_number)}&type=link',
            'note': 'Forum discussions and user mentions'
        }
    
    def _pinterest_search(self, phone_number: str) -> Dict[str, Any]:
        """Pinterest search"""
        return {
            'platform': 'Pinterest',
            'search_url': f'https://www.pinterest.com/search/pins/?q={quote(phone_number)}',
            'people_url': f'https://www.pinterest.com/search/people/?q={quote(phone_number)}',
            'note': 'Visual search platform'
        }
    
    def _snapchat_search(self, phone_number: str) -> Dict[str, Any]:
        """Snapchat search"""
        return {
            'platform': 'Snapchat',
            'search_url': f'https://www.snapchat.com/add/{quote(phone_number)}',
            'note': 'Direct add by phone number (if user allows)'
        }
    
    def _discord_search(self, phone_number: str) -> Dict[str, Any]:
        """Discord search"""
        return {
            'platform': 'Discord',
            'search_url': f'https://discord.com/users/{quote(phone_number)}',
            'note': 'Direct user lookup by phone number (limited)'
        }
    
    # Photo search methods
    def _sogou_photo_search(self, image_path: str) -> Dict[str, Any]:
        """Sogou photo search"""
        _ = image_path
        return {
            'engine': 'Sogou Images',
            'search_url': 'https://pic.sogou.com/',
            'upload_url': 'https://pic.sogou.com/pics',
            'note': 'Chinese image search engine'
        }
    
    def _yandex_reverse_search(self, image_path: str) -> Dict[str, Any]:
        """Yandex reverse image search"""
        _ = image_path
        return {
            'engine': 'Yandex Reverse',
            'search_url': 'https://yandex.com/images/',
            'upload_url': 'https://yandex.com/images/app/?_url=',
            'features': ['Face detection', 'Object recognition', 'Similar images']
        }
    
    def _iqdb_search(self, image_path: str) -> Dict[str, Any]:
        """IQDB anime image search"""
        _ = image_path
        return {
            'engine': 'IQDB',
            'search_url': 'https://iqdb.org/',
            'upload_url': 'https://iqdb.org/',
            'note': 'Specialized for anime/manga images'
        }
    
    def _saucenao_search(self, image_path: str) -> Dict[str, Any]:
        """SauceNAO image search"""
        _ = image_path
        return {
            'engine': 'SauceNAO',
            'search_url': 'https://saucenao.com/',
            'upload_url': 'https://saucenao.com/search.php',
            'note': 'Specialized for anime/manga images'
        }
    
    # API services
    def _twilio_lookup(self, phone_number: str) -> Dict[str, Any]:
        """Twilio Lookup API"""
        _ = phone_number
        return {
            'service': 'Twilio Lookup',
            'requires_api_key': True,
            'api_url': 'https://lookups.twilio.com/v1/PhoneNumbers/',
            'api_key_required': 'Get API key from https://www.twilio.com/',
            'features': ['Carrier info', 'Type detection', 'Country code', 'Caller ID']
        }
    
    def _infobel_lookup(self, phone_number: str) -> Dict[str, Any]:
        """Infobel phone lookup"""
        _ = phone_number
        return {
            'service': 'Infobel',
            'requires_api_key': True,
            'api_url': 'https://api.infobel.com/v1/lookup/',
            'api_key_required': 'Get API key from https://www.infobel.com/',
            'features': ['International coverage', 'Business listings', 'Historical data']
        }
    
    def _globalphone_lookup(self, phone_number: str) -> Dict[str, Any]:
        """Global Phone lookup"""
        _ = phone_number
        return {
            'service': 'Global Phone',
            'requires_api_key': True,
            'api_url': 'https://api.globalphoneapi.com/lookup',
            'api_key_required': 'Get API key from https://globalphoneapi.com/',
            'features': ['Global coverage', 'Real-time data', 'Advanced filtering']
        }
    
    # Facial recognition services
    def _amazon_rekognition_analysis(self, image_path: str) -> Dict[str, Any]:
        """Amazon Rekognition analysis"""
        _ = image_path
        return {
            'service': 'Amazon Rekognition',
            'requires_api_key': True,
            'api_url': 'https://rekognition.amazonaws.com/',
            'api_key_required': 'Get AWS credentials from https://aws.amazon.com/',
            'features': [
                'Face detection',
                'Age estimation',
                'Gender detection',
                'Emotion analysis',
                'Celebrity recognition',
                'Face comparison'
            ],
            'note': 'AWS cloud-based facial analysis'
        }
    
    def _azure_face_analysis(self, image_path: str) -> Dict[str, Any]:
        """Azure Face API analysis"""
        _ = image_path
        return {
            'service': 'Azure Face API',
            'requires_api_key': True,
            'api_url': 'https://westcentralus.api.cognitive.microsoft.com/face/v1.0/',
            'api_key_required': 'Get API key from https://azure.microsoft.com/',
            'features': [
                'Face detection',
                'Age estimation',
                'Gender detection',
                'Emotion recognition',
                'Face landmarks',
                'Similar face matching'
            ],
            'note': 'Microsoft Azure cognitive services'
        }
    
    def _google_vision_analysis(self, image_path: str) -> Dict[str, Any]:
        """Google Vision API analysis"""
        _ = image_path
        return {
            'service': 'Google Vision API',
            'requires_api_key': True,
            'api_url': 'https://vision.googleapis.com/v1/',
            'api_key_required': 'Get API key from https://cloud.google.com/vision/',
            'features': [
                'Face detection',
                'Label detection',
                'Text extraction',
                'Logo detection',
                'Landmark recognition',
                'Web detection'
            ],
            'note': 'Google Cloud Vision AI services'
        }
    
    # Include previous methods (simplified versions)
    def _facebook_search(self, phone_number: str) -> Dict[str, Any]:
        """Facebook search"""
        return {
            'platform': 'Facebook',
            'search_url': f'https://www.facebook.com/search/people/?q={quote(phone_number)}',
            'messenger_url': f'https://m.me/{quote(phone_number)}',
            'note': 'May require login for full results'
        }
    
    def _instagram_search(self, phone_number: str) -> Dict[str, Any]:
        """Instagram search"""
        return {
            'platform': 'Instagram',
            'search_url': f'https://www.instagram.com/explore/tags/{quote(phone_number)}/',
            'note': 'Phone numbers rarely used as usernames'
        }
    
    def _twitter_search(self, phone_number: str) -> Dict[str, Any]:
        """Twitter search"""
        return {
            'platform': 'Twitter',
            'search_url': f'https://twitter.com/search?q={quote(phone_number)}',
            'note': 'Real-time search capabilities'
        }
    
    def _linkedin_search(self, phone_number: str) -> Dict[str, Any]:
        """LinkedIn search"""
        return {
            'platform': 'LinkedIn',
            'search_url': f'https://www.linkedin.com/search/results/all/?keywords={quote(phone_number)}',
            'note': 'Professional network search'
        }
    
    def _telegram_search(self, phone_number: str) -> Dict[str, Any]:
        """Telegram search"""
        return {
            'platform': 'Telegram',
            'direct_url': f'https://t.me/{quote(phone_number)}',
            'note': 'Phone numbers are private on Telegram'
        }
    
    def _whatsapp_search(self, phone_number: str) -> Dict[str, Any]:
        """WhatsApp search"""
        return {
            'platform': 'WhatsApp',
            'chat_url': f'https://wa.me/{quote(phone_number)}',
            'note': 'Direct WhatsApp chat link'
        }
    
    def _vk_search(self, phone_number: str) -> Dict[str, Any]:
        """VK.com search"""
        return {
            'platform': 'VK.com',
            'search_url': f'https://vk.com/search?c[section]=people&c[q]={quote(phone_number)}',
            'note': 'Russian social network'
        }
    
    def _ok_search(self, phone_number: str) -> Dict[str, Any]:
        """OK.ru search"""
        return {
            'platform': 'OK.ru',
            'search_url': f'https://ok.ru/search?st.query={quote(phone_number)}',
            'note': 'Russian social network'
        }
    
    def _google_photo_search(self, image_path: str) -> Dict[str, Any]:
        """Google photo search"""
        _ = image_path
        return {
            'engine': 'Google Images',
            'search_url': 'https://images.google.com/',
            'upload_url': 'https://images.google.com/searchbyimage/upload',
            'note': 'Upload image for reverse search'
        }
    
    def _bing_photo_search(self, image_path: str) -> Dict[str, Any]:
        """Bing photo search"""
        _ = image_path
        return {
            'engine': 'Bing Visual Search',
            'search_url': 'https://www.bing.com/visualsearch',
            'note': 'Microsoft visual search capabilities'
        }
    
    def _tineye_search(self, image_path: str) -> Dict[str, Any]:
        """TinEye photo search"""
        _ = image_path
        return {
            'engine': 'TinEye',
            'search_url': 'https://tineye.com/',
            'note': 'Best for finding exact image copies'
        }
    
    def _baidu_photo_search(self, image_path: str) -> Dict[str, Any]:
        """Baidu photo search"""
        _ = image_path
        return {
            'engine': 'Baidu Images',
            'search_url': 'https://image.baidu.com/',
            'note': 'Chinese image search engine'
        }
    
    def _yandex_photo_search(self, image_path: str) -> Dict[str, Any]:
        """Yandex photo search"""
        _ = image_path
        return {
            'engine': 'Yandex Images',
            'search_url': 'https://yandex.com/images/',
            'note': 'Russian image search with face detection'
        }
    
    def _numverify_api(self, phone_number: str) -> Dict[str, Any]:
        """Numverify API"""
        _ = phone_number
        return {
            'service': 'Numverify',
            'requires_api_key': True,
            'api_url': 'http://apilayer.net/api/validate',
            'note': 'Phone validation and carrier info'
        }
    
    def _abstract_api(self, phone_number: str) -> Dict[str, Any]:
        """Abstract API"""
        _ = phone_number
        return {
            'service': 'Abstract API',
            'requires_api_key': True,
            'api_url': 'https://phonevalidation.abstractapi.com/v1/',
            'note': 'Advanced phone validation'
        }
    
    def _ipapi_lookup(self, phone_number: str) -> Dict[str, Any]:
        """IPAPI lookup"""
        _ = phone_number
        return {
            'service': 'IPAPI',
            'requires_api_key': True,
            'api_url': 'https://ipapi.com/phone_api.json',
            'note': 'Phone validation and location data'
        }
    
    def _face_recognition_analysis(self, image_path: str) -> Dict[str, Any]:
        """Local face recognition"""
        _ = image_path
        return {
            'service': 'Local Face Recognition',
            'requires_library': 'face_recognition',
            'note': 'Local processing with face_recognition library'
        }
    
    def _facepp_analysis(self, image_path: str) -> Dict[str, Any]:
        """Face++ analysis"""
        _ = image_path
        return {
            'service': 'Face++',
            'requires_api_key': True,
            'api_url': 'https://api-cn.faceplusplus.com/facepp/v3/detect',
            'note': 'Advanced facial analysis'
        }
    
    def _kairos_analysis(self, image_path: str) -> Dict[str, Any]:
        """Kairos analysis"""
        _ = image_path
        return {
            'service': 'Kairos',
            'requires_api_key': True,
            'api_url': 'https://api.kairos.com/v2/api/detect',
            'note': 'Professional facial analysis'
        }
    
    # Phone validation
    def validate_and_parse(self, phone_number: str) -> Optional[phonenumbers.PhoneNumber]:
        """Validate and parse phone number"""
        phone_number = phone_number.replace(' ', '').replace('-', '')
        
        try:
            parsed = parse(phone_number, None)
            if is_valid_number(parsed):
                return parsed
        except NumberParseException:
            pass
        
        for region in ['US', 'GB', 'FR', 'DE', 'RU', 'IN', 'BR', 'CN', 'JP']:
            try:
                parsed = parse(phone_number, region)
                if is_valid_number(parsed):
                    return parsed
            except NumberParseException:
                continue
        
        return None
    
    def get_basic_phone_info(self, phone_number: phonenumbers.PhoneNumber) -> Dict[str, Any]:
        """Get basic phone information"""
        # `geocoder.description_for_number` returns a human-friendly location (region/city),
        # not necessarily the country name. We expose both for convenience.
        location_en = geocoder.description_for_number(phone_number, 'en')
        # Russian localization for UI/CLI users in this workspace.
        location_ru = geocoder.description_for_number(phone_number, 'ru')

        return {
            'country_code': phone_number.country_code,
            'national_number': phone_number.national_number,
            'international_format': phonenumbers.format_number(phone_number, PhoneNumberFormat.INTERNATIONAL),
            'national_format': phonenumbers.format_number(phone_number, PhoneNumberFormat.NATIONAL),
            'e164_format': phonenumbers.format_number(phone_number, PhoneNumberFormat.E164),
            # Backward compatible keys:
            'country': location_en,
            'carrier': carrier.name_for_number(phone_number, 'en'),
            # New fields:
            'region_code': phonenumbers.region_code_for_number(phone_number),
            'location_en': location_en,
            'location_ru': location_ru,
        }

    def _owner_lookup_local(self, phone_e164: str, limit: int = 5) -> Dict[str, Any]:
        """Try to find an owner/name for a phone using local, user-provided datasets.

        IMPORTANT:
        - We do NOT attempt any network scraping here.
        - This relies only on `business_directory.db` (see `directory_db.py`).
        - Returned values are "candidates" because directory datasets may contain
          organizations, departments, call-centers, etc. It's not always a private owner.
        """
        try:
            result = search_records(query=phone_e164, field='phone', limit=limit, offset=0)
            items = result.get('items') or []

            candidates = []
            for item in items[:limit]:
                if not isinstance(item, dict):
                    continue
                candidates.append(
                    {
                        'name': item.get('name') or '',
                        'legal_form': item.get('legal_form') or '',
                        'category': item.get('category') or '',
                        'city': item.get('city') or '',
                        'dataset_name': item.get('dataset_name') or '',
                        # Keep the phone itself for traceability; other fields are available
                        # via /api/directory/search if needed.
                        'phones': item.get('phones') or '',
                    }
                )

            return {
                'service': 'Local Directory (business_directory.db)',
                'found': bool(candidates),
                'matches': int(result.get('total') or 0),
                'candidates': candidates,
                'note': 'Owner/name is derived from local directory datasets; it may be an organization, not a private individual.'
            }
        except (sqlite3.Error, OSError, ValueError, TypeError) as e:
            return {
                'service': 'Local Directory (business_directory.db)',
                'found': False,
                'error': str(e),
            }
    
    def extract_photo_metadata(self, image_path: str) -> Dict[str, Any]:
        """Extract photo metadata"""
        try:
            with Image.open(image_path) as img:
                metadata = {
                    'filename': os.path.basename(image_path),
                    'format': img.format,
                    'size': img.size,
                    'file_size': os.path.getsize(image_path)
                }
                
                # EXIF data (use public Pillow API)
                try:
                    exif = img.getexif()
                except (AttributeError, OSError, ValueError):
                    exif = None

                if exif:
                    exif_data = {}
                    for tag, value in exif.items():
                        if tag in ExifTags.TAGS:
                            exif_data[ExifTags.TAGS[tag]] = str(value) if not isinstance(value, (str, int, float)) else value
                    metadata['exif'] = exif_data
                
                return metadata
        except (OSError, ValueError) as e:
            return {'error': f'Failed to extract metadata: {str(e)}'}
    
    def allowed_file(self, filename: str) -> bool:
        """Check if file is allowed"""
        return '.' in filename and \
               filename.rsplit('.', 1)[1].lower() in self.allowed_extensions
    
    def universal_phone_search(self, phone_number: str, search_types: List[str] = None) -> Dict[str, Any]:
        """Universal phone search across all sources"""
        if search_types is None:
            search_types = ['basic', 'google', 'social']
        
        parsed = self.validate_and_parse(phone_number)
        if not parsed:
            return {
                'error': 'Invalid phone number',
                'input': phone_number,
                'valid': False
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
        
        # Basic info
        if 'basic' in search_types:
            results['results']['basic'] = self.get_basic_phone_info(parsed)

        # Local owner/name lookup (offline, from directory DB only)
        if 'owner' in search_types or 'all' in search_types:
            results['results']['owner'] = self._owner_lookup_local(formatted)
        
        # Search engines
        if 'search_engines' in search_types or 'google' in search_types or 'all' in search_types:
            results['results']['search_engines'] = {}
            for engine in self.phone_search_engines:
                results['results']['search_engines'][engine] = self.phone_search_engines[engine](formatted)
        
        # Social platforms
        if 'social' in search_types or 'all' in search_types:
            results['results']['social_platforms'] = {}
            for platform in self.social_platforms:
                results['results']['social_platforms'][platform] = self.social_platforms[platform](formatted)
        
        # API services
        if 'api' in search_types or 'all' in search_types:
            results['results']['api_services'] = {}
            for service in self.api_services:
                results['results']['api_services'][service] = self.api_services[service](formatted)

        # Local breach databases (redacted)
        if 'data_breaches' in search_types or 'all' in search_types:
            results['results']['data_breaches'] = self.data_breaches_search(formatted)

        # Optional external vendor lookups (requires user-supplied API keys)
        if 'xosint_phone' in search_types:
            results['results']['xosint_phone'] = self.xosint.phone_external_lookup(formatted)
        
        return results
    
    def universal_photo_search(self, image_path: str, search_types: List[str] = None) -> Dict[str, Any]:
        """Universal photo search across all sources"""
        if search_types is None:
            search_types = ['metadata', 'google', 'yandex']
        
        results = {
            'image_path': image_path,
            'search_types': search_types,
            'timestamp': time.strftime('%Y-%m-%d %H:%M:%S'),
            'results': {}
        }
        
        # Metadata
        if 'metadata' in search_types:
            results['results']['metadata'] = self.extract_photo_metadata(image_path)
        
        # Image search engines
        if 'search_engines' in search_types or 'google' in search_types or 'all' in search_types:
            results['results']['image_search'] = {}
            for engine in self.photo_search_engines:
                results['results']['image_search'][engine] = self.photo_search_engines[engine](image_path)
        
        # Facial recognition
        if 'facial' in search_types or 'face' in search_types or 'all' in search_types:
            results['results']['facial_recognition'] = {}
            for service in self.facial_services:
                results['results']['facial_recognition'][service] = self.facial_services[service](image_path)
        
        return results

# Flask Application
app = Flask(__name__)
app.config['MAX_CONTENT_LENGTH'] = 16 * 1024 * 1024
app.config['UPLOAD_FOLDER'] = 'uploads'

universal_search = UniversalSearchSystem()

@app.route('/')
def index():
    """Main universal search interface"""
    return render_template('universal_search.html')

@app.route('/api/phone_search', methods=['POST'])
def api_phone_search():
    """Phone search API"""
    data = request.get_json()
    
    if not data or 'phone' not in data:
        return jsonify({'error': 'Phone number is required'}), 400
    
    phone_number = data['phone']
    search_types = data.get('search_types', ['basic', 'google', 'social'])
    
    result = universal_search.universal_phone_search(phone_number, search_types)
    return jsonify(result)

@app.route('/api/photo_search', methods=['POST'])
def api_photo_search():
    """Photo search API"""
    if 'file' not in request.files:
        return jsonify({'error': 'No file selected'}), 400
    
    file = request.files['file']
    if file.filename == '':
        return jsonify({'error': 'No file selected'}), 400
    
    if file and universal_search.allowed_file(file.filename):
        filename = secure_filename(file.filename)
        unique_filename = f"{uuid.uuid4().hex}_{filename}"
        filepath = os.path.join(app.config['UPLOAD_FOLDER'], unique_filename)
        file.save(filepath)
        
        search_types = request.form.getlist('search_types')
        if not search_types:
            search_types = ['metadata', 'google', 'yandex']
        
        result = universal_search.universal_photo_search(filepath, search_types)
        result['filename'] = filename
        result['unique_filename'] = unique_filename
        
        return jsonify(result)
    
    return jsonify({'error': 'File type not allowed'}), 400

@app.route('/api/sources', methods=['GET'])
def api_sources():
    """Get available search sources"""
    return jsonify({
        'phone_search_engines': list(universal_search.phone_search_engines.keys()),
        'social_platforms': list(universal_search.social_platforms.keys()),
        'photo_search_engines': list(universal_search.photo_search_engines.keys()),
        'facial_services': list(universal_search.facial_services.keys()),
        'api_services': list(universal_search.api_services.keys()),
        'local_databases': ['data_breaches', 'business_directory'],
        'osint_tools': ['phone_check', 'ip_lookup', 'email_check', 'xosint_phone', 'directory_search', 'directory_stats'],
        'phone_search_types': ['basic', 'search_engines', 'social', 'api', 'data_breaches', 'xosint_phone', 'owner', 'all'],
        'allowed_extensions': list(universal_search.allowed_extensions)
    })


@app.route('/api/phone_check', methods=['GET'])
def api_phone_check():
    """Fast phone check endpoint.

    Returns:
    - validation + formatting
    - basic phone info
    - optional breach summary (redacted)
    - optional external vendor lookups (requires API keys)

    Query params:
    - phone: required
    - breaches: (true/false), default true
    - external: (true/false), default false
    """
    phone = (request.args.get('phone') or '').strip()
    if not phone:
        return jsonify({'error': 'Phone number is required'}), 400

    breaches_flag = (request.args.get('breaches') or 'true').strip().lower() in {'1', 'true', 'yes', 'y'}
    external_flag = (request.args.get('external') or 'false').strip().lower() in {'1', 'true', 'yes', 'y'}

    parsed = universal_search.validate_and_parse(phone)
    if not parsed:
        return jsonify({'error': 'Invalid phone number', 'input': phone, 'valid': False}), 400

    formatted = phonenumbers.format_number(parsed, PhoneNumberFormat.E164)
    payload = {
        'input': phone,
        'formatted': formatted,
        'valid': True,
        'timestamp': time.strftime('%Y-%m-%d %H:%M:%S'),
        'results': {
            'basic': universal_search.get_basic_phone_info(parsed),
        },
    }

    if breaches_flag:
        payload['results']['data_breaches'] = universal_search.data_breaches_search(formatted)

    if external_flag:
        payload['results']['xosint_phone'] = universal_search.xosint.phone_external_lookup(formatted)

    return jsonify(payload)


@app.route('/api/ip_lookup', methods=['GET'])
def api_ip_lookup():
    """IP lookup (safe OSINT)."""
    ip = (request.args.get('ip') or '').strip()
    if not ip:
        return jsonify({'error': 'IP is required'}), 400

    payload = universal_search.xosint.ip_lookup(ip)
    if payload.get('valid') is False:
        return jsonify(payload), 400
    return jsonify(payload)


@app.route('/api/email_check', methods=['GET'])
def api_email_check():
    """Basic email checks (safe OSINT)."""
    email = (request.args.get('email') or '').strip()
    if not email:
        return jsonify({'error': 'Email is required'}), 400
    return jsonify(universal_search.xosint.email_check(email))


@app.route('/api/directory/search', methods=['GET'])
def api_directory_search():
    """Search business directory by phone/name/address or across all fields."""
    query = (request.args.get('query') or '').strip()
    field = (request.args.get('field') or 'all').strip().lower()
    limit = request.args.get('limit', '50')
    offset = request.args.get('offset', '0')

    if not query:
        return jsonify({'error': 'query is required'}), 400

    result = search_records(query=query, field=field, limit=limit, offset=offset)
    return jsonify(result)


@app.route('/api/directory/stats', methods=['GET'])
def api_directory_stats():
    """Stats by city and category for the business directory."""
    top_n = request.args.get('top', '20')
    result = stats_by_city_and_category(top_n=top_n)
    return jsonify(result)

if __name__ == '__main__':
    print("Universal Search System")
    print("=====================")
    print("Comprehensive OSINT tool for phones and photos")
    print("Starting web server on http://localhost:5000")
    print()
    print("Available endpoints:")
    print("  POST /api/phone_search - Universal phone search")
    print("  POST /api/photo_search - Universal photo search")
    print("  GET /api/sources - Available search sources")
    print()
    
    app.run(host='0.0.0.0', port=5000, debug=True)
