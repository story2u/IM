from typing import Any, Protocol
from uuid import UUID

from pydantic import BaseModel, Field

from app.domain.enums import IMChannel, Priority, RuleType


class InboundMessage(BaseModel):
    owner_user_id: UUID | None = None
    channel: IMChannel
    external_message_id: str
    conversation_id: str
    sender_external_id: str | None = None
    sender_display_name: str | None = None
    text: str | None = None
    source_type: str = "private"
    group_name: str | None = None
    raw_message_links: list[str] = Field(default_factory=list)
    raw_payload: dict[str, Any] = Field(default_factory=dict)


class SendReceipt(BaseModel):
    provider_message_id: str | None = None
    raw_response: dict[str, Any] = Field(default_factory=dict)


class DetectionRule(BaseModel):
    id: UUID
    name: str
    rule_type: RuleType
    pattern: str
    score: float
    priority: int


class DetectionResult(BaseModel):
    is_opportunity: bool
    confidence: float = 0.0
    title: str | None = None
    summary: str | None = None
    reason: str | None = None
    matched_keywords: list[str] = Field(default_factory=list)
    priority: Priority = Priority.NORMAL


class IMAdapter(Protocol):
    channel: IMChannel

    async def parse_webhook(
        self,
        payload: dict[str, Any],
        headers: dict[str, str],
        query: dict[str, str] | None = None,
    ) -> InboundMessage | None:
        ...

    async def send_message(self, conversation_id: str, text: str) -> SendReceipt:
        ...


class OpportunityAIClassifier(Protocol):
    async def classify(self, text: str, rule_score: float) -> DetectionResult:
        ...


class ReplyGenerator(Protocol):
    async def generate_reply(self, opportunity_id: UUID) -> str:
        ...


class TaskQueue(Protocol):
    def enqueue_ai_reply(self, opportunity_id: UUID) -> None:
        ...

    def notify_reviewers(self, opportunity_id: UUID) -> None:
        ...
