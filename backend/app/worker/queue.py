from uuid import UUID

import structlog

from app.worker.celery_app import celery_app

logger = structlog.get_logger(__name__)


class CeleryTaskQueue:
    def enqueue_ai_reply(self, opportunity_id: UUID) -> None:
        celery_app.send_task("ai.generate_and_send_reply", args=[str(opportunity_id)])

    def notify_reviewers(self, opportunity_id: UUID) -> None:
        logger.info("opportunity.pending_review", opportunity_id=str(opportunity_id))
