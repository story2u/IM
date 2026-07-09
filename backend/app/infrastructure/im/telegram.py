from typing import Any

import httpx

from app.core.config import Settings
from app.core.security import require_secret
from app.domain.enums import IMChannel
from app.domain.ports import InboundMessage, SendReceipt


class TelegramAdapter:
    channel = IMChannel.TELEGRAM

    def __init__(self, settings: Settings) -> None:
        self.settings = settings
        self.api_base_url = f"https://api.telegram.org/bot{settings.telegram_bot_token}"

    async def parse_webhook(
        self,
        payload: dict[str, Any],
        headers: dict[str, str],
        query: dict[str, str] | None = None,
    ) -> InboundMessage | None:
        secret = headers.get("x-telegram-bot-api-secret-token", "")
        require_secret(secret, self.settings.telegram_webhook_secret, "invalid telegram secret")

        message = payload.get("message") or payload.get("edited_message")
        if not message:
            return None

        text = message.get("text") or message.get("caption")
        if not text:
            return None

        chat = message.get("chat") or {}
        sender = message.get("from") or {}
        external_message_id = str(message.get("message_id"))
        conversation_id = str(chat.get("id"))

        sender_name = " ".join(
            item
            for item in [sender.get("first_name"), sender.get("last_name")]
            if item
        ) or sender.get("username")

        return InboundMessage(
            channel=self.channel,
            external_message_id=external_message_id,
            conversation_id=conversation_id,
            sender_external_id=str(sender.get("id")) if sender.get("id") else None,
            sender_display_name=sender_name,
            text=text,
            raw_payload=payload,
        )

    async def send_message(self, conversation_id: str, text: str) -> SendReceipt:
        if not self.settings.im_send_enabled:
            return SendReceipt(
                provider_message_id=None,
                raw_response={"dry_run": True, "channel": self.channel, "chat_id": conversation_id},
            )

        async with httpx.AsyncClient(timeout=10.0) as client:
            response = await client.post(
                f"{self.api_base_url}/sendMessage",
                json={"chat_id": conversation_id, "text": text},
            )
            response.raise_for_status()
            data = response.json()
            message_id = data.get("result", {}).get("message_id")
            return SendReceipt(
                provider_message_id=str(message_id) if message_id else None,
                raw_response=data,
            )
