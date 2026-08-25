"""Lightweight integration of selected *safe* X-osint-like capabilities.

This project intentionally does NOT vendor/copy the upstream X-osint script.
Instead, it provides a small, testable adapter layer that exposes a subset of
benign OSINT utilities (IP lookup, basic email checks, optional phone lookups
via user-supplied API keys).

Rationale:
- Upstream X-osint is a single large interactive script with many optional,
  platform-specific features and runtime package installation.
- This codebase is an API-first Flask app; we keep integrations modular.

No hacking/pentesting features are included here.
"""

from __future__ import annotations

import ipaddress
import os
import re
import socket
from typing import Any, Dict, Optional

import requests


_EMAIL_RE = re.compile(r"^[^@\s]+@[^@\s]+\.[^@\s]+$")


class XOsintToolkit:
    """Small OSINT toolkit inspired by X-osint (safe subset)."""

    def __init__(self, session: Optional[requests.Session] = None):
        self.session = session or requests.Session()

    def ip_lookup(self, ip: str) -> Dict[str, Any]:
        """Return a combined IP intel payload from public sources.

        Sources:
        - ip-api.com (geo/ASN/ISP). No key.
        - ipinfo.io (geo/org). Optional token via IPINFO_TOKEN.
        - ARIN WHOIS REST. No key.
        - Reverse DNS (local resolver).
        """

        try:
            ip_obj = ipaddress.ip_address(ip)
        except ValueError:
            return {"input": ip, "valid": False, "error": "Invalid IP address"}

        ip_str = str(ip_obj)
        out: Dict[str, Any] = {
            "input": ip_str,
            "valid": True,
            "reverse_dns": None,
            "providers": {},
        }

        # Reverse DNS (best-effort)
        try:
            host, _, _ = socket.gethostbyaddr(ip_str)
            out["reverse_dns"] = host
        except OSError:
            out["reverse_dns"] = None

        # ip-api.com (best-effort)
        try:
            fields = (
                "status,message,country,countryCode,region,regionName,city,zip,lat,lon,"
                "timezone,isp,org,as,query,proxy,hosting,mobile"
            )
            resp = self.session.get(
                f"http://ip-api.com/json/{ip_str}",
                params={"fields": fields},
                timeout=10,
            )
            if resp.ok:
                out["providers"]["ip_api"] = resp.json()
            else:
                out["providers"]["ip_api"] = {"error": f"HTTP {resp.status_code}"}
        except (requests.RequestException, ValueError) as e:
            out["providers"]["ip_api"] = {"error": str(e)}

        # ipinfo.io (best-effort, optional token)
        try:
            params = {}
            token = os.getenv("IPINFO_TOKEN")
            if token:
                params["token"] = token
            resp = self.session.get(f"https://ipinfo.io/{ip_str}/json", params=params, timeout=10)
            if resp.ok:
                out["providers"]["ipinfo"] = resp.json()
            else:
                out["providers"]["ipinfo"] = {"error": f"HTTP {resp.status_code}"}
        except (requests.RequestException, ValueError) as e:
            out["providers"]["ipinfo"] = {"error": str(e)}

        # ARIN WHOIS REST (best-effort)
        try:
            resp = self.session.get(f"https://whois.arin.net/rest/ip/{ip_str}.json", timeout=10)
            if resp.ok:
                out["providers"]["arin_whois"] = resp.json()
            else:
                out["providers"]["arin_whois"] = {"error": f"HTTP {resp.status_code}"}
        except (requests.RequestException, ValueError) as e:
            out["providers"]["arin_whois"] = {"error": str(e)}

        return out

    def email_check(self, email: str) -> Dict[str, Any]:
        """Basic email sanity checks.

        - Validates format
        - Checks if the domain resolves (A/AAAA)
        - Checks if the domain is disposable (Kickbox public endpoint)

        This does not attempt account enumeration or any intrusive checks.
        """

        email = (email or "").strip()
        if not email:
            return {"input": email, "valid": False, "error": "Email is required"}

        format_ok = bool(_EMAIL_RE.match(email))
        domain = email.split("@")[-1].lower() if "@" in email else ""

        out: Dict[str, Any] = {
            "input": email,
            "format_valid": format_ok,
            "domain": domain or None,
            "domain_resolves": None,
            "disposable_domain": None,
        }

        # Domain resolves (best-effort)
        if domain:
            try:
                socket.getaddrinfo(domain, None)
                out["domain_resolves"] = True
            except OSError:
                out["domain_resolves"] = False

        # Disposable domain (best-effort)
        if domain:
            try:
                resp = self.session.get(
                    f"https://open.kickbox.com/v1/disposable/{domain}",
                    timeout=10,
                )
                if resp.ok:
                    data = resp.json()
                    out["disposable_domain"] = bool(data.get("disposable"))
                else:
                    out["disposable_domain"] = None
            except (requests.RequestException, ValueError):
                out["disposable_domain"] = None

        # Overall verdict
        out["valid"] = bool(format_ok and domain)
        return out

    def phone_external_lookup(self, phone_e164: str) -> Dict[str, Any]:
        """Optional phone intel via 3rd-party APIs.

        Requires user-provided API keys via environment variables:
        - IPQUALITYSCORE_API_KEY
        - VONAGE_API_KEY, VONAGE_API_SECRET (optional)

        Security note: this function only calls well-known vendor APIs; it does
        not scrape private sources.
        """

        phone_e164 = (phone_e164 or "").strip()
        if not phone_e164:
            return {"input": phone_e164, "ok": False, "error": "Phone is required"}

        out: Dict[str, Any] = {"input": phone_e164, "ok": True, "providers": {}}

        ipqs_key = os.getenv("IPQUALITYSCORE_API_KEY")
        if ipqs_key:
            try:
                resp = self.session.get(
                    f"https://ipqualityscore.com/api/json/phone/{ipqs_key}/{phone_e164}",
                    timeout=15,
                )
                out["providers"]["ipqualityscore"] = resp.json() if resp.ok else {"error": f"HTTP {resp.status_code}"}
            except (requests.RequestException, ValueError) as e:
                out["providers"]["ipqualityscore"] = {"error": str(e)}
        else:
            out["providers"]["ipqualityscore"] = {"not_configured": True}

        vonage_key = os.getenv("VONAGE_API_KEY")
        vonage_secret = os.getenv("VONAGE_API_SECRET")
        if vonage_key and vonage_secret:
            try:
                resp = self.session.get(
                    "https://api.nexmo.com/ni/advanced/async/json",
                    params={"api_key": vonage_key, "api_secret": vonage_secret, "number": phone_e164},
                    timeout=15,
                )
                out["providers"]["vonage"] = resp.json() if resp.ok else {"error": f"HTTP {resp.status_code}"}
            except (requests.RequestException, ValueError) as e:
                out["providers"]["vonage"] = {"error": str(e)}
        else:
            out["providers"]["vonage"] = {"not_configured": True}

        return out
