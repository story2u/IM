from __future__ import annotations

from datetime import datetime, timedelta, timezone
from uuid import UUID

from sqlmodel import col, select
from sqlmodel.ext.asyncio.session import AsyncSession

from app.domain.enums import (
    FrontendOpportunityStatus,
    IMChannel,
    MessageDirection,
    MessageSource,
    OpportunityStatus,
    Priority,
)
from app.domain.ports import DetectionRule, InboundMessage
from app.infrastructure.db.models import AppConfig, Message, Opportunity, ReplyTemplate, Rule, utc_now


FRONTEND_STATUS_MAP: dict[FrontendOpportunityStatus, set[OpportunityStatus]] = {
    FrontendOpportunityStatus.PENDING: {
        OpportunityStatus.PENDING_HUMAN,
        OpportunityStatus.AI_AUTO_REPLY,
    },
    FrontendOpportunityStatus.REPLIED: {
        OpportunityStatus.REPLIED,
        OpportunityStatus.FOLLOWING,
    },
    FrontendOpportunityStatus.IGNORED: {
        OpportunityStatus.IGNORED,
        OpportunityStatus.CLOSED,
    },
}


class MessageRepository:
    def __init__(self, session: AsyncSession) -> None:
        self.session = session

    async def get_by_external_id(self, channel: IMChannel, external_message_id: str) -> Message | None:
        statement = select(Message).where(
            Message.channel == channel,
            Message.external_message_id == external_message_id,
        )
        result = await self.session.exec(statement)
        return result.first()

    async def create_incoming(self, inbound: InboundMessage) -> Message:
        message = Message(
            channel=inbound.channel,
            external_message_id=inbound.external_message_id,
            conversation_id=inbound.conversation_id,
            sender_external_id=inbound.sender_external_id,
            sender_display_name=inbound.sender_display_name,
            direction=MessageDirection.INCOMING,
            text=inbound.text,
            raw_payload=inbound.raw_payload,
        )
        self.session.add(message)
        await self.session.commit()
        await self.session.refresh(message)
        return message

    async def create_outgoing(
        self,
        *,
        channel: IMChannel,
        conversation_id: str,
        text: str,
        source: MessageSource,
        opportunity_id: UUID,
        external_message_id: str,
        raw_payload: dict,
    ) -> Message:
        message = Message(
            channel=channel,
            external_message_id=external_message_id,
            conversation_id=conversation_id,
            sender_display_name="商机助手",
            direction=MessageDirection.OUTGOING,
            source=source,
            text=text,
            raw_payload=raw_payload,
            opportunity_id=opportunity_id,
        )
        self.session.add(message)
        await self.session.commit()
        await self.session.refresh(message)
        return message

    async def attach_opportunity(self, message_id: UUID, opportunity_id: UUID) -> None:
        message = await self.session.get(Message, message_id)
        if not message:
            return
        message.opportunity_id = opportunity_id
        message.processed_at = utc_now()
        message.updated_at = utc_now()
        self.session.add(message)
        await self.session.commit()

    async def mark_processed(self, message_id: UUID) -> None:
        message = await self.session.get(Message, message_id)
        if not message:
            return
        message.processed_at = utc_now()
        message.updated_at = utc_now()
        self.session.add(message)
        await self.session.commit()

    async def list_by_opportunity(self, opportunity_id: UUID) -> list[Message]:
        statement = (
            select(Message)
            .where(Message.opportunity_id == opportunity_id)
            .order_by(col(Message.sent_at).asc(), col(Message.created_at).asc())
        )
        result = await self.session.exec(statement)
        return list(result.all())

    async def list_by_conversation(
        self,
        channel: IMChannel,
        conversation_id: str,
        limit: int = 20,
    ) -> list[Message]:
        statement = (
            select(Message)
            .where(Message.channel == channel, Message.conversation_id == conversation_id)
            .order_by(col(Message.sent_at).desc())
            .limit(limit)
        )
        result = await self.session.exec(statement)
        return list(reversed(result.all()))


class OpportunityRepository:
    def __init__(self, session: AsyncSession) -> None:
        self.session = session

    async def create(
        self,
        *,
        channel: IMChannel,
        conversation_id: str,
        customer_external_id: str | None,
        contact_name: str | None,
        source_type: str,
        group_name: str | None,
        source_message_id: UUID,
        title: str,
        summary: str | None,
        matched_keywords: list[str],
        raw_message_links: list[str],
        confidence: float,
        priority: Priority,
        detection_reason: str | None,
        status: OpportunityStatus,
        last_message_preview: str,
    ) -> Opportunity:
        opportunity = Opportunity(
            channel=channel,
            conversation_id=conversation_id,
            customer_external_id=customer_external_id,
            contact_name=contact_name or customer_external_id or "未知联系人",
            source_type=source_type,
            group_name=group_name,
            source_message_id=source_message_id,
            title=title,
            summary=summary,
            matched_keywords=matched_keywords,
            raw_message_links=raw_message_links,
            trust_score=80 if not raw_message_links else 65,
            confidence=confidence,
            priority=priority,
            detection_reason=detection_reason,
            status=status,
            last_message_preview=last_message_preview,
        )
        self.session.add(opportunity)
        await self.session.commit()
        await self.session.refresh(opportunity)
        return opportunity

    async def get(self, opportunity_id: UUID) -> Opportunity | None:
        return await self.session.get(Opportunity, opportunity_id)

    async def list(
        self,
        *,
        frontend_status: FrontendOpportunityStatus | None = None,
        channel: IMChannel | None = None,
        limit: int = 100,
        offset: int = 0,
    ) -> list[Opportunity]:
        statement = select(Opportunity)
        if frontend_status:
            statement = statement.where(Opportunity.status.in_(FRONTEND_STATUS_MAP[frontend_status]))
        if channel:
            statement = statement.where(Opportunity.channel == channel)
        statement = statement.order_by(col(Opportunity.last_message_at).desc()).offset(offset).limit(limit)
        result = await self.session.exec(statement)
        return list(result.all())

    async def update_status(
        self,
        opportunity: Opportunity,
        status: OpportunityStatus,
        *,
        final_reply: str | None = None,
        assigned_to: str | None = None,
    ) -> Opportunity:
        opportunity.status = status
        if final_reply is not None:
            opportunity.final_reply = final_reply
            opportunity.last_message_preview = final_reply
            opportunity.last_message_at = utc_now()
        if assigned_to is not None:
            opportunity.assigned_to = assigned_to
        opportunity.updated_at = utc_now()
        self.session.add(opportunity)
        await self.session.commit()
        await self.session.refresh(opportunity)
        return opportunity

    async def save_ai_draft(self, opportunity: Opportunity, draft: str) -> Opportunity:
        opportunity.ai_reply_draft = draft
        opportunity.updated_at = utc_now()
        self.session.add(opportunity)
        await self.session.commit()
        await self.session.refresh(opportunity)
        return opportunity

    async def pending_human_older_than(self, minutes: int) -> list[Opportunity]:
        cutoff = datetime.now(timezone.utc) - timedelta(minutes=minutes)
        statement = select(Opportunity).where(
            Opportunity.status == OpportunityStatus.PENDING_HUMAN,
            Opportunity.created_at <= cutoff,
        )
        result = await self.session.exec(statement)
        return list(result.all())


class RuleRepository:
    def __init__(self, session: AsyncSession) -> None:
        self.session = session

    async def enabled_detection_rules(self) -> list[DetectionRule]:
        statement = select(Rule).where(Rule.enabled.is_(True)).order_by(col(Rule.priority).asc())
        result = await self.session.exec(statement)
        return [
            DetectionRule(
                id=rule.id,
                name=rule.name,
                rule_type=rule.rule_type,
                pattern=rule.pattern,
                score=rule.score,
                priority=rule.priority,
            )
            for rule in result.all()
        ]

    async def list(self) -> list[Rule]:
        result = await self.session.exec(select(Rule).order_by(col(Rule.priority).asc()))
        return list(result.all())

    async def create(self, rule: Rule) -> Rule:
        self.session.add(rule)
        await self.session.commit()
        await self.session.refresh(rule)
        return rule

    async def get(self, rule_id: UUID) -> Rule | None:
        return await self.session.get(Rule, rule_id)

    async def save(self, rule: Rule) -> Rule:
        rule.updated_at = utc_now()
        self.session.add(rule)
        await self.session.commit()
        await self.session.refresh(rule)
        return rule

    async def delete(self, rule: Rule) -> None:
        await self.session.delete(rule)
        await self.session.commit()


class ConfigRepository:
    def __init__(self, session: AsyncSession) -> None:
        self.session = session

    async def get_value(self, key: str) -> dict | None:
        config = await self.session.get(AppConfig, key)
        return config.value if config else None

    async def set_value(self, key: str, value: dict, description: str | None = None) -> AppConfig:
        config = await self.session.get(AppConfig, key)
        if config:
            config.value = value
            config.description = description or config.description
            config.updated_at = utc_now()
        else:
            config = AppConfig(key=key, value=value, description=description)
        self.session.add(config)
        await self.session.commit()
        await self.session.refresh(config)
        return config


class ReplyTemplateRepository:
    def __init__(self, session: AsyncSession) -> None:
        self.session = session

    async def list(self, enabled_only: bool = True) -> list[ReplyTemplate]:
        statement = select(ReplyTemplate).order_by(col(ReplyTemplate.created_at).desc())
        if enabled_only:
            statement = statement.where(ReplyTemplate.enabled.is_(True))
        result = await self.session.exec(statement)
        return list(result.all())

    async def create(self, template: ReplyTemplate) -> ReplyTemplate:
        self.session.add(template)
        await self.session.commit()
        await self.session.refresh(template)
        return template

    async def get(self, template_id: UUID) -> ReplyTemplate | None:
        return await self.session.get(ReplyTemplate, template_id)

    async def save(self, template: ReplyTemplate) -> ReplyTemplate:
        template.updated_at = utc_now()
        self.session.add(template)
        await self.session.commit()
        await self.session.refresh(template)
        return template
