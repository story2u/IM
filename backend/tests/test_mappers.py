from app.application.mappers import frontend_status
from app.domain.enums import FrontendOpportunityStatus, OpportunityStatus


def test_frontend_status_mapping() -> None:
    assert frontend_status(OpportunityStatus.PENDING_HUMAN) == FrontendOpportunityStatus.PENDING
    assert frontend_status(OpportunityStatus.AI_AUTO_REPLY) == FrontendOpportunityStatus.PENDING
    assert frontend_status(OpportunityStatus.REPLIED) == FrontendOpportunityStatus.REPLIED
    assert frontend_status(OpportunityStatus.IGNORED) == FrontendOpportunityStatus.IGNORED
