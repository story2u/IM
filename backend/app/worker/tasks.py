import asyncio
from uuid import UUID

import structlog

from app.application.use_cases.ai_reply import AIAutoReplyUseCase, transition_pending_to_ai
from app.core.config import get_settings
from app.core.time_window import WorkTimeConfig, WorkTimeService
from app.infrastructure.ai.litellm_client import LiteLLMReplyGenerator
from app.infrastructure.db.repositories import ConfigRepository, MessageRepository, OpportunityRepository
from app.infrastructure.db.session import AsyncSessionLocal
from app.infrastructure.im.base import AdapterRegistry
from app.infrastructure.im.telegram import TelegramAdapter
from app.infrastructure.im.wecom import WeComAdapter
from app.worker.celery_app import celery_app

logger = structlog.get_logger(__name__)


@celery_app.task(
    name="ai.generate_and_send_reply",
    queue="ai",
    autoretry_for=(Exception,),
    retry_backoff=True,
    retry_kwargs={"max_retries": 3},
)
def generate_and_send_reply(opportunity_id: str) -> None:
    asyncio.run(_generate_and_send_reply(UUID(opportunity_id)))


@celery_app.task(name="opportunity.sweep_pending_for_ai", queue="default")
def sweep_pending_for_ai() -> None:
    asyncio.run(_sweep_pending_for_ai())


async def _generate_and_send_reply(opportunity_id: UUID) -> None:
    settings = get_settings()
    async with AsyncSessionLocal() as session:
        opportunity_repo = OpportunityRepository(session)
        message_repo = MessageRepository(session)
        opportunity = await opportunity_repo.get(opportunity_id)
        if not opportunity:
            logger.warning("opportunity.not_found", opportunity_id=str(opportunity_id))
            return

        reply_generator = LiteLLMReplyGenerator(
            settings=settings,
            opportunity_repo=opportunity_repo,
            message_repo=message_repo,
        )
        adapters = AdapterRegistry(
            [
                TelegramAdapter(settings),
                WeComAdapter(settings),
            ]
        )
        use_case = AIAutoReplyUseCase(
            opportunity_repo=opportunity_repo,
            message_repo=message_repo,
            adapters=adapters,
            reply_generator=reply_generator,
        )
        await use_case.execute(opportunity)


async def _sweep_pending_for_ai() -> None:
    settings = get_settings()
    async with AsyncSessionLocal() as session:
        config_repo = ConfigRepository(session)
        raw_config = await config_repo.get_value("working_hours")
        config = (
            WorkTimeConfig.model_validate(raw_config)
            if raw_config
            else WorkTimeConfig.from_settings(settings)
        )
        work_time = WorkTimeService(config)
        if work_time.is_working_time() or not config.auto_reply_after_hours:
            return

        opportunity_repo = OpportunityRepository(session)
        stale = await opportunity_repo.pending_human_older_than(settings.pending_human_sla_minutes)
        for opportunity in stale:
            updated = await transition_pending_to_ai(opportunity_repo, opportunity.id)
            if updated:
                celery_app.send_task("ai.generate_and_send_reply", args=[str(updated.id)])
