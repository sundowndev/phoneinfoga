"""Unified CLI for this workspace.

Goal: a single, testable entrypoint that ties together:
- Flask API (`universal_search_system.app`)
- Safe OSINT utilities (see `xosint_toolkit.py`)

We intentionally do NOT include features that facilitate hacking / abuse.
"""

from __future__ import annotations

import argparse
import asyncio
import json
import sys
from typing import List

import phonenumbers
from phonenumbers.phonenumberutil import PhoneNumberFormat

from directory_db import search_records, stats_by_city_and_category
from universal_search_system import app, universal_search


def _print_json(data) -> None:
    json.dump(data, sys.stdout, ensure_ascii=False, indent=2)
    sys.stdout.write("\n")


def cmd_serve(args: argparse.Namespace) -> int:
    # Keep debug off by default for safety/reproducibility.
    app.run(host=args.host, port=args.port, debug=args.debug)
    return 0


def cmd_ip(args: argparse.Namespace) -> int:
    payload = universal_search.xosint.ip_lookup(args.ip)
    _print_json(payload)
    return 0 if payload.get("valid") else 2


def cmd_email(args: argparse.Namespace) -> int:
    payload = universal_search.xosint.email_check(args.email)
    _print_json(payload)
    return 0 if payload.get("valid") else 2


def cmd_phone(args: argparse.Namespace) -> int:
    search_types: List[str] = args.search_types or ["basic", "google", "social"]
    payload = universal_search.universal_phone_search(args.phone, search_types)
    _print_json(payload)
    return 0 if payload.get("valid") else 2


def cmd_phone_check(args: argparse.Namespace) -> int:
    parsed = universal_search.validate_and_parse(args.phone)
    if not parsed:
        _print_json({"input": args.phone, "valid": False, "error": "Invalid phone number"})
        return 2

    formatted = phonenumbers.format_number(parsed, PhoneNumberFormat.E164)

    out = {
        "input": args.phone,
        "formatted": formatted,
        "valid": True,
        "results": {
            "basic": universal_search.get_basic_phone_info(parsed),
        },
    }

    if not args.no_breaches:
        out["results"]["data_breaches"] = universal_search.data_breaches_search(formatted)
    if args.external:
        out["results"]["xosint_phone"] = universal_search.xosint.phone_external_lookup(formatted)

    _print_json(out)
    return 0


def cmd_directory_search(args: argparse.Namespace) -> int:
    result = search_records(
        query=args.query,
        field=args.field,
        limit=args.limit,
        offset=args.offset,
        db_path=args.db,
    )
    _print_json(result)
    return 0


def cmd_directory_stats(args: argparse.Namespace) -> int:
    result = stats_by_city_and_category(db_path=args.db, top_n=args.top)
    _print_json(result)
    return 0


def cmd_telegram_bot(args: argparse.Namespace) -> int:
    _ = args
    from telegram_bot import run

    asyncio.run(run())
    return 0


def build_parser() -> argparse.ArgumentParser:
    p = argparse.ArgumentParser(prog="osint_cli", description="Unified OSINT CLI (safe subset).")
    sub = p.add_subparsers(dest="cmd", required=True)

    s = sub.add_parser("serve", help="Run Flask API server")
    s.add_argument("--host", default="0.0.0.0")
    s.add_argument("--port", type=int, default=5000)
    s.add_argument("--debug", action="store_true", help="Enable Flask debug mode")
    s.set_defaults(func=cmd_serve)

    ip = sub.add_parser("ip", help="IP lookup (safe)")
    ip.add_argument("ip")
    ip.set_defaults(func=cmd_ip)

    em = sub.add_parser("email", help="Email checks (safe)")
    em.add_argument("email")
    em.set_defaults(func=cmd_email)

    ph = sub.add_parser("phone", help="Phone search via universal engine")
    ph.add_argument("phone", help="Phone number (ideally E.164, e.g. +7915...) ")
    ph.add_argument(
        "--search-types",
        nargs="+",
        default=None,
        help="Search types, e.g. basic google social data_breaches xosint_phone",
    )
    ph.set_defaults(func=cmd_phone)

    chk = sub.add_parser("phone-check", help="Fast phone validation + basic info (+ optional breaches/external)")
    chk.add_argument("phone", help="Phone number (ideally E.164, e.g. +7915...)")
    chk.add_argument("--no-breaches", action="store_true", help="Do not query local breach DB")
    chk.add_argument("--external", action="store_true", help="Include external vendor lookups (requires API keys)")
    chk.set_defaults(func=cmd_phone_check)

    ds = sub.add_parser("directory-search", help="Search business directory")
    ds.add_argument("query", help="Search string (phone, name, address, etc.)")
    ds.add_argument("--field", default="all", choices=["all", "phone", "name", "address"])
    ds.add_argument("--limit", type=int, default=50)
    ds.add_argument("--offset", type=int, default=0)
    ds.add_argument("--db", default=None, help="Optional path to business_directory.db")
    ds.set_defaults(func=cmd_directory_search)

    st = sub.add_parser("directory-stats", help="Directory statistics by city/category")
    st.add_argument("--top", type=int, default=20)
    st.add_argument("--db", default=None, help="Optional path to business_directory.db")
    st.set_defaults(func=cmd_directory_stats)

    tg = sub.add_parser("telegram-bot", help="Run Telegram bot integration")
    tg.set_defaults(func=cmd_telegram_bot)

    return p


def main(argv: List[str] | None = None) -> int:
    parser = build_parser()
    args = parser.parse_args(argv)
    return int(args.func(args))


if __name__ == "__main__":
    raise SystemExit(main())
