from enum import StrEnum


class IMChannel(StrEnum):
    TELEGRAM = "telegram"
    WECOM = "wecom"


class MessageDirection(StrEnum):
    INCOMING = "incoming"
    OUTGOING = "outgoing"


class MessageSource(StrEnum):
    HUMAN = "human"
    AI = "ai"


class OpportunityStatus(StrEnum):
    PENDING_HUMAN = "pending_human"
    AI_AUTO_REPLY = "ai_auto_reply"
    REPLIED = "replied"
    FOLLOWING = "following"
    IGNORED = "ignored"
    CLOSED = "closed"


class FrontendOpportunityStatus(StrEnum):
    PENDING = "pending"
    REPLIED = "replied"
    IGNORED = "ignored"


class Priority(StrEnum):
    LOW = "low"
    NORMAL = "normal"
    HIGH = "high"
    URGENT = "urgent"


class RuleType(StrEnum):
    KEYWORD = "keyword"
    REGEX = "regex"
    AI_HINT = "ai_hint"
