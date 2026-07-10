import asyncio
import getpass
import os

from telethon import TelegramClient
from telethon.errors import SessionPasswordNeededError
from telethon.sessions import StringSession


async def prompt(message: str) -> str:
    return (await asyncio.to_thread(input, message)).strip()


async def prompt_password(message: str) -> str:
    return await asyncio.to_thread(getpass.getpass, message)


async def main() -> None:
    api_id = int(os.getenv("TELEGRAM_API_ID") or await prompt("Telegram API ID: "))
    api_hash = os.getenv("TELEGRAM_API_HASH") or await prompt("Telegram API hash: ")
    phone = await prompt("Telegram phone number, for example +8613800000000: ")
    client = TelegramClient(
        StringSession(),
        api_id,
        api_hash,
    )
    await client.connect()
    if not await client.is_user_authorized():
        await client.send_code_request(phone)
        code = await prompt("Login code: ")
        try:
            await client.sign_in(phone=phone, code=code)
        except SessionPasswordNeededError:
            password = await prompt_password("Two-step verification password, if enabled: ")
            await client.sign_in(password=password)

    print("\nTELEGRAM_SESSION_STRING=" + client.session.save())
    await client.disconnect()


if __name__ == "__main__":
    asyncio.run(main())
