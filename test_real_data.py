"""Manual smoke-check against a *running* local API server.

This file is intentionally NOT a pytest test module.

To run it manually:
  1) start the API server (e.g. `python osint_cli.py serve`)
  2) run: `python test_real_data.py`
"""

# Prevent pytest from collecting/executing this module.
__test__ = False


def main() -> None:
    import requests
    from requests.exceptions import RequestException
    import sqlite3

    # Тестируем все источники данных по номеру +79156129531
    phone = '+79156129531'

    print('=== ПРОВЕРКА РЕАЛЬНЫХ ДАННЫХ ===')

    # 1. Basic info
    try:
        r1 = requests.post(
            'http://localhost:5000/api/phone_search',
            json={'phone': phone, 'search_types': ['basic']},
            timeout=15,
        )
        basic = r1.json()
        print(f'Basic info статус: {r1.status_code}')
        print(f'Basic info найден: {"basic" in basic.get("results", {})}')

        if 'basic' in basic.get('results', {}):
            basic_data = basic['results']['basic']
            print(f'- Страна: {basic_data.get("country", "N/A")}')
            print(f'- Оператор: {basic_data.get("carrier", "N/A")}')
            print(f'- Валидный: {basic_data.get("valid", False)}')
    except RequestException as e:
        print(f'Basic info ошибка: {e}')

    print()

    # 2. Data breaches
    try:
        r2 = requests.post(
            'http://localhost:5000/api/phone_search',
            json={'phone': phone, 'search_types': ['data_breaches']},
            timeout=15,
        )
        breaches = r2.json()
        print(f'Data breaches статус: {r2.status_code}')
        print(f'Data breaches найден: {"data_breaches" in breaches.get("results", {})}')

        if 'data_breaches' in breaches.get('results', {}):
            breaches_data = breaches['results']['data_breaches']
            print(f'- Найдено записей: {breaches_data.get("matches", 0)}')
            print(f'- Есть данные: {breaches_data.get("found", False)}')
    except RequestException as e:
        print(f'Data breaches ошибка: {e}')

    print()

    # 3. Sherlock report
    try:
        r3 = requests.post(
            'http://localhost:5000/api/phone_search',
            json={'phone': phone, 'search_types': ['sherlock']},
            timeout=15,
        )
        sherlock = r3.json()
        print(f'Sherlock статус: {r3.status_code}')
        print(f'Sherlock report найден: {"sherlock_report" in sherlock.get("results", {})}')

        if 'sherlock_report' in sherlock.get('results', {}):
            sherlock_data = sherlock['results']['sherlock_report']
            print(f'- Профилей в Sherlock: {sherlock_data.get("total_profiles", 0)}')
            print(f'- Есть данные: {sherlock_data.get("found", False)}')
            print(f'- Источников: {sherlock_data.get("total_sources", 0)}')
    except RequestException as e:
        print(f'Sherlock ошибка: {e}')

    print()

    # 4. Проверяем базу данных напрямую
    try:
        conn = sqlite3.connect('sherlock_reports.db')
        cursor = conn.cursor()
        cursor.execute('SELECT COUNT(*) FROM sherlock_profiles WHERE phone = ?', (phone,))
        count = cursor.fetchone()[0]
        print(f'Прямая проверка БД - записей в sherlock_profiles: {count}')
        conn.close()
    except sqlite3.Error as e:
        print(f'Ошибка БД: {e}')

    # 5. Проверяем data breaches базу
    try:
        conn = sqlite3.connect('data_breaches.db')
        cursor = conn.cursor()
        cursor.execute('SELECT COUNT(*) FROM users WHERE phone = ?', (phone,))
        count = cursor.fetchone()[0]
        print(f'Прямая проверка БД - записей в data_breaches: {count}')
        conn.close()
    except sqlite3.Error as e:
        print(f'Ошибка БД data_breaches: {e}')


if __name__ == '__main__':
    main()
