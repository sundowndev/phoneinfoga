import json
import os
import sqlite3
import tempfile


def test_api_phone_search_basic_success():
    # Import inside the test so pytest collection remains fast/failures are localized.
    from universal_search_system import app

    client = app.test_client()
    resp = client.post(
        "/api/phone_search",
        data=json.dumps({"phone": "+79156129531", "search_types": ["basic"]}),
        content_type="application/json",
    )

    assert resp.status_code == 200
    payload = resp.get_json()
    assert isinstance(payload, dict)
    assert payload.get("valid") is True
    assert payload.get("results", {}).get("basic")


def test_api_phone_search_missing_phone_is_400():
    from universal_search_system import app

    client = app.test_client()
    resp = client.post("/api/phone_search", json={})

    assert resp.status_code == 400
    payload = resp.get_json()
    assert payload and "error" in payload


def test_api_phone_search_data_breaches_redacted():
    from universal_search_system import app

    client = app.test_client()
    # This phone exists in data_breaches.py sample dataset.
    resp = client.post(
        "/api/phone_search",
        json={"phone": "+79991234567", "search_types": ["data_breaches"]},
    )

    assert resp.status_code == 200
    payload = resp.get_json()
    assert payload.get("valid") is True

    breaches = payload.get("results", {}).get("data_breaches")
    assert isinstance(breaches, dict)
    assert breaches.get("service") == "Data Breaches Database"
    assert breaches.get("found") is True
    assert breaches.get("matches", 0) >= 1

    # Privacy guardrails: no raw leaked records are returned.
    assert breaches.get("data_redacted") is True
    assert breaches.get("data") == []
    # And no obvious PII fields should appear at top-level in the breach response.
    forbidden_keys = {"email", "name", "address", "password", "password_hash", "username"}
    assert not (forbidden_keys & set(breaches.keys()))


def test_api_ip_lookup_invalid_ip_is_400():
    from universal_search_system import app

    client = app.test_client()
    resp = client.get("/api/ip_lookup?ip=not-an-ip")

    assert resp.status_code == 400
    payload = resp.get_json()
    assert payload
    assert payload.get("valid") is False


def test_api_email_check_invalid_format_is_ok_and_non_network():
    from universal_search_system import app

    client = app.test_client()
    # Intentionally invalid: no domain => implementation won't attempt DNS/Kickbox.
    resp = client.get("/api/email_check?email=not-an-email")

    assert resp.status_code == 200
    payload = resp.get_json()
    assert payload
    assert payload.get("format_valid") is False


def test_api_phone_check_invalid_is_400():
    from universal_search_system import app

    client = app.test_client()
    resp = client.get("/api/phone_check?phone=123")

    assert resp.status_code == 400
    payload = resp.get_json()
    assert payload
    assert payload.get("valid") is False


def test_api_phone_check_valid_includes_basic_and_redacted_breaches():
    from universal_search_system import app

    client = app.test_client()
    # This phone exists in data_breaches.py sample dataset.
    resp = client.get("/api/phone_check?phone=%2B79991234567")

    assert resp.status_code == 200
    payload = resp.get_json()
    assert payload.get("valid") is True
    assert payload.get("results", {}).get("basic")

    breaches = payload.get("results", {}).get("data_breaches")
    assert isinstance(breaches, dict)
    assert breaches.get("found") is True
    assert breaches.get("data_redacted") is True
    assert breaches.get("data") == []


def _seed_directory_db(tmp_path: str) -> None:
    conn = sqlite3.connect(tmp_path)
    cursor = conn.cursor()
    cursor.execute(
        """
        CREATE TABLE datasets (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT,
            source_file TEXT UNIQUE,
            city TEXT,
            header_json TEXT,
            encoding TEXT,
            imported_at TEXT
        )
        """
    )
    cursor.execute(
        """
        CREATE TABLE records (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            dataset_id INTEGER,
            oid TEXT,
            name TEXT,
            legal_form TEXT,
            category TEXT,
            subcategory TEXT,
            services TEXT,
            description TEXT,
            phones TEXT,
            email TEXT,
            site TEXT,
            address TEXT,
            postal_code TEXT,
            hours TEXT,
            extra_1 TEXT,
            extra_2 TEXT,
            extra_3 TEXT,
            vk TEXT,
            facebook TEXT,
            skype TEXT,
            twitter TEXT,
            instagram TEXT,
            icq TEXT,
            jabber TEXT,
            raw_json TEXT
        )
        """
    )
    cursor.execute(
        "INSERT INTO datasets (name, source_file, city, header_json, encoding, imported_at) VALUES (?, ?, ?, ?, ?, ?)",
        ("TestCity", "test.csv", "TestCity", "[]", "utf-8", "2025-01-01"),
    )
    dataset_id = cursor.lastrowid
    cursor.execute(
        """
        INSERT INTO records (dataset_id, oid, name, legal_form, category, subcategory, services, description, phones,
                             email, site, address, postal_code, hours, extra_1, extra_2, extra_3, vk, facebook, skype,
                             twitter, instagram, icq, jabber, raw_json)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
        """,
        (
            dataset_id,
            "1",
            "Test Org",
            "LLC",
            "Services",
            "Repair",
            "Some services",
            "Desc",
            "+7-999-111-22-33",
            "test@example.com",
            "example.com",
            "Test address 1",
            "123000",
            "09:00-18:00",
            "",
            "",
            "",
            "",
            "",
            "",
            "",
            "",
            "",
            "",
            "[]",
        ),
    )
    conn.commit()
    conn.close()


def test_api_directory_search_and_stats_with_temp_db(monkeypatch):
    from universal_search_system import app

    with tempfile.TemporaryDirectory() as tmpdir:
        db_path = os.path.join(tmpdir, "business_directory.db")
        _seed_directory_db(db_path)
        monkeypatch.setenv("DIRECTORY_DB_PATH", db_path)

        client = app.test_client()

        resp = client.get("/api/directory/search?query=Test%20Org&field=name")
        assert resp.status_code == 200
        payload = resp.get_json()
        assert payload
        assert payload.get("total") == 1
        assert payload.get("items")

        resp2 = client.get("/api/directory/stats?top=5")
        assert resp2.status_code == 200
        payload2 = resp2.get_json()
        assert payload2
        assert payload2.get("total") == 1
        assert payload2.get("by_city")
