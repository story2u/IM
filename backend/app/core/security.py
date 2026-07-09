import hmac

from fastapi import HTTPException, status


def constant_time_equals(left: str, right: str) -> bool:
    return hmac.compare_digest(left.encode("utf-8"), right.encode("utf-8"))


def require_secret(actual: str, expected: str, detail: str = "invalid signature") -> None:
    if not expected or not constant_time_equals(actual, expected):
        raise HTTPException(status_code=status.HTTP_401_UNAUTHORIZED, detail=detail)
