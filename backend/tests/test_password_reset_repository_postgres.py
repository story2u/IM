import os
from collections.abc import AsyncIterator
from datetime import timedelta

import pytest
from sqlalchemy import delete
from sqlalchemy.ext.asyncio import async_sessionmaker, create_async_engine
from sqlmodel.ext.asyncio.session import AsyncSession

from app.core.security import hash_password, verify_password
from app.infrastructure.db.models import PasswordResetChallenge, User, utc_now
from app.infrastructure.db.repositories import PasswordResetRepository

TEST_DATABASE_URL = os.getenv("SUBSCRIPTION_TEST_DATABASE_URL")
pytestmark = pytest.mark.skipif(
    not TEST_DATABASE_URL,
    reason="SUBSCRIPTION_TEST_DATABASE_URL is required for PostgreSQL password reset tests",
)


@pytest.fixture
async def reset_subject() -> AsyncIterator[tuple[async_sessionmaker[AsyncSession], User]]:
    assert TEST_DATABASE_URL
    engine = create_async_engine(TEST_DATABASE_URL)
    factory = async_sessionmaker(engine, class_=AsyncSession, expire_on_commit=False)
    user = User(email=f"password-reset-{os.urandom(8).hex()}@example.test")
    async with factory() as session:
        session.add(user)
        await session.commit()

    yield factory, user

    async with factory() as session:
        await session.exec(
            delete(PasswordResetChallenge).where(PasswordResetChallenge.user_id == user.id)
        )
        await session.exec(delete(User).where(User.id == user.id))
        await session.commit()
    await engine.dispose()


async def test_new_challenge_invalidates_old_and_password_change_consumes_all(reset_subject) -> None:
    factory, user = reset_subject
    async with factory() as session:
        repo = PasswordResetRepository(session)
        first = await repo.create(
            user_id=user.id,
            token_digest="a" * 64,
            code_digest="b" * 64,
            expires_at=utc_now() + timedelta(minutes=15),
        )
        second = await repo.create(
            user_id=user.id,
            token_digest="c" * 64,
            code_digest="d" * 64,
            expires_at=utc_now() + timedelta(minutes=15),
        )
        await session.refresh(first)
        assert first.used_at is not None
        assert await repo.active_by_token(first.token_digest) is None

        locked = await repo.active_by_token(second.token_digest)
        assert locked is not None
        loaded_user = await repo.user_for_challenge(locked)
        assert loaded_user is not None
        await repo.replace_password(
            user=loaded_user,
            password_hash=hash_password("new-password-123"),
            challenge=locked,
        )
        assert loaded_user.auth_version == 1
        assert verify_password("new-password-123", loaded_user.password_hash or "")
        assert await repo.active_by_token(second.token_digest) is None
