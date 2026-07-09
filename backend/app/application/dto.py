from datetime import datetime
from uuid import UUID

from pydantic import BaseModel, ConfigDict, Field

from app.core.time_window import WorkTimeConfig
from app.domain.enums import (
    FrontendOpportunityStatus,
    IMChannel,
    MessageSource,
    OpportunityStatus,
    Priority,
    RuleType,
)


class OpportunityRead(BaseModel):
    id: UUID
    platform: IMChannel
    contactName: str
    contactAvatar: str
    summary: str
    matchedKeywords: list[str]
    confidenceScore: float
    status: FrontendOpportunityStatus
    internalStatus: OpportunityStatus
    priority: Priority
    lastMessagePreview: str
    createdAt: datetime
    updatedAt: datetime


class OpportunityDetailRead(OpportunityRead):
    aiReplyDraft: str | None = None
    finalReply: str | None = None
    detectionReason: str | None = None


class ChatMessageRead(BaseModel):
    id: UUID
    senderName: str
    content: str
    isFromContact: bool
    sentAt: datetime
    source: MessageSource | None


class ManualReplyRequest(BaseModel):
    text: str = Field(min_length=1, max_length=4000)
    operator_id: str = Field(default="operator", min_length=1, max_length=128)
    mark_following: bool = True


class AIDraftResponse(BaseModel):
    opportunity_id: UUID
    draft: str


class OpportunityStatusUpdate(BaseModel):
    status: OpportunityStatus


class RuleCreate(BaseModel):
    name: str = Field(min_length=1, max_length=128)
    rule_type: RuleType
    pattern: str = Field(min_length=1, max_length=500)
    score: float = Field(default=0.5, ge=0.0, le=1.0)
    priority: int = Field(default=100, ge=1, le=1000)
    enabled: bool = True


class RuleUpdate(BaseModel):
    name: str | None = Field(default=None, min_length=1, max_length=128)
    rule_type: RuleType | None = None
    pattern: str | None = Field(default=None, min_length=1, max_length=500)
    score: float | None = Field(default=None, ge=0.0, le=1.0)
    priority: int | None = Field(default=None, ge=1, le=1000)
    enabled: bool | None = None


class RuleRead(BaseModel):
    model_config = ConfigDict(from_attributes=True)

    id: UUID
    name: str
    enabled: bool
    priority: int
    rule_type: RuleType
    pattern: str
    score: float
    created_at: datetime
    updated_at: datetime


class ConfigRead(BaseModel):
    key: str
    value: dict
    description: str | None = None
    updated_at: datetime | None = None


class ConfigUpdate(BaseModel):
    value: dict
    description: str | None = None


class WorkModeRead(BaseModel):
    mode: str
    is_working_time: bool
    work_time: WorkTimeConfig


class ReplyTemplateCreate(BaseModel):
    title: str = Field(min_length=1, max_length=128)
    content: str = Field(min_length=1, max_length=4000)
    category: str = Field(default="通用", max_length=64)


class ReplyTemplateUpdate(BaseModel):
    title: str | None = Field(default=None, min_length=1, max_length=128)
    content: str | None = Field(default=None, min_length=1, max_length=4000)
    category: str | None = Field(default=None, max_length=64)
    enabled: bool | None = None


class ReplyTemplateRead(BaseModel):
    id: UUID
    title: str
    content: str
    category: str


class StatsSummaryRead(BaseModel):
    total: int
    pending: int
    replied: int
    ignored: int
    avgConfidence: float
