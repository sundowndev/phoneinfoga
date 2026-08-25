"""CLI tool to import organization CSVs into business_directory.db.

Usage:
  python import_directory_data.py --files "C:\\path\\Data_*.csv"

Notes:
- Data is assumed to be owned by the user with consent for processing.
- Imports all fields; also stores raw JSON per row.
"""

from __future__ import annotations

import argparse
import glob
import os

from directory_db import import_many, DEFAULT_DB_PATH


def build_parser() -> argparse.ArgumentParser:
    p = argparse.ArgumentParser(description="Import organization CSVs into business_directory.db")
    p.add_argument(
        "--files",
        required=True,
        help='Glob for CSV files, e.g. "C:\\Users\\test\\Downloads\\Data_*.csv"',
    )
    p.add_argument("--db", default=DEFAULT_DB_PATH)
    return p


def main(argv=None) -> int:
    args = build_parser().parse_args(argv)
    paths = glob.glob(args.files)
    if not paths:
        print("No files matched:", args.files)
        return 2

    results = import_many(paths, db_path=args.db)
    total = sum(count for _, count in results)

    print("Imported files:")
    for path, count in results:
        print(f"  {os.path.basename(path)}: {count} records")

    print(f"Total records added: {total}")
    print(f"Database: {os.path.abspath(args.db)}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
