from app.application.dto import (
    ChatMessageRead,
    OpportunityDetailRead,
    OpportunityRead,
    ReplyTemplateRead,
)
from app.domain.enums import (
    FrontendOpportunityStatus,
    MessageDirection,
    OpportunityStatus,
)
from app.infrastructure.db.models import Message, Opportunity, ReplyTemplate


def frontend_status(status: OpportunityStatus) -> FrontendOpportunityStatus:
    if status in {OpportunityStatus.PENDING_HUMAN, OpportunityStatus.AI_AUTO_REPLY}:
        return FrontendOpportunityStatus.PENDING
    if status in {OpportunityStatus.REPLIED, OpportunityStatus.FOLLOWING}:
        return FrontendOpportunityStatus.REPLIED
    return FrontendOpportunityStatus.IGNORED


def to_opportunity_read(opportunity: Opportunity) -> OpportunityRead:
    return OpportunityRead(
        id=opportunity.id,
        platform=opportunity.channel,
        contactName=opportunity.contact_name,
        contactAvatar=opportunity.contact_avatar,
        summary=opportunity.summary or opportunity.title,
        matchedKeywords=opportunity.matched_keywords,
        confidenceScore=opportunity.confidence,
        status=frontend_status(opportunity.status),
        internalStatus=opportunity.status,
        priority=opportunity.priority,
        lastMessagePreview=opportunity.last_message_preview,
        createdAt=opportunity.created_at,
        updatedAt=opportunity.updated_at,
    )


def to_opportunity_detail(opportunity: Opportunity) -> OpportunityDetailRead:
    base = to_opportunity_read(opportunity)
    return OpportunityDetailRead(
        **base.model_dump(),
        aiReplyDraft=opportunity.ai_reply_draft,
        finalReply=opportunity.final_reply,
        detectionReason=opportunity.detection_reason,
    )


def to_chat_message_read(message: Message) -> ChatMessageRead:
    return ChatMessageRead(
        id=message.id,
        senderName=message.sender_display_name or "客户",
        content=message.text or "",
        isFromContact=message.direction == MessageDirection.INCOMING,
        sentAt=message.sent_at,
        source=message.source,
    )


def to_reply_template_read(template: ReplyTemplate) -> ReplyTemplateRead:
    return ReplyTemplateRead(
        id=template.id,
        title=template.title,
        content=template.content,
        category=template.category,
    )
