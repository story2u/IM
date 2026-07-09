from app.application.dto import StatsSummaryRead
from app.application.mappers import frontend_status
from app.domain.enums import FrontendOpportunityStatus
from app.infrastructure.db.models import Opportunity


def build_stats_summary(opportunities: list[Opportunity]) -> StatsSummaryRead:
    total = len(opportunities)
    pending = 0
    replied = 0
    ignored = 0
    confidence_sum = 0.0

    for opportunity in opportunities:
        status = frontend_status(opportunity.status)
        pending += status == FrontendOpportunityStatus.PENDING
        replied += status == FrontendOpportunityStatus.REPLIED
        ignored += status == FrontendOpportunityStatus.IGNORED
        confidence_sum += opportunity.confidence

    return StatsSummaryRead(
        total=total,
        pending=pending,
        replied=replied,
        ignored=ignored,
        avgConfidence=round(confidence_sum / total, 4) if total else 0.0,
    )
