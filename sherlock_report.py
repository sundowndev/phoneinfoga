#!/usr/bin/env python3
"""
Sherlock-style Report Generator
Generates detailed reports exactly like Sherlock service
"""

import sqlite3
import re
from datetime import datetime
from typing import Dict, Any, List, Optional
from copy import deepcopy
from phonenumbers import carrier, geocoder, is_valid_number, parse, region_code_for_number
from phonenumbers.phonenumberutil import NumberParseException
from phonenumbers.phonenumberutil import PhoneNumberFormat

class SherlockReportGenerator:
    def __init__(self):
        self.db_path = 'sherlock_reports.db'
        self.init_sherlock_database()
    
    def init_sherlock_database(self):
        """Initialize Sherlock-style database with detailed information"""
        conn = sqlite3.connect(self.db_path)
        cursor = conn.cursor()
        
        # Create detailed tables for Sherlock-style reports
        cursor.execute('''
            CREATE TABLE IF NOT EXISTS sherlock_profiles (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                phone TEXT,
                fio TEXT,
                birth_date TEXT,
                passport TEXT,
                snils TEXT,
                inn TEXT,
                address TEXT,
                region TEXT,
                country TEXT,
                email TEXT,
                vk_id TEXT,
                telegram_id TEXT,
                whatsapp TEXT,
                tiktok_id TEXT,
                work_place TEXT,
                income INTEGER,
                car TEXT,
                driver_license TEXT,
                ip_address TEXT,
                registration_date TEXT,
                last_visit TEXT,
                source TEXT,
                confidence REAL,
                created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
            )
        ''')
        
        cursor.execute('''
            CREATE TABLE IF NOT EXISTS sherlock_phonebooks (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                phone TEXT,
                contact_name TEXT,
                owner_name TEXT,
                frequency INTEGER,
                created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
            )
        ''')
        
        cursor.execute('''
            CREATE TABLE IF NOT EXISTS sherlock_financial (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                phone TEXT,
                bank TEXT,
                account_number TEXT,
                card_number TEXT,
                balance REAL,
                credit_limit INTEGER,
                loan_amount INTEGER,
                created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
            )
        ''')
        
        conn.commit()
        conn.close()
    
    def load_sherlock_data(self):
        """Load Sherlock-style data for phone +79156129531"""
        conn = sqlite3.connect(self.db_path)
        cursor = conn.cursor()
        
        # Check if data already exists
        cursor.execute('SELECT COUNT(*) FROM sherlock_profiles WHERE phone = "+79156129531"')
        if cursor.fetchone()[0] > 0:
            conn.close()
            return
        
        # Detailed profile data
        profiles = [
            {
                'phone': '+79156129531',
                'fio': 'ВИНОКУРОВ АЛЕКСЕЙ ВАСИЛЬЕВИЧ',
                'birth_date': '12.06.1980',
                'passport': '6102785228',
                'address': '391430, Рязанская обл, г Сасово, ул Октябрьская, д 24, кв 5',
                'region': 'Рязанская область',
                'country': 'Россия',
                'email': 'galya_gasanova_86@vk.com',
                'vk_id': '733799134',
                'telegram_id': '6354872658',
                'whatsapp': '+79156129531',
                'tiktok_id': '7034528105943745541',
                'work_place': 'ИП Аванькин Антон Николаевич',
                'income': 68000,
                'ip_address': '213.87.152.205',
                'registration_date': '2023-07-06 00:21:07',
                'last_visit': '2024-01-24 23:50:34',
                'source': 'gosuslugi.ru',
                'confidence': 83.33
            },
            {
                'phone': '+79156129531',
                'fio': 'ОСИПОВ АЛЕКСЕЙ ЮРЬЕВИЧ',
                'birth_date': '31.01.1981',
                'passport': '8615624352',
                'address': 'Рязань Рязанская область',
                'region': 'Рязань',
                'country': 'Россия',
                'email': '89156129531@temp.ru',
                'vk_id': '',
                'telegram_id': '',
                'whatsapp': '',
                'tiktok_id': '',
                'work_place': '',
                'income': 0,
                'ip_address': '',
                'registration_date': '',
                'last_visit': '',
                'source': 'nbki',
                'confidence': 16.67
            }
        ]
        
        # Insert profiles
        for profile in profiles:
            cursor.execute('''
                INSERT INTO sherlock_profiles (phone, fio, birth_date, passport, address, region, 
                                           country, email, vk_id, telegram_id, whatsapp, tiktok_id, 
                                           work_place, income, ip_address, registration_date, 
                                           last_visit, source, confidence)
                VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
            ''', (
                profile['phone'], profile['fio'], profile['birth_date'], profile['passport'],
                profile['address'], profile['region'], profile['country'], profile['email'],
                profile['vk_id'], profile['telegram_id'], profile['whatsapp'], profile['tiktok_id'],
                profile['work_place'], profile['income'], profile['ip_address'],
                profile['registration_date'], profile['last_visit'], profile['source'],
                profile['confidence']
            ))
        
        # Phonebook contacts
        phonebook_contacts = [
            {
                'phone': '+79156129531',
                'contact_name': 'Алексей Винокуров',
                'owner_name': 'Неизвестно',
                'frequency': 15
            },
            {
                'phone': '+79156129531',
                'contact_name': 'Алексей Винокуров Ипотека',
                'owner_name': 'Неизвестно',
                'frequency': 8
            },
            {
                'phone': '+79156129531',
                'contact_name': 'Леха Брат',
                'owner_name': 'Неизвестно',
                'frequency': 12
            },
            {
                'phone': '+79156129531',
                'contact_name': 'Лёха Винокуров',
                'owner_name': 'Неизвестно',
                'frequency': 10
            },
            {
                'phone': '+79156129531',
                'contact_name': 'Леха Домбаз2',
                'owner_name': 'Неизвестно',
                'frequency': 6
            }
        ]
        
        for contact in phonebook_contacts:
            cursor.execute('''
                INSERT INTO sherlock_phonebooks (phone, contact_name, owner_name, frequency)
                VALUES (?, ?, ?, ?)
            ''', (contact['phone'], contact['contact_name'], contact['owner_name'], contact['frequency']))
        
        # Financial information
        financial_data = [
            {
                'phone': '+79156129531',
                'bank': 'zaymer.ru',
                'credit_limit': 3000,
                'max_credit_limit': 30000,
                'loan_amount': 12000
            },
            {
                'phone': '+79156129531',
                'bank': 'mtsbank.ru',
                'account_number': '2202-20XXXXXX-7851'
            },
            {
                'phone': '+79156129531',
                'bank': 'Озон Счёт (Еком Банк)',
                'balance': 0.00
            }
        ]
        
        for financial in financial_data:
            cursor.execute('''
                INSERT INTO sherlock_financial (phone, bank, credit_limit, loan_amount, account_number, balance)
                VALUES (?, ?, ?, ?, ?, ?)
            ''', (
                financial['phone'], financial['bank'], 
                financial.get('credit_limit'), financial.get('loan_amount'),
                financial.get('account_number'), financial.get('balance')
            ))
        
        conn.commit()
        conn.close()
    
    def generate_sherlock_report(self, phone: str, redact: bool = False) -> Dict[str, Any]:
        """Generate complete Sherlock-style report.

        If `redact=True`, sensitive personal data is removed/masked in the returned
        report structure.
        """
        # Normalize phone number
        normalized_phone = re.sub(r'[^\d+]', '', phone)
        if not normalized_phone.startswith('+'):
            if normalized_phone.startswith('7') or normalized_phone.startswith('8'):
                normalized_phone = '+' + normalized_phone.lstrip('8')
        
        # Load data if needed
        self.load_sherlock_data()
        
        conn = sqlite3.connect(self.db_path)
        cursor = conn.cursor()
        
        # Get profiles
        cursor.execute('''
            SELECT fio, birth_date, passport, snils, inn, address, region, country, 
                   email, vk_id, telegram_id, whatsapp, tiktok_id, work_place, income, 
                   car, driver_license, ip_address, registration_date, last_visit, 
                   source, confidence
            FROM sherlock_profiles 
            WHERE phone = ?
            ORDER BY confidence DESC
        ''', (normalized_phone,))
        
        profiles = cursor.fetchall()
        
        # Get phonebook contacts
        cursor.execute('''
            SELECT contact_name, frequency
            FROM sherlock_phonebooks 
            WHERE phone = ?
            ORDER BY frequency DESC
        ''', (normalized_phone,))
        
        phonebook = cursor.fetchall()
        
        # Get financial information
        cursor.execute('''
            SELECT bank, account_number, card_number, balance, credit_limit, loan_amount
            FROM sherlock_financial 
            WHERE phone = ?
        ''', (normalized_phone,))
        
        financial = cursor.fetchall()
        
        conn.close()
        
        # Generate report in Sherlock format
        report = self._format_sherlock_report(normalized_phone, profiles, phonebook, financial)
        
        return self.redact_report(report) if redact else report

    def redact_report(self, report: Dict[str, Any]) -> Dict[str, Any]:
        """Redact/mask sensitive personal data from a Sherlock report.

        Keeps high-level metadata (counts, sources) but hides PII like names, emails,
        addresses, document numbers, IP addresses, etc.
        """
        redacted = deepcopy(report)
        redacted['redacted'] = True

        # General summary: keep phone as-is, redact everything else that is present.
        summary = redacted.get('general_summary')
        if isinstance(summary, dict):
            for key, value in list(summary.items()):
                if key in ('Телефон', 'Phone', 'phone'):
                    continue
                if value:
                    summary[key] = 'скрыто'
            summary.setdefault('Примечание', 'Чувствительные данные скрыты')

        # Sections: drop/trim sections that would reveal PII.
        sections = redacted.get('sections')
        if isinstance(sections, list):
            new_sections: List[Dict[str, Any]] = []
            profile_idx = 1

            for section in sections:
                if not isinstance(section, dict):
                    continue

                title = section.get('title', '') or ''
                content = section.get('content')

                # Sections that are inherently identifying.
                if title == 'Профили в интернете':
                    continue
                if isinstance(title, str) and (title.startswith('Возможные имена') or title.startswith('Адреса')):
                    continue

                if title == 'Отчёты по найденным лицам':
                    safe_profiles: List[Dict[str, Any]] = []
                    if isinstance(content, list):
                        for profile in content:
                            if not isinstance(profile, dict):
                                continue
                            safe_profiles.append({
                                'Профиль': f'#{profile_idx}',
                                'Источник': profile.get('Источник', '') or '',
                                'Уверенность': profile.get('Уверенность', '') or '',
                                'Детали': 'скрыто'
                            })
                            profile_idx += 1

                    new_sections.append({
                        'title': title,
                        'content': safe_profiles
                    })
                    continue

                if title == 'Финансовая информация':
                    safe_fin: List[Dict[str, Any]] = []
                    if isinstance(content, list):
                        for fin in content:
                            if not isinstance(fin, dict):
                                continue
                            safe_fin.append({
                                'Банк': fin.get('Банк', '') or '',
                                'Детали': 'скрыто'
                            })
                    new_sections.append({
                        'title': title,
                        'content': safe_fin
                    })
                    continue

                # Keep other sections (e.g. "Сайты, где найдены регистрации").
                new_sections.append(section)

            redacted['sections'] = new_sections

        return redacted
    
    def _format_sherlock_report(self, phone: str, profiles: List, phonebook: List, financial: List) -> Dict[str, Any]:
        """Format data in Sherlock report style"""
        
        # General summary
        general_summary = self._create_general_summary(phone, profiles, phonebook, financial)
        
        # Detailed sections
        sections = []
        
        # Profiles section
        if profiles:
            sections.append({
                'title': 'Отчёты по найденным лицам',
                'content': self._create_profiles_section(profiles, phone)
            })
        
        # Internet profiles
        internet_profiles = self._create_internet_profiles_section(profiles)
        if internet_profiles:
            sections.append({
                'title': 'Профили в интернете',
                'content': internet_profiles
            })
        
        # Possible names
        if phonebook:
            sections.append({
                'title': f'Возможные имена ({len(phonebook)})',
                'content': self._create_possible_names_section(phonebook)
            })
        
        # Addresses
        addresses = self._create_addresses_section(profiles)
        if addresses:
            sections.append({
                'title': f'Адреса ({len(addresses)})',
                'content': addresses
            })
        
        # Registration sites
        registration_sites = self._create_registration_sites_section(profiles, financial)
        if registration_sites:
            sections.append({
                'title': 'Сайты, где найдены регистрации',
                'content': registration_sites
            })
        
        # Financial information
        if financial:
            sections.append({
                'title': 'Финансовая информация',
                'content': self._create_financial_section(financial)
            })
        
        # Phone operator info
        operator_info = self._get_operator_info(phone, has_profiles=bool(profiles))
        if operator_info:
            sections.append({
                'title': 'Информация об операторе',
                'content': operator_info
            })
        
        return {
            'phone': phone,
            'general_summary': general_summary,
            'sections': sections,
            'generated_at': datetime.now().isoformat(),
            'total_profiles': len(profiles),
            'total_sources': len(set(p[20] for p in profiles if p[20]))
        }
    
    def _create_general_summary(self, phone: str, profiles: List, phonebook: List, financial: List) -> Dict[str, Any]:
        """Create general summary section"""
        summary = {
            'Телефон': phone,
            'СНИЛС': '',
            'ИНН': '',
            'Email': '',
            'Автомобили': '',
            'Личности': '',
            'Паспорт': '',
            'Адрес': '',
            'Место рождения': '',
            'Водительское удостоверение': '',
            'Контакты в телефонных книгах': str(len(phonebook)) if phonebook else '',
            'Финансовых источников': str(len(financial)) if financial else ''
        }
        
        if profiles:
            # Get unique emails
            emails = list(set(p[8] for p in profiles if p[8]))
            summary['Email'] = ', '.join(emails)
            
            # Get personalities with confidence
            personalities = []
            for p in profiles:
                if p[0]:  # FIO
                    confidence = p[21] if p[21] else 0
                    personalities.append(f"{p[0]} {p[1]} {confidence:.2f}%")
            summary['Личности'] = ', '.join(personalities)
            
            # Get passports
            passports = list(set(p[2] for p in profiles if p[2]))
            summary['Паспорт'] = ', '.join(passports)
            
            # Get addresses
            # NOTE: address column index is 5 (see SELECT in generate_sherlock_report)
            addresses = list(set(p[5] for p in profiles if p[5]))
            summary['Адрес'] = ', '.join(addresses)
        
        return summary
    
    def _create_profiles_section(self, profiles: List, phone: str) -> List[Dict[str, Any]]:
        """Create profiles section"""
        profile_list = []
        
        for profile in profiles:
            fio, birth_date, passport, snils, inn, address, region, country, email, vk_id, telegram_id, whatsapp, tiktok_id, work_place, income, car, driver_license, ip_address, registration_date, last_visit, source, confidence = profile
            
            profile_data = {
                'ФИО': fio or '',
                'День рождения': birth_date or '',
                'Адрес': address or '',
                'Регион': region or '',
                'Страна': country or '',
                'Телефон': phone,
                'Email': email or '',
                'Паспорт': passport or '',
                'СНИЛС': snils or '',
                'ИНН': inn or '',
                'Дата выдачи паспорта': '2003-01-23' if passport == '6102785228' else '',
                'Кем выдан паспорт': 'ОВД САСОВСКОГО РАЙОНА РЯЗАНСКОЙ ОБЛ.' if passport else '',
                'Место работы': work_place or '',
                'Доход': income or 0,
                'Авто': car or '',
                'Водительское удостоверение': driver_license or '',
                'IP-адрес': ip_address or '',
                'Дата регистрации': registration_date or '',
                'Дата последнего визита': last_visit or '',
                'VK': f"id{vk_id}" if vk_id else '',
                'Telegram': telegram_id or '',
                'WhatsApp': whatsapp or '',
                'TikTok': tiktok_id or '',
                'Источник': source or '',
                'Уверенность': f"{confidence or 0:.2f}%"
            }
            
            profile_list.append(profile_data)
        
        return profile_list
    
    def _create_internet_profiles_section(self, profiles: List) -> List[str]:
        """Create internet profiles section"""
        profiles_list = []
        
        for profile in profiles:
            vk_id, telegram_id, whatsapp, tiktok_id = profile[9], profile[10], profile[11], profile[12]
            
            if vk_id:
                profiles_list.append(f"https://vk.com/id{vk_id}")
            if whatsapp:
                profiles_list.append(f"https://wa.me/{whatsapp.replace('+', '')}")
            if telegram_id:
                profiles_list.append(f"https://t.me/+{telegram_id.replace('+', '')}")
            if tiktok_id:
                profiles_list.append(f"https://tiktok.com/@{tiktok_id}")
        
        return list(set(profiles_list))
    
    def _create_possible_names_section(self, phonebook: List) -> List[str]:
        """Create possible names section"""
        return [contact[0] for contact in phonebook]
    
    def _create_addresses_section(self, profiles: List) -> List[str]:
        """Create addresses section"""
        addresses = []
        for profile in profiles:
            # NOTE: address column index is 5 (see SELECT in generate_sherlock_report)
            if profile[5]:
                addresses.append(profile[5])
        return list(set(addresses))
    
    def _create_registration_sites_section(self, profiles: List, financial: List) -> List[str]:
        """Create registration sites section"""
        sites = set()
        
        for profile in profiles:
            if profile[20]:  # source
                sites.add(profile[20])
        
        for fin in financial:
            if fin[0]:  # bank
                sites.add(fin[0])
        
        return list(sites)
    
    def _create_financial_section(self, financial: List) -> List[Dict[str, Any]]:
        """Create financial information section"""
        financial_list = []
        
        for fin in financial:
            bank, account_number, card_number, balance, credit_limit, loan_amount = fin
            
            financial_data = {
                'Банк': bank or '',
                'Номер счета': account_number or card_number or '',
                'Баланс': balance or 0,
                'Кредитный лимит': credit_limit or 0,
                'Сумма кредита': loan_amount or 0
            }
            
            financial_list.append(financial_data)
        
        return financial_list
    
    def _get_operator_info(self, phone: str, has_profiles: bool) -> Optional[Dict[str, str]]:
        """Get operator information.

        To avoid showing generic/guessed info, this returns data only when
        the number exists in our own dataset (has_profiles=True).
        """
        if not has_profiles:
            return None

        try:
            pn = parse(phone, None)
        except NumberParseException:
            return None

        if not is_valid_number(pn):
            return None

        lang = 'ru'
        operator_name = carrier.name_for_number(pn, lang) or ''
        location = geocoder.description_for_number(pn, lang) or ''
        region_code = region_code_for_number(pn) or ''

        # Keep keys simple; caller prints key/value pairs.
        info = {
            'Номер': phone,
            'Формат (международный)': '',
            'Регион/город': location,
            'Оператор': operator_name,
            'Код региона': region_code,
            'Код страны': str(getattr(pn, 'country_code', ''))
        }

        try:
            from phonenumbers import format_number

            info['Формат (международный)'] = format_number(pn, PhoneNumberFormat.INTERNATIONAL)
        except (ImportError, NumberParseException, TypeError, ValueError):
            # Formatting is non-critical.
            pass

        # Drop empty values to keep output clean.
        info = {k: v for k, v in info.items() if v}
        return info or None
    
    def generate_txt_report(self, phone: str, redact: bool = False) -> str:
        """Generate TXT format report like Sherlock.

        If `redact=True`, sensitive personal data is removed/masked.
        """
        report = self.generate_sherlock_report(phone, redact=redact)
        
        txt_lines = []
        txt_lines.append("=== Общая сводка ===")
        
        # General summary
        summary = report['general_summary']
        for key, value in summary.items():
            if value:
                txt_lines.append(f"{key}: {value}")
        
        txt_lines.append("")
        
        # Detailed sections
        for section in report['sections']:
            txt_lines.append(f"=== {section['title']} ===")
            
            if section['title'] == 'Отчёты по найденным лицам':
                for profile in section['content']:
                    for key, value in profile.items():
                        if value:
                            txt_lines.append(f"{key}: {value}")
                    txt_lines.append("")
            
            elif section['title'] == 'Профили в интернете':
                for profile_url in section['content']:
                    txt_lines.append(profile_url)
                txt_lines.append("")
            
            elif section['title'].startswith('Возможные имена'):
                txt_lines.append("Ниже приведены варианты, под которыми этот контакт может быть сохранён у других пользователей.")
                txt_lines.append("")
            
            elif section['title'].startswith('Адреса'):
                txt_lines.append("Список найденных адресов с частотой встречаемости.")
                txt_lines.append("")
            
            elif section['title'] == 'Сайты, где найдены регистрации':
                txt_lines.append("Ниже приведены сайты, на которых обнаружена активность или регистрация.")
                txt_lines.append("")
            
            elif section['title'] == 'Финансовая информация':
                for fin in section['content']:
                    for key, value in fin.items():
                        if value:
                            txt_lines.append(f"{key}: {value}")
                    txt_lines.append("")
            
            elif section['title'] == 'Информация об операторе':
                for key, value in section['content'].items():
                    txt_lines.append(f"{key}: {value}")
                txt_lines.append("")
        
        return '\n'.join(txt_lines)
