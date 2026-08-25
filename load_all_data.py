#!/usr/bin/env python3
"""Seed/import all local datasets into SQLite databases.

This workspace uses a few local SQLite databases:
- data_breaches.db: simulated leak dataset (see `data_breaches.py`)
- sherlock_reports.db: simulated Sherlock-style dataset (see `sherlock_report.py`)
- business_directory.db: user-provided organization directories imported from CSV (see `directory_db.py`)

This script ensures the DBs exist and are populated with their built-in sample data.
Optionally, it imports *your own* directory CSV files into business_directory.db.

Examples:
  python load_all_data.py
  python load_all_data.py --directory-glob "C:\\Users\\test\\Downloads\\Data_*.csv"
  python load_all_data.py --directory-glob "D:\\datasets\\*.csv" --directory-db "D:\\db\\business_directory.db"

Notes:
- Directory CSV imports are assumed to be data you own and are allowed to process.
- Re-running is safe: built-in datasets are inserted only if missing; directory datasets
  are de-duplicated by absolute source_file path.
"""

from __future__ import annotations

import argparse
import glob
import os
import sqlite3
from typing import Optional, Sequence

from data_breaches import DataBreachesParser
from directory_db import DEFAULT_DB_PATH, import_many, init_db
from sherlock_report import SherlockReportGenerator


def _count_rows(db_path: str, table: str) -> Optional[int]:
    if not os.path.exists(db_path):
        return None
    conn = sqlite3.connect(db_path)
    try:
        cur = conn.cursor()
        cur.execute(f"SELECT COUNT(*) FROM {table}")
        return int(cur.fetchone()[0])
    except sqlite3.Error:
        return None
    finally:
        conn.close()


def _print_db_summary(directory_db_path: str = DEFAULT_DB_PATH) -> None:
    summaries = [
        (
            "data_breaches.db",
            ["users", "breaches"],
        ),
        (
            "sherlock_reports.db",
            ["sherlock_profiles", "sherlock_phonebooks", "sherlock_financial"],
        ),
        (
            directory_db_path,
            ["datasets", "records"],
        ),
    ]

    print("\n=== Database summary ===")
    for db_path, tables in summaries:
        abs_path = os.path.abspath(db_path)
        exists = os.path.exists(db_path)
        print(f"{db_path} ({abs_path}) - exists={exists}")
        for t in tables:
            cnt = _count_rows(db_path, t)
            if cnt is None:
                print(f"  {t}: n/a")
            else:
                print(f"  {t}: {cnt}")


def _parse_args(argv: Optional[Sequence[str]]) -> argparse.Namespace:
    p = argparse.ArgumentParser(description="Seed/import all local datasets into SQLite DBs")
    p.add_argument(
        "--directory-glob",
        default=None,
        help='Optional glob for directory CSVs, e.g. "C:\\Users\\me\\Downloads\\Data_*.csv"',
    )
    p.add_argument(
        "--directory-db",
        default=DEFAULT_DB_PATH,
        help=f"Path to directory DB (default: {DEFAULT_DB_PATH})",
    )
    p.add_argument(
        "--no-breaches",
        action="store_true",
        help="Skip seeding data_breaches.db",
    )
    p.add_argument(
        "--no-sherlock",
        action="store_true",
        help="Skip seeding sherlock_reports.db",
    )
    return p.parse_args(argv)


def main(argv: Optional[Sequence[str]] = None) -> int:
    args = _parse_args(argv)

    if not args.no_breaches:
        # Seeds itself (only if empty)
        DataBreachesParser()
        print("Seeded: data_breaches.db")

    if not args.no_sherlock:
        # Create tables + insert sample records (only if missing)
        s = SherlockReportGenerator()
        s.load_sherlock_data()
        print("Seeded: sherlock_reports.db")

    # Ensure directory DB exists even if no CSVs were provided.
    init_db(args.directory_db)

    imported_total = 0
    if args.directory_glob:
        paths = glob.glob(args.directory_glob)
        if not paths:
            print("No files matched directory glob:", args.directory_glob)
            _print_db_summary(directory_db_path=args.directory_db)
            return 2

        results = import_many(paths, db_path=args.directory_db)
        imported_total = sum(cnt for _, cnt in results)
        print("\nImported directory CSVs:")
        for path, cnt in results:
            print(f"  {os.path.basename(path)}: {cnt} records")
        print(f"Total directory records added: {imported_total}")
        print(f"Directory DB: {os.path.abspath(args.directory_db)}")

    _print_db_summary(directory_db_path=args.directory_db)
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
