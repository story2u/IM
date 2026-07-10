from uuid import uuid4

from app.infrastructure.im.telegram_user import TelegramUserClient, TelegramUserClientConfig


def test_telegram_user_client_normalizes_chat_ids() -> None:
    config = TelegramUserClientConfig(
        user_id=uuid4(),
        api_id=1,
        api_hash="hash",
        session_string="",
        chats=["-1001234567890", "public_jobs_channel", 42, ""],
    )

    client = TelegramUserClient(config)

    assert client.normalized_chats() == [-1001234567890, "public_jobs_channel", 42]


def test_telegram_user_client_extracts_unique_links() -> None:
    config = TelegramUserClientConfig(
        user_id=uuid4(),
        api_id=1,
        api_hash="hash",
        session_string="",
        chats=[],
    )

    client = TelegramUserClient(config)

    assert client._extract_links("Apply: https://example.com/job https://example.com/job", []) == [
        "https://example.com/job"
    ]
