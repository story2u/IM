import asyncio
from datetime import datetime, timezone

from app.domain.enums import IMChannel, MessageDirection, MessageSource, OpportunityStatus, Priority, RuleType
from app.infrastructure.db.models import Message, Opportunity, ReplyTemplate, Rule
from app.infrastructure.db.repositories import ConfigRepository
from app.infrastructure.db.session import AsyncSessionLocal
from sqlmodel import select


RULES = [
    {
        "name": "采购与报价意图",
        "rule_type": RuleType.KEYWORD,
        "pattern": "报价,价格,采购,批量采购,折扣,合同,续约,quote,pricing",
        "score": 0.45,
        "priority": 10,
    },
    {
        "name": "企业能力咨询",
        "rule_type": RuleType.KEYWORD,
        "pattern": "企业版,API,私有化部署,数据安全,等保,SSO,SLA,坐席",
        "score": 0.4,
        "priority": 20,
    },
    {
        "name": "试用转化",
        "rule_type": RuleType.KEYWORD,
        "pattern": "试用,免费试用,demo,trial,开通",
        "score": 0.35,
        "priority": 30,
    },
]

TEMPLATES = [
    {
        "title": "首次咨询欢迎语",
        "content": "您好 {{联系人姓名}}！感谢您的关注。请问您目前主要想解决什么业务问题？我可以为您做针对性介绍。",
        "category": "开场白",
    },
    {
        "title": "企业版功能介绍",
        "content": "{{联系人姓名}} 您好，企业版包含 API 接入、SSO 单点登录、客户成功经理与 SLA 保障。需要我发一份功能对比表吗？",
        "category": "产品介绍",
    },
    {
        "title": "非工作时间自动回复",
        "content": "您好 {{联系人姓名}}，感谢您的消息！现在是非工作时间，我已记录您的需求，工作时间内会第一时间详细答复。",
        "category": "自动回复",
    },
]

OPPORTUNITIES = [
    {
        "channel": IMChannel.TELEGRAM,
        "conversation_id": "10001",
        "customer_external_id": "tg-10001",
        "contact_name": "Michael Chen",
        "title": "企业版 API 接入与批量采购报价",
        "summary": "询问企业版 API 接入方案与批量采购价格，预计团队规模 200 人",
        "matched_keywords": ["API 接入", "批量采购", "企业版"],
        "confidence": 0.94,
        "priority": Priority.URGENT,
        "status": OpportunityStatus.PENDING_HUMAN,
        "last_message_preview": "我们下周需要给管理层一个方案，能尽快发一份报价吗？",
        "messages": [
            ("Michael Chen", "你好，我们是一家做 SaaS 的公司，团队大概 200 人。", True, None),
            ("Michael Chen", "想问一下企业版支持 API 接入吗？批量采购有折扣吗？", True, None),
            ("Michael Chen", "我们下周需要给管理层一个方案，能尽快发一份报价吗？", True, None),
        ],
    },
    {
        "channel": IMChannel.WECOM,
        "conversation_id": "wecom-lina",
        "customer_external_id": "wecom-lina",
        "contact_name": "王丽娜",
        "title": "私有化部署与数据安全合规",
        "summary": "对私有化部署感兴趣，关注数据安全合规认证情况",
        "matched_keywords": ["私有化部署", "数据安全"],
        "confidence": 0.87,
        "priority": Priority.HIGH,
        "status": OpportunityStatus.PENDING_HUMAN,
        "last_message_preview": "你们有等保三级认证吗？私有化部署大概什么周期？",
        "messages": [
            ("王丽娜", "您好，我是华信金融的采购负责人。我们对你们的产品比较感兴趣。", True, None),
            ("王丽娜", "你们有等保三级认证吗？私有化部署大概什么周期？", True, None),
        ],
    },
]


async def main() -> None:
    async with AsyncSessionLocal() as session:
        config_repo = ConfigRepository(session)
        await config_repo.set_value(
            "working_hours",
            {
                "timezone": "Asia/Shanghai",
                "weekdays": [1, 2, 3, 4, 5],
                "start": "09:00",
                "end": "18:30",
                "auto_reply_after_hours": True,
            },
            description="工作时间配置",
        )

        for item in RULES:
            result = await session.exec(select(Rule).where(Rule.name == item["name"]))
            if not result.first():
                session.add(Rule(**item))

        for item in TEMPLATES:
            result = await session.exec(select(ReplyTemplate).where(ReplyTemplate.title == item["title"]))
            if not result.first():
                session.add(ReplyTemplate(**item))

        await session.commit()

        for item in OPPORTUNITIES:
            data = dict(item)
            messages = data.pop("messages")
            result = await session.exec(
                select(Opportunity).where(
                    Opportunity.channel == data["channel"],
                    Opportunity.conversation_id == data["conversation_id"],
                    Opportunity.title == data["title"],
                )
            )
            if result.first():
                continue

            opportunity = Opportunity(**data)
            session.add(opportunity)
            await session.commit()
            await session.refresh(opportunity)

            for index, (sender, text, incoming, source) in enumerate(messages, start=1):
                session.add(
                    Message(
                        channel=opportunity.channel,
                        external_message_id=f"seed-{opportunity.id}-{index}",
                        conversation_id=opportunity.conversation_id,
                        sender_external_id=opportunity.customer_external_id if incoming else None,
                        sender_display_name=sender,
                        direction=MessageDirection.INCOMING if incoming else MessageDirection.OUTGOING,
                        source=source or (None if incoming else MessageSource.HUMAN),
                        text=text,
                        raw_payload={"seed": True},
                        opportunity_id=opportunity.id,
                        sent_at=datetime.now(timezone.utc),
                    )
                )
            await session.commit()


if __name__ == "__main__":
    asyncio.run(main())
