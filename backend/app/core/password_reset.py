import hashlib
import hmac
import secrets

from app.core.config import Settings

RESET_CODE_ALPHABET = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"


def generate_reset_token() -> str:
    return secrets.token_urlsafe(32)


def generate_reset_code() -> str:
    return "".join(secrets.choice(RESET_CODE_ALPHABET) for _ in range(10))


def reset_credential_digest(value: str, settings: Settings) -> str:
    secret = settings.jwt_secret_key or settings.admin_api_token
    return hmac.new(
        secret.encode("utf-8"),
        f"password-reset\0{value}".encode("utf-8"),
        hashlib.sha256,
    ).hexdigest()
