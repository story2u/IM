from uuid import uuid4

import pytest

from app.domain.enums import RuleType
from app.domain.ports import DetectionRule
from app.domain.services.detection_policy import OpportunityDetector


@pytest.mark.asyncio
async def test_detector_marks_high_intent_message_as_opportunity() -> None:
    detector = OpportunityDetector()
    result = await detector.detect(
        "我们想了解企业版 API 接入和 200 人批量采购报价",
        [
            DetectionRule(
                id=uuid4(),
                name="报价",
                rule_type=RuleType.KEYWORD,
                pattern="报价,批量采购",
                score=0.5,
                priority=10,
            ),
            DetectionRule(
                id=uuid4(),
                name="企业能力",
                rule_type=RuleType.KEYWORD,
                pattern="企业版,API",
                score=0.4,
                priority=20,
            ),
        ],
    )

    assert result.is_opportunity
    assert result.confidence >= 0.75
    assert "报价" in result.matched_keywords


@pytest.mark.asyncio
async def test_detector_marks_recruiting_group_message_as_opportunity() -> None:
    detector = OpportunityDetector()
    result = await detector.detect(
        "招聘 Python 后端工程师，远程全职，薪资 25k-35k，简历发 @hr_jobs",
        [],
    )

    assert result.is_opportunity
    assert result.confidence >= 0.75
    assert "招聘" in result.matched_keywords
