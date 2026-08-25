"""Manual Sherlock report generation example.

This is NOT a pytest test module.
"""

__test__ = False


def main() -> None:
    from sherlock_report import SherlockReportGenerator

    sherlock = SherlockReportGenerator()
    report = sherlock.generate_sherlock_report('+79156129531')

    print('=== SHERLOCK ОТЧЕТ ===')
    print('Phone:', report['phone'])
    print('Total Profiles:', report['total_profiles'])
    print('Sections:', len(report['sections']))
    for section in report['sections']:
        print(f'- {section["title"]}')

    print('\n=== ОБЩАЯ СВОДКА ===')
    for key, value in report['general_summary'].items():
        if value:
            print(f'{key}: {value}')


if __name__ == '__main__':
    main()
