"""Manual full-search example against a *running* local API server.

This is NOT a pytest test module.
"""

__test__ = False


def main() -> None:
    import requests

    # Полный поиск с реальными данными
    r = requests.post(
        'http://localhost:5000/api/phone_search',
        json={'phone': '+79156129531', 'search_types': ['all']},
        timeout=30,
    )
    data = r.json()

    print('=== ПОЛНЫЙ ПОИСК С РЕАЛЬНЫМИ ДАННЫМИ ===')
    print(f'Всего результатов: {len(data.get("results", {}))}')
    print(f'Keys: {list(data.get("results", {}).keys())}')

    # Sherlock данные
    sherlock = data.get("results", {}).get("sherlock_report", {})
    print(f'Sherlock профилей: {sherlock.get("total_profiles", 0)}')
    print(f'Sherlock источников: {sherlock.get("total_sources", 0)}')
    print(f'Sherlock найден: {sherlock.get("found", False)}')

    # Data breaches данные
    breaches = data.get("results", {}).get("data_breaches", {})
    print(f'Data breaches записей: {breaches.get("matches", 0)}')
    print(f'Data breaches найден: {breaches.get("found", False)}')

    # Basic данные
    basic = data.get("results", {}).get("basic", {})
    print(f'Basic страна: {basic.get("country", "N/A")}')
    print(f'Basic валидный: {basic.get("valid", False)}')

    print('\n=== РЕАЛЬНЫЕ ДАННЫЕ ИЗ SHERLOCK ===')
    if sherlock.get("found", False):
        report = sherlock.get("report", {})
        summary = report.get("general_summary", {})
        print(f'Телефон: {summary.get("Телефон", "N/A")}')
        print(f'Email: {summary.get("Email", "N/A")}')
        print(f'Личности: {summary.get("Личности", "N/A")}')
        print(f'Паспорт: {summary.get("Паспорт", "N/A")}')

        # Показываем первые 2 профиля
        sections = report.get("sections", [])
        for section in sections:
            if section["title"] == "Отчёты по найденным лицам":
                profiles = section["content"][:2]  # Первые 2 профиля
                for i, profile in enumerate(profiles, 1):
                    print(f'\nПрофиль {i}:')
                    print(f'  ФИО: {profile.get("ФИО", "N/A")}')
                    print(f'  Дата рождения: {profile.get("День рождения", "N/A")}')
                    print(f'  Адрес: {profile.get("Адрес", "N/A")}')
                    print(f'  Email: {profile.get("Email", "N/A")}')
                    print(f'  Работа: {profile.get("Место работы", "N/A")}')
                    print(f'  Доход: {profile.get("Доход", 0)}')
                break


if __name__ == '__main__':
    main()
