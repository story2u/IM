import re

from app.domain.enums import Priority, RuleType
from app.domain.ports import DetectionResult, DetectionRule, OpportunityAIClassifier


HIGH_INTENT_WORDS = {
    "报价",
    "价格",
    "采购",
    "私有化",
    "部署",
    "企业版",
    "API",
    "合同",
    "续约",
    "折扣",
    "试用",
    "demo",
    "pricing",
    "quote",
    "enterprise",
    "trial",
}


class OpportunityDetector:
    def __init__(self, ai_classifier: OpportunityAIClassifier | None = None) -> None:
        self.ai_classifier = ai_classifier

    async def detect(self, text: str, rules: list[DetectionRule]) -> DetectionResult:
        normalized = text.strip()
        if not normalized:
            return DetectionResult(is_opportunity=False)

        score = 0.0
        reasons: list[str] = []
        matched_keywords: list[str] = []

        for rule in sorted(rules, key=lambda item: item.priority):
            matched = self._match_rule(rule, normalized)
            if not matched:
                continue

            score = min(1.0, score + rule.score)
            reasons.append(f"{rule.name}:{matched}")
            if rule.rule_type in {RuleType.KEYWORD, RuleType.AI_HINT}:
                matched_keywords.append(matched)

        for word in HIGH_INTENT_WORDS:
            if word.lower() in normalized.lower() and word not in matched_keywords:
                score = min(1.0, score + 0.12)
                matched_keywords.append(word)

        if score >= 0.75:
            return self._build_positive_result(normalized, score, reasons, matched_keywords)

        if score < 0.35 or not self.ai_classifier:
            return DetectionResult(
                is_opportunity=score >= 0.45,
                confidence=score,
                title=self._title(normalized) if score >= 0.45 else None,
                summary=normalized if score >= 0.45 else None,
                reason="; ".join(reasons) if reasons else None,
                matched_keywords=matched_keywords,
                priority=self._priority(score, matched_keywords),
            )

        ai_result = await self.ai_classifier.classify(text=normalized, rule_score=score)
        if matched_keywords:
            ai_result.matched_keywords = sorted(
                {*ai_result.matched_keywords, *matched_keywords},
                key=lambda item: normalized.lower().find(item.lower())
                if item.lower() in normalized.lower()
                else 999,
            )
        return ai_result

    def _match_rule(self, rule: DetectionRule, text: str) -> str | None:
        if rule.rule_type in {RuleType.KEYWORD, RuleType.AI_HINT}:
            candidates = [item.strip() for item in rule.pattern.split(",") if item.strip()]
            for candidate in candidates:
                if candidate.lower() in text.lower():
                    return candidate
            return None

        if rule.rule_type == RuleType.REGEX:
            match = re.search(rule.pattern, text, flags=re.IGNORECASE)
            return match.group(0) if match else None

        return None

    def _build_positive_result(
        self,
        text: str,
        score: float,
        reasons: list[str],
        matched_keywords: list[str],
    ) -> DetectionResult:
        return DetectionResult(
            is_opportunity=True,
            confidence=score,
            title=self._title(text),
            summary=text[:240],
            reason="; ".join(reasons) if reasons else "high intent keyword match",
            matched_keywords=matched_keywords[:8],
            priority=self._priority(score, matched_keywords),
        )

    def _priority(self, score: float, matched_keywords: list[str]) -> Priority:
        urgent_words = {"报价", "采购", "合同", "urgent", "quote"}
        if score >= 0.92 or urgent_words.intersection({item.lower() for item in matched_keywords}):
            return Priority.URGENT
        if score >= 0.8:
            return Priority.HIGH
        if score >= 0.55:
            return Priority.NORMAL
        return Priority.LOW

    def _title(self, text: str) -> str:
        return text[:42] + ("..." if len(text) > 42 else "")
