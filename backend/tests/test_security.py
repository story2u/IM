from uuid import UUID, uuid4

from app.core.config import Settings
from app.core.security import (
    create_access_token,
    decode_access_token,
    decrypt_secret,
    encrypt_secret,
    hash_password,
    verify_password,
)


def test_password_hash_round_trip() -> None:
    password_hash = hash_password("correct horse battery staple")

    assert verify_password("correct horse battery staple", password_hash)
    assert not verify_password("wrong password", password_hash)


def test_access_token_round_trip() -> None:
    settings = Settings(
        database_url="postgresql+asyncpg://user:password@localhost:5432/im",
        admin_api_token="admin-secret",
        jwt_secret_key="jwt-secret",
    )
    user_id = uuid4()

    token = create_access_token(subject=user_id, settings=settings)
    payload = decode_access_token(token, settings)

    assert UUID(payload["sub"]) == user_id


def test_secret_encryption_round_trip() -> None:
    settings = Settings(
        database_url="postgresql+asyncpg://user:password@localhost:5432/im",
        admin_api_token="admin-secret",
        jwt_secret_key="jwt-secret",
    )

    encrypted = encrypt_secret("telegram-session", settings)

    assert encrypted != "telegram-session"
    assert decrypt_secret(encrypted, settings) == "telegram-session"
