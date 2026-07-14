from app.core.config import Settings
from app.core.password_reset import reset_credential_digest
from app.core.security import constant_time_equals, hash_password, verify_password
from app.domain.ports import TaskQueue
from app.infrastructure.db.models import User
from app.infrastructure.db.repositories import PasswordResetRepository, UserRepository


class PasswordManagementError(Exception):
    pass


class CurrentPasswordInvalid(PasswordManagementError):
    pass


class PasswordResetInvalid(PasswordManagementError):
    pass


class PasswordResetRequired(PasswordManagementError):
    pass


class PasswordUnchanged(PasswordManagementError):
    pass


class PasswordManagementUseCase:
    def __init__(
        self,
        *,
        settings: Settings,
        user_repo: UserRepository,
        reset_repo: PasswordResetRepository,
        task_queue: TaskQueue,
    ) -> None:
        self.settings = settings
        self.user_repo = user_repo
        self.reset_repo = reset_repo
        self.task_queue = task_queue

    async def request_reset(self, email: str) -> bool:
        return self.task_queue.enqueue_password_reset(email)

    async def change_password(
        self, *, user: User, current_password: str, new_password: str
    ) -> None:
        if not user.password_hash:
            raise PasswordResetRequired
        if not verify_password(current_password, user.password_hash):
            raise CurrentPasswordInvalid
        if verify_password(new_password, user.password_hash):
            raise PasswordUnchanged
        await self.reset_repo.replace_password(user=user, password_hash=hash_password(new_password))

    async def confirm_reset(
        self,
        *,
        new_password: str,
        token: str | None = None,
        email: str | None = None,
        code: str | None = None,
    ) -> None:
        challenge = None
        user = None
        if token:
            challenge = await self.reset_repo.active_by_token(
                reset_credential_digest(token, self.settings)
            )
            if challenge:
                user = await self.reset_repo.user_for_challenge(challenge)
        elif email and code:
            record = await self.reset_repo.latest_active_for_email(email)
            if record:
                challenge, user = record
                supplied_digest = reset_credential_digest(code, self.settings)
                if not constant_time_equals(supplied_digest, challenge.code_digest):
                    await self.reset_repo.register_failed_attempt(
                        challenge,
                        max_attempts=self.settings.password_reset_max_attempts,
                    )
                    raise PasswordResetInvalid
        if not challenge or not user or not user.is_active:
            raise PasswordResetInvalid
        if challenge.failed_attempts >= self.settings.password_reset_max_attempts:
            raise PasswordResetInvalid
        if user.password_hash and verify_password(new_password, user.password_hash):
            raise PasswordUnchanged
        await self.reset_repo.replace_password(
            user=user,
            password_hash=hash_password(new_password),
            challenge=challenge,
        )
