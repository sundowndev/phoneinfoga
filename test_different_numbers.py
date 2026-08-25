"""Manual multi-number check against a *running* local API server.

This is NOT a pytest test module.
"""

__test__ = False


def main() -> None:
    import requests

    # Тестируем разные номера
    numbers = ['+79156129531', '+79991234567', '+78005553535']

    for phone in numbers:
        r = requests.post(
            'http://localhost:5000/api/phone_search',
            json={'phone': phone, 'search_types': ['basic', 'sherlock']},
            timeout=20,
        )
        data = r.json()

        print(f'=== НОМЕР: {phone} ===')
        print(f'Status: {r.status_code}')
        print(f'Valid: {data.get("valid", False)}')
        print(f'Sherlock found: {data.get("results", {}).get("sherlock_report", {}).get("found", False)}')
        print(f'Sherlock profiles: {data.get("results", {}).get("sherlock_report", {}).get("total_profiles", 0)}')
        print()


if __name__ == '__main__':
    main()
