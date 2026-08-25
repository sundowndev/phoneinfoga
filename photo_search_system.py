#!/usr/bin/env python3
"""
Photo Search System
Reverse image search and facial recognition for OSINT
"""

import os
import base64
import hashlib
import json
import requests
from io import BytesIO
from flask import Flask, render_template, request, jsonify, redirect, url_for
from werkzeug.utils import secure_filename
from PIL import Image, ExifTags
from typing import Dict, Any, List, Optional
import logging
import time
import uuid

# Configure logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

class PhotoSearchSystem:
    """Photo analysis and reverse image search system"""
    
    def __init__(self):
        self.session = requests.Session()
        self.session.headers.update({
            'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36'
        })
        
        self.upload_folder = 'uploads'
        self.allowed_extensions = {'png', 'jpg', 'jpeg', 'gif', 'bmp', 'webp'}
        self.max_file_size = 16 * 1024 * 1024  # 16MB
        
        # Create upload folder if it doesn't exist
        if not os.path.exists(self.upload_folder):
            os.makedirs(self.upload_folder)
        
        self.search_engines = {
            'google': self._google_images_search,
            'yandex': self._yandex_images_search,
            'bing': self._bing_images_search,
            'tineye': self._tineye_search,
            'baidu': self._baidu_search
        }
        
        self.facial_recognition = {
            'face_recognition': self._face_recognition_analysis,
            'facepp': self._facepp_analysis,
            'kairos': self._kairos_analysis
        }
    
    def allowed_file(self, filename: str) -> bool:
        """Check if file extension is allowed"""
        return '.' in filename and \
               filename.rsplit('.', 1)[1].lower() in self.allowed_extensions
    
    def extract_metadata(self, image_path: str) -> Dict[str, Any]:
        """Extract metadata from image file"""
        try:
            with Image.open(image_path) as img:
                metadata = {
                    'filename': os.path.basename(image_path),
                    'format': img.format,
                    'mode': img.mode,
                    'size': img.size,
                    'file_size': os.path.getsize(image_path)
                }
                
                # Extract EXIF data
                exif_data = {}
                if hasattr(img, '_getexif') and img._getexif() is not None:
                    exif = img._getexif()
                    for tag, value in exif.items():
                        if tag in ExifTags.TAGS:
                            exif_data[ExifTags.TAGS[tag]] = str(value) if not isinstance(value, (str, int, float)) else value
                
                metadata['exif'] = exif_data
                
                # Extract useful EXIF information
                if 'DateTime' in exif_data:
                    metadata['date_taken'] = exif_data['DateTime']
                if 'Make' in exif_data:
                    metadata['camera_make'] = exif_data['Make']
                if 'Model' in exif_data:
                    metadata['camera_model'] = exif_data['Model']
                if 'GPSInfo' in exif_data:
                    metadata['gps_info'] = 'GPS data available'
                
                # Calculate image hash
                with open(image_path, 'rb') as f:
                    img_data = f.read()
                    metadata['md5_hash'] = hashlib.md5(img_data).hexdigest()
                    metadata['sha256_hash'] = hashlib.sha256(img_data).hexdigest()
                
                return metadata
                
        except Exception as e:
            return {'error': f'Failed to extract metadata: {str(e)}'}
    
    def _google_images_search(self, image_path: str) -> Dict[str, Any]:
        """Google reverse image search"""
        try:
            # For Google Images, we need to upload the image
            # This is a simplified version - in production you'd use proper API
            search_url = "https://images.google.com/searchbyimage/upload"
            
            with open(image_path, 'rb') as f:
                files = {'encoded_image': f}
                response = self.session.post(search_url, files=files)
            
            return {
                'engine': 'Google Images',
                'search_url': 'https://images.google.com/',
                'upload_url': search_url,
                'manual_search': 'Upload image manually at Google Images',
                'note': 'Automatic upload requires Google API key',
                'features': ['Visual similarity', 'Web pages containing image', 'Similar images']
            }
        except Exception as e:
            return {
                'engine': 'Google Images',
                'error': str(e),
                'search_url': 'https://images.google.com/',
                'manual_search': 'Upload image manually at Google Images'
            }
    
    def _yandex_images_search(self, image_path: str) -> Dict[str, Any]:
        """Yandex reverse image search"""
        return {
            'engine': 'Yandex Images',
            'search_url': 'https://yandex.com/images/',
            'upload_url': 'https://yandex.com/images/app/?_url=',
            'manual_search': 'Upload image manually at Yandex Images',
            'features': ['Visual similarity', 'Face detection', 'Object recognition'],
            'note': 'Yandex has excellent face detection capabilities'
        }
    
    def _bing_images_search(self, image_path: str) -> Dict[str, Any]:
        """Bing reverse image search"""
        return {
            'engine': 'Bing Visual Search',
            'search_url': 'https://www.bing.com/visualsearch',
            'upload_url': 'https://www.bing.com/visualsearch',
            'manual_search': 'Upload image manually at Bing Visual Search',
            'features': ['Visual similarity', 'Object detection', 'Text extraction'],
            'note': 'Microsoft Bing has advanced visual search capabilities'
        }
    
    def _tineye_search(self, image_path: str) -> Dict[str, Any]:
        """TinEye reverse image search"""
        return {
            'engine': 'TinEye',
            'search_url': 'https://tineye.com/',
            'upload_url': 'https://tineye.com/search',
            'manual_search': 'Upload image manually at TinEye',
            'features': ['Exact matches', 'Modified versions', 'Usage tracking'],
            'note': 'Best for finding exact copies and usage history'
        }
    
    def _baidu_search(self, image_path: str) -> Dict[str, Any]:
        """Baidu reverse image search"""
        return {
            'engine': 'Baidu Images',
            'search_url': 'https://image.baidu.com/',
            'upload_url': 'https://image.baidu.com/pcdutu',
            'manual_search': 'Upload image manually at Baidu Images',
            'features': ['Face recognition', 'Celebrity matching', 'Visual similarity'],
            'note': 'Good for Asian face recognition'
        }
    
    def _face_recognition_analysis(self, image_path: str) -> Dict[str, Any]:
        """Local face recognition analysis (placeholder)"""
        return {
            'service': 'Local Face Recognition',
            'requires_library': 'face_recognition',
            'installation': 'pip install face_recognition',
            'features': [
                'Face detection',
                'Face encoding comparison',
                'Similarity matching',
                'Multiple face detection'
            ],
            'note': 'Requires face_recognition library and dlib',
            'setup_commands': [
                'pip install face_recognition',
                'pip install dlib',
                'pip install cmake'
            ]
        }
    
    def _facepp_analysis(self, image_path: str) -> Dict[str, Any]:
        """Face++ API analysis"""
        return {
            'service': 'Face++',
            'requires_api_key': True,
            'api_url': 'https://api-cn.faceplusplus.com/facepp/v3/detect',
            'api_key_required': 'Get API key from https://www.faceplusplus.com/',
            'features': [
                'Face detection',
                'Age estimation',
                'Gender detection',
                'Emotion recognition',
                'Beauty analysis',
                'Ethnicity detection'
            ],
            'sample_request': f'POST /facepp/v3/detect with image file',
            'note': 'Advanced facial analysis with Chinese AI technology'
        }
    
    def _kairos_analysis(self, image_path: str) -> Dict[str, Any]:
        """Kairos facial analysis"""
        return {
            'service': 'Kairos',
            'requires_api_key': True,
            'api_url': 'https://api.kairos.com/v2/api/detect',
            'api_key_required': 'Get API key from https://www.kairos.com/',
            'features': [
                'Face detection',
                'Age estimation',
                'Gender detection',
                'Emotion analysis',
                'Facial landmarks'
            ],
            'sample_request': f'POST /v2/api/detect with image file',
            'note': 'Professional facial analysis service'
        }
    
    def analyze_photo(self, image_path: str, analysis_types: List[str] = None) -> Dict[str, Any]:
        """Comprehensive photo analysis"""
        if analysis_types is None:
            analysis_types = ['metadata', 'google', 'yandex']
        
        if not os.path.exists(image_path):
            return {
                'error': 'Image file not found',
                'path': image_path
            }
        
        results = {
            'image_path': image_path,
            'analysis_types': analysis_types,
            'timestamp': time.strftime('%Y-%m-%d %H:%M:%S'),
            'results': {}
        }
        
        # Extract metadata
        if 'metadata' in analysis_types or 'all' in analysis_types:
            results['results']['metadata'] = self.extract_metadata(image_path)
        
        # Reverse image search
        if 'google' in analysis_types or 'search' in analysis_types or 'all' in analysis_types:
            results['results']['image_search'] = {}
            for engine in ['google', 'yandex', 'bing', 'tineye', 'baidu']:
                results['results']['image_search'][engine] = self.search_engines[engine](image_path)
        
        # Facial recognition
        if 'face' in analysis_types or 'facial' in analysis_types or 'all' in analysis_types:
            results['results']['facial_recognition'] = {}
            for service in ['face_recognition', 'facepp', 'kairos']:
                results['results']['facial_recognition'][service] = self.facial_recognition[service](image_path)
        
        return results

# Flask Web Application
app = Flask(__name__)
app.config['MAX_CONTENT_LENGTH'] = 16 * 1024 * 1024  # 16MB max file size
app.config['UPLOAD_FOLDER'] = 'uploads'

photo_search = PhotoSearchSystem()

@app.route('/')
def index():
    """Main photo search interface"""
    return render_template('photo_search.html')

@app.route('/upload', methods=['POST'])
def upload_file():
    """Handle file upload"""
    if 'file' not in request.files:
        return jsonify({'error': 'No file selected'}), 400
    
    file = request.files['file']
    if file.filename == '':
        return jsonify({'error': 'No file selected'}), 400
    
    if file and photo_search.allowed_file(file.filename):
        filename = secure_filename(file.filename)
        unique_filename = f"{uuid.uuid4().hex}_{filename}"
        filepath = os.path.join(app.config['UPLOAD_FOLDER'], unique_filename)
        file.save(filepath)
        
        # Get analysis types
        analysis_types = request.form.getlist('analysis_types')
        if not analysis_types:
            analysis_types = ['metadata', 'google']
        
        # Analyze the photo
        result = photo_search.analyze_photo(filepath, analysis_types)
        result['filename'] = filename
        result['unique_filename'] = unique_filename
        
        return jsonify(result)
    
    return jsonify({'error': 'File type not allowed'}), 400

@app.route('/api/analyze', methods=['POST'])
def api_analyze():
    """API endpoint for photo analysis"""
    data = request.get_json()
    
    if not data or 'image_path' not in data:
        return jsonify({'error': 'Image path is required'}), 400
    
    image_path = data['image_path']
    analysis_types = data.get('analysis_types', ['metadata', 'google'])
    
    result = photo_search.analyze_photo(image_path, analysis_types)
    return jsonify(result)

@app.route('/api/search_engines', methods=['GET'])
def api_search_engines():
    """Get available search engines"""
    return jsonify({
        'image_search': list(photo_search.search_engines.keys()),
        'facial_recognition': list(photo_search.facial_recognition.keys()),
        'allowed_extensions': list(photo_search.allowed_extensions),
        'max_file_size': photo_search.max_file_size
    })

if __name__ == '__main__':
    print("Photo Search System")
    print("==================")
    print("Starting web server on http://localhost:5000")
    print("Available endpoints:")
    print("  POST /upload - Upload and analyze photo")
    print("  POST /api/analyze - Analyze photo by path")
    print("  GET /api/search_engines - Available search engines")
    print()
    
    app.run(host='0.0.0.0', port=5000, debug=True)
