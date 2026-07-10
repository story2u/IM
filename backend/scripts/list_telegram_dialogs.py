import asyncio
import os

from telethon import TelegramClient
from telethon.sessions import StringSession


async def prompt(message: str) -> str:
    return (await asyncio.to_thread(input, message)).strip()


async def main() -> None:
    api_id = int(os.getenv("TELEGRAM_API_ID") or await prompt("Telegram API ID: "))
    api_hash = os.getenv("TELEGRAM_API_HASH") or await prompt("Telegram API hash: ")
    session_string = os.getenv("TELEGRAM_SESSION_STRING") or await prompt("Telegram session string: ")
    client = TelegramClient(
        StringSession(session_string),
        api_id,
        api_hash,
    )
    await client.start()
    async for dialog in client.iter_dialogs():
        entity = dialog.entity
        username = getattr(entity, "username", None)
        print(f"{dialog.id}\t{dialog.name}\tusername={username or '-'}")
    await client.disconnect()


if __name__ == "__main__":
    asyncio.run(main())
