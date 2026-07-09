import json
from uuid import UUID

from langchain_core.messages import HumanMessage, SystemMessage
from langchain_litellm import ChatLiteLLM

from app.core.config import Settings
from app.domain.enums import Priority
from app.domain.ports import DetectionResult
from app.infrastructure.db.repositories import MessageRepository, OpportunityRepository


class LiteLLMOpportunityClassifier:
    def __init__(self, settings: Settings) -> None:
        self.settings = settings

    async def classify(self, text: str, rule_score: float) -> DetectionResult:
        if not self.settings.ai_enabled:
            return DetectionResult(
                is_opportunity=rule_score >= 0.45,
                confidence=rule_score,
                title=text[:42],
                summary=text[:240],
                reason="rule score fallback",
                priority=Priority.NORMAL,
            )

        llm = ChatLiteLLM(model=self.settings.litellm_model, temperature=0.1)
        response = await llm.ainvoke(
            [
                SystemMessage(
                    content=(
                        "你是B2B商机识别器。只输出JSON，字段："
                        "is_opportunity(bool), confidence(0-1), title, summary, "
                        "matched_keywords(array), priority(low|normal|high|urgent), reason。"
                    )
                ),
                HumanMessage(content=f"规则分={rule_score}\n客户消息：{text}"),
            ]
        )
        try:
            data = self._loads_json(str(response.content))
        except (json.JSONDecodeError, TypeError, ValueError):
            return DetectionResult(
                is_opportunity=rule_score >= 0.45,
                confidence=rule_score,
                title=text[:42],
                summary=text[:240],
                reason="ai classifier returned non-json; rule score fallback",
                priority=Priority.NORMAL,
            )
        return DetectionResult(
            is_opportunity=bool(data.get("is_opportunity")),
            confidence=float(data.get("confidence", rule_score)),
            title=data.get("title") or text[:42],
            summary=data.get("summary") or text[:240],
            reason=data.get("reason") or "ai classifier",
            matched_keywords=list(data.get("matched_keywords") or []),
            priority=Priority(data.get("priority") or Priority.NORMAL),
        )

    def _loads_json(self, content: str) -> dict:
        content = content.strip()
        if content.startswith("```"):
            content = content.strip("`")
            content = content.removeprefix("json").strip()
        return json.loads(content)


class LiteLLMReplyGenerator:
    def __init__(
        self,
        *,
        settings: Settings,
        opportunity_repo: OpportunityRepository,
        message_repo: MessageRepository,
    ) -> None:
        self.settings = settings
        self.opportunity_repo = opportunity_repo
        self.message_repo = message_repo

    async def generate_reply(self, opportunity_id: UUID) -> str:
        opportunity = await self.opportunity_repo.get(opportunity_id)
        if not opportunity:
            raise ValueError("opportunity not found")

        messages = await self.message_repo.list_by_conversation(
            opportunity.channel,
            opportunity.conversation_id,
            limit=12,
        )
        history = "\n".join(
            f"{message.sender_display_name or '客户'}: {message.text}"
            for message in messages
            if message.text
        )

        if not self.settings.ai_enabled:
            keywords = "、".join(opportunity.matched_keywords) or "您的需求"
            return (
                f"{opportunity.contact_name} 您好！关于您提到的「{keywords}」，"
                "我们可以为您整理一份针对性的方案。方便的话，我想先确认一下使用场景、"
                "团队规模和期望上线时间，然后安排顾问进一步沟通。"
            )

        llm = ChatLiteLLM(model=self.settings.litellm_model, temperature=0.3)
        response = await llm.ainvoke(
            [
                SystemMessage(
                    content=(
                        "你是B2B商机助手。回复要自然、专业、简洁。"
                        "不要承诺最低价、合同条款、绝对交付结果。"
                        "目标是确认需求并推动下一步沟通。"
                    )
                ),
                HumanMessage(
                    content=(
                        f"联系人：{opportunity.contact_name}\n"
                        f"商机摘要：{opportunity.summary}\n"
                        f"关键词：{opportunity.matched_keywords}\n"
                        f"聊天历史：\n{history}\n\n"
                        "请生成一条可直接发送的回复，长度控制在120字内。"
                    )
                ),
            ]
        )
        return str(response.content).strip()
