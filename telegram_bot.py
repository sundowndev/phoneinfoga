#!/usr/bin/env python3
"""Telegram bot for quick phone checks using this workspace search engine."""

from __future__ import annotations

import asyncio
import logging
import os
from typing import Any, Dict, List, Set

from dotenv import load_dotenv
from telegram import Update
from telegram.constants import ParseMode
from telegram.ext import (
    Application,
    CommandHandler,
    ContextTypes,
    MessageHandler,
    filters,
)

from universal_search_system import universal_search


logging.basicConfig(
    format="%(asctime)s | %(levelname)s | %(name)s | %(message)s",
    level=logging.INFO,
)
logger = logging.getLogger(__name__)


def _parse_allowed_chat_ids(raw: str) -> Set[int]:
    """Parse comma-separated chat IDs from env variable."""
    allowed: Set[int] = set()
    for part in (raw or "").split(","):
        value = part.strip()
        if not value:
            continue
        try:
            allowed.add(int(value))
        except ValueError:
            logger.warning("Invalid TELEGRAM_ALLOWED_CHAT_IDS value ignored: %s", value)
    return allowed


def _is_chat_allowed(update: Update, allowed_chat_ids: Set[int]) -> bool:
    if not allowed_chat_ids:
        return True
    if not update.effective_chat:
        return False
    return int(update.effective_chat.id) in allowed_chat_ids


async def _deny_if_not_allowed(update: Update, allowed_chat_ids: Set[int]) -> bool:
    if _is_chat_allowed(update, allowed_chat_ids):
        return False
    if update.message:
        await update.message.reply_text("⛔ Доступ к боту ограничен для этого чата.")
    return True


def _compact_lines(lines: List[str], max_len: int = 4000) -> str:
    """Join lines and cap output to Telegram safe text length."""
    text = "\n".join(lines)
    if len(text) <= max_len:
        return text
    return f"{text[: max_len - 15]}\n...\n(обрезано)"


def _safe_get(d: Dict[str, Any], *keys: str, default: Any = "—") -> Any:
    current: Any = d
    for key in keys:
        if not isinstance(current, dict) or key not in current:
            return default
        current = current[key]
    return current


def _format_result(payload: Dict[str, Any]) -> str:
    if not payload.get("valid"):
        err = payload.get("error", "Не удалось обработать номер")
        return f"❌ {err}"

    basic = _safe_get(payload, "results", "basic", default={})
    owner = _safe_get(payload, "results", "owner", default={})
    breaches = _safe_get(payload, "results", "data_breaches", default={})

    lines = [
        "✅ <b>Результат проверки</b>",
        f"📞 Ввод: <code>{payload.get('input', '—')}</code>",
        f"🌍 E.164: <code>{payload.get('formatted', '—')}</code>",
        "",
        "<b>Базовая информация</b>",
        f"• Страна/регион: {_safe_get(basic, 'country')} ({_safe_get(basic, 'region_code')})",
        f"• Оператор: {_safe_get(basic, 'carrier')}",
        f"• Международный формат: <code>{_safe_get(basic, 'international_format')}</code>",
    ]

    if isinstance(owner, dict) and owner.get("found"):
        lines.extend([
            "",
            "<b>Локальный справочник</b>",
            f"• Найдено совпадений: {owner.get('matches', 0)}",
        ])
        candidates = owner.get("candidates") or []
        for i, item in enumerate(candidates[:3], start=1):
            name = item.get("name") or "—"
            city = item.get("city") or "—"
            category = item.get("category") or "—"
            lines.append(f"  {i}. {name} | {city} | {category}")
    else:
        lines.extend([
            "",
            "<b>Локальный справочник</b>",
            "• Совпадений не найдено",
        ])

    if isinstance(breaches, dict):
        found = breaches.get("found")
        matches = breaches.get("matches", 0)
        lines.extend([
            "",
            "<b>Проверка утечек (редактированный вывод)</b>",
            f"• Обнаружено: {'да' if found else 'нет'}",
            f"• Количество записей: {matches}",
        ])

    return _compact_lines(lines)


def _format_ip_result(payload: Dict[str, Any]) -> str:
    if payload.get("valid") is False:
        return f"❌ {payload.get('error', 'Некорректный IP')}"

    lines = [
        "✅ <b>IP lookup</b>",
        f"🌐 IP: <code>{payload.get('ip', '—')}</code>",
        f"📍 Страна: {payload.get('country', '—')}",
        f"🏙 Город: {payload.get('city', '—')}",
        f"🛰 Провайдер/ASN: {payload.get('org') or payload.get('asn', '—')}",
    ]
    return _compact_lines(lines)


def _format_email_result(payload: Dict[str, Any]) -> str:
    if payload.get("valid") is False:
        return f"❌ {payload.get('error', 'Некорректный email')}"

    lines = [
        "✅ <b>Email check</b>",
        f"📧 Email: <code>{payload.get('email', '—')}</code>",
        f"• Синтаксис: {'OK' if payload.get('valid_format') else 'ошибка'}",
        f"• Домен: <code>{payload.get('domain', '—')}</code>",
        f"• MX записи: {'есть' if payload.get('has_mx') else 'нет/неизвестно'}",
    ]
    return _compact_lines(lines)


def _arg_from_context(context: ContextTypes.DEFAULT_TYPE) -> str:
    return " ".join(context.args).strip() if context.args else ""


async def start_cmd(update: Update, context: ContextTypes.DEFAULT_TYPE) -> None:
    allowed_chat_ids: Set[int] = context.application.bot_data.get("allowed_chat_ids", set())
    if await _deny_if_not_allowed(update, allowed_chat_ids):
        return

    if not update.message:
        return

    text = (
        "Привет! Я бот для быстрой проверки телефонных номеров в формате OSINT.\n\n"
        "Команды:\n"
        "• /start — старт\n"
        "• /help — помощь\n"
        "• /search <номер> — проверить номер\n"
        "• /ip <адрес> — информация по IP\n"
        "• /email <адрес> — базовая проверка email\n\n"
        "Пример: <code>/search +79001234567</code>\n"
        "Также можно просто отправить номер сообщением."
    )
    await update.message.reply_text(text, parse_mode=ParseMode.HTML)


async def help_cmd(update: Update, context: ContextTypes.DEFAULT_TYPE) -> None:
    allowed_chat_ids: Set[int] = context.application.bot_data.get("allowed_chat_ids", set())
    if await _deny_if_not_allowed(update, allowed_chat_ids):
        return

    if not update.message:
        return

    await update.message.reply_text(
        "Доступные команды:\n"
        "• <code>/search +79001234567</code> — проверка телефона\n"
        "• <code>/ip 8.8.8.8</code> — IP lookup\n"
        "• <code>/email test@example.com</code> — проверка email\n\n"
        "Или просто отправьте номер телефона отдельным сообщением.",
        parse_mode=ParseMode.HTML,
    )


async def _run_search_and_reply(update: Update, phone: str) -> None:
    if not update.message:
        return

    await update.message.reply_text("⏳ Ищу информацию, секунду...")
    payload = universal_search.universal_phone_search(
        phone,
        ["basic", "owner", "data_breaches"],
    )
    await update.message.reply_text(_format_result(payload), parse_mode=ParseMode.HTML)


async def search_cmd(update: Update, context: ContextTypes.DEFAULT_TYPE) -> None:
    allowed_chat_ids: Set[int] = context.application.bot_data.get("allowed_chat_ids", set())
    if await _deny_if_not_allowed(update, allowed_chat_ids):
        return

    if not update.message:
        return

    if not context.args:
        await update.message.reply_text("Укажи номер: <code>/search +79001234567</code>", parse_mode=ParseMode.HTML)
        return

    phone = _arg_from_context(context)
    await _run_search_and_reply(update, phone)


async def ip_cmd(update: Update, context: ContextTypes.DEFAULT_TYPE) -> None:
    allowed_chat_ids: Set[int] = context.application.bot_data.get("allowed_chat_ids", set())
    if await _deny_if_not_allowed(update, allowed_chat_ids):
        return

    if not update.message:
        return

    value = _arg_from_context(context)
    if not value:
        await update.message.reply_text("Укажи IP: <code>/ip 8.8.8.8</code>", parse_mode=ParseMode.HTML)
        return

    await update.message.reply_text("⏳ Проверяю IP...")
    payload = universal_search.xosint.ip_lookup(value)
    await update.message.reply_text(_format_ip_result(payload), parse_mode=ParseMode.HTML)


async def email_cmd(update: Update, context: ContextTypes.DEFAULT_TYPE) -> None:
    allowed_chat_ids: Set[int] = context.application.bot_data.get("allowed_chat_ids", set())
    if await _deny_if_not_allowed(update, allowed_chat_ids):
        return

    if not update.message:
        return

    value = _arg_from_context(context)
    if not value:
        await update.message.reply_text(
            "Укажи email: <code>/email test@example.com</code>",
            parse_mode=ParseMode.HTML,
        )
        return

    await update.message.reply_text("⏳ Проверяю email...")
    payload = universal_search.xosint.email_check(value)
    await update.message.reply_text(_format_email_result(payload), parse_mode=ParseMode.HTML)


async def text_fallback(update: Update, context: ContextTypes.DEFAULT_TYPE) -> None:
    allowed_chat_ids: Set[int] = context.application.bot_data.get("allowed_chat_ids", set())
    if await _deny_if_not_allowed(update, allowed_chat_ids):
        return

    if not update.message or not update.message.text:
        return

    phone = update.message.text.strip()
    await _run_search_and_reply(update, phone)


def build_app(token: str) -> Application:
    app = Application.builder().token(token).build()
    app.add_handler(CommandHandler("start", start_cmd))
    app.add_handler(CommandHandler("help", help_cmd))
    app.add_handler(CommandHandler("search", search_cmd))
    app.add_handler(CommandHandler("ip", ip_cmd))
    app.add_handler(CommandHandler("email", email_cmd))
    app.add_handler(MessageHandler(filters.TEXT & ~filters.COMMAND, text_fallback))
    return app


def run_sync() -> None:
    asyncio.run(run())


async def run() -> None:
    load_dotenv()
    token = os.getenv("TELEGRAM_BOT_TOKEN", "").strip()
    if not token:
        raise RuntimeError("TELEGRAM_BOT_TOKEN is not set. Add it to .env file.")

    allowed_chat_ids = _parse_allowed_chat_ids(os.getenv("TELEGRAM_ALLOWED_CHAT_IDS", ""))

    app = build_app(token)
    app.bot_data["allowed_chat_ids"] = allowed_chat_ids

    if allowed_chat_ids:
        logger.info("Telegram bot access control enabled for %d chat(s)", len(allowed_chat_ids))
    else:
        logger.info("Telegram bot access control disabled (all chats allowed)")

    logger.info("Telegram bot started")
    await app.initialize()
    await app.start()
    await app.updater.start_polling()

    try:
        while True:
            await asyncio.sleep(3600)
    except (KeyboardInterrupt, SystemExit):
        logger.info("Stopping Telegram bot")
    finally:
        await app.updater.stop()
        await app.stop()
        await app.shutdown()


if __name__ == "__main__":
    run_sync()