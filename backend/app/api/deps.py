from collections.abc import AsyncGenerator
from uuid import UUID

from fastapi import Depends, HTTPException, status
from fastapi.security import HTTPAuthorizationCredentials, HTTPBearer
from redis.asyncio import Redis
from sqlmodel.ext.asyncio.session import AsyncSession

from app.core.config import Settings, get_settings
from app.core.security import constant_time_equals
from app.core.time_window import WorkTimeConfig, WorkTimeService
from app.domain.services.detection_policy import OpportunityDetector
from app.infrastructure.ai.litellm_client import LiteLLMOpportunityClassifier, LiteLLMReplyGenerator
from app.infrastructure.db.models import Opportunity
from app.infrastructure.db.repositories import (
    ConfigRepository,
    MessageRepository,
    OpportunityRepository,
    ReplyTemplateRepository,
    RuleRepository,
)
from app.infrastructure.db.session import get_session
from app.infrastructure.im.base import AdapterRegistry
from app.infrastructure.im.telegram import TelegramAdapter
from app.infrastructure.im.wecom import WeComAdapter
from app.worker.queue import CeleryTaskQueue

bearer = HTTPBearer(auto_error=False)


async def require_admin(
    credentials: HTTPAuthorizationCredentials | None = Depends(bearer),
    settings: Settings = Depends(get_settings),
) -> None:
    if not credentials or credentials.scheme.lower() != "bearer":
        raise HTTPException(status_code=status.HTTP_401_UNAUTHORIZED, detail="missing token")
    if not constant_time_equals(credentials.credentials, settings.admin_api_token):
        raise HTTPException(status_code=status.HTTP_401_UNAUTHORIZED, detail="invalid token")


async def get_redis_client(settings: Settings = Depends(get_settings)) -> AsyncGenerator[Redis, None]:
    redis = Redis.from_url(settings.redis_url, decode_responses=True)
    try:
        yield redis
    finally:
        await redis.aclose()


def get_message_repo(session: AsyncSession = Depends(get_session)) -> MessageRepository:
    return MessageRepository(session)


def get_opportunity_repo(session: AsyncSession = Depends(get_session)) -> OpportunityRepository:
    return OpportunityRepository(session)


def get_rule_repo(session: AsyncSession = Depends(get_session)) -> RuleRepository:
    return RuleRepository(session)


def get_template_repo(session: AsyncSession = Depends(get_session)) -> ReplyTemplateRepository:
    return ReplyTemplateRepository(session)


async def get_work_time_service(
    session: AsyncSession = Depends(get_session),
    settings: Settings = Depends(get_settings),
) -> WorkTimeService:
    config_repo = ConfigRepository(session)
    raw_config = await config_repo.get_value("working_hours")
    config = (
        WorkTimeConfig.model_validate(raw_config)
        if raw_config
        else WorkTimeConfig.from_settings(settings)
    )
    return WorkTimeService(config)


def get_detector(settings: Settings = Depends(get_settings)) -> OpportunityDetector:
    classifier = LiteLLMOpportunityClassifier(settings)
    return OpportunityDetector(ai_classifier=classifier)


def get_adapter_registry(
    settings: Settings = Depends(get_settings),
    redis: Redis = Depends(get_redis_client),
) -> AdapterRegistry:
    return AdapterRegistry(
        [
            TelegramAdapter(settings),
            WeComAdapter(settings, redis=redis),
        ]
    )


def get_task_queue() -> CeleryTaskQueue:
    return CeleryTaskQueue()


def get_reply_generator(
    settings: Settings = Depends(get_settings),
    opportunity_repo: OpportunityRepository = Depends(get_opportunity_repo),
    message_repo: MessageRepository = Depends(get_message_repo),
) -> LiteLLMReplyGenerator:
    return LiteLLMReplyGenerator(
        settings=settings,
        opportunity_repo=opportunity_repo,
        message_repo=message_repo,
    )


async def get_opportunity_or_404(
    opportunity_id: UUID,
    repo: OpportunityRepository = Depends(get_opportunity_repo),
) -> Opportunity:
    opportunity = await repo.get(opportunity_id)
    if not opportunity:
        raise HTTPException(status_code=status.HTTP_404_NOT_FOUND, detail="opportunity not found")
    return opportunity
